package products

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/products")

func (c *component) GetProduct(ctx context.Context, in products.GetProductInput) (products.GetProductOutput, error) {
	ctx, span := tracer.Start(ctx, "Products.Component.GetProduct")
	defer span.End()

	p, err := c.getUC.Execute(ctx, in.ID)
	if err != nil {
		return products.GetProductOutput{}, err
	}

	return products.GetProductOutput{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}, nil
}
