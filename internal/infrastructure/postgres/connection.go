package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrruke12/payment-service/internal/infrastructure/config"
)

func ConnString(cfg config.DBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSL,
	)
}

func Connect(ctx context.Context, cfg config.DBConfig) *pgxpool.Pool {
	_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	connString := ConnString(cfg)

	connConfig, err := pgxpool.ParseConfig(connString)

	if err != nil {
		panic(err)
	}

	connConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(
		_ctx,
		connConfig,
	)

	if err != nil {
		panic(err)
	}

	return pool
}

func Health(ctx context.Context, pool *pgxpool.Pool) error {
	_ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	return pool.Ping(_ctx)
}
