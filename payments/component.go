package payments

import "github.com/emersonmatsumoto/clean-go/payments/internal/usecases"

type Component interface {
	ProcessPayment(in ProcessPaymentInput) (ProcessPaymentOutput, error)
}

type component struct {
	payUC *usecases.ProcessPaymentUseCase
}
