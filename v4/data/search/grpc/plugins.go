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

// Package grpc provides the Pydio grpc service for querying indexer.
//
// Insertion in the index is not performed directly but via events broadcasted by the broker.
package grpc

import (
	"context"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/service"
)

var (
	Name = common.ServiceGrpcNamespace_ + common.ServiceSearch
)

func init() {

	config.RegisterExposedConfigs(Name, ExposedConfigs)

	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(Name),
			service.Context(ctx),
			service.Tag(common.ServiceTagData),
			service.Description("Search Engine"),
			service.Fork(true),
			/*
				service.RouterDependencies(),
				service.AutoRestart(true),
				service.WithMicro(func(m micro.Service) error {

					ctx := m.Options().Context
					cfg := servicecontext.GetConfig(ctx)

					indexContent := cfg.Val("indexContent").Bool()
					if indexContent {
						log.Logger(m.Options().Context).Info("Enabling content indexation in search engine")
					} else {
						log.Logger(m.Options().Context).Info("disabling content indexation in search engine")
					}

					dir, _ := config.ServiceDataDir(Name)
					bleve.BleveIndexPath = filepath.Join(dir, "searchengine.bleve")
					bleveConfs := make(map[string]interface{})
					bleveConfs["basenameAnalyzer"] = cfg.Val("basenameAnalyzer").String()
					bleveConfs["contentAnalyzer"] = cfg.Val("contentAnalyzer").String()

					bleveEngine, err := bleve.NewBleveEngine(indexContent, bleveConfs)
					if err != nil {
						return err
					}

					server := &SearchServer{
						Engine:           bleveEngine,
						TreeClient:       tree.NewNodeProviderClient(defaults.NewClientConn(common.ServiceTree)),
						ReIndexThrottler: make(chan struct{}, 5),
					}

					tree.RegisterSearcherHandler(m.Options().Server, server)
					sync.RegisterSyncEndpointHandler(m.Options().Server, server)

					m.Init(
						micro.BeforeStop(bleveEngine.Close),
					)

					// Register Subscribers
					if err := m.Options().Server.Subscribe(
						m.Options().Server.NewSubscriber(
							common.TopicMetaChanges,
							server.CreateNodeChangeSubscriber(),
						),
					); err != nil {
						return err
					}

					return nil
				}),
			*/
		)
	})
}
