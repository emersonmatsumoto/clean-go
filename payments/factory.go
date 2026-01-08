package payments

import (
	"github.com/emersonmatsumoto/clean-go/payments/internal/external"
	"github.com/emersonmatsumoto/clean-go/payments/internal/usecases"
)

func NewComponent(stripeKey string) Component {
	gateway := external.NewStripeAdapter(stripeKey)

	return &component{
		payUC: usecases.NewProcessPaymentUseCase(gateway),
	}
}
