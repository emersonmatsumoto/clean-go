package products

import (
	"context"

	"go.opentelemetry.io/otel"
)

type GetProductInput struct {
	ID string
}

type GetProductOutput struct {
	ID    string
	Name  string
	Price float64
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/products")

func (c *component) GetProduct(ctx context.Context, in GetProductInput) (GetProductOutput, error) {
	ctx, span := tracer.Start(ctx, "Products.Component.GetProduct")
	defer span.End()

	p, err := c.getUC.Execute(ctx, in.ID)
	if err != nil {
		return GetProductOutput{}, err
	}

	return GetProductOutput{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}, nil
}
