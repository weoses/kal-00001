package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/weoses/memelo/gen/proto/v1/v1connect"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/weoses/memelo/auth-service/api"
	"github.com/weoses/memelo/auth-service/conf"
	"github.com/weoses/memelo/auth-service/middleware"
	"github.com/weoses/memelo/auth-service/service"
	"github.com/weoses/memelo/auth-service/storage"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(storage.NewGormDb),
		fx.Provide(storage.NewUserStorage),
		fx.Provide(storage.NewIntegrationTelegramStorage),
		fx.Provide(storage.NewIntegrationWebappBasicStorage),

		fx.Provide(service.NewAuthService),
		fx.Provide(service.NewIntegrationTelegramService),
		fx.Provide(service.NewIntegrationWebappBasicService),

		fx.Provide(api.NewAuthGrpcApi),
		fx.Provide(api.NewIntegrationTelegramGrpcApi),
		fx.Provide(api.NewIntegrationWebappBasicGrpcApi),

		fx.Provide(func(c *conf.Config) (net.Listener, error) {
			return net.Listen("tcp", c.Server.ListenAddress)
		}),

		fx.Invoke(func(cfg *conf.Config) error {
			return storage.RunMigrations(cfg, slog.With("service", "Migrator"))
		}),
		fx.Invoke(startup),
	)
}

func startup(
	lc fx.Lifecycle,
	ln net.Listener,
	authApi v1connect.AuthServiceHandler,
	telegramApi v1connect.IntegrationTelegramServiceHandler,
	webappBasicApi v1connect.IntegrationWebappBasicServiceHandler,
) {
	interceptors := connect.WithInterceptors(middleware.NewLoggingInterceptor(slog.With("service", "router")))

	mux := http.NewServeMux()
	mux.Handle(v1connect.NewAuthServiceHandler(authApi, interceptors))
	mux.Handle(v1connect.NewIntegrationTelegramServiceHandler(telegramApi, interceptors))
	mux.Handle(v1connect.NewIntegrationWebappBasicServiceHandler(webappBasicApi, interceptors))
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		WriteTimeout: time.Second * 30,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() { _ = srv.Serve(ln) }()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
