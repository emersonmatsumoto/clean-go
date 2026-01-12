package payments

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/payments/internal/usecases"
)

type Component interface {
	ProcessPayment(ctx context.Context, in ProcessPaymentInput) (ProcessPaymentOutput, error)
}

type component struct {
	payUC *usecases.ProcessPaymentUseCase
}
