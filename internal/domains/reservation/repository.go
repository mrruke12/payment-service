package reservation

import (
	"context"
)

type Repository interface {
	CreateReservation(ctx context.Context, reservation *Reservation) error
	TidyReservations(ctx context.Context, batchSize int) error
}
