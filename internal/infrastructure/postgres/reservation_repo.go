package postgres

import (
	"context"

	"github.com/mrruke12/payment-service/internal/domains/reservation"
)

type ReservationRepository struct {
	db Querier
}

func NewReservationRepository(db Querier) reservation.Repository {
	return &ReservationRepository{
		db: db,
	}
}

func (r *ReservationRepository) CreateReservation(ctx context.Context, reservation *reservation.Reservation) error {
	sql, args, err := builder.Insert("reservations").
		Columns("id", "order_id", "status").
		Values(reservation.ID, reservation.OrderID, reservation.Status).
		ToSql()

	_, err = r.db.Exec(ctx, sql, args...)

	return err
}
