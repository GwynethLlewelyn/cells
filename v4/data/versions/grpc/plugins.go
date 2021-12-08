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

// Package grpc is in charge of storing versions metadata
package grpc

import (
	"context"
	"path"
	"time"

	grpc2 "github.com/pydio/cells/v4/common/client/grpc"

	"google.golang.org/grpc"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/proto/docstore"
	"github.com/pydio/cells/v4/common/proto/jobs"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/service"
	"github.com/pydio/cells/v4/common/utils/std"
	"github.com/pydio/cells/v4/data/versions"
	json "github.com/pydio/cells/v4/x/jsonx"
)

var (
	Name = common.ServiceGrpcNamespace_ + common.ServiceVersions
)

func init() {

	plugins.Register("main", func(ctx context.Context) {
		config.RegisterExposedConfigs(Name, ExposedConfigs)

		service.NewService(
			service.Name(Name),
			service.Context(ctx),
			service.Tag(common.ServiceTagData),
			service.Description("Versioning service"),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceJobs, []string{}),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceDocStore, []string{}),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceTree, []string{}),
			service.Migrations([]*service.Migration{
				{
					TargetVersion: service.FirstRun(),
					Up:            InitDefaults,
				},
			}),
			service.Unique(true),
			service.WithGRPC(func(ctx context.Context, server *grpc.Server) error {

				serviceDir, e := config.ServiceDataDir(common.ServiceGrpcNamespace_ + common.ServiceVersions)
				if e != nil {
					return e
				}
				store, err := versions.NewBoltStore(path.Join(serviceDir, "versions.db"))
				if err != nil {
					return err
				}

				engine := &Handler{
					db: store,
				}

				tree.RegisterNodeVersionerServer(server, engine)

				jobsClient := jobs.NewJobServiceClient(grpc2.NewClientConn(common.ServiceJobs))
				// TODO V4
				// to := registry.ShortRequestTimeout()
				// Migration from old prune-versions-job : delete if exists, replaced by composed job
				var reinsert bool
				if _, e := jobsClient.GetJob(ctx, &jobs.GetJobRequest{JobID: "prune-versions-job"}); e == nil {
					jobsClient.DeleteJob(ctx, &jobs.DeleteJobRequest{JobID: "prune-versions-job"})
					reinsert = true
				}
				vJob := getVersioningJob()
				if _, err := jobsClient.GetJob(ctx, &jobs.GetJobRequest{JobID: vJob.ID} /*, to*/); err != nil || reinsert {
					log.Logger(ctx).Info("Inserting versioning job")
					jobsClient.PutJob(ctx, &jobs.PutJobRequest{Job: vJob} /*, to*/)
				}

				return nil
			}),
		)
	})
}

func InitDefaults(ctx context.Context) error {

	// todo V4

	log.Logger(ctx).Info("Inserting default versioning policies")

	max30, _ := json.Marshal(&tree.VersioningPolicy{
		Uuid:        "default-policy",
		Name:        "30 days max",
		Description: "Remove versions after 30 days",
		KeepPeriods: []*tree.VersioningKeepPeriod{
			{IntervalStart: "0", MaxNumber: -1},
			{IntervalStart: "15d", MaxNumber: 10},
			{IntervalStart: "30d"},
		},
		VersionsDataSourceName: "versions",
		NodeDeletedStrategy:    tree.VersioningNodeDeletedStrategy_KeepLast,
	})

	keepAll, _ := json.Marshal(&tree.VersioningPolicy{
		Uuid:        "keep-all",
		Name:        "Keep all",
		Description: "Keep all versions forever",
		KeepPeriods: []*tree.VersioningKeepPeriod{
			{IntervalStart: "0", MaxNumber: -1},
		},
		VersionsDataSourceName: "versions",
		NodeDeletedStrategy:    tree.VersioningNodeDeletedStrategy_KeepAll,
	})

	regular, _ := json.Marshal(&tree.VersioningPolicy{
		Uuid:        "regular-pruning",
		Name:        "Regular Pruning",
		Description: "Prune versions regularly at various intervals and remove after 30 days",
		KeepPeriods: []*tree.VersioningKeepPeriod{
			{IntervalStart: "0", MaxNumber: 10},
			{IntervalStart: "10m", MaxNumber: 10},
			{IntervalStart: "3h", MaxNumber: 10},
			{IntervalStart: "1d", MaxNumber: 10},
			{IntervalStart: "10d", MaxNumber: 10},
			{IntervalStart: "30d", MaxNumber: 0},
		},
		VersionsDataSourceName: "versions",
		NodeDeletedStrategy:    tree.VersioningNodeDeletedStrategy_KeepLast,
	})

	return std.Retry(ctx, func() error {

		dc := docstore.NewDocStoreClient(grpc2.NewClientConn(common.ServiceDocStore))
		_, e := dc.PutDocument(ctx, &docstore.PutDocumentRequest{
			StoreID:    common.DocStoreIdVersioningPolicies,
			DocumentID: "default-policy",
			Document: &docstore.Document{
				ID:    "default-policy",
				Owner: common.PydioSystemUsername,
				Type:  docstore.DocumentType_JSON,
				Data:  string(max30),
			},
		})
		if e != nil {
			return e
		}

		_, e = dc.PutDocument(ctx, &docstore.PutDocumentRequest{
			StoreID:    common.DocStoreIdVersioningPolicies,
			DocumentID: "keep-all",
			Document: &docstore.Document{
				ID:    "keep-all",
				Owner: common.PydioSystemUsername,
				Type:  docstore.DocumentType_JSON,
				Data:  string(keepAll),
			},
		})
		if e != nil {
			return e
		}

		_, e = dc.PutDocument(ctx, &docstore.PutDocumentRequest{
			StoreID:    common.DocStoreIdVersioningPolicies,
			DocumentID: "regular-pruning",
			Document: &docstore.Document{
				ID:    "regular-pruning",
				Owner: common.PydioSystemUsername,
				Type:  docstore.DocumentType_JSON,
				Data:  string(regular),
			},
		})

		return e
	}, 3*time.Second, 10*time.Second)

}
