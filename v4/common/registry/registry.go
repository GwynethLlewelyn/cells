package registry

import (
	"context"
	"fmt"
	"github.com/micro/micro/v3/service/registry"
)

type Registry interface{
	Register(Service) error
	Deregister(Service) error
	GetService(string) ([]Service, error)
	ListServices() ([]Service, error)
	Watch(...WatchOption) (Watcher, error)
}

type WatchOption func(WatchOptions) error

type WatchOptions interface{}

type Context interface {
	Context(context.Context)
}

type Watcher interface {
	Next() (Result, error)
	Stop()
}

type Result interface {
	Action() string
	Service() Service
}

type reg struct{
	registry.Registry
}

func New(r registry.Registry) Registry {
	return &reg{
		Registry: r,
	}
}

func (r *reg) Register(s Service) error {
	var p *registry.Service
	if ok := s.As(&p); !ok {
		return fmt.Errorf("not a service")
	}
	return r.Registry.Register(p)
}

func (r *reg) Deregister(s Service) error {
	var p *registry.Service
	if ok := s.As(&p); !ok {
		return fmt.Errorf("not a service")
	}
	return r.Registry.Deregister(p)
}

func (r *reg) ListServices() ([]Service, error) {
	var services []Service

	ss, err := r.Registry.ListServices()
	if err != nil {
		return nil, err
	}
	
	for _, s := range ss {
		services = append(services, &service{
			s: s,
		})
	}
	return services, nil
}

func (r *reg) GetService(name string) ([]Service, error) {
	var services []Service

	ss, err := r.Registry.GetService(name)
	if err != nil {
		return nil, err
	}
	for _, s := range ss {
		services = append(services, &service{
			s: s,
		})
	}
	return services, nil
}

func (r *reg) Watch(opts ...WatchOption) (Watcher, error) {
	w, err := r.Registry.Watch()
	if err != nil {
		return nil, err
	}

	return &watcher{w}, nil
}

type watcher struct {
	w registry.Watcher
}

func (w *watcher) Next() (Result, error) {
	res, err := w.w.Next()
	if err != nil {
		return nil, err
	}

	return &result{r: res}, nil
}

func (w *watcher) Stop() {
	w.w.Stop()
}

type result struct {
	r *registry.Result
}

func (r *result) Action() string {
	return r.r.Action
}

func (r *result) Service() Service {
	return &service{
		s: r.r.Service,
	}
}
