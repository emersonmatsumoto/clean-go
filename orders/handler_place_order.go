package orders

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/orders"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/orders")

func (c *component) PlaceOrder(ctx context.Context, in orders.PlaceOrderInput) (orders.PlaceOrderOutput, error) {
	ctx, span := tracer.Start(ctx, "Orders.Component.PlaceOrder")
	defer span.End()

	return c.createUC.Execute(ctx, in)
}
