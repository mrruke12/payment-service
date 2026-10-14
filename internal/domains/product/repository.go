package product

import (
	"context"
	"uuid"
)

type Repository interface {
	SubtractQuantities(ctx context.Context, ids []uuid.UUID, deltas []int) error
}
