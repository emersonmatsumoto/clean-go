package local

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
)

type productGateway struct {
	comp products.Component
}

func NewProductGateway(comp products.Component) ports.ProductGateway {
	return &productGateway{comp: comp}
}

func (a *productGateway) GetProduct(ctx context.Context, in ports.GetProductInput) (ports.GetProductOutput, error) {
	return a.comp.GetProduct(ctx, in)
}
