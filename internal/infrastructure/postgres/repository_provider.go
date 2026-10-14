package postgres

import (
	"github.com/mrruke12/payment-service/internal/domains/order"
	"github.com/mrruke12/payment-service/internal/domains/product"
	"github.com/mrruke12/payment-service/internal/domains/reservation"
)

type RepositoryProvider struct {
	tx Querier
}

func (r *RepositoryProvider) ProductRepository() product.Repository {
	return NewProductRepository(r.tx)
}

func (r *RepositoryProvider) OrderRepository() order.Repository {
	return NewOrderRepository(r.tx)
}

func (r *RepositoryProvider) ReservationRepository() reservation.Repository {
	return NewReservationRepository(r.tx)
}
