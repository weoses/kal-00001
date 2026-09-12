package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/weoses/memelo/auth-service/storage"
)

type IntegrationTelegramService interface {
	// Authorize resolves a Telegram numeric user id to a user. An id not
	// linked to any user comes back as AuthStatusNotFound (anonymous),
	// even if the same person is a registered user via a different linked id.
	Authorize(ctx context.Context, telegramId int64) (*AuthResult, error)
}

type IntegrationTelegramServiceImpl struct {
	links storage.IntegrationTelegramStorage
	users storage.UserStorage
}

func NewIntegrationTelegramService(links storage.IntegrationTelegramStorage, users storage.UserStorage) IntegrationTelegramService {
	return &IntegrationTelegramServiceImpl{links: links, users: users}
}

func (s *IntegrationTelegramServiceImpl) Authorize(ctx context.Context, telegramId int64) (*AuthResult, error) {
	link, err := s.links.FindByTelegramId(ctx, telegramId)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return notFoundResult(), nil
		}
		return nil, fmt.Errorf("find telegram link failed: %w", err)
	}

	return resultForUser(ctx, s.users, link.UserId)
}
