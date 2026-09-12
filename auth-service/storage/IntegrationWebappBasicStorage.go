package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/weoses/memelo/auth-service/entity"
	"gorm.io/gorm"
)

type IntegrationWebappBasicStorage interface {
	FindByUsername(ctx context.Context, username string) (*entity.IntegrationWebappBasic, error)
}

type IntegrationWebappBasicStorageImpl struct {
	db *gorm.DB
}

func NewIntegrationWebappBasicStorage(db *gorm.DB) IntegrationWebappBasicStorage {
	return &IntegrationWebappBasicStorageImpl{db: db}
}

func (s *IntegrationWebappBasicStorageImpl) FindByUsername(ctx context.Context, username string) (*entity.IntegrationWebappBasic, error) {
	var link entity.IntegrationWebappBasic
	err := s.db.WithContext(ctx).First(&link, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find by username failed: %w", err)
	}
	return &link, nil
}
