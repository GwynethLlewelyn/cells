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

// Package grpc provides a pydio GRPC service for managing files encryption keys
package grpc

import (
	"context"

	"github.com/pydio/cells/v4/common/proto/encryption"
	"github.com/pydio/cells/v4/data/key"
	"google.golang.org/grpc"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/micro/broker"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/service"
)

func init() {
	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(common.ServiceGrpcNamespace_+common.ServiceEncKey),
			service.Context(ctx),
			service.Tag(common.ServiceTagData),
			service.Description("Encryption Keys server"),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceMeta, []string{}),
			service.WithStorage(key.NewDAO, "data_key"),
			service.WithGRPC(func(c context.Context, srv *grpc.Server) error {
				h := &NodeKeyManagerHandler{}
				encryption.RegisterNodeKeyManagerServer(srv, h)
				sub, e := broker.Subscribe(common.TopicTreeChanges, func(message broker.Message) error {
					msg := &tree.NodeChangeEvent{}
					if ctx, e := message.Unmarshal(msg); e == nil {
						return h.HandleTreeChanges(ctx, msg)
					}
					return nil
				})
				if e != nil {
					return e
				}
				go func() {
					<-c.Done()
					_ = sub()
				}()
				return nil
			}),
		)
	})
}
