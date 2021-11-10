/*
 * Copyright (c) 2018. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package grpc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"go.uber.org/zap"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/crypto"
	"github.com/pydio/cells/v4/common/log"
	enc "github.com/pydio/cells/v4/common/proto/encryption"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/idm/key"
)

type userKeyStore struct{}

// NewUserKeyStore creates a master password based
func NewUserKeyStore() (enc.UserKeyStoreHandler, error) {
	return &userKeyStore{}, nil
}

func (ukm *userKeyStore) getDAO(ctx context.Context) (key.DAO, error) {
	dao := servicecontext.GetDAO(ctx)
	if dao == nil {
		return nil, errors.InternalServerError(common.ServiceUserKey, "No DAO found Wrong initialization")
	}

	keyDao, ok := dao.(key.DAO)
	if !ok {
		return nil, errors.New(common.ServiceUserKey, "unsupported dao", 500)
	}
	return keyDao, nil
}

func (ukm *userKeyStore) AddKey(ctx context.Context, req *enc.AddKeyRequest, rsp *enc.AddKeyResponse) error {
	dao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}

	err = seal(req.Key, []byte(req.StrPassword))
	if err != nil {
		return err
	}

	return dao.SaveKey(req.Key)
}

func (ukm *userKeyStore) GetKey(ctx context.Context, req *enc.GetKeyRequest, rsp *enc.GetKeyResponse) error {

	dao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}

	// TODO: Extract user / password info from Context
	user := common.PydioSystemUsername
	pwd := config.Vault().Val("masterPassword").Bytes()

	rsp.Key, err = dao.GetKey(user, req.KeyID)
	if err != nil {
		return err
	}

	if rsp.Key == nil {
		return nil
	}
	return open(rsp.Key, pwd)
}

func (ukm *userKeyStore) AdminListKeys(ctx context.Context, req *enc.AdminListKeysRequest, rsp *enc.AdminListKeysResponse) error {

	keyDao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}

	rsp.Keys, err = keyDao.ListKeys(common.PydioSystemUsername)
	return err
}

func (ukm *userKeyStore) AdminCreateKey(ctx context.Context, req *enc.AdminCreateKeyRequest, rsp *enc.AdminCreateKeyResponse) error {

	keyDao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}

	if _, err := keyDao.GetKey(common.PydioSystemUsername, req.KeyID); err != nil {
		if errors.Parse(err.Error()).Code == 404 {
			return createSystemKey(keyDao, req.KeyID, req.Label)
		} else {
			return err
		}
	} else {
		return errors.BadRequest(common.ServiceEncKey, "Key already exists with this id!")
	}
}

func (ukm *userKeyStore) AdminDeleteKey(ctx context.Context, req *enc.AdminDeleteKeyRequest, rsp *enc.AdminDeleteKeyResponse) error {

	dao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}
	return dao.DeleteKey(common.PydioSystemUsername, req.KeyID)
}

func (ukm *userKeyStore) AdminImportKey(ctx context.Context, req *enc.AdminImportKeyRequest, rsp *enc.AdminImportKeyResponse) error {

	dao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}

	log.Logger(ctx).Debug("Received request", zap.Any("Data", req))

	var k *enc.Key
	k, err = dao.GetKey(common.PydioSystemUsername, req.Key.ID)
	if err != nil {
		if errors.Parse(err.Error()).Code != 404 {
			return err
		}
	} else if k != nil && !req.Override {
		return errors.BadRequest(common.ServiceEncKey, fmt.Sprintf("Key already exists with [%s] id", req.Key.ID))
	}

	log.Logger(ctx).Debug("Opening sealed key with imported password")
	err = open(req.Key, []byte(req.StrPassword))
	if err != nil {
		rsp.Success = false
		return errors.InternalServerError(common.ServiceEncKey, "unable to decrypt %s for import, cause: %s", req.Key.ID, err.Error())
	}

	log.Logger(ctx).Debug("Sealing with master key")
	err = sealWithMasterKey(req.Key)
	if err != nil {
		rsp.Success = false
		return errors.InternalServerError(common.ServiceEncKey, "unable to encrypt %s.%s for export, cause: %s", common.PydioSystemUsername, req.Key.ID, err.Error())
	}

	if req.Key.CreationDate == 0 {
		if k != nil {
			req.Key.CreationDate = k.CreationDate
		} else {
			req.Key.CreationDate = int32(time.Now().Unix())
		}
	}

	if len(req.Key.Owner) == 0 {
		req.Key.Owner = common.PydioSystemUsername
	}

	log.Logger(ctx).Debug("Received request", zap.Any("Data", req))

	log.Logger(ctx).Debug("Adding import info")
	// We set import info

	if req.Key.Info == nil {
		req.Key.Info = &enc.KeyInfo{}
	}

	if req.Key.Info.Imports == nil {
		req.Key.Info.Imports = []*enc.Import{}
	}

	req.Key.Info.Imports = append(req.Key.Info.Imports, &enc.Import{
		By:   common.PydioSystemUsername,
		Date: int32(time.Now().Unix()),
	})

	log.Logger(ctx).Debug("Saving new key")
	err = dao.SaveKey(req.Key)
	if err != nil {
		rsp.Success = false
		return errors.InternalServerError(common.ServiceEncKey, "failed to save imported key, cause: %s", err.Error())
	}

	log.Logger(ctx).Debug("Returning response")
	rsp.Success = true
	return nil
}

func (ukm *userKeyStore) AdminExportKey(ctx context.Context, req *enc.AdminExportKeyRequest, rsp *enc.AdminExportKeyResponse) error {

	//Get key from dao
	dao, err := ukm.getDAO(ctx)
	if err != nil {
		return err
	}

	rsp.Key, err = dao.GetKey(common.PydioSystemUsername, req.KeyID)
	if err != nil {
		return err
	}
	log.Logger(ctx).Debug(fmt.Sprintf("Exporting key %s", rsp.Key.Content))

	// We set export info
	if rsp.Key.Info.Exports == nil {
		rsp.Key.Info.Exports = []*enc.Export{}
	}
	rsp.Key.Info.Exports = append(rsp.Key.Info.Exports, &enc.Export{
		By:   common.PydioSystemUsername,
		Date: int32(time.Now().Unix()),
	})

	// We update the key
	err = dao.SaveKey(rsp.Key)
	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "failed to update key info, cause: %s", err.Error())
	}

	err = openWithMasterKey(rsp.Key)
	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "unable to decrypt for %s with key %s, cause: %s", common.PydioSystemUsername, req.KeyID, err)
	}

	err = seal(rsp.Key, []byte(req.StrPassword))
	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "unable to encrypt for %s with key %s for export, cause: %s", common.PydioSystemUsername, req.KeyID, err)
	}
	return nil
}

// Create a default key or create a system key with a given ID
func createSystemKey(dao key.DAO, keyID string, keyLabel string) error {
	systemKey := &enc.Key{
		ID:           keyID,
		Owner:        common.PydioSystemUsername,
		Label:        keyLabel,
		CreationDate: int32(time.Now().Unix()),
	}

	keyContentBytes := make([]byte, 32)
	_, err := rand.Read(keyContentBytes)
	if err != nil {
		return err
	}

	masterPasswordBytes, err := getMasterPassword()
	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "failed to get password. Make sure you have the system keyring installed. Cause: %s", err.Error())
	}

	masterKey := crypto.KeyFromPassword(masterPasswordBytes, 32)
	encryptedKeyContentBytes, err := crypto.Seal(masterKey, keyContentBytes)
	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "failed to encrypt the default key. Cause: %s", err.Error())
	}
	systemKey.Content = base64.StdEncoding.EncodeToString(encryptedKeyContentBytes)
	log.Logger(context.Background()).Debug(fmt.Sprintf("Saving default key %s", systemKey.Content))
	return dao.SaveKey(systemKey)
}

func sealWithMasterKey(k *enc.Key) error {
	masterPasswordBytes, err := getMasterPassword()
	if len(masterPasswordBytes) == 0 {
		return errors.InternalServerError(common.ServiceEncKey, "failed to get %s password, cause: %s", common.PydioSystemUsername, err.Error())
	}
	return seal(k, masterPasswordBytes)
}

func openWithMasterKey(k *enc.Key) error {
	masterPasswordBytes, err := getMasterPassword()
	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "failed to get %s password, cause: %s", common.PydioSystemUsername, err.Error())
	}
	return open(k, masterPasswordBytes)
}

func seal(k *enc.Key, passwordBytes []byte) error {
	keyContentBytes, err := base64.StdEncoding.DecodeString(k.Content)
	if err != nil {
		return errors.New(common.ServiceUserKey, "unable to decode key", 400)
	}

	passwordKey := crypto.KeyFromPassword(passwordBytes, 32)
	encryptedKeyContentBytes, err := crypto.Seal(passwordKey, keyContentBytes)

	if err != nil {
		return errors.InternalServerError(common.ServiceEncKey, "failed to encrypt the default key, cause: %s", err.Error())
	}
	k.Content = base64.StdEncoding.EncodeToString(encryptedKeyContentBytes)
	return nil
}

func open(k *enc.Key, passwordBytes []byte) error {
	sealedContentBytes, err := base64.StdEncoding.DecodeString(k.Content)
	if err != nil {
		return errors.New(common.ServiceUserKey, "unable to decode key", 400)
	}

	passwordKey := crypto.KeyFromPassword(passwordBytes, 32)

	nonce := sealedContentBytes[:12]
	keySealContentBytes := sealedContentBytes[12:]

	keyPlainContentBytes, err := crypto.Open(passwordKey, nonce, keySealContentBytes)
	if err != nil {
		return err
	}

	k.Content = base64.StdEncoding.EncodeToString(keyPlainContentBytes)
	return nil
}

func getMasterPassword() ([]byte, error) {
	var masterPasswordBytes []byte
	masterPassword := config.Vault().Val("masterPassword").String()
	if masterPassword == "" {
		return masterPasswordBytes, errors.InternalServerError("master.key.load", "cannot get master password")
	}
	masterPasswordBytes = []byte(masterPassword)
	return masterPasswordBytes, nil
}
