package reservation

import (
	"time"
	"uuid"
)

type Reservation struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	Status    Status
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewReservation(orderID uuid.UUID) *Reservation {
	return &Reservation{
		ID:      uuid.NewV7(),
		OrderID: orderID,
		Status:  StatusActive,
	}
}
