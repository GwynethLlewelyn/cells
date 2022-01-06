package hooks

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7/pkg/signer"
	"github.com/minio/minio/cmd"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/auth"
	"github.com/pydio/cells/v4/common/auth/claim"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/common/service/context/metadata"
	"github.com/pydio/cells/v4/common/utils/permissions"
)

// authHandler - handles all the incoming authorization headers and validates them if possible.
type pydioAuthHandler struct {
	handler         http.Handler
	jwtVerifier     *auth.JWTVerifier
	globalAccessKey string
}

// GetPydioAuthHandlerFunc validates Pydio authorization headers for the incoming request.
func GetPydioAuthHandlerFunc(globalAccessKey string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return pydioAuthHandler{
			handler:         h,
			jwtVerifier:     auth.DefaultJWTVerifier(),
			globalAccessKey: globalAccessKey,
		}
	}
}

// handler for validating incoming authorization headers.
func (a pydioAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//var md map[string]string
	var userName string
	ctx := r.Context()
	ctx = servicecontext.HttpRequestInfoToMetadata(ctx, r)
	ctx = servicecontext.WithServiceName(ctx, common.ServiceGatewayData)

	resignRequestV4 := false
	jwt := r.URL.Query().Get("pydio_jwt")

	if len(jwt) > 0 {
		//logger.Info("Found JWT in URL: replace by header and remove from URL")
		r.Header.Set("X-Pydio-Bearer", jwt)
		r.URL.RawQuery = strings.Replace(r.URL.RawQuery, "&pydio_jwt="+jwt, "", 1)

	} else if bearer, ok := r.Header["X-Pydio-Bearer"]; !ok || len(bearer) == 0 {
		// Copy request.
		req := *r
		// Save authorization header.
		v4Auth := req.Header.Get("Authorization")
		// Parse signature version '4' header.
		signedKey, err := cmd.ExposedParseSignV4(v4Auth)
		if err == nil {
			if signedKey != a.globalAccessKey {
				log.Logger(ctx).Debug("Use AWS Api Key as JWT: " + signedKey)
				resignRequestV4 = true
				r.Header.Set("X-Pydio-Bearer", signedKey)
			}
		}
	}

	if bearer, ok := r.Header["X-Pydio-Bearer"]; ok && len(bearer) > 0 {

		rawIDToken := strings.Join(bearer, "")
		var err error
		var claims claim.Claims
		ctx, claims, err = a.jwtVerifier.Verify(ctx, rawIDToken)
		if err != nil {
			cmd.ExposedWriteErrorResponse(ctx, w, cmd.ErrAccessDenied, r.URL)
			return
		}
		userName = claims.Name
		if resignRequestV4 {
			// User is OK, override signature with service account ID/Secret
			r = signer.SignV4(*r, common.S3GatewayRootUser, common.S3GatewayRootPassword, "", common.S3GatewayDefaultRegion)
		}

	} else if agent, aOk := r.Header["User-Agent"]; aOk && strings.Contains(strings.Join(agent, ""), "pydio.sync.client.s3") {

		userName = common.PydioSystemUsername

	} else {

		if user, er := permissions.SearchUniqueUser(ctx, common.PydioS3AnonUsername, ""); er == nil {
			userName = common.PydioS3AnonUsername
			var s []string
			for _, role := range user.Roles {
				if role.UserRole { // Just append the User Role
					s = append(s, role.Uuid)
				}
			}
			anonClaim := claim.Claims{
				Name:      common.PydioS3AnonUsername,
				Roles:     strings.Join(s, ","),
				Profile:   "anon",
				GroupPath: "/",
			}
			ctx = context.WithValue(ctx, claim.ContextKey, anonClaim)

		} else {
			cmd.ExposedWriteErrorResponse(ctx, w, cmd.ErrAccessDenied, r.URL)
			return
		}

	}

	ctx = metadata.WithUserNameMetadata(ctx, userName)
	newRequest := r.WithContext(ctx)
	a.handler.ServeHTTP(w, newRequest)

}
