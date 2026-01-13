package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
	"github.com/emersonmatsumoto/clean-go/users"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type CreateOrderUseCase struct {
	repo     ports.OrderRepository
	prodComp products.Component
	payComp  payments.Component
	userComp users.Component
}

func NewCreateOrderUseCase(r ports.OrderRepository, p products.Component, pay payments.Component, user users.Component) *CreateOrderUseCase {
	return &CreateOrderUseCase{repo: r, prodComp: p, payComp: pay, userComp: user}
}

var (
	tracer          = otel.Tracer("github.com/emersonmatsumoto/clean-go/orders/internal/usecases")
	meter           = otel.Meter("orders-component")
	orderCounter, _ = meter.Int64Counter(
		"orders.created.total",
		metric.WithDescription("Total de pedidos criados"),
	)
	orderValue, _ = meter.Float64Histogram(
		"orders.value",
		metric.WithUnit("BRL"),
		metric.WithDescription("Distribuição dos valores dos pedidos"),
	)
)

func (uc *CreateOrderUseCase) Execute(ctx context.Context, userID string, itemsInput []entities.OrderItem, cardToken string) (*entities.Order, error) {
	ctx, span := tracer.Start(ctx, "Orders.CreateOrderUseCase.Execute")
	defer span.End()

	var domainItems []entities.OrderItem

	for _, item := range itemsInput {
		p, err := uc.prodComp.GetProduct(ctx, products.GetProductInput{ID: item.ProductID})
		if err != nil {
			return nil, fmt.Errorf("produto %s não encontrado", item.ProductID)
		}

		domainItems = append(domainItems, entities.OrderItem{
			ProductID: p.ID,
			Price:     p.Price,
			Quantity:  item.Quantity,
		})
	}

	userData, err := uc.userComp.GetUser(ctx, users.GetUserInput{ID: userID})
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	addressStr := fmt.Sprintf("%s, %s - %s", userData.Address.Street, userData.Address.City, userData.Address.ZipCode)
	order := entities.NewOrder(userID, domainItems, addressStr)

	span.SetAttributes(
		attribute.String("user.id", order.UserID),
		attribute.Float64("order.total", order.Total),
	)

	payRes, err := uc.payComp.ProcessPayment(ctx, payments.ProcessPaymentInput{
		OrderID:  order.ID,
		Amount:   order.Total,
		TokenID:  cardToken,
		Currency: "BRL",
	})

	if err != nil || payRes.Status != "SUCCESS" {
		return nil, errors.New("falha no pagamento")
	}

	order.MarkAsPaid(payRes.TransactionID)

	orderID, err := uc.repo.Save(ctx, order)
	if err != nil {
		return nil, err
	}

	order.SetID(orderID)
	span.SetAttributes(
		attribute.String("order.id", order.ID),
	)

	orderCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("status", "success"),
		attribute.String("city", userData.Address.City),
	))
	orderValue.Record(ctx, order.Total)

	return order, err
}
