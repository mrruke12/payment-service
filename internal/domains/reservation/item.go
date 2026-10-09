package reservation

import "uuid"

type Item struct {
	ID                 uuid.UUID
	ItemID             uuid.UUID
	OrderReservationID uuid.UUID
}
