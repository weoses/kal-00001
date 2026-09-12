package api

import (
	"context"
	"log/slog"

	v1 "github.com/weoses/memelo/gen/proto/v1"
	"github.com/weoses/memelo/gen/proto/v1/v1connect"

	"github.com/weoses/memelo/auth-service/service"
)

type AuthGrpcApiImpl struct {
	slogger     *slog.Logger
	authService service.AuthService
}

func NewAuthGrpcApi(authService service.AuthService) v1connect.AuthServiceHandler {
	return &AuthGrpcApiImpl{
		slogger:     slog.With("service", "AuthGrpcApi"),
		authService: authService,
	}
}

func (a *AuthGrpcApiImpl) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	result, err := a.authService.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserResponse{Result: toProtoAuthResult(result)}, nil
}
