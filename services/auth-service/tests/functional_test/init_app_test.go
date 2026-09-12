package functional_test

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/weoses/memelo/auth-service/app"
	"github.com/weoses/memelo/auth-service/conf"
	"go.uber.org/fx"
)

func startTestServer(cfg *conf.Config) (addr string, stop func(), err error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", nil, fmt.Errorf("listen: %w", err)
	}
	addr = fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port)

	fxapp := fx.New(
		fx.Supply(cfg),
		app.Module(),
		fx.Decorate(func() net.Listener { return ln }),
	)

	if err := fxapp.Start(context.Background()); err != nil {
		return "", nil, fmt.Errorf("start app: %w", err)
	}

	stop = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = fxapp.Stop(ctx)
	}
	return addr, stop, nil
}
