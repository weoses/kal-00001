package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/weoses/memelo/auth-service/entity"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type UserStorage interface {
	GetById(ctx context.Context, id string) (*entity.User, error)
	GetPermissionCodes(ctx context.Context, userId string) ([]string, error)
}

type UserStorageImpl struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return &UserStorageImpl{db: db}
}

func (s *UserStorageImpl) GetById(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user by id failed: %w", err)
	}
	return &user, nil
}

func (s *UserStorageImpl) GetPermissionCodes(ctx context.Context, userId string) ([]string, error) {
	var codes []string
	err := s.db.WithContext(ctx).
		Model(&entity.Permission{}).
		Joins("JOIN user_permissions up ON up.permission_id = permissions.id").
		Where("up.user_id = ?", userId).
		Pluck("permissions.code", &codes).Error
	if err != nil {
		return nil, fmt.Errorf("get permission codes failed: %w", err)
	}
	return codes, nil
}
