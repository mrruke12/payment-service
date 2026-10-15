package orderusecase

import (
	"context"
	"uuid"

	"github.com/mrruke12/payment-service/internal/domains/order"
	"github.com/mrruke12/payment-service/internal/domains/reservation"
	"github.com/mrruke12/payment-service/pkg/aggregation"
)

func (s *Service) CreateOrder(ctx context.Context, userID uuid.UUID, items []order.Item) (*order.Order, error) {
	o := order.NewOrder(userID) // create an order

	oItems := order.ItemsToOrderItems(o.ID, items)                                                 // convert object-ish items to persistence order items
	products := aggregation.Map(items, func(v order.Item, _ int) uuid.UUID { return v.ProductID }) // get product ids
	deltas := aggregation.Map(items, func(v order.Item, _ int) int { return v.Quantity })          // get deltas

	reserv := reservation.NewReservation(o.ID) // create reservation

	err := s.uow.Do(ctx, func(ctx context.Context, r RepositoryProvider) error {
		productRepo := r.ProductRepository()
		orderRepo := r.OrderRepository()
		reservationRepo := r.ReservationRepository()

		err := orderRepo.CreateOrder(ctx, o, oItems) // create order and items in database

		if err != nil {
			return err
		}

		err = productRepo.SubtractQuantities(ctx, products, deltas) // reserve order item products

		if err != nil {
			return err
		}

		return reservationRepo.CreateReservation(ctx, reserv) // create reservation
	})

	if err != nil {
		return nil, err
	}

	return o, nil
}
