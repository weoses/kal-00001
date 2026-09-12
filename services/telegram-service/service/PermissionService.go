package service

import (
	"context"
	"errors"
	"log/slog"
)

const (
	PermissionCreate    = "create"
	PermissionDelete    = "delete"
	PermissionRecompute = "recompute"
	PermissionSearch    = "search"
)

// permissionCodes maps a telegram-service action to the permission code
// auth-service grants a user. PermissionSearch has no entry: search stays
// open to any resolved (non-blocked) user, matching the read-only nature of
// browsing the (possibly public) dataset.
var permissionCodes = map[string]string{
	PermissionCreate:    "CREATE",
	PermissionDelete:    "DELETE",
	PermissionRecompute: "UPDATE",
}

var ErrForbidden = errors.New("forbidden")

type PermissionService interface {
	IsAllowed(ctx context.Context, userId int64, permission string) bool
}

// InvokeWithPermission checks the named permission for userId before calling fn.
// Returns ErrForbidden if the user is not allowed.
func InvokeWithPermission[T any](
	ctx context.Context,
	svc PermissionService,
	userId int64,
	permission string,
	fn func() (T, error),
) (T, error) {
	if !svc.IsAllowed(ctx, userId, permission) {
		var zero T
		return zero, ErrForbidden
	}
	return fn()
}

type PermissionServiceImpl struct {
	auth AuthConnector
	log  *slog.Logger
}

func NewPermissionService(auth AuthConnector) PermissionService {
	return &PermissionServiceImpl{auth: auth, log: slog.With("service", "PermissionService")}
}

func (p *PermissionServiceImpl) IsAllowed(ctx context.Context, userId int64, permission string) bool {
	code, needsPermission := permissionCodes[permission]
	if !needsPermission {
		return true
	}

	result, err := p.auth.AuthorizeTelegram(ctx, userId)
	if err != nil {
		p.log.ErrorContext(ctx, "authorize failed, denying by default", "userId", userId, "permission", permission, "error", err)
		return false
	}

	return result.HasPermission(code)
}
