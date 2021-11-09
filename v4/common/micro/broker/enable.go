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

package broker

import (
	"fmt"

	"github.com/pydio/cells/common"

	"github.com/micro/micro/v3/service/broker/memory"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/server"

	defaults "github.com/pydio/cells/v4/common/micro"
	"github.com/pydio/cells/v4/common/micro/broker/service"
	"github.com/pydio/cells/v4/common/registry"
)

/* TODO v4 still needed ?
// EnableNATS enables the nats broker
func EnableNATS() {
	addr := viper.GetString("nats_address")
	b := NewBroker(
		nats.NewBroker(
			broker.Addrs(addr),
			broker.Codec(proto.NewCodec()),
		),
		BeforeDisconnect(func() error {
			s, err := registry.ListRunningServices()
			if err != nil {
				return err
			}
			if len(s) > 0 {
				return fmt.Errorf("services are still running")
			}

			return nil
		}),
	)

	defaults.InitServer(func() server.Option {
		return server.Broker(b)
	})

	defaults.InitClient(func() client.Option {
		return client.Broker(b)
	})

	broker.DefaultBroker = b

	// Establishing connectin
	broker.Connect()
}

// EnableSTAN enables the stan broker
func EnableSTAN() {
	addr := viper.GetString("nats_address")
	b := NewBroker(
		stan.NewBroker(
			stan.ClusterID(viper.GetString("nats_streaming_cluster_id")),
			broker.Addrs(addr),
			broker.Codec(proto.NewCodec()),
		),
		BeforeDisconnect(func() error {
			s, err := registry.ListRunningServices()
			if err != nil {
				return err
			}
			if len(s) > 0 {
				return fmt.Errorf("services are still running")
			}

			return nil
		}),
	)

	defaults.InitServer(func() server.Option {
		return server.Broker(b)
	})

	defaults.InitClient(func() client.Option {
		return client.Broker(b)
	})

	broker.DefaultBroker = b

	// Establishing connectin
	broker.Connect()
}

// EnableSTAN enables the http broker
func EnableHTTP() {

	addr := viper.GetString("broker_address")
	b := http.NewBroker(broker.Addrs(addr))

	defaults.InitServer(func() server.Option {
		return server.Broker(b)
	})

	defaults.InitClient(func() client.Option {
		return client.Broker(b)
	})

	broker.DefaultBroker = b

	// Establishing connectin
	broker.Connect()
}
*/

func EnableMemory() {
	b := NewBroker(
		memory.NewBroker(),
		BeforeDisconnect(func() error {
			s, err := registry.ListServices()
			if err != nil {
				return err
			}
			if len(s) > 0 {
				return fmt.Errorf("services are still running")
			}

			return nil
		}),
	)

	defaults.InitServer(func() server.Option {
		return server.Broker(b)
	})

	defaults.InitClient(func() client.Option {
		return client.Broker(b)
	})

	b.Connect()
}

func EnableService(hostname, port string) {
	b := service.NewBroker(
		service.WithClientConn(
			defaults.NewClientConn(common.ServiceBroker),
		),
	)

	defaults.InitServer(func() server.Option {
		return server.Broker(b)
	})

	defaults.InitClient(func() client.Option {
		return client.Broker(b)
	})

	// Establishing connectin
	b.Connect()
}
