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

// Package grpc provides persistence layer for user-defined metadata
package grpc

import (
	"context"

	"github.com/go-sql-driver/mysql"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/proto/idm"
	service2 "github.com/pydio/cells/v4/common/proto/service"
	"github.com/pydio/cells/v4/common/service"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/idm/meta"
)

var (
	Name = common.ServiceGrpcNamespace_ + common.ServiceUserMeta
)

func init() {
	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(Name),
			service.Context(ctx),
			service.Tag(common.ServiceTagIdm),
			service.Description("User-defined Metadata"),
			/*
				service.WithStorage(meta.NewDAO, "idm_usr_meta"),
				service.Unique(true),
				service.Migrations([]*service.Migration{
					{
						TargetVersion: service.FirstRun(),
						Up:            defaultMetas,
					},
				}),
				service.WithMicro(func(m micro.Service) error {
					ctx := m.Options().Context
					server := NewHandler()

					service.AddMicroMeta(m, meta2.ServiceMetaProvider, "stream")
					service.AddMicroMeta(m, meta2.ServiceMetaNsProvider, "list")

					m.Init(micro.BeforeStop(func() error {
						server.Stop()
						return nil
					}))
					idm.RegisterUserMetaServiceHandler(m.Options().Server, server)
					tree.RegisterNodeProviderStreamerHandler(m.Options().Server, server)

					// Clean role on user deletion
					cleaner := NewCleaner(servicecontext.GetDAO(ctx))
					if err := m.Options().Server.Subscribe(m.Options().Server.NewSubscriber(common.TopicIdmEvent, cleaner)); err != nil {
						return err
					}
					return nil
				}),

			*/
		)
	})
}

func defaultMetas(ctx context.Context) error {
	log.Logger(ctx).Info("Inserting default namespace for metadata")
	dao := servicecontext.GetDAO(ctx).(meta.DAO)
	err := dao.GetNamespaceDao().Add(&idm.UserMetaNamespace{
		Namespace:      "usermeta-tags",
		Label:          "Tags",
		Indexable:      true,
		JsonDefinition: "{\"type\":\"tags\"}",
		Policies: []*service2.ResourcePolicy{
			{Action: service2.ResourcePolicyAction_READ, Subject: "*", Effect: service2.ResourcePolicy_allow},
			{Action: service2.ResourcePolicyAction_WRITE, Subject: "*", Effect: service2.ResourcePolicy_allow},
		},
	})

	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return err
	}
	if me.Number == 1062 {
		// This is a duplicate error, we ignore it
		return nil
	}
	return err
}
