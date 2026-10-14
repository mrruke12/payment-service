package order

import "uuid"

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}

func NewOrderItem(orderID, productID uuid.UUID, quantity int) *OrderItem {
	return &OrderItem{
		ID:        uuid.NewV7(),
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
	}
}

func ItemsToOrderItems(orderID uuid.UUID, items []Item) []OrderItem {
	res := make([]OrderItem, len(items))

	for i := range items {
		res[i] = *NewOrderItem(
			orderID,
			items[i].ProductID,
			items[i].Quantity,
		)
	}

	return res
}

type Item struct {
	ProductID uuid.UUID
	Quantity  int
}
