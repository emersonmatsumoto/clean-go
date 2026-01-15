package payments

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/payments")

func (c *component) ProcessPayment(ctx context.Context, in payments.ProcessPaymentInput) (payments.ProcessPaymentOutput, error) {
	ctx, span := tracer.Start(ctx, "Payments.Component.ProcessPayment")
	defer span.End()

	return c.payUC.Execute(ctx, in)
}
