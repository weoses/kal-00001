package storage

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/weoses/memelo/auth-service/conf"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations applies all pending golang-migrate migrations embedded in
// this binary. It is invoked once at startup, mirroring storage-service's
// auto-migrate-on-boot behavior.
func RunMigrations(cfg *conf.Config, logger *slog.Logger) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("load embedded migrations failed: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, cfg.Postgres.Dsn)
	if err != nil {
		return fmt.Errorf("create migrator failed: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations failed: %w", err)
	}

	logger.Info("migrations applied")
	return nil
}
