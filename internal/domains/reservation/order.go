package reservation

import (
	"time"
	"uuid"
)

type Order struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
