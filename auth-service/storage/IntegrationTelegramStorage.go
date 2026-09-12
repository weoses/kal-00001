package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/weoses/memelo/auth-service/entity"
	"gorm.io/gorm"
)

type IntegrationTelegramStorage interface {
	// FindByTelegramId returns the link row whose TelegramIds array contains
	// telegramId, or ErrNotFound if that id isn't linked to any user.
	FindByTelegramId(ctx context.Context, telegramId int64) (*entity.IntegrationTelegram, error)
}

type IntegrationTelegramStorageImpl struct {
	db *gorm.DB
}

func NewIntegrationTelegramStorage(db *gorm.DB) IntegrationTelegramStorage {
	return &IntegrationTelegramStorageImpl{db: db}
}

func (s *IntegrationTelegramStorageImpl) FindByTelegramId(ctx context.Context, telegramId int64) (*entity.IntegrationTelegram, error) {
	var link entity.IntegrationTelegram
	err := s.db.WithContext(ctx).
		Where("telegram_ids @> ARRAY[?]::bigint[]", telegramId).
		First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find by telegram id failed: %w", err)
	}
	return &link, nil
}
