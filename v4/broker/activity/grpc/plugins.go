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

// Package grpc is the persistence service for all activities.
//
// It is listening to many events broadcasted by the application and storing them inside associated feeds, depending
// on the event context, owner, object type, etc...
// Persistence is implemented on a Bolt database.
package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/pydio/cells/v4/broker/activity"
	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/log"
	defaults "github.com/pydio/cells/v4/common/micro"
	"github.com/pydio/cells/v4/common/micro/broker"
	"github.com/pydio/cells/v4/common/nodes/meta"
	"github.com/pydio/cells/v4/common/plugins"
	proto "github.com/pydio/cells/v4/common/proto/activity"
	"github.com/pydio/cells/v4/common/proto/idm"
	"github.com/pydio/cells/v4/common/proto/jobs"
	serviceproto "github.com/pydio/cells/v4/common/proto/service"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/service"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/common/utils/cache"
	"github.com/pydio/cells/v4/common/utils/std"
)

var (
	Name = common.ServiceGrpcNamespace_ + common.ServiceActivity
)

func init() {
	plugins.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(Name),
			service.Context(ctx),
			service.Tag(common.ServiceTagBroker),
			service.Description("Activity Service is collecting activity for users and nodes"),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceJobs, []string{}),
			service.Dependency(common.ServiceGrpcNamespace_+common.ServiceTree, []string{}),
			service.Metadata(meta.ServiceMetaProvider, "stream"),
			service.Migrations([]*service.Migration{
				{
					TargetVersion: service.FirstRun(),
					Up:            RegisterDigestJob,
				},
			}),
			service.WithStorage(activity.NewDAO, "broker_activity"),
			service.Unique(true),
			service.WithGRPC(func(c context.Context, srv *grpc.Server) error {

				dao := servicecontext.GetDAO(c).(activity.DAO)
				// Register Subscribers
				subscriber := NewEventsSubscriber(dao)
				batcher := cache.NewEventsBatcher(c, 3*time.Second, 20*time.Second, 2000, true, func(ctx context.Context, msg ...*tree.NodeChangeEvent) {
					subscriber.HandleNodeChange(ctx, msg[0])
				})

				u1, e := broker.Subscribe(common.TopicTreeChanges, func(message broker.Message) error {
					msg := &tree.NodeChangeEvent{}
					if ctx, e := message.Unmarshal(msg); e == nil {
						if msg.Target != nil && (msg.Target.Etag == common.NodeFlagEtagTemporary || msg.Target.HasMetaKey(common.MetaNamespaceDatasourceInternal)) {
							return nil
						}
						if msg.Source != nil && msg.Source.HasMetaKey(common.MetaNamespaceDatasourceInternal) {
							return nil
						}
						if msg.Optimistic {
							return nil
						}
						batcher.Events <- &cache.EventWithContext{NodeChangeEvent: msg, Ctx: ctx}
					}
					return nil
				})
				if e != nil {
					return e
				}

				u2, e := broker.Subscribe(common.TopicMetaChanges, func(message broker.Message) error {
					msg := &tree.NodeChangeEvent{}
					if ctx, e := message.Unmarshal(msg); e == nil {
						if msg.Optimistic || msg.Type != tree.NodeChangeEvent_UPDATE_USER_META {
							return nil
						}
						batcher.Events <- &cache.EventWithContext{NodeChangeEvent: msg, Ctx: ctx}
					}
					return nil
				})
				if e != nil {
					return e
				}

				u3, e := broker.Subscribe(common.TopicIdmEvent, func(message broker.Message) error {
					msg := &idm.ChangeEvent{}
					if ctx, e := message.Unmarshal(msg); e == nil {
						return subscriber.HandleIdmChange(ctx, msg)
					}
					return nil
				})
				if e != nil {
					return e
				}

				proto.RegisterActivityServiceServer(srv, new(Handler))
				tree.RegisterNodeProviderStreamerServer(srv, new(MetaProvider))

				/*
					// TODO V4 Replaced by below ?
					m.Init(micro.BeforeStop(func() error {
						batcher.Stop()
						return nil
					}))
				*/

				go func() {
					<-c.Done()
					batcher.Stop()
					// Unsubscribes
					_ = u1()
					_ = u2()
					_ = u3()
				}()

				return nil
			}),
		)
	})
}

func RegisterDigestJob(ctx context.Context) error {

	log.Logger(ctx).Info("Registering default job for creating activities digests")
	// Build queries for standard users
	q1, _ := anypb.New(&idm.UserSingleQuery{NodeType: idm.NodeType_USER})
	q2, _ := anypb.New(&idm.UserSingleQuery{AttributeName: idm.UserAttrHidden, AttributeAnyValue: true, Not: true})
	job := &jobs.Job{
		ID:             "users-activity-digest",
		Label:          "Users activities digest",
		Owner:          common.PydioSystemUsername,
		MaxConcurrency: 1,
		AutoStart:      false,
		Schedule: &jobs.Schedule{
			Iso8601Schedule: "R/2012-06-04T19:25:16.828696-07:00/PT15M", // every 5 mn
		},
		Actions: []*jobs.Action{
			{
				ID: "broker.activity.actions.mail-digest",
				UsersSelector: &jobs.UsersSelector{
					Label: "All users except hidden",
					Query: &serviceproto.Query{
						SubQueries: []*anypb.Any{q1, q2},
						Operation:  serviceproto.OperationType_AND,
					},
				},
			},
		},
	}

	return std.Retry(ctx, func() error {
		cliJob := jobs.NewJobServiceClient(defaults.NewClientConn(common.ServiceJobs))
		_, e := cliJob.PutJob(ctx, &jobs.PutJobRequest{Job: job} /*, registry.ShortRequestTimeout() - TODO V4*/)
		return e
	}, 5*time.Second, 20*time.Second)

}
