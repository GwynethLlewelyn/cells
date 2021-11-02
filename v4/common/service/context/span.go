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

package servicecontext

import (
	"context"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/pborman/uuid"
	"google.golang.org/grpc"
)

type spanContextKey struct{}

const (
	SpanMetadataId           = "x-pydio-span-id"
	SpanMetadataRootParentId = "x-pydio-span-root-id"
	OperationMetadataId      = "x-pydio-operation-id"
)

type Span struct {
	SpanId       string
	ParentId     string
	RootParentId string
}

func NewSpan() *Span {
	return &Span{
		SpanId: uuid.NewUUID().String(),
	}
}

func NewSpanFromParent(s *Span) *Span {
	root := s.RootParentId
	if len(root) == 0 {
		root = s.SpanId
	}
	return &Span{
		SpanId:       uuid.NewUUID().String(),
		ParentId:     s.SpanId,
		RootParentId: root,
	}
}

func WithSpan(ctx context.Context, s *Span) context.Context {
	md := metadata.Metadata{}
	if meta, ok := metadata.FromContext(ctx); ok {
		for k, v := range meta {
			md[k] = v
		}
	}
	md[SpanMetadataId] = s.SpanId
	if len(s.RootParentId) > 0 {
		md[SpanMetadataRootParentId] = s.RootParentId
	}
	ctx = metadata.NewContext(ctx, md)
	return context.WithValue(ctx, spanContextKey{}, s)
}

func SpanFromContext(ctx context.Context) (*Span, bool) {
	if val := ctx.Value(spanContextKey{}); val != nil {
		return val.(*Span), true
	} else {
		return nil, false
	}
}

func SpanFromHeader(md metadata.Metadata) (*Span, bool) {
	if md == nil {
		return nil, false
	}
	if id, ok := md[SpanMetadataId]; ok {
		s := &Span{
			SpanId: id,
		}
		if parent, ok := md[SpanMetadataRootParentId]; ok {
			s.RootParentId = parent
		}
		return s, true
	}
	return nil, false
}

func childOrNewSpan(ctx context.Context) context.Context {
	var s *Span
	if parent, ok := SpanFromContext(ctx); ok {
		s = NewSpanFromParent(parent)
	} else if md, ok := metadata.FromContext(ctx); ok {
		if headSpan, okk := SpanFromHeader(md); okk {
			s = NewSpanFromParent(headSpan)
		}
	}
	if s == nil {
		s = NewSpan()
	}
	return WithSpan(ctx, s)
}

func ctxWithOpIdFromMeta(ctx context.Context) context.Context {
	if md, ok := metadata.FromContext(ctx); ok {
		if opId, o := md[OperationMetadataId]; o {
			ctx = WithOperationID(ctx, opId)
		}
	}
	return ctx
}

func SpanUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		//fmt.Println("Client Call", method)
		if _, ok := SpanFromContext(ctx); !ok {
			s := NewSpan()
			ctx = WithSpan(ctx, s)
		}
		if opID, _ := GetOperationID(ctx); opID != "" {
			md := metadata.Metadata{}
			if meta, ok := metadata.FromContext(ctx); ok {
				for k, v := range meta {
					md[k] = v
				}
			}
			md[OperationMetadataId] = opID
			ctx = metadata.NewContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func SpanStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		//fmt.Println("Client Stream", method)
		if _, ok := SpanFromContext(ctx); !ok {
			s := NewSpan()
			ctx = WithSpan(ctx, s)
		}
		if opID, _ := GetOperationID(ctx); opID != "" {
			md := metadata.Metadata{}
			if meta, ok := metadata.FromContext(ctx); ok {
				for k, v := range meta {
					md[k] = v
				}
			}
			md[OperationMetadataId] = opID
			ctx = metadata.NewContext(ctx, md)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func SpanUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// fmt.Println("Server", req.Method())
		ctx = childOrNewSpan(ctx)
		ctx = ctxWithOpIdFromMeta(ctx)
		return handler(ctx, req)
	}
}

func SpanStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		ctx = childOrNewSpan(ctx)
		ctx = ctxWithOpIdFromMeta(ctx)

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

/*
// SpanSubscriberWrapper wraps a db connection for each subscriber
func SpanSubscriberWrapper(subscriberFunc server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		ctx = childOrNewSpan(ctx)
		ctx = ctxWithOpIdFromMeta(ctx)
		return subscriberFunc(ctx, msg)
	}
}
*/

// HttpSpanHandlerWrapper extracts data from request and put it in context Metadata field
func HttpSpanHandlerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(childOrNewSpan(r.Context()))
		h.ServeHTTP(w, r)
	})
}
