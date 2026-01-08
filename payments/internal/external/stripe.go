package external

import (
	"github.com/emersonmatsumoto/clean-go/payments/internal/ports"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/charge"
)

type stripeAdapter struct {
	key string
}

func NewStripeAdapter(key string) ports.PaymentGateway {
	return &stripeAdapter{key}
}

func (a *stripeAdapter) Charge(amount float64, token string, currency string) (string, error) {
	stripe.Key = a.key
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(int64(amount * 100)),
		Currency: stripe.String(currency),
	}
	params.SetSource(token)

	ch, err := charge.New(params)
	if err != nil {
		return "", err
	}

	return ch.ID, nil
}
