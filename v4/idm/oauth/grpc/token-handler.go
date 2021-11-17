package grpc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
	"time"

	"github.com/ory/fosite/token/hmac"
	"github.com/ory/fosite/token/jwt"
	"go.uber.org/zap"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/proto/auth"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/common/service/errors"
	"github.com/pydio/cells/v4/common/utils/permissions"
	"github.com/pydio/cells/v4/common/utils/uuid"
	"github.com/pydio/cells/v4/idm/oauth"
	json "github.com/pydio/cells/v4/x/jsonx"
)

var tokensKey []byte

type PatScopeClaims struct {
	Scopes []string `json:"scopes"`
}

type PatHandler struct {
	strategy *hmac.HMACStrategy
}

func (p *PatHandler) getDao(ctx context.Context) oauth.DAO {
	return servicecontext.GetDAO(ctx).(oauth.DAO)
}

func (p *PatHandler) getStrategy() *hmac.HMACStrategy {
	if p.strategy == nil {
		p.strategy = &hmac.HMACStrategy{
			TokenEntropy:         32,
			GlobalSecret:         p.getKey(),
			RotatedGlobalSecrets: nil,
			Mutex:                sync.Mutex{},
		}
	}
	return p.strategy
}

func (p *PatHandler) getKey() []byte {

	if len(tokensKey) > 0 {
		return tokensKey
	}

	cVal := config.Get("defaults", "personalTokens", "secureKey")
	if cVal.String() == "" {
		tokensKey = p.generateRandomKey(32)
		strKey := base64.StdEncoding.EncodeToString(tokensKey)
		cVal.Set(strKey)
		config.Save(common.PydioSystemUsername, "Creating random key for personal tokens service")
	} else if t, e := base64.StdEncoding.DecodeString(cVal.String()); e == nil {
		tokensKey = t
	} else {
		log.Logger(context.Background()).Error("Could not read generated key for personal tokens!", zap.Error(e))
	}
	return tokensKey
}

func (p *PatHandler) Verify(ctx context.Context, request *auth.VerifyTokenRequest, response *auth.VerifyTokenResponse) error {
	dao := p.getDao(ctx)
	if err := p.getStrategy().Validate(request.Token); err != nil {
		return errors.Unauthorized("token.invalid", "Cannot validate token")
	}
	pat, e := dao.Load(request.Token)
	if e != nil {
		return errors.Unauthorized("token.not.found", "Cannot find corresponding Personal Access Token")
	}
	// Check Expiration Date
	if time.Unix(pat.ExpiresAt, 0).Before(time.Now()) {
		return errors.Unauthorized("token.expired", "Personal token is expired")
	}
	if pat.AutoRefreshWindow > 0 {
		// Recompute expire date
		pat.ExpiresAt = time.Now().Add(time.Duration(pat.AutoRefreshWindow) * time.Second).Unix()
		if er := dao.Store(request.Token, pat, true); er != nil {
			return errors.BadRequest("internal.error", "Cannot store updated token "+er.Error())
		}
	}

	cl := jwt.IDTokenClaims{
		Subject:   pat.UserUuid,
		Issuer:    "local",
		ExpiresAt: time.Unix(pat.ExpiresAt, 0),
		Audience:  []string{common.ServiceGrpcNamespace_ + common.ServiceToken},
	}
	if len(pat.Scopes) > 0 {
		cl.Extra = map[string]interface{}{
			"scopes": pat.Scopes,
		}
	}
	m, _ := json.Marshal(cl)
	response.Success = true
	response.Data = m
	return nil
}

func (p *PatHandler) Generate(ctx context.Context, request *auth.PatGenerateRequest, response *auth.PatGenerateResponse) error {
	dao := p.getDao(ctx)
	token := &auth.PersonalAccessToken{
		Uuid:              uuid.New(),
		Type:              request.Type,
		Label:             request.Label,
		UserUuid:          request.UserUuid,
		UserLogin:         request.UserLogin,
		Scopes:            request.Scopes,
		AutoRefreshWindow: request.AutoRefreshWindow,
		ExpiresAt:         request.ExpiresAt,
	}
	if request.AutoRefreshWindow > 0 {
		request.ExpiresAt = time.Now().Add(time.Duration(request.AutoRefreshWindow) * time.Second).Unix()
		token.ExpiresAt = request.ExpiresAt
	} else if request.ExpiresAt > 0 {
		token.ExpiresAt = request.ExpiresAt
	} else {
		return errors.BadRequest("missing.parameters", "Please provide one of ExpiresAt or AutoRefreshWindow")
	}
	token.CreatedAt = time.Now().Unix()
	if uName, _ := permissions.FindUserNameInContext(ctx); uName != "" {
		token.CreatedBy = uName
	}
	accessToken, _, err := p.getStrategy().Generate()
	if err != nil {
		return err
	}
	if err := dao.Store(accessToken, token, false); err != nil {
		return err
	}
	response.TokenUuid = token.Uuid
	response.AccessToken = accessToken
	return nil
}

func (p *PatHandler) Revoke(ctx context.Context, request *auth.PatRevokeRequest, response *auth.PatRevokeResponse) error {
	dao := p.getDao(ctx)
	return dao.Delete(request.GetUuid())
}

func (p *PatHandler) List(ctx context.Context, request *auth.PatListRequest, response *auth.PatListResponse) error {
	dao := p.getDao(ctx)
	tt, er := dao.List(request.Type, request.ByUserLogin)
	if er != nil {
		return nil
	}
	response.Tokens = tt
	return nil
}

func (p *PatHandler) PruneTokens(ctx context.Context, request *auth.PruneTokensRequest, response *auth.PruneTokensResponse) error {
	i, e := p.getDao(ctx).PruneExpired()
	if e != nil {
		return e
	}
	response.Count = int32(i)
	return nil
}

func (p *PatHandler) generateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
