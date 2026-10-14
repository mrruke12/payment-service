package order

import (
	"context"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *Order, items []OrderItem) error
}
