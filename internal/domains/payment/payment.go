package payment

import (
	"time"
	"uuid"
)

type Payment struct {
	PaymentID uuid.UUID
	UserID    uuid.UUID
	OrderID   uuid.UUID
	Amount    int
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
