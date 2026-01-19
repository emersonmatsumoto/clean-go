package ports

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/payments"
)

type ProcessPaymentInput = payments.ProcessPaymentInput
type ProcessPaymentOutput = payments.ProcessPaymentOutput

type PaymentGateway interface {
	ProcessPayment(ctx context.Context, in ProcessPaymentInput) (ProcessPaymentOutput, error)
}
