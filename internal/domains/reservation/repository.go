package reservation

import (
	"context"
)

type Repository interface {
	CreateReservation(ctx context.Context, reservation *Reservation) error
}
