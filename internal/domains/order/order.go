package order

import (
	"time"
	"uuid"
)

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(userID uuid.UUID) *Order {
	return &Order{
		ID:     uuid.NewV7(),
		UserID: userID,
		Status: StatusPending,
	}
}
