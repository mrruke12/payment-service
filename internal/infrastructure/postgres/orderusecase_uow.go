package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrruke12/payment-service/internal/application/orderusecase"
)

func NewRepositoryProvider(tx Querier) orderusecase.RepositoryProvider {
	return &RepositoryProvider{
		tx: tx,
	}
}

func NewCreateOrderUnitOfWork(pool *pgxpool.Pool) orderusecase.UnitOfWork {
	return &UnitOfWork[orderusecase.RepositoryProvider]{
		pool:    pool,
		factory: NewRepositoryProvider,
	}
}
