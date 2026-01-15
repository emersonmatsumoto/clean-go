package payments

import (
	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"github.com/emersonmatsumoto/clean-go/payments/internal/external"
	"github.com/emersonmatsumoto/clean-go/payments/internal/usecases"
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
