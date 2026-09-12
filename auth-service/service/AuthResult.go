package service

import (
	"context"
	"fmt"

	"github.com/weoses/memelo/auth-service/storage"
)

type AuthStatus int

const (
	AuthStatusNotFound AuthStatus = iota
	AuthStatusOk
)

// AuthResult is the service layer's own DTO — proto types must not leak in
// here, and this type must not be returned directly from the api/ layer.
type AuthResult struct {
	Status      AuthStatus
	UserId      string
	AccountId   string
	Permissions []string
}

func notFoundResult() *AuthResult {
	return &AuthResult{Status: AuthStatusNotFound}
}

// resultForUser loads a user's permissions and assembles the AuthResult
// shared by every integration's Authorize implementation.
func resultForUser(ctx context.Context, users storage.UserStorage, userId string) (*AuthResult, error) {
	user, err := users.GetById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	permissions, err := users.GetPermissionCodes(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("get permission codes failed: %w", err)
	}

	return &AuthResult{
		Status:      AuthStatusOk,
		UserId:      user.Id,
		AccountId:   user.AccountId,
		Permissions: permissions,
	}, nil
}
