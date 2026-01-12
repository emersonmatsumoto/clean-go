package payments

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ProcessPaymentInput struct {
	OrderID  string
	Amount   float64
	TokenID  string
	Currency string
}

type ProcessPaymentOutput struct {
	TransactionID string
	Status        string
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/payments")

func (c *component) ProcessPayment(ctx context.Context, in ProcessPaymentInput) (ProcessPaymentOutput, error) {
	ctx, span := tracer.Start(ctx, "Payments.Component.ProcessPayment")
	defer span.End()

	res, err := c.payUC.Execute(in.Amount, in.TokenID, in.Currency)

	status := "SUCCESS"
	if err != nil {
		status = "FAILED"
	}

	span.SetAttributes(
		attribute.String("payment.status", status),
		attribute.String("payment.transaction_id", res.TransactionID),
	)

	return ProcessPaymentOutput{
		TransactionID: res.TransactionID,
		Status:        status,
	}, err
}
