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

package cmd

import (
	"github.com/pydio/cells/v4/common/broker"
	"golang.org/x/sync/errgroup"

	"github.com/pydio/cells/v4/common/config/runtime"
	"github.com/pydio/cells/v4/common/plugins"
	pb "github.com/pydio/cells/v4/common/proto/registry"
	"github.com/pydio/cells/v4/common/registry"
	"github.com/pydio/cells/v4/common/server"
	"github.com/pydio/cells/v4/common/server/caddy"
	servercontext "github.com/pydio/cells/v4/common/server/context"
	"github.com/pydio/cells/v4/common/server/fork"
	"github.com/pydio/cells/v4/common/server/generic"
	"github.com/pydio/cells/v4/common/server/grpc"
	"github.com/pydio/cells/v4/common/server/http"
	"github.com/pydio/cells/v4/common/service"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FilterStartTags    []string
	FilterStartExclude []string
)

// StartCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("args", args)

		bindViperFlags(cmd.Flags(), map[string]string{
			//	"log":  "logs_level",
			"fork": "is_fork",
		})

		initLogLevel()

		initConfig()

		// Making sure we capture the signals
		handleSignals()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// broker.Connect()

		pluginsReg, err := registry.OpenRegistry(ctx, "memory:///")
		if err != nil {
			return err
		}

		reg, err := registry.OpenRegistry(ctx, viper.GetString("registry"))
		if err != nil {
			return err
		}

		broker.Register(broker.NewBroker(viper.GetString("broker")))

		ctx = servercontext.WithRegistry(ctx, reg)
		ctx = servicecontext.WithRegistry(ctx, pluginsReg)

		/*
			watcher, err := reg.Watch()
			if err != nil {
				return err
			}

			go func() {
				for {
					w, err := watcher.Next()
					if err != nil {
						fmt.Println("And the error is ? ", err)
					}

					if w.Action() == "start_request" {
						fmt.Println("Received start request for ? ", w.Item().Name())
						var node registry.Node

						if w.Item().As(&node) {
							serverType, ok := node.Metadata()["type"]
							if !ok {
								continue
							}

							fmt.Println("Starting ", node.Name())
							switch serverType {
							case "grpc":
								grpc.New(ctx)
							}
						}

						var sss registry.Service
						if w.Item().As(&sss) {
							ss, err := pluginsReg.Get(sss.Name(), registry.WithType(pb.ItemType_SERVICE))
							if err != nil {
								fmt.Println(err)
								continue
							}

							var s service.Service
							if ss.As(&s) {
								opts := s.Options()

								opts.Context = ctx

								s.Start()
							}
						}
					}
				}
			}()



			srvGRPC := grpc.New(ctx)
			var srvHTTP server.Server
			if !runtime.IsFork() {
				if h, err := caddy.New(ctx, ""); err != nil {
					return err
				} else {
					srvHTTP = h
				}
			} else {
				srvHTTP = http.New(ctx)
			}
			if err != nil {
				return err
			}
			srvGeneric := generic.New(ctx)

		*/

		//ctx = servicecontext.WithServer(ctx, "grpc", srvGRPC)
		//ctx = servicecontext.WithServer(ctx, "http", srvHTTP)
		//ctx = servicecontext.WithServer(ctx, "generic", srvGeneric)

		plugins.Init(ctx, "main")

		services, err := pluginsReg.List(registry.WithType(pb.ItemType_SERVICE))
		if err != nil {
			return err
		}

		var (
			srvGRPC    server.Server
			srvHTTP    server.Server
			srvGeneric server.Server
			srvs       []server.Server
		)

		for _, ss := range services {
			if !runtime.IsRequired(ss.Name()) {
				continue
			}

			var s service.Service
			if ss.As(&s) {
				opts := s.Options()

				opts.Context = servicecontext.WithRegistry(opts.Context, reg)

				if opts.Fork && !runtime.IsFork() {
					if !opts.AutoStart {
						continue
					}

					srvFork := fork.NewServer(opts.Context)
					var srvForkAs *fork.ForkServer
					if srvFork.As(&srvForkAs) {
						srvForkAs.RegisterForkParam(opts.Name)
					}

					srvs = append(srvs, srvFork)

					opts.Server = srvFork

					continue
				}

				if s.IsGRPC() {
					if srvGRPC == nil {
						srvGRPC = grpc.New(ctx)
						srvs = append(srvs, srvGRPC)
					}

					opts.Server = srvGRPC
				}

				if s.IsREST() {
					if srvHTTP == nil {
						if runtime.IsFork() {
							srvHTTP = http.New(ctx)
						} else {
							srvHTTP, _ = caddy.New(opts.Context, "")
						}
						srvs = append(srvs, srvHTTP)
					}

					opts.Server = srvHTTP
				}

				if s.IsGeneric() {
					if srvGeneric == nil {
						srvGeneric = generic.New(ctx)
						srvs = append(srvs, srvGeneric)
					}

					opts.Server = srvGeneric
				}

				// Checking which service is needed
				bs, ok := opts.Server.(server.WrappedServer)
				if ok {
					bs.RegisterBeforeServe(s.Start)
					bs.RegisterAfterServe(func() error {
						// Register service again to update nodes information
						if err := reg.Register(s); err != nil {
							return err
						}
						return nil
					})
					bs.RegisterBeforeStop(s.Stop)
				}
			}
		}

		var g errgroup.Group
		for _, srv := range srvs {
			g.Go(srv.Serve)
		}

		if err := g.Wait(); err != nil {
			return err
		}

		select {
		case <-cmd.Context().Done():
		}

		return nil
	},
}

func init() {
	// Flags for selecting / filtering services
	StartCmd.Flags().StringArrayVarP(&FilterStartTags, "tags", "t", []string{}, "Select services to start by tags, possible values are 'broker', 'data', 'datasource', 'discovery', 'frontend', 'gateway', 'idm', 'scheduler'")
	StartCmd.Flags().StringArrayVarP(&FilterStartExclude, "exclude", "x", []string{}, "Select services to start by filtering out some specific ones by name")

	StartCmd.Flags().String("grpc.address", ":8001", "gRPC Server Address")
	StartCmd.Flags().String("http.address", ":8002", "HTTP Server Address")

	StartCmd.Flags().Bool("fork", false, "Used internally by application when forking processes")

	addRegistryFlags(StartCmd.Flags())

	StartCmd.Flags().MarkHidden("fork")
	StartCmd.Flags().MarkHidden("registry")
	StartCmd.Flags().MarkHidden("broker")

	RootCmd.AddCommand(StartCmd)
}
