package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/weoses/memelo/auth-service/storage"
	"golang.org/x/crypto/bcrypt"
)

type IntegrationWebappBasicService interface {
	Authorize(ctx context.Context, username string, password string) (*AuthResult, error)
}

type IntegrationWebappBasicServiceImpl struct {
	links storage.IntegrationWebappBasicStorage
	users storage.UserStorage
}

func NewIntegrationWebappBasicService(links storage.IntegrationWebappBasicStorage, users storage.UserStorage) IntegrationWebappBasicService {
	return &IntegrationWebappBasicServiceImpl{links: links, users: users}
}

func (s *IntegrationWebappBasicServiceImpl) Authorize(ctx context.Context, username string, password string) (*AuthResult, error) {
	link, err := s.links.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return notFoundResult(), nil
		}
		return nil, fmt.Errorf("find webapp basic link failed: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(link.PasswordHash), []byte(password)) != nil {
		return notFoundResult(), nil
	}

	return resultForUser(ctx, s.users, link.UserId)
}
