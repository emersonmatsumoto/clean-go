package local

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
)

type paymentGateway struct {
	comp payments.Component
}

func NewPaymentGateway(comp payments.Component) ports.PaymentGateway {
	return &paymentGateway{comp: comp}
}

func (a *paymentGateway) ProcessPayment(ctx context.Context, in ports.ProcessPaymentInput) (ports.ProcessPaymentOutput, error) {
	return a.comp.ProcessPayment(ctx, in)
}
