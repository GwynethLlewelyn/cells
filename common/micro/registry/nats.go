package registry

import (
	"context"
	"net"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	gonats "github.com/nats-io/nats.go"
	gostan "github.com/nats-io/stan.go"

	defaults "github.com/pydio/cells/common/micro"
	"github.com/pydio/cells/common/micro/registry/nats"
	"github.com/pydio/cells/common/micro/registry/stan"
	"github.com/pydio/cells/common/micro/selector/cache"
	"github.com/spf13/viper"
)

type customDialer struct {
	ctx             context.Context
	nc              *gonats.Conn
	connectTimeout  time.Duration
	connectTimeWait time.Duration
}

func (cd *customDialer) Dial(network, address string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(cd.ctx, cd.connectTimeout)
	defer cancel()

	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		select {
		case <-cd.ctx.Done():
			return nil, cd.ctx.Err()
		default:
			d := &net.Dialer{}
			if conn, err := d.DialContext(ctx, network, address); err == nil {
				return conn, nil
			} else {
				time.Sleep(cd.connectTimeWait)
			}
		}
	}
}

func EnableNats() {
	addr := viper.GetString("nats_address")
	r := nats.NewRegistry(
		registry.Addrs(addr),
	)

	s := cache.NewSelector(selector.Registry(r))

	defaults.InitServer(func() server.Option {
		return server.Registry(r)
	})

	defaults.InitClient(
		func() client.Option {
			return client.Selector(s)
		}, func() client.Option {
			return client.Registry(r)
		}, func() client.Option {
			return client.Retries(5)
		}, func() client.Option {
			return client.Retry(RetryOnError)
		},
	)

	registry.DefaultRegistry = r
}

func EnableStan() {
	addr := viper.GetString("nats_address")

	r := stan.NewRegistry(
		stan.ClusterID(viper.GetString("nats_streaming_cluster_id")),
		registry.Addrs(addr),
		stan.Options(
			gostan.NatsURL(addr),
		),
	)

	s := cache.NewSelector(selector.Registry(r))

	defaults.InitServer(func() server.Option {
		return server.Registry(r)
	})

	defaults.InitClient(
		func() client.Option {
			return client.Selector(s)
		}, func() client.Option {
			return client.Registry(r)
		}, func() client.Option {
			return client.Retries(5)
		}, func() client.Option {
			return client.Retry(RetryOnError)
		},
	)

	registry.DefaultRegistry = r
}

// RetryOnError retries a request on a 500 or timeout error
func RetryOnError(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	e := errors.Parse(err.Error())
	if e == nil {
		return false, nil
	}

	switch e.Code {
	// retry on timeout or internal server error
	case 408, 500:
		return true, nil
	default:
		return false, nil
	}
}
