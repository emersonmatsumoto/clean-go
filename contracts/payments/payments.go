package payments

import "context"

type Component interface {
	ProcessPayment(ctx context.Context, in ProcessPaymentInput) (ProcessPaymentOutput, error)
}
