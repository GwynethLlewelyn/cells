package grpc

import (
	"context"
	"strings"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	metadata2 "google.golang.org/grpc/metadata"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/auth"
	"github.com/pydio/cells/v4/common/config"
	"github.com/pydio/cells/v4/common/crypto/providers"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/plugins"
	"github.com/pydio/cells/v4/common/proto/install"
	"github.com/pydio/cells/v4/common/proto/tree"
	"github.com/pydio/cells/v4/common/server"
	grpc2 "github.com/pydio/cells/v4/common/server/grpc"
	"github.com/pydio/cells/v4/common/service"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"github.com/pydio/cells/v4/common/service/context/metadata"
)

func init() {

	handlersRegister := func(runtimeCtx context.Context, g *grpc.Server, clear bool) {
		h := &TreeHandler{runtimeCtx: runtimeCtx}
		if clear {
			h.name = common.ServiceGatewayGrpcClear
		} else {
			h.name = common.ServiceGatewayGrpc
		}
		// Do not use Enhanced here
		tree.RegisterNodeProviderServer(g, h)
		tree.RegisterNodeReceiverServer(g, h)
		tree.RegisterNodeChangesStreamerServer(g, h)
		tree.RegisterNodeProviderStreamerServer(g, h)
		tree.RegisterNodeReceiverStreamServer(g, h)
	}

	// Build options - optionally force port
	baseOpts := []service.ServiceOption{
		service.Tag(common.ServiceTagGateway),
		service.Dependency(common.ServiceGrpcNamespace_+common.ServiceTree, []string{}),
		service.Dependency(common.ServiceGatewayProxy, []string{}),
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
			clearOpts = append(clearOpts,
				service.Context(ctx),
				service.WithServerProvider(createServerProvider(false)),
				service.WithGRPC(func(runtimeCtx context.Context, srv *grpc.Server) error {
					handlersRegister(runtimeCtx, srv, true)
					return nil
				}),
			)
			service.NewService(clearOpts...)
		}
		if hasTls {
			tlsOpts = append(tlsOpts,
				service.Context(ctx),
				service.WithServerProvider(createServerProvider(true)),
				service.WithGRPC(func(runtimeCtx context.Context, srv *grpc.Server) error {
					handlersRegister(runtimeCtx, srv, false)
					return nil
				}),
			)
			service.NewService(tlsOpts...)
		}
	})

}

func createServerProvider(tls bool) service.ServerProvider {
	return func(ctx context.Context) (server.Server, error) {
		grpcOptions := []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				servicecontext.ContextUnaryServerInterceptor(servicecontext.SpanIncomingContext),
				servicecontext.MetricsUnaryServerInterceptor(),
				servicecontext.ContextUnaryServerInterceptor(servicecontext.MetaIncomingContext),
				servicecontext.ContextUnaryServerInterceptor(jwtCtxModifier),
				servicecontext.ContextUnaryServerInterceptor(grpcMetaCtxModifier),
			),
			grpc.ChainStreamInterceptor(
				servicecontext.ContextStreamServerInterceptor(servicecontext.SpanIncomingContext),
				servicecontext.MetricsStreamServerInterceptor(),
				servicecontext.ContextStreamServerInterceptor(servicecontext.MetaIncomingContext),
				servicecontext.ContextStreamServerInterceptor(jwtCtxModifier),
				servicecontext.ContextStreamServerInterceptor(grpcMetaCtxModifier),
			),
		}
		if tls {
			localConfig := &install.ProxyConfig{
				Binds:     []string{"localhost"},
				TLSConfig: &install.ProxyConfig_SelfSigned{SelfSigned: &install.TLSSelfSigned{}},
			}
			tlsConfig, e := providers.LoadTLSServerConfig(localConfig)
			if e != nil {
				return nil, e
			}
			grpcOptions = append(grpcOptions, grpc.Creds(credentials.NewTLS(tlsConfig)))
		}

		srv := grpc.NewServer(grpcOptions...)
		addr := ":0" // Will pick a random port
		if !tls {
			if port := viper.GetString("grpc_external"); port != "" {
				addr = ":" + port
			}
			logCtx := servicecontext.WithServiceName(ctx, common.ServiceGatewayGrpcClear)
			log.Logger(logCtx).Info("Starting HTTP only gRPC gateway. Will be accessed directly through " + addr)
		} else {
			logCtx := servicecontext.WithServiceName(ctx, common.ServiceGatewayGrpc)
			log.Logger(logCtx).Info("Activating self-signed configuration for gRPC gateway to allow full TLS chain.")
		}
		return grpc2.NewWithServer(ctx, srv, addr), nil
	}
}

// jwtCtxModifier extracts x-pydio-bearer metadata to validate authentication
func jwtCtxModifier(ctx context.Context) (context.Context, bool, error) {

	jwtVerifier := auth.DefaultJWTVerifier()
	meta, ok := metadata2.FromIncomingContext(ctx)
	if !ok {
		return ctx, false, nil
	}
	bearer := strings.Join(meta.Get("x-pydio-bearer"), "")
	if bearer == "" {
		return ctx, false, nil
	}

	if ct, _, err := jwtVerifier.Verify(ctx, bearer); err != nil {
		log.Auditer(ctx).Error("Blocked invalid JWT", log.GetAuditId(common.AuditInvalidJwt))
		return ctx, false, err
	} else {
		log.Logger(ctx).Debug("Got valid Claims from Bearer!")
		return ct, true, nil
	}

}

// grpcMetaCtxModifier extracts specific meta from IncomingContext
func grpcMetaCtxModifier(ctx context.Context) (context.Context, bool, error) {

	meta := map[string]string{}
	if existing, ok := metadata2.FromIncomingContext(ctx); ok {
		translate := map[string]string{
			"user-agent":      servicecontext.HttpMetaUserAgent,
			"content-type":    servicecontext.HttpMetaContentType,
			"x-forwarded-for": servicecontext.HttpMetaRemoteAddress,
			"x-pydio-span-id": servicecontext.SpanMetadataId,
		}
		for k, v := range existing {
			if k == ":authority" { // Ignore grpc-specific meta
				continue
			}
			if newK, ok := translate[k]; ok {
				meta[newK] = strings.Join(v, "")
			} else {
				meta[k] = strings.Join(v, "")
			}
		}
		// Override with specific header
		if ua, ok := existing["x-pydio-grpc-user-agent"]; ok {
			meta[servicecontext.HttpMetaUserAgent] = strings.Join(ua, "")
		}
	}
	meta[servicecontext.HttpMetaExtracted] = servicecontext.HttpMetaExtracted
	layout := "2006-01-02T15:04-0700"
	t := time.Now()
	meta[servicecontext.ServerTime] = t.Format(layout)
	// We currently use server time instead of client time. TODO: Retrieve client time and locale and set it here.
	meta[servicecontext.ClientTime] = t.Format(layout)

	return metadata.NewContext(ctx, meta), true, nil
}
