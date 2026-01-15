package usecases

import (
	"context"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"github.com/emersonmatsumoto/clean-go/payments/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type ProcessPaymentUseCase struct {
	gateway ports.PaymentGateway
}

func NewProcessPaymentUseCase(g ports.PaymentGateway) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{gateway: g}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/payments")

func (uc *ProcessPaymentUseCase) Execute(ctx context.Context, in payments.ProcessPaymentInput) (payments.ProcessPaymentOutput, error) {
	ctx, span := tracer.Start(ctx, "Payments.ProcessPaymentUseCase.Execute")
	defer span.End()

	if in.Amount <= 0 {
		err := fmt.Errorf("Amount deve ser maior que zero")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return payments.ProcessPaymentOutput{Status: "FAILED"}, err
	}

	txID, err := uc.gateway.Charge(in.Amount, in.TokenID, in.Currency)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "erro ao processar pagamento")
		return payments.ProcessPaymentOutput{Status: "FAILED"}, err
	}

	span.SetAttributes(
		attribute.String("payment.status", "PAID"),
		attribute.String("payment.transaction_id", txID),
	)

	return payments.ProcessPaymentOutput{
		TransactionID: txID,
		Status:        "PAID",
	}, nil
}
