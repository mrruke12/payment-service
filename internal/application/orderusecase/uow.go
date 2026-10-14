package orderusecase

import (
	"context"

	"github.com/mrruke12/payment-service/internal/domains/order"
	"github.com/mrruke12/payment-service/internal/domains/product"
	"github.com/mrruke12/payment-service/internal/domains/reservation"
)

type RepositoryProvider interface {
	ProductRepository() product.Repository
	OrderRepository() order.Repository
	ReservationRepository() reservation.Repository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, r RepositoryProvider) error) error
}
