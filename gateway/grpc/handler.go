package grpc

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/micro/go-micro/errors"

	"github.com/pydio/cells/common/nodes"
	"github.com/pydio/cells/common/proto/tree"
)

type TreeHandler struct {
	router nodes.Client
}

func (t *TreeHandler) fixMode(n *tree.Node) {
	if n.IsLeaf() {
		n.Mode = int32(os.ModePerm)
	} else {
		n.Mode = int32(os.ModePerm & os.ModeDir)
		n.Size = 0
		n.MTime = time.Now().Unix()
	}
}

func (t *TreeHandler) ReadNodeStream(ctx context.Context, s tree.NodeProviderStreamer_ReadNodeStreamStream) error {
	router := t.getRouter()
	var err error
	for {
		r, e := s.Recv()
		if e != nil {
			if e != io.EOF && e != io.ErrUnexpectedEOF {
				s.SendMsg(e)
				err = e
			}
			break
		}
		resp, _ := router.ReadNode(ctx, r)
		if resp == nil {
			resp = &tree.ReadNodeResponse{Success: false}
		} else {
			resp.Success = true
		}
		sE := s.Send(resp)
		if sE != nil {
			// Error while sending
			break
		}
	}
	s.Close()
	return err
}

func (t *TreeHandler) CreateNodeStream(ctx context.Context, s tree.NodeReceiverStream_CreateNodeStreamStream) error {
	router := t.getRouter()
	var err error
	for {
		r, e := s.Recv()
		if e != nil {
			if e != io.EOF {
				s.SendMsg(e)
				err = e
			}
			break
		}
		t.fixMode(r.Node)
		resp, er := router.CreateNode(ctx, r)
		if er != nil {
			s.SendMsg(er)
			err = er
			break
		}
		if err = s.Send(resp); err != nil {
			break
		}
	}
	s.Close()
	return err //errors.BadRequest("not.implemented", "CreateNodeStream not implemented yet")
}

func (t *TreeHandler) UpdateNodeStream(context.Context, tree.NodeReceiverStream_UpdateNodeStreamStream) error {
	return errors.BadRequest("not.implemented", "UpdateNodeStream not implemented yet")
}

func (t *TreeHandler) DeleteNodeStream(context.Context, tree.NodeReceiverStream_DeleteNodeStreamStream) error {
	return errors.BadRequest("not.implemented", "DeleteNodeStream not implemented yet")
}

// ReadNode forwards to router
func (t *TreeHandler) ReadNode(ctx context.Context, request *tree.ReadNodeRequest, response *tree.ReadNodeResponse) error {
	r, e := t.getRouter().ReadNode(ctx, request)
	if e != nil {
		response.Success = false
		return e
	} else {
		response.Node = r.Node
		response.Success = r.Success
		return nil
	}
}

// ListNodes forwards to router
func (t *TreeHandler) ListNodes(ctx context.Context, request *tree.ListNodesRequest, stream tree.NodeProvider_ListNodesStream) error {

	st, e := t.getRouter().ListNodes(ctx, request)
	if e != nil {
		return e
	}
	defer st.Close()

	for {
		r, e := st.Recv()
		if e == io.EOF || e == io.ErrUnexpectedEOF {
			break
		}
		if e != nil {
			return e
		}
		stream.Send(r)
	}

	return nil
}

// StreamChanges sends events to the client
func (t *TreeHandler) StreamChanges(ctx context.Context, req *tree.StreamChangesRequest, resp tree.NodeChangesStreamer_StreamChangesStream) error {

	streamer, err := t.getRouter().StreamChanges(ctx, req)
	if err != nil {
		return err
	}
	defer streamer.Close()
	for {
		r, e := streamer.Recv()
		if e != nil {
			break
		}
		resp.Send(r)
	}

	return nil
}

// CreateNode is used for creating folders
func (t *TreeHandler) CreateNode(ctx context.Context, req *tree.CreateNodeRequest, resp *tree.CreateNodeResponse) error {
	t.fixMode(req.Node)
	r, e := t.getRouter().CreateNode(ctx, req)
	if e != nil {
		resp.Success = false
		return e
	}
	resp.Node = r.Node
	resp.Success = r.Success
	return nil
}

// UpdateNode is used for moving nodes paths
func (t *TreeHandler) UpdateNode(ctx context.Context, req *tree.UpdateNodeRequest, resp *tree.UpdateNodeResponse) error {
	r, e := t.getRouter().UpdateNode(ctx, req)
	if e != nil {
		resp.Success = false
		return e
	}
	resp.Success = r.Success
	resp.Node = r.Node
	return nil
}

// DeleteNode is used to delete nodes
func (t *TreeHandler) DeleteNode(ctx context.Context, req *tree.DeleteNodeRequest, resp *tree.DeleteNodeResponse) error {
	r, e := t.getRouter().DeleteNode(ctx, req)
	if e != nil {
		resp.Success = false
		return e
	}
	resp.Success = r.Success
	return nil
}

func (t *TreeHandler) getRouter() nodes.Client {
	if t.router != nil {
		return t.router
	}
	t.router = nodes.NewStandardRouter(nodes.RouterOptions{
		AdminView:        false,
		WatchRegistry:    true,
		LogReadEvents:    false,
		AuditEvent:       false,
		SynchronousTasks: true,
	})
	return t.router
}
