package ports

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/products"
)

type GetProductInput = products.GetProductInput
type GetProductOutput = products.GetProductOutput

type ProductGateway interface {
	GetProduct(ctx context.Context, input GetProductInput) (GetProductOutput, error)
}
