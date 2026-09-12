package storage

import (
	"fmt"

	"github.com/weoses/memelo/auth-service/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDb(cfg *conf.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres failed: %w", err)
	}
	return db, nil
}
