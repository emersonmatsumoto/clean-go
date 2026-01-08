package usecases

import (
	"errors"

	"github.com/emersonmatsumoto/clean-go/payments/internal/entities"
	"github.com/emersonmatsumoto/clean-go/payments/internal/ports"
)

type ProcessPaymentUseCase struct {
	gateway ports.PaymentGateway
}

func NewProcessPaymentUseCase(g ports.PaymentGateway) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{gateway: g}
}

func (uc *ProcessPaymentUseCase) Execute(amount float64, token string, currency string) (*entities.Payment, error) {
	if amount <= 0 {
		return nil, errors.New("valor inválido")
	}

	txID, err := uc.gateway.Charge(amount, token, currency)
	if err != nil {
		return &entities.Payment{Status: "FAILED"}, err
	}

	return &entities.Payment{
		TransactionID: txID,
		Status:        "PAID",
	}, nil
}
