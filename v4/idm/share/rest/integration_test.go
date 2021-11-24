/*
 * Copyright (c) 2019-2021. Abstrium SAS <team (at) pydio.com>
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

package rest_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"path"
	"testing"

	"github.com/pydio/cells/v4/common/nodes"

	"github.com/pydio/cells/v4/common/nodes/compose"
	"github.com/pydio/cells/v4/common/proto/service"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/pydio/cells/v4/common/auth"
	"github.com/pydio/cells/v4/common/utils/permissions"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/emicklei/go-restful"
	"github.com/pydio/cells/v4/common/server/stubs/idmtest"
	rest2 "github.com/pydio/cells/v4/idm/share/rest"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/client/grpc"
	"github.com/pydio/cells/v4/common/proto/idm"
	"github.com/pydio/cells/v4/common/proto/rest"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/server/stubs/datatest"
	"github.com/pydio/cells/v4/idm/share"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {

	testData, er := idmtest.GetStartData()
	if er != nil {
		log.Fatal(er)
	}

	ds, er := datatest.NewDocStoreService()
	if er != nil {
		log.Fatal(er)
	}
	grpc.RegisterMock(common.ServiceDocStore, ds)

	if er := datatest.RegisterTreeAndDatasources(); er != nil {
		log.Fatal(er)
	}
	if er := idmtest.RegisterIdmMocksWithData(testData); er != nil {
		log.Fatal(er)
	}

	nodes.UseMockStorageClientType()

	m.Run()
}

type reqBody struct {
	bytes.Buffer
}

func (r *reqBody) Close() error {
	return nil
}

type respWriter struct {
	bytes.Buffer
	statusCode int
	hh         http.Header
}

func (r *respWriter) Header() http.Header {
	return r.hh
}

func (r *respWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func TestShareLinks(t *testing.T) {

	Convey("Test CRUD Share Link on File", t, func() {

		ctx := context.Background()
		u, e := permissions.SearchUniqueUser(ctx, "admin", "")
		So(e, ShouldBeNil)
		ctx = auth.WithImpersonate(ctx, u)

		newNode := &tree.Node{Path: "pydiods1/file.ex", Type: tree.NodeType_LEAF, Size: 24}
		nc := tree.NewNodeReceiverClient(grpc.NewClientConn(common.ServiceTree))
		cR, e := nc.CreateNode(ctx, &tree.CreateNodeRequest{Node: newNode})
		So(e, ShouldBeNil)
		newNode = cR.GetNode()
		So(newNode.Uuid, ShouldNotBeEmpty)

		inputData, _ := protojson.Marshal(&rest.PutShareLinkRequest{
			ShareLink: &rest.ShareLink{
				Label:     "Link to File.ex",
				RootNodes: []*tree.Node{{Uuid: newNode.Uuid}},
				Permissions: []rest.ShareLinkAccessType{
					rest.ShareLinkAccessType_Download, rest.ShareLinkAccessType_Preview,
				},
			},
		})
		input := bytes.NewBuffer(inputData)
		body := &reqBody{Buffer: *input}
		req := &restful.Request{
			Request: (&http.Request{
				Body: body,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
					"Accept":       {"application/json"},
				},
			}).WithContext(ctx),
		}
		restful.DefaultResponseContentType(restful.MIME_JSON)
		output := bytes.NewBuffer([]byte{})
		resp := restful.NewResponse(&respWriter{Buffer: *output, hh: map[string][]string{}})

		h := rest2.NewSharesHandler()
		h.PutShareLink(req, resp)
		sCode := resp.ResponseWriter.(*respWriter).statusCode
		sContent := resp.ResponseWriter.(*respWriter).String()
		So(sCode, ShouldEqual, http.StatusOK)
		So(sContent, ShouldNotBeEmpty)
		outputLink := &rest.ShareLink{}
		So(protojson.Unmarshal([]byte(sContent), outputLink), ShouldBeNil)
		t.Log("Created a shared link", outputLink)

		// Now try to access link as the new user
		hiddenUser, e := permissions.SearchUniqueUser(context.Background(), outputLink.UserLogin, "")
		So(e, ShouldBeNil)
		So(hiddenUser.Attributes, ShouldContainKey, "hidden")
		hiddenCtx := auth.WithImpersonate(context.Background(), hiddenUser)
		wsClient := idm.NewWorkspaceServiceClient(grpc.NewClientConn(common.ServiceWorkspace))
		q, _ := anypb.New(&idm.WorkspaceSingleQuery{Uuid: outputLink.Uuid})
		sr, e := wsClient.SearchWorkspace(hiddenCtx, &idm.SearchWorkspaceRequest{Query: &service.Query{SubQueries: []*anypb.Any{q}}})
		So(e, ShouldBeNil)
		srw, _ := sr.Recv()
		slugRoot := srw.GetWorkspace().GetSlug()

		// Create slug/
		hash := md5.New()
		hash.Write([]byte(newNode.Uuid))
		rand := hex.EncodeToString(hash.Sum(nil))
		rootKey := rand[0:8] + "-" + path.Base(newNode.GetPath())

		read, e := compose.PathClient().ReadNode(hiddenCtx, &tree.ReadNodeRequest{Node: &tree.Node{Path: path.Join(slugRoot, rootKey)}})
		So(e, ShouldBeNil)
		So(read, ShouldNotBeEmpty)
		t.Log("Router Accessed File from Hidden User", read.Node)

	})
}

func TestBasicMocks(t *testing.T) {
	bg := context.Background()
	Convey("Test Basic Docstore Mock", t, func() {
		e := share.StoreHashDocument(bg, &idm.User{Uuid: "uuid", Login: "login"}, &rest.ShareLink{
			Uuid:             "link-uuid",
			LinkHash:         "hash",
			Label:            "My Link",
			Description:      "My Description",
			PasswordRequired: false,
		})
		So(e, ShouldBeNil)
		loadLink := &rest.ShareLink{Uuid: "link-uuid"}
		e = share.LoadHashDocumentData(bg, loadLink, []*idm.ACL{})
		So(e, ShouldBeNil)
		So(loadLink.LinkHash, ShouldEqual, "hash")
	})

	Convey("Test Index Mock", t, func() {
		cl := tree.NewNodeReceiverClient(grpc.NewClientConn(common.ServiceDataIndex_ + "pydiods1"))
		resp, e := cl.CreateNode(bg, &tree.CreateNodeRequest{Node: &tree.Node{Path: "/test", Type: tree.NodeType_COLLECTION, Size: 24, Etag: "etag"}})
		So(e, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Node.Uuid, ShouldNotBeEmpty)

		cl2 := tree.NewNodeProviderClient(grpc.NewClientConn(common.ServiceDataIndex_ + "pydiods1"))
		st, e := cl2.ListNodes(bg, &tree.ListNodesRequest{Node: &tree.Node{Path: "/"}})
		So(e, ShouldBeNil)
		var nn []*tree.Node
		for {
			r, e := st.Recv()
			if e != nil {
				break
			}
			nn = append(nn, r.GetNode())
		}
		So(nn, ShouldHaveLength, 1)
	})

	Convey("Test Tree Mock", t, func() {
		conn := grpc.NewClientConn(common.ServiceTree)
		conn2 := grpc.NewClientConn(common.ServiceMeta)
		cl := tree.NewNodeReceiverClient(conn)
		resp, e := cl.CreateNode(bg, &tree.CreateNodeRequest{Node: &tree.Node{Path: "/pydiods1/test", Type: tree.NodeType_COLLECTION, Size: 24, Etag: "etag"}})
		So(e, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(resp.Node.Uuid, ShouldNotBeEmpty)
		clM := tree.NewNodeReceiverClient(conn2)
		clone := resp.Node.Clone()
		clone.MustSetMeta("namespace", "\"value\"")
		_, e = clM.CreateNode(bg, &tree.CreateNodeRequest{Node: clone})
		So(e, ShouldBeNil)

		cl2 := tree.NewNodeProviderClient(conn)
		st, e := cl2.ListNodes(bg, &tree.ListNodesRequest{Node: &tree.Node{Path: "/"}, Recursive: true})
		So(e, ShouldBeNil)
		var nn []*tree.Node
		var cloneRes *tree.Node
		for {
			r, e := st.Recv()
			if e != nil {
				break
			}
			if r.GetNode().GetUuid() == clone.GetUuid() {
				cloneRes = r.GetNode()
			}
			nn = append(nn, r.GetNode())
		}
		So(nn, ShouldHaveLength, 6) // All DSS Roots + New Node
		So(cloneRes, ShouldNotBeEmpty)
		So(cloneRes.HasMetaKey("namespace"), ShouldBeTrue)
	})
}
