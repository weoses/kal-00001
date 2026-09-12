package functional_test

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	v1 "github.com/weoses/memelo/gen/proto/v1"
	"github.com/weoses/memelo/gen/proto/v1/v1connect"
	"golang.org/x/net/http2"
)

// h2cClient creates an HTTP client that speaks HTTP/2 cleartext (h2c).
func h2cClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
}

type TestClient struct {
	telegram    v1connect.IntegrationTelegramServiceClient
	webappBasic v1connect.IntegrationWebappBasicServiceClient
}

func NewTestClient(addr string) *TestClient {
	c := h2cClient()
	return &TestClient{
		telegram:    v1connect.NewIntegrationTelegramServiceClient(c, addr),
		webappBasic: v1connect.NewIntegrationWebappBasicServiceClient(c, addr),
	}
}

func (c *TestClient) AuthorizeTelegram(ctx context.Context, telegramId int64) (*v1.TelegramAuthorizeResponse, error) {
	return c.telegram.Authorize(ctx, &v1.TelegramAuthorizeRequest{TelegramId: telegramId})
}

func (c *TestClient) AuthorizeWebappBasic(ctx context.Context, username, password string) (*v1.WebappBasicAuthorizeResponse, error) {
	return c.webappBasic.Authorize(ctx, &v1.WebappBasicAuthorizeRequest{Username: username, Password: password})
}
