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

package nodes

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	servercontext "github.com/pydio/cells/v4/common/server/context"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"

	"github.com/pydio/cells/v4/common"
	clientgrpc "github.com/pydio/cells/v4/common/client/grpc"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/proto/object"
	pb "github.com/pydio/cells/v4/common/proto/registry"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/registry"
	"github.com/pydio/cells/v4/common/utils/configx"
)

type sourceAlias struct {
	dataSource string
	bucket     string
}

type LoadedSource struct {
	*object.DataSource
	Client StorageClient
}

type SourcesPool interface {
	Close()
	GetRuntimeCtx() context.Context
	GetTreeClient() tree.NodeProviderClient
	GetTreeClientWrite() tree.NodeReceiverClient
	GetDataSourceInfo(dsName string, retries ...int) (LoadedSource, error)
	GetDataSources() map[string]LoadedSource
	LoadDataSources()
}

func (s LoadedSource) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	return encoder.AddObject("DataSource", s.DataSource)
}

// ClientsPool is responsible for discovering available datasources and
// keeping an up to date registry that is used by the routers.
type ClientsPool struct {
	ctx context.Context

	sync.Mutex
	sources map[string]LoadedSource
	aliases map[string]sourceAlias

	// Statically set for testing
	treeClient      tree.NodeProviderClient
	treeClientWrite tree.NodeReceiverClient

	regWatcher  registry.Watcher
	confWatcher configx.Receiver
}

// NewClientsPool creates a client Pool and initialises it by calling the registry.
func NewClientsPool(ctx context.Context, watchRegistry bool) *ClientsPool {

	pool := &ClientsPool{
		ctx:     ctx,
		sources: make(map[string]LoadedSource),
		aliases: make(map[string]sourceAlias),
	}

	if IsUnitTestEnv {
		// Workaround the fact that no registry is present when doing unit tests
		return pool
	}

	go pool.LoadDataSources()
	if watchRegistry {
		reg := servercontext.GetRegistry(ctx)
		if reg == nil {
			debug.PrintStack()
			log.Logger(context.Background()).Warn("Starting clients pool registry watcher with empty registry, will use default")
		}
		go func() {
			e := pool.watchRegistry(reg)
			if e != nil {
				log.Logger(context.Background()).Warn("Cannot watch registry in client pool", zap.Error(e))
			}
		}()
		go pool.watchConfigChanges()
	}

	return pool
}

// Close stops the underlying watcher if defined.
func (p *ClientsPool) Close() {
	if p.regWatcher != nil {
		p.regWatcher.Stop()
	}
	if p.confWatcher != nil {
		p.confWatcher.Stop()
	}
}

func (p *ClientsPool) GetRuntimeCtx() context.Context {
	return p.ctx
}

// GetTreeClient returns the internal NodeProviderClient pointing to the TreeService.
func (p *ClientsPool) GetTreeClient() tree.NodeProviderClient {
	if p.treeClient != nil {
		return p.treeClient
	}
	return tree.NewNodeProviderClient(clientgrpc.GetClientConnFromCtx(p.ctx, common.ServiceGrpcNamespace_+common.ServiceTree))
}

// GetTreeClientWrite returns the internal NodeReceiverClient pointing to the TreeService.
func (p *ClientsPool) GetTreeClientWrite() tree.NodeReceiverClient {
	if p.treeClientWrite != nil {
		return p.treeClientWrite
	}
	return tree.NewNodeReceiverClient(clientgrpc.GetClientConnFromCtx(p.ctx, common.ServiceGrpcNamespace_+common.ServiceTree))
}

// GetDataSourceInfo tries to find information about a DataSource, eventually retrying as DataSource
// could be currently starting and not yet registered in the ClientsPool.
func (p *ClientsPool) GetDataSourceInfo(dsName string, retries ...int) (LoadedSource, error) {

	if dsName == "default" {
		dsName = config.Get("defaults", "datasource").Default("default").String()
	}

	if cl, ok := p.sources[dsName]; ok {

		return cl, nil

	} else if alias, aOk := p.aliases[dsName]; aOk {

		if dsi, e := p.GetDataSourceInfo(alias.dataSource); e != nil {

			return dsi, e

		} else {

			ds := LoadedSource{}
			ds.DataSource = proto.Clone(dsi.DataSource).(*object.DataSource)
			ds.DataSource.ObjectsBucket = alias.bucket
			ds.Client = dsi.Client
			return ds, nil

		}

	} else if len(retries) == 0 || retries[0] <= 5 {

		var retry int
		if len(retries) > 0 {
			retry = retries[0]
		}
		delay := (retry + 1) * 2

		log.Logger(context.Background()).Debug(fmt.Sprintf("[ClientsPool] cannot find datasource, retrying in %ds...", delay), zap.String("ds", dsName), zap.Any("retries", retry))

		<-time.After(time.Duration(delay) * time.Second)
		p.LoadDataSources()
		return p.GetDataSourceInfo(dsName, retry+1)

	} else {

		e := fmt.Errorf("Could not find DataSource " + dsName)
		var keys []string
		for k := range p.sources {
			keys = append(keys, k)
		}
		log.Logger(context.Background()).Error(e.Error(), zap.Strings("currentSources", keys))
		return LoadedSource{}, e

	}

}

// GetDataSources returns currently loaded datasources
func (p *ClientsPool) GetDataSources() map[string]LoadedSource {
	return p.sources
}

// LoadDataSources queries the registry to reload available datasources
func (p *ClientsPool) LoadDataSources() {
	if IsUnitTestEnv {
		// Workaround the fact that no registry is present when doing unit tests
		return
	}

	sources := config.Get("services", common.ServiceGrpcNamespace_+common.ServiceDataSync, "sources").StringArray()
	sources = config.SourceNamesFiltered(sources)

	for _, source := range sources {
		endpointClient := object.NewDataSourceEndpointClient(clientgrpc.GetClientConnFromCtx(p.ctx, common.ServiceGrpcNamespace_+common.ServiceDataSync_+source))
		response, err := endpointClient.GetDataSourceConfig(context.Background(), &object.GetDataSourceConfigRequest{})
		if err == nil && response.DataSource != nil {
			log.Logger(context.Background()).Debug("Creating client for datasource " + source)
			if e := p.CreateClientsForDataSource(source, response.DataSource); e != nil {
				log.Logger(context.Background()).Warn("Cannot create clients for datasource "+source, zap.Error(e))
			}
		} else {
			log.Logger(context.Background()).Warn("no answer from endpoint, maybe not ready yet? "+common.ServiceGrpcNamespace_+common.ServiceDataSync_+source, zap.Any("r", response), zap.Error(err))
		}
	}

	if e := p.registerAlternativeClient(common.PydioThumbstoreNamespace); e != nil {
		log.Logger(context.Background()).Warn("Cannot register alternative client "+common.PydioThumbstoreNamespace, zap.Error(e))
	}
	if e := p.registerAlternativeClient(common.PydioDocstoreBinariesNamespace); e != nil {
		log.Logger(context.Background()).Warn("Cannot register alternative client "+common.PydioDocstoreBinariesNamespace, zap.Error(e))
	}
	if e := p.registerAlternativeClient(common.PydioVersionsNamespace); e != nil {
		log.Logger(context.Background()).Warn("Cannot register alternative client "+common.PydioVersionsNamespace, zap.Error(e))
	}
}

func (p *ClientsPool) registerAlternativeClient(namespace string) error {
	dataSource, bucket, err := GetGenericStoreClientConfig(namespace)
	if err != nil {
		return err
	}
	p.Lock()
	defer p.Unlock()
	p.aliases[namespace] = sourceAlias{
		dataSource: dataSource,
		bucket:     bucket,
	}
	return nil
}

func (p *ClientsPool) watchRegistry(reg registry.Registry) error {
	if reg == nil {
		defaultReg, err := registry.OpenRegistry(context.Background(), viper.GetString("registry"))
		if err != nil {
			return err
		}
		reg = defaultReg
	}

	w, err := reg.Watch(registry.WithType(pb.ItemType_SERVICE))
	if err != nil {
		return err
	}
	p.regWatcher = w

	prefix := common.ServiceGrpcNamespace_ + common.ServiceDataSync_

	for {
		r, err := w.Next()
		if err != nil {
			return err
		}

		for _, item := range r.Items() {
			var s registry.Service
			if !item.As(&s) {
				continue
			}
			if !strings.HasPrefix(s.Name(), prefix) {
				continue
			}
			dsName := strings.TrimPrefix(s.Name(), prefix)
			if _, ok := p.sources[dsName]; ok && r.Action() == pb.ActionType_DELETE {
				p.Lock()
				delete(p.sources, dsName)
				p.Unlock()
			}
		}
		p.LoadDataSources()
	}

}

func (p *ClientsPool) watchConfigChanges() {
	for {
		watcher, err := config.Watch("services", common.ServiceGrpcNamespace_+common.ServiceDataSync, "sources")
		if err != nil {
			// Cool-off period
			time.Sleep(1 * time.Second)
			continue
		}

		p.confWatcher = watcher
		for {
			event, err := watcher.Next()
			if err != nil {
				break
			}

			if event != nil {
				p.LoadDataSources()
			}
		}

		watcher.Stop()

		// Cool-off period
		time.Sleep(1 * time.Second)
	}
}

func (p *ClientsPool) CreateClientsForDataSource(dataSourceName string, dataSource *object.DataSource) error {

	log.Logger(context.Background()).Debug("Adding dataSource", zap.String("dsname", dataSourceName))
	loaded := LoadedSource{
		DataSource: dataSource,
	}
	client, err := NewStorageClient(dataSource.ClientConfig())
	if err != nil {
		return err
	}
	loaded.Client = client

	p.Lock()
	p.sources[dataSourceName] = loaded
	p.Unlock()

	return nil
}

func MakeFakeClientsPool(tc tree.NodeProviderClient, tw tree.NodeReceiverClient) *ClientsPool {
	IsUnitTestEnv = true
	c := NewClientsPool(context.TODO(), false)

	c.treeClient = tc
	c.treeClientWrite = tw

	mockDatasource := &object.DataSource{
		Name:          "datasource",
		ObjectsHost:   "localhost",
		ObjectsPort:   9078,
		ApiKey:        "access",
		ApiSecret:     "secret",
		ObjectsSecure: false,
		ObjectsBucket: "bucket",
	}

	// mockClient, err := mockDatasource.CreateClient()
	// if err != nil {
	// 	return nil, err
	// }

	loaded := LoadedSource{
		DataSource: mockDatasource,
	}
	cfg := configx.New()
	_ = cfg.Val("type").Set("mock")
	client, _ := NewStorageClient(cfg)
	loaded.Client = client
	c.sources = map[string]LoadedSource{
		"datasource": loaded,
	}
	/*
		_ = c.CreateClientsForDataSource("datasource", mockDatasource)
		c.aliases["datasource"] = sourceAlias{
			dataSource: "localhost:9078",
			bucket:     "bucket",
		}

	*/

	// coreClient, _ := minio.NewCore("localhost:9078", "access", "secret", false)
	// c.Clients = map[string]*minio.Core{
	// 	"datasource": coreClient,
	// }
	// c.dsBuckets = map[string]string{
	// 	"datasource": "bucket",
	// }
	// c.dsEncrypted = map[string]bool{
	// 	"datasource": true,
	// }
	return c
}
