package service

import (
	"context"
	"crypto/tls"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/dao"
	"github.com/pydio/cells/v4/common/server"
	"github.com/pydio/cells/v4/common/service/frontend"
	"github.com/pydio/cells/v4/common/utils/uuid"
)

type ServerProvider func(ctx context.Context) (server.Server, error)

// ServiceOptions stores all options for a pydio service
type ServiceOptions struct {
	Name string
	ID   string
	Tags []string

	Version     string
	Description string
	Source      string
	Metadata    map[string]string

	Context context.Context
	Cancel  context.CancelFunc

	DAO        func(dao.DAO) dao.DAO
	Prefix     interface{}
	Migrations []*Migration

	// Port      string
	TLSConfig *tls.Config

	Server         server.Server
	ServerProvider ServerProvider
	serverType     server.ServerType
	serverStart    func() error
	serverStop     func() error

	Dependencies []*dependency

	// Starting options
	AutoStart   bool
	AutoRestart bool
	Fork        bool
	Unique      bool

	// Before and After funcs
	BeforeStart []func(context.Context) error
	AfterStart  []func(context.Context) error
	AfterServe  []func(context.Context) error

	BeforeStop []func(context.Context) error
	AfterStop  []func(context.Context) error

	UseWebSession      bool
	WebSessionExcludes []string
}

type dependency struct {
	Name string
	Tag  []string
}

//
type ServiceOption func(*ServiceOptions)

// Name option for a service
func Name(n string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Name = n
	}
}

// Tag option for a service
func Tag(t ...string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Tags = append(o.Tags, t...)
	}
}

// Description option for a service
func Description(d string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Description = d
	}
}

// Source option for a service
func Source(s string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Source = s
	}
}

// Context option for a service
func Context(c context.Context) ServiceOption {
	return func(o *ServiceOptions) {
		o.Context = c
	}
}

// Cancel option for a service
func Cancel(c context.CancelFunc) ServiceOption {
	return func(o *ServiceOptions) {
		o.Cancel = c
	}
}

// WithTLSConfig option for a service
func WithTLSConfig(c *tls.Config) ServiceOption {
	return func(o *ServiceOptions) {
		o.TLSConfig = c
	}
}

// WithServer directly presets the server.Server instance
func WithServer(s server.Server) ServiceOption {
	return func(o *ServiceOptions) {
		o.Server = s
	}
}

// WithServerProvider passes a callback producing a new server.Server
func WithServerProvider(provider ServerProvider) ServiceOption {
	return func(o *ServiceOptions) {
		o.ServerProvider = provider
	}
}

// AutoStart option for a service
func AutoStart(b bool) ServiceOption {
	return func(o *ServiceOptions) {
		o.AutoStart = b
	}
}

// AfterStart registers a callback to be run after service.Start (blocking)
func AfterStart(f func(ctx context.Context) error) ServiceOption {
	return func(o *ServiceOptions) {
		o.AfterStart = append(o.AfterStart, f)
	}
}

// AfterServe registers a callback that is run after Server is finally started (non-blocking)
func AfterServe(f func(ctx context.Context) error) ServiceOption {
	return func(o *ServiceOptions) {
		o.AfterServe = append(o.AfterServe, f)
	}
}

// Fork option for a service
func Fork(f bool) ServiceOption {
	return func(o *ServiceOptions) {
		o.Fork = f
	}
}

// Unique option for a service
func Unique(b bool) ServiceOption {
	return func(o *ServiceOptions) {
		o.Unique = b
	}
}

// Migrations option for a service
func Migrations(migrations []*Migration) ServiceOption {
	return func(o *ServiceOptions) {
		o.Migrations = migrations
	}
}

func Metadata(name, value string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Metadata[name] = value
	}
}

// Dependency option for a service
func Dependency(n string, t []string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Dependencies = append(o.Dependencies, &dependency{n, t})
	}
}

// PluginBoxes option for a service
func PluginBoxes(boxes ...frontend.PluginBox) ServiceOption {
	return func(o *ServiceOptions) {
		o.Dependencies = append(o.Dependencies, &dependency{common.ServiceWebNamespace_ + common.ServiceFrontStatics, []string{}})
		frontend.RegisterPluginBoxes(boxes...)
	}
}

func WithWebSession(excludes ...string) ServiceOption {
	return func(o *ServiceOptions) {
		o.UseWebSession = true
		o.WebSessionExcludes = excludes
	}
}

func newOptions(opts ...ServiceOption) *ServiceOptions {
	opt := &ServiceOptions{}

	opt.ID = uuid.New()
	opt.Metadata = make(map[string]string)
	opt.Version = common.Version().String()
	opt.AutoStart = true

	for _, o := range opts {
		if o == nil {
			continue
		}

		o(opt)
	}

	return opt
}
