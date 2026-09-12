package api

import (
	"context"
	"log/slog"

	v1 "github.com/weoses/memelo/gen/proto/v1"
	"github.com/weoses/memelo/gen/proto/v1/v1connect"

	"github.com/weoses/memelo/auth-service/service"
)

type IntegrationTelegramGrpcApiImpl struct {
	slogger *slog.Logger
	svc     service.IntegrationTelegramService
}

func NewIntegrationTelegramGrpcApi(svc service.IntegrationTelegramService) v1connect.IntegrationTelegramServiceHandler {
	return &IntegrationTelegramGrpcApiImpl{
		slogger: slog.With("service", "IntegrationTelegramGrpcApi"),
		svc:     svc,
	}
}

func (a *IntegrationTelegramGrpcApiImpl) Authorize(ctx context.Context, req *v1.TelegramAuthorizeRequest) (*v1.TelegramAuthorizeResponse, error) {
	result, err := a.svc.Authorize(ctx, req.TelegramId)
	if err != nil {
		return nil, err
	}
	return &v1.TelegramAuthorizeResponse{Result: toProtoAuthResult(result)}, nil
}
