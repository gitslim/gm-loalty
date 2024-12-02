package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/gitslim/gophermart/internal/conf"
	"github.com/gitslim/gophermart/internal/logging"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

func NewConnPool(config *conf.Config) (*sqlx.DB, error) {
	// Создаем подключение через sqlx
	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настраиваем пул
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

func RegisterPoolHooks(lc fx.Lifecycle, cfg *conf.Config, log logging.Logger, db *sqlx.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// пингуем бд
			if err := db.PingContext(context.Background()); err != nil {
				return fmt.Errorf("postgres connection error: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// закрываем соединение
			db.Close()
			return nil
		},
	})
}
