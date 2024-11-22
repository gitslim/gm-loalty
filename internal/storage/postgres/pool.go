package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/gitslim/gophermart/internal/conf"
	"github.com/gitslim/gophermart/internal/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewConnPool(config *conf.Config) (*pgxpool.Pool, error) {
	// По дефолту запросы подготавливаются и кэшируются: default_query_exec_mode=cache_statement
	c, err := pgxpool.ParseConfig(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres database dsn: %w", err)
	}

	// настройка пула
	c.MaxConns = 10
	c.MinConns = 2
	c.MaxConnIdleTime = 5 * time.Minute

	// Подключение к бд
	pool, err := pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("failed create postgres connection pool: %w", err)
	}
	return pool, nil
}

func RegisterPoolHooks(lc fx.Lifecycle, cfg *conf.Config, log logging.Logger, pool *pgxpool.Pool) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// пингуем бд
			if err := pool.Ping(context.Background()); err != nil {
				return fmt.Errorf("postgres connection error: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// закрываем соединение
			pool.Close()
			return nil
		},
	})
}
