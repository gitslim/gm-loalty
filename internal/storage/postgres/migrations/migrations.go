package migrations

import (
	"errors"
	"fmt"

	"github.com/gitslim/gophermart/internal/conf"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations запускает миграции базы данных
func RunMigrations(config *conf.Config) error {
	m, err := migrate.New(
		"file://migrations",
		config.DatabaseURI,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return nil
}
