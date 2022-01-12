package registry

import (
	"context"

	pb "github.com/pydio/cells/v4/common/proto/registry"
)

const (
	ServiceMetaOverride = "service-override"
)

type Option func(*Options) error

type Options struct {
	Context context.Context
	Name    string
	Type    pb.ItemType
	Filter  func(item Item) bool
}

func WithName(n string) Option {
	return func(o *Options) error {
		o.Name = n
		return nil
	}
}

func WithType(t pb.ItemType) Option {
	return func(o *Options) error {
		o.Type = t
		return nil
	}
}

func WithMeta(name, value string) Option {
	return func(options *Options) error {
		options.Filter = func(item Item) bool {
			mm := item.Metadata()
			val, has := mm[name]
			if !has {
				return false
			}
			if len(value) > 0 && val != value {
				return false
			}
			return true
		}
		return nil
	}
}
