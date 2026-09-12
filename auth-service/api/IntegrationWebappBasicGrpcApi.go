package api

import (
	"context"
	"log/slog"

	v1 "github.com/weoses/memelo/gen/proto/v1"
	"github.com/weoses/memelo/gen/proto/v1/v1connect"

	"github.com/weoses/memelo/auth-service/service"
)

type IntegrationWebappBasicGrpcApiImpl struct {
	slogger *slog.Logger
	svc     service.IntegrationWebappBasicService
}

func NewIntegrationWebappBasicGrpcApi(svc service.IntegrationWebappBasicService) v1connect.IntegrationWebappBasicServiceHandler {
	return &IntegrationWebappBasicGrpcApiImpl{
		slogger: slog.With("service", "IntegrationWebappBasicGrpcApi"),
		svc:     svc,
	}
}

func (a *IntegrationWebappBasicGrpcApiImpl) Authorize(ctx context.Context, req *v1.WebappBasicAuthorizeRequest) (*v1.WebappBasicAuthorizeResponse, error) {
	result, err := a.svc.Authorize(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.WebappBasicAuthorizeResponse{Result: toProtoAuthResult(result)}, nil
}
