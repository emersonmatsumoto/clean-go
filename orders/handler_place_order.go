package orders

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/orders")

type PlaceOrderInput = usecases.PlaceOrderInput
type PlaceOrderOutput = usecases.PlaceOrderOutput
type OrderItemInput = usecases.OrderItemInput

func (c *component) PlaceOrder(ctx context.Context, in usecases.PlaceOrderInput) (usecases.PlaceOrderOutput, error) {
	ctx, span := tracer.Start(ctx, "Orders.Component.PlaceOrder")
	defer span.End()

	return c.createUC.Execute(ctx, in)
}
