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

func (r *ReservationRepository) TidyReservations(ctx context.Context, batchSize int) error {
	sql := `
		with to_update as (
			select r.id as reservation_id, r.order_id as order_id 
			from reservations as r
			where r.expires_at < now()
			and status = 'active'
			limit $1
			for update
		),
		upd_res as (
			update reservations r
			set status = 'expired'
			from to_update t
			where r.id = t.reservation_id
		)
		update orders o
		set status = 'cancelled'
		from to_update t
		where o.id = t.order_id
	`

	_, err := r.db.Exec(ctx, sql, batchSize)

	return err
}
