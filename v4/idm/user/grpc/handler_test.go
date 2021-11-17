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
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/proto/idm"
	"github.com/pydio/cells/v4/common/proto/service"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/common/sql"
	"github.com/pydio/cells/v4/common/utils/cache"
	"github.com/pydio/cells/v4/idm/user"
	"github.com/pydio/cells/v4/x/configx"
)

var (
	ctx context.Context
	wg  sync.WaitGroup
)

func TestMain(m *testing.M) {

	// Use the cache mechanism to avoid trying to retrieve the role service
	autoAppliesCache = cache.NewShort(cache.WithEviction(3600*time.Second), cache.WithCleanWindow(7200*time.Second))
	autoAppliesCache.Set("autoApplies", map[string][]*idm.Role{
		"autoApplyProfile": {{Uuid: "auto-apply", AutoApplies: []string{"autoApplyProfile"}}},
	})

	sqlDao := sql.NewDAO("sqlite3", "file::memory:?mode=memory&cache=shared", "idm_user")
	if sqlDao == nil {
		log.Fatal("unable to open sqlite3 DB file, could not start test")
		return
	}

	mockDAO := user.NewDAO(sqlDao)
	var options = configx.New()
	if err := mockDAO.Init(options); err != nil {
		log.Fatal("could not start test: unable to initialise DAO, error: ", err)
		return
	}

	ctx = servicecontext.WithDAO(context.Background(), mockDAO)

	m.Run()
	wg.Wait()
}

func TestUser(t *testing.T) {

	h := new(Handler)

	Convey("Create one user", t, func() {
		resp, err := h.CreateUser(ctx, &idm.CreateUserRequest{User: &idm.User{Login: "user1"}})

		So(err, ShouldBeNil)
		So(resp.GetUser().GetLogin(), ShouldEqual, "user1")
	})

	Convey("Create a second user with name attribute", t, func() {
		resp, err := h.CreateUser(ctx, &idm.CreateUserRequest{User: &idm.User{Login: "user2", Attributes: map[string]string{"name": "User 2"}}})

		So(err, ShouldBeNil)
		So(resp.GetUser().GetLogin(), ShouldEqual, "user2")
	})

	Convey("Get User", t, func() {
		mock := &userStreamMock{ctx: ctx}
		err := h.StreamUser(mock)

		So(err, ShouldBeNil)
		So(len(mock.InternalBuffer), ShouldEqual, 0)
	})

	Convey("Search User", t, func() {
		mock := &userStreamMock{ctx: ctx}
		userQuery := &idm.UserSingleQuery{
			Login: "user1",
		}
		userQueryAny, _ := anypb.New(userQuery)
		request := &idm.SearchUserRequest{
			Query: &service.Query{
				SubQueries: []*anypb.Any{userQueryAny},
			},
		}
		err := h.SearchUser(request, mock)

		So(err, ShouldBeNil)
		So(len(mock.InternalBuffer), ShouldEqual, 1)

		resp, err := h.CountUser(ctx, request)
		So(err, ShouldBeNil)
		So(resp.Count, ShouldEqual, 1)
	})

	Convey("Create a user with auto apply role", t, func() {
		resp, err := h.CreateUser(ctx, &idm.CreateUserRequest{User: &idm.User{Login: "user3", Attributes: map[string]string{"profile": "autoApplyProfile"}}})

		So(err, ShouldBeNil)
		So(resp.GetUser().GetLogin(), ShouldEqual, "user3")
		mock := &userStreamMock{ctx: ctx}
		userQuery := &idm.UserSingleQuery{
			Login: "user3",
		}
		userQueryAny, _ := anypb.New(userQuery)
		request := &idm.SearchUserRequest{
			Query: &service.Query{
				SubQueries: []*anypb.Any{userQueryAny},
			},
		}
		err = h.SearchUser(request, mock)
		So(err, ShouldBeNil)
		So(len(mock.InternalBuffer), ShouldEqual, 1)
		fmt.Println(mock.InternalBuffer[0].Roles)
		So(mock.InternalBuffer[0].Roles, ShouldHaveLength, 3)
		So(mock.InternalBuffer[0].Roles[0].Uuid, ShouldEqual, "ROOT_GROUP")
		So(mock.InternalBuffer[0].Roles[0].GroupRole, ShouldBeTrue)
		So(mock.InternalBuffer[0].Roles[0].UserRole, ShouldBeFalse)

		So(mock.InternalBuffer[0].Roles[1].Uuid, ShouldEqual, "auto-apply")

		So(mock.InternalBuffer[0].Roles[2].Uuid, ShouldEqual, mock.InternalBuffer[0].Uuid)
		So(mock.InternalBuffer[0].Roles[2].UserRole, ShouldBeTrue)
		So(mock.InternalBuffer[0].Roles[2].GroupRole, ShouldBeFalse)

	})

	Convey("Del User", t, func() {
		_, err := h.DeleteUser(ctx, &idm.DeleteUserRequest{})
		So(err, ShouldNotBeNil)
	})

	Convey("Del User", t, func() {
		singleQ1 := new(idm.UserSingleQuery)
		singleQ1.Login = "user1"
		singleQ1Any, err := anypb.New(singleQ1)
		So(err, ShouldBeNil)

		query := &service.Query{
			SubQueries: []*anypb.Any{singleQ1Any},
		}

		_, err = h.DeleteUser(ctx, &idm.DeleteUserRequest{Query: query})
		So(err, ShouldBeNil)
	})

	Convey("List all Users", t, func() {
		mock := &userStreamMock{ctx: ctx}
		err := h.SearchUser(&idm.SearchUserRequest{}, mock)

		So(err, ShouldBeNil)
		So(len(mock.InternalBuffer), ShouldEqual, 3)
	})

	Convey("Create and bind user", t, func() {
		resp := new(idm.CreateUserResponse)
		resp, err := h.CreateUser(ctx, &idm.CreateUserRequest{User: &idm.User{Login: "john", Password: "f00"}})

		So(err, ShouldBeNil)
		So(resp.GetUser().GetLogin(), ShouldEqual, "john")

		// Correct bind
		{
			bindResp, err := h.BindUser(ctx, &idm.BindUserRequest{UserName: "john", Password: "f00"})
			So(err, ShouldBeNil)
			So(bindResp.User, ShouldNotBeNil)
		}
		// Wrong user name
		{
			_, err := h.BindUser(ctx, &idm.BindUserRequest{UserName: "johnFAIL", Password: "f00"})
			So(err, ShouldNotBeNil)
		}
		// Wrong Password
		{
			_, err = h.BindUser(ctx, &idm.BindUserRequest{UserName: "john", Password: "f00FAIL"})
			So(err, ShouldNotBeNil)
		}

	})

	Convey("Create and bind user with a legacy password", t, func() {
		attributes := make(map[string]string, 1)
		attributes[idm.UserAttrPassHashed] = "true"
		resp, err := h.CreateUser(ctx, &idm.CreateUserRequest{User: &idm.User{
			Login:      "legacy",
			Password:   "sha256:1000:ojGf2O2ELslNab+PZ/CkbVddgHQnSwx/FcMZ0Pa4+EE=:d6mVgF+fS+wg7X+0lmSn1T1IOU7DLZhz",
			Attributes: attributes,
		}})

		So(err, ShouldBeNil)
		So(resp.GetUser().GetLogin(), ShouldEqual, "legacy")

		bindResp, err := h.BindUser(ctx, &idm.BindUserRequest{UserName: "legacy", Password: "P@ssw0rd"})
		So(err, ShouldBeNil)
		So(bindResp.User, ShouldNotBeNil)
	})

	Convey("Test password change lock", t, func() {

		resp, err := h.CreateUser(ctx, &idm.CreateUserRequest{
			User: &idm.User{Login: "emma",
				Password: "oldpassword",
				Attributes: map[string]string{
					"locks": `["pass_change","other"]`,
				}}})
		So(err, ShouldBeNil)

		updatedContext := context.WithValue(ctx, common.PydioContextUserKey, "emma")
		resp, err = h.CreateUser(updatedContext, &idm.CreateUserRequest{
			User: &idm.User{
				Login:       "emma",
				OldPassword: "oldpassword",
				Password:    "oldpassword",
				Attributes: map[string]string{
					"locks": `["pass_change","other"]`,
				}}})
		So(err, ShouldNotBeNil)

		resp, err = h.CreateUser(updatedContext, &idm.CreateUserRequest{
			User: &idm.User{Login: "emma",
				Password: "newpassword",
				Attributes: map[string]string{
					"locks": `["pass_change","other"]`,
				}}})
		So(err, ShouldBeNil)
		So(resp.User.Attributes, ShouldNotBeNil)
		So(resp.User.Attributes["locks"], ShouldEqual, `["other"]`)

	})

}

// =================================================
// * Mock *
// =================================================

type userStreamMock struct {
	ctx            context.Context
	InternalBuffer []*idm.User
}

func (x *userStreamMock) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (x *userStreamMock) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (x *userStreamMock) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (x *userStreamMock) Context() context.Context {
	return x.ctx
}

func (x *userStreamMock) Close() error {
	return nil
}

func (x *userStreamMock) SendMsg(m interface{}) error {
	return nil
}

func (x *userStreamMock) RecvMsg(m interface{}) error {
	return nil
}

func (x *userStreamMock) Recv() (*idm.SearchUserRequest, error) {
	return nil, nil
}

func (x *userStreamMock) Send(m *idm.SearchUserResponse) error {
	x.InternalBuffer = append(x.InternalBuffer, m.User)
	return nil
}
