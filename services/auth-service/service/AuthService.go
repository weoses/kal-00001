package service

import (
	"context"
	"errors"

	"github.com/weoses/memelo/auth-service/storage"
)

// AuthService is a generic cross-integration user lookup, used for
// admin/debug tooling rather than on the hot authorize-a-request path.
type AuthService interface {
	GetUser(ctx context.Context, userId string) (*AuthResult, error)
}

type AuthServiceImpl struct {
	users storage.UserStorage
}

func NewAuthService(users storage.UserStorage) AuthService {
	return &AuthServiceImpl{users: users}
}

func (s *AuthServiceImpl) GetUser(ctx context.Context, userId string) (*AuthResult, error) {
	result, err := resultForUser(ctx, s.users, userId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return notFoundResult(), nil
		}
		return nil, err
	}
	return result, nil
}
