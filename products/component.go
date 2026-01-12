package products

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/products/internal/usecases"
)

type Component interface {
	GetProduct(ctx context.Context, input GetProductInput) (GetProductOutput, error)
}

type component struct {
	getUC *usecases.GetProductUseCase
}
