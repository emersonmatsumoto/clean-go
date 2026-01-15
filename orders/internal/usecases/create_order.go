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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

type OrderItemInput struct {
	ProductID string
	Quantity  int
}

type PlaceOrderInput struct {
	UserID    string
	Items     []OrderItemInput
	CardToken string
}

type PlaceOrderOutput struct {
	OrderID string
	Total   float64
	Status  string
}

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

func (uc *CreateOrderUseCase) Execute(ctx context.Context, in PlaceOrderInput) (PlaceOrderOutput, error) {
	ctx, span := tracer.Start(ctx, "Orders.CreateOrderUseCase.Execute")
	defer span.End()

	if len(in.Items) == 0 {
		err := fmt.Errorf("o pedido deve ter pelo menos um item")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return PlaceOrderOutput{}, err
	}

	if in.CardToken == "" {
		err := fmt.Errorf("token do cartão é obrigatório")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return PlaceOrderOutput{}, err
	}

	var domainItems []entities.OrderItem

	for _, item := range in.Items {
		p, err := uc.prodComp.GetProduct(ctx, products.GetProductInput{ID: item.ProductID})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(
				codes.Error,
				fmt.Sprintf("produto %s não encontrado", item.ProductID),
			)
			return PlaceOrderOutput{}, fmt.Errorf("produto %s não encontrado", item.ProductID)
		}

		domainItems = append(domainItems, entities.OrderItem{
			ProductID: p.ID,
			Price:     p.Price,
			Quantity:  item.Quantity,
		})
	}

	userData, err := uc.userComp.GetUser(ctx, users.GetUserInput{ID: in.UserID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "usuário não encontrado")
		return PlaceOrderOutput{}, errors.New("usuário não encontrado")
	}

	addressStr := fmt.Sprintf("%s, %s - %s", userData.Address.Street, userData.Address.City, userData.Address.ZipCode)
	order := entities.NewOrder(in.UserID, domainItems, addressStr)

	span.SetAttributes(
		attribute.String("user.id", order.UserID),
		attribute.Float64("order.total", order.Total),
	)

	payRes, err := uc.payComp.ProcessPayment(ctx, payments.ProcessPaymentInput{
		OrderID:  order.ID,
		Amount:   order.Total,
		TokenID:  in.CardToken,
		Currency: "BRL",
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "erro ao processar pagamento")
		return PlaceOrderOutput{}, errors.New("falha no pagamento")
	}

	if payRes.Status != "SUCCESS" {
		err := fmt.Errorf("pagamento recusado: %s", payRes.Status)
		span.RecordError(err)
		span.SetStatus(codes.Error, "pagamento recusado")
		return PlaceOrderOutput{}, errors.New("falha no pagamento")
	}

	order.MarkAsPaid(payRes.TransactionID)

	orderID, err := uc.repo.Save(ctx, order)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "erro ao salvar pedido")
		return PlaceOrderOutput{}, err
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

	span.SetStatus(codes.Ok, "ordem criada com sucesso")

	return PlaceOrderOutput{
		OrderID: order.ID,
		Total:   order.Total,
		Status:  order.Status,
	}, nil
}
