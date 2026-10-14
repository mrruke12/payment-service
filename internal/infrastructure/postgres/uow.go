package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type repositoryProviderFactory[RP any] func(tx Querier) RP

type UnitOfWork[RP any] struct {
	pool    *pgxpool.Pool
	factory repositoryProviderFactory[RP]
}

func (u *UnitOfWork[RP]) Do(
	ctx context.Context,
	fn func(ctx context.Context, r RP) error,
) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	repos := u.factory(tx)

	if err := fn(ctx, repos); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
