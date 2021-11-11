package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/auth"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/crypto/providers"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/proto/install"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/service"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/common/utils/net"
)

func init() {

	grpcServerWithLog := func(ctx context.Context, g *grpc.Server, msg string) error {
		log.Logger(ctx).Info(msg)
		/*
			// TODO V4
			m.Init(micro.WrapHandler(jwtWrapper(m.Options().Context), httpMetaWrapper()))
		*/
		h := &TreeHandler{}
		tree.RegisterNodeProviderServer(g, h)
		/*
			TODO V4
			tree.RegisterNodeReceiverServer(g, h)
			tree.RegisterNodeChangesStreamerServer(g, h)
			tree.RegisterNodeProviderStreamerServer(g, h)
			tree.RegisterNodeReceiverStreamServer(g, h)
		*/
		return nil
	}

	// Build options - optionally force port
	baseOpts := []service.ServiceOption{
		service.Tag(common.ServiceTagGateway),
		service.Dependency(common.ServiceGrpcNamespace_+common.ServiceTree, []string{}),
		service.Dependency(common.ServiceGatewayProxy, []string{}),
		/*
			service.WithMicro(func(m micro.Service) error {
				m.Init(micro.WrapHandler(jwtWrapper(m.Options().Context), httpMetaWrapper()))
				h := &TreeHandler{}
				srv := m.Options().Server
				tree.RegisterNodeProviderHandler(srv, h)
				tree.RegisterNodeReceiverHandler(srv, h)
				tree.RegisterNodeChangesStreamerHandler(srv, h)
				tree.RegisterNodeProviderStreamerHandler(srv, h)
				tree.RegisterNodeReceiverStreamHandler(srv, h)
				return nil
			}),

		*/
	}
	tlsOpts := append(baseOpts,
		service.Name(common.ServiceGatewayGrpc),
		service.Description("External gRPC Access (tls)"),
	)
	clearOpts := append(baseOpts,
		service.Name(common.ServiceGatewayGrpcClear),
		service.Description("External gRPC Access (clear)"),
	)
	plugins.Register("main", func(ctx context.Context) {
		ss, _ := config.LoadSites()
		var hasClear, hasTls bool
		for _, s := range ss {
			if s.HasTLS() {
				hasTls = true
			} else {
				hasClear = true
			}
		}
		if hasClear {
			var p string
			if port := viper.GetString("grpc_external"); port != "" {
				p = port
			} else {
				p = fmt.Sprintf("%d", net.GetAvailablePort())
			}
			logCtx := servicecontext.WithServiceName(ctx, common.ServiceGatewayGrpcClear)
			clearOpts = append(clearOpts,
				// service.Port(p), TODO V4
				service.Context(ctx),
				service.WithGRPC(func(g *grpc.Server) error {
					return grpcServerWithLog(logCtx, g, "Starting HTTP only gRPC gateway. Will be accessed directly through port "+p)
				}),
				//service.WithMicro(microServiceWithLog(logCtx, "Starting HTTP only gRPC gateway. Will be accessed directly through port "+p)),
			)
			service.NewService(clearOpts...)
		}
		if hasTls {
			logCtx := servicecontext.WithServiceName(ctx, common.ServiceGatewayGrpc)
			tlsOpts = append(tlsOpts,
				service.Context(ctx),
				service.WithGRPC(func(g *grpc.Server) error {
					return grpcServerWithLog(logCtx, g, "Activating self-signed configuration for gRPC gateway to allow full TLS chain.")
				}),
			)
			localConfig := &install.ProxyConfig{
				Binds:     []string{common.ServiceGatewayGrpc},
				TLSConfig: &install.ProxyConfig_SelfSigned{SelfSigned: &install.TLSSelfSigned{}},
			}
			if tls, e := providers.LoadTLSServerConfig(localConfig); e == nil {
				tlsOpts = append(tlsOpts, service.WithTLSConfig(tls))
			}
			service.NewService(tlsOpts...)
		}
		/*
			if len(ss) == 1 && !ss[0].HasTLS() {
				// This is a simple config without TLS - Access will be direct not through proxy
				//fmt.Println("[NO-TLS] " + common.ServiceGatewayGrpc + " served as HTTP and should be accessed directly (no TLS)")
				if port := viper.Get("grpc_external"); port != nil {
					log.Logger(ctx).Info("Using HTTP configuration for gRPC gateway. Should be accessed directly through port " + port.(string))
					tlsOpts = append(tlsOpts, service.Port(port.(string)))
				} else {
					log.Logger(ctx).Info("Using HTTP configuration for gRPC gateway. Should be accessed directly on this service port")
				}
			} else {
				log.Logger(ctx).Info("Activating self-signed configuration for gRPC gateway to allow full TLS chain.")
				localConfig := &install.ProxyConfig{
					Binds:     []string{common.ServiceGatewayGrpc},
					TLSConfig: &install.ProxyConfig_SelfSigned{SelfSigned: &install.TLSSelfSigned{}},
				}
				if tls, e := providers.LoadTLSServerConfig(localConfig); e == nil {
					//fmt.Println("[TLS] Activating self-signed TLS on " + common.ServiceGatewayGrpc)
					tlsOpts = append(tlsOpts, service.WithTLSConfig(tls))
				}
			}
			tlsOpts = append(tlsOpts, service.Context(ctx))
			service.NewService(tlsOpts...)
		*/
	})

}

// jwtWrapper extracts x-pydio-bearer metadata to validate authentication
func jwtWrapper(serviceCtx context.Context) func(handlerFunc server.HandlerFunc) server.HandlerFunc {

	jwtVerifier := auth.DefaultJWTVerifier()

	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {

		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			if meta, ok := metadata.FromContext(ctx); ok {

				bearer, o := meta["x-pydio-bearer"] //strings.Join(meta.Get("x-pydio-bearer"), "")
				if o {
					var err error
					ctx, _, err = jwtVerifier.Verify(ctx, bearer)
					if err != nil {
						log.Auditer(serviceCtx).Error(
							"Blocked invalid JWT",
							log.GetAuditId(common.AuditInvalidJwt),
						)
						return err
					} else {
						log.Logger(serviceCtx).Debug("Got valid Claims from Bearer!")
					}
				}
			}

			return handlerFunc(ctx, req, rsp)

		}
	}

}

// httpMetaWrapper translates gRPC meta headers (lowercase x-header-name) to standard cells metadata
func httpMetaWrapper() func(handlerFunc server.HandlerFunc) server.HandlerFunc {

	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			return handlerFunc(ctxRequestInfoToMetadata(ctx), req, rsp)

		}
	}
}

func ctxRequestInfoToMetadata(ctx context.Context) context.Context {

	meta := metadata.Metadata{}
	if existing, ok := metadata.FromContext(ctx); ok {
		if _, already := existing[servicecontext.HttpMetaExtracted]; already {
			return ctx
		}
		translate := map[string]string{
			"user-agent":      servicecontext.HttpMetaUserAgent,
			"content-type":    servicecontext.HttpMetaContentType,
			"x-forwarded-for": servicecontext.HttpMetaRemoteAddress,
			"x-pydio-span-id": servicecontext.SpanMetadataId,
		}
		for k, v := range existing {
			if newK, ok := translate[k]; ok {
				meta[newK] = v
			} else {
				meta[k] = v
			}
		}
		// Override with specific header
		if ua, ok := existing["x-pydio-grpc-user-agent"]; ok {
			meta[servicecontext.HttpMetaUserAgent] = ua
		}
	}
	meta[servicecontext.HttpMetaExtracted] = servicecontext.HttpMetaExtracted
	layout := "2006-01-02T15:04-0700"
	t := time.Now()
	meta[servicecontext.ServerTime] = t.Format(layout)
	// We currently use server time instead of client time. TODO: Retrieve client time and locale and set it here.
	meta[servicecontext.ClientTime] = t.Format(layout)

	return metadata.NewContext(ctx, meta)
}
