package postgres

import (
	"context"
	"errors"
	"uuid"

	"github.com/mrruke12/payment-service/internal/domains/product"
)

type ProductRepository struct {
	db Querier
}

func NewProductRepository(db Querier) product.Repository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) SubtractQuantities(ctx context.Context, ids []uuid.UUID, deltas []int) error {
	sql := `
		update products p
		set quantity = p.quantity - s.delta
		from unnest(
			$1::uuid[],
			$2::int[]
		) as s(id, delta)
		where p.id = s.id
		and p.quantity >= s.delta
	`
	tag, err := r.db.Exec(ctx, sql, ids, deltas)

	if err != nil {
		return err
	}

	if tag.RowsAffected() != int64(len(ids)) {
		return errors.New("affected rows mismatch")
	}

	return nil
}
