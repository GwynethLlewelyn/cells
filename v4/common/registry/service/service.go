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

package service

import (
	"context"
	"net/url"
	"time"

	cgrpc "github.com/pydio/cells/v4/common/client/grpc"

	"google.golang.org/grpc"

	pb "github.com/pydio/cells/v4/common/proto/registry"
	"github.com/pydio/cells/v4/common/registry"
	"github.com/pydio/cells/v4/common/utils/std"
)

var scheme = "grpc"

type URLOpener struct {
	grpc.ClientConnInterface
}

func init() {
	o := &URLOpener{}
	registry.DefaultURLMux().Register(scheme, o)
}

func (o *URLOpener) OpenURL(ctx context.Context, u *url.URL) (registry.Registry, error) {
	// We use WithBlock, shall we timeout and retry here ?
	var conn grpc.ClientConnInterface
	err := std.Retry(ctx, func() error {
		// c, can := context.WithTimeout(ctx, 1*time.Minute)
		// defer can()
		// var e error

		address := u.Hostname()
		if port := u.Port(); port != "" {
			address = address + ":" + port
		}

		// TODO v4 error handling
		cli, err := grpc.Dial(u.Hostname()+":"+u.Port(), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return err
		}

		conn = cgrpc.NewClientConn("pydio.grpc.registry", cgrpc.WithClientConn(cli))

		return nil
	}, 30*time.Second, 5*time.Minute)
	if err != nil {
		return nil, err
	}
	/*
		conn, err := grpc.Dial(u.Hostname()+":"+u.Port(), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return nil, err
		}
	*/

	return NewRegistry(
		WithConn(conn),
	), nil
}

var (
	// The default service name
	DefaultService = "go.micro.registry"
)

type serviceRegistry struct {
	opts Options
	// name of the registry
	name string
	// address
	address []string
	// client to call registry
	client pb.RegistryClient
}

func (s *serviceRegistry) callOpts() []grpc.CallOption {
	var opts []grpc.CallOption

	// set registry address
	//if len(s.address) > 0 {
	//	opts = append(opts, client.WithAddress(s.address...))
	//}

	// set timeout
	if s.opts.Timeout > time.Duration(0) {
		// opts = append(opts, grpc.client.WithRequestTimeout(s.opts.Timeout))
	}

	// add retries
	// TODO : charles' GUTS feeling :-)
	// opts = append(opts, client.WithRetries(10))

	return opts
}

func (s *serviceRegistry) Init(opts ...Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	if len(s.opts.Addrs) > 0 {
		s.address = s.opts.Addrs
	}

	// extract the client from the context, fallback to grpc
	var conn *grpc.ClientConn
	if c, ok := s.opts.Context.Value(connKey{}).(*grpc.ClientConn); ok {
		conn = c
	} else {
		conn, _ = grpc.Dial(":8000")
	}

	s.client = pb.NewRegistryClient(conn)

	return nil
}

func (s *serviceRegistry) Options() Options {
	return s.opts
}

func (s *serviceRegistry) Start(item registry.Item) error {
	_, err := s.client.Start(s.opts.Context, ToProtoItem(item), s.callOpts()...)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceRegistry) Stop(item registry.Item) error {
	_, err := s.client.Stop(s.opts.Context, ToProtoItem(item), s.callOpts()...)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceRegistry) Register(item registry.Item) error {
	_, err := s.client.Register(s.opts.Context, ToProtoItem(item), s.callOpts()...)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceRegistry) Deregister(item registry.Item) error {
	_, err := s.client.Deregister(s.opts.Context, ToProtoItem(item), s.callOpts()...)
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceRegistry) Get(name string, opts ...registry.Option) (registry.Item, error) {
	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	rsp, err := s.client.Get(s.opts.Context, &pb.GetRequest{
		Name: name,
		Options: &pb.Options{
			Type: options.Type,
		},
	}, s.callOpts()...)
	if err != nil {
		return nil, err
	}

	return ToItem(rsp.Item), nil
}

func (s *serviceRegistry) List(opts ...registry.Option) ([]registry.Item, error) {
	var options registry.Options
	for _, o := range opts {
		o(&options)
	}
	rsp, err := s.client.List(s.opts.Context, &pb.ListRequest{
		Options: &pb.Options{
			Type: options.Type,
		},
	}, s.callOpts()...)
	if err != nil {
		return nil, err
	}

	items := make([]registry.Item, 0, len(rsp.Items))
	for _, item := range rsp.Items {
		casted := ToItem(item)
		if options.Filter != nil && !options.Filter(casted) {
			continue
		}
		items = append(items, casted)
	}

	return items, nil
}

func (s *serviceRegistry) Watch(opts ...registry.Option) (registry.Watcher, error) {
	var options registry.Options
	for _, o := range opts {
		o(&options)
	}

	ctx := context.TODO()
	req := &pb.WatchRequest{
		Options: &pb.Options{
			Type: options.Type,
		},
	}

	if options.Context != nil {
		ctx = options.Context
	}

	req.Name = options.Name

	stream, err := s.client.Watch(ctx, req, s.callOpts()...)

	if err != nil {
		return nil, err
	}

	return newWatcher(stream), nil
}

func (s *serviceRegistry) As(interface{}) bool {
	return false
}

func (s *serviceRegistry) String() string {
	return "service"
}

// NewRegistry returns a new registry service client
func NewRegistry(opts ...Option) registry.Registry {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	var ctx context.Context
	var cancel context.CancelFunc

	ctx = options.Context
	if ctx == nil {
		ctx = context.TODO()
	}

	ctx, cancel = context.WithCancel(ctx)

	options.Context = ctx

	// extract the client from the context, fallback to grpc
	var conn grpc.ClientConnInterface
	conn, ok := options.Context.Value(connKey{}).(grpc.ClientConnInterface)
	if !ok {
		conn, _ = grpc.Dial(":8000")
	}

	// service name. TODO: accept option
	name := DefaultService

	r := &serviceRegistry{
		opts:   options,
		name:   name,
		client: pb.NewRegistryClient(conn),
	}

	go func() {
		// Check the stream has a connection to the registry
		watcher, err := r.Watch()
		if err != nil {
			cancel()
			return
		}

		for {
			_, err := watcher.Next()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	return registry.NewRegistry(r)
}
