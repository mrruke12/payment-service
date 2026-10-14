package postgres

import (
	"context"

	"github.com/mrruke12/payment-service/internal/domains/order"
)

type OrderRepository struct {
	db Querier
}

func NewOrderRepository(db Querier) order.Repository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *order.Order, items []order.OrderItem) error {
	sql, args, err := builder.Insert("orders").
		Columns("id", "user_id", "status").
		Values(order.ID, order.UserID, order.Status).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	itemsQuery := builder.Insert("order_items").
		Columns("id", "order_id", "product_id", "quantity")

	for i := range items {
		item := &items[i]
		itemsQuery = itemsQuery.Values(item.ID, item.OrderID, item.ProductID, item.Quantity)
	}

	sql, args, err = itemsQuery.ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)

	return err
}
