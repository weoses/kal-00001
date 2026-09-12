package api

import (
	v1 "github.com/weoses/memelo/gen/proto/v1"

	"github.com/weoses/memelo/auth-service/service"
)

func toProtoAuthResult(r *service.AuthResult) *v1.AuthResult {
	return &v1.AuthResult{
		Status:      toProtoStatus(r.Status),
		UserId:      r.UserId,
		AccountId:   r.AccountId,
		Permissions: r.Permissions,
	}
}

func toProtoStatus(s service.AuthStatus) v1.AuthStatus {
	if s == service.AuthStatusOk {
		return v1.AuthStatus_AUTH_STATUS_OK
	}
	return v1.AuthStatus_AUTH_STATUS_NOT_FOUND
}
