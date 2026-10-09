package product

import "uuid"

type Product struct {
	ID       uuid.UUID
	name     string
	Quantity int
	Price    int
}
