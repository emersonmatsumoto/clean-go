package payments

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"github.com/emersonmatsumoto/clean-go/payments/internal/external"
	"github.com/emersonmatsumoto/clean-go/payments/internal/usecases"
	"go.opentelemetry.io/otel"
)

type component struct {
	payUC *usecases.ProcessPaymentUseCase
}

func NewComponent(stripeKey string) payments.Component {
	gateway := external.NewStripeAdapter(stripeKey)

	return &component{
		payUC: usecases.NewProcessPaymentUseCase(gateway),
	}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/payments")

func (c *component) ProcessPayment(ctx context.Context, in payments.ProcessPaymentInput) (payments.ProcessPaymentOutput, error) {
	ctx, span := tracer.Start(ctx, "Payments.Component.ProcessPayment")
	defer span.End()

	return c.payUC.Execute(ctx, in)
}
