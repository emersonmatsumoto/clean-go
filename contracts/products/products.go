package products

import "context"

type Component interface {
	GetProduct(ctx context.Context, input GetProductInput) (GetProductOutput, error)
}
