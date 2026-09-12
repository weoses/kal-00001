package storage

import (
	"fmt"
	"log/slog"

	"github.com/weoses/memelo/auth-service/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewGormDb(cfg *conf.Config) (*gorm.DB, error) {
	// A missing Telegram/webapp link is an expected "not found" result, not
	// an error — don't let GORM log it at error level on every anonymous lookup.
	gormLog := gormlogger.NewSlogLogger(slog.With("service", "gorm"), gormlogger.Config{
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: true,
	})
	db, err := gorm.Open(postgres.Open(cfg.Postgres.Dsn), &gorm.Config{Logger: gormLog})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres failed: %w", err)
	}
	return db, nil
}
