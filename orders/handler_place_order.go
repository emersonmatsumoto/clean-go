package orders

import (
	"context"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"go.opentelemetry.io/otel"
)

type OrderItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type PlaceOrderInput struct {
	UserID    string           `json:"user_id"`
	Items     []OrderItemInput `json:"items"`
	CardToken string           `json:"card_token"`
}

type PlaceOrderOutput struct {
	OrderID string  `json:"order_id"`
	Total   float64 `json:"total"`
	Status  string  `json:"status"`
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/orders")

func (c *component) PlaceOrder(ctx context.Context, in PlaceOrderInput) (PlaceOrderOutput, error) {
	ctx, span := tracer.Start(ctx, "Orders.Component.PlaceOrder")
	defer span.End()

	if len(in.Items) == 0 {
		return PlaceOrderOutput{}, fmt.Errorf("o pedido deve ter pelo menos um item")
	}
	if in.CardToken == "" {
		return PlaceOrderOutput{}, fmt.Errorf("token do cartão é obrigatório")
	}

	var itemsToProcess []entities.OrderItem
	for _, item := range in.Items {
		itemsToProcess = append(itemsToProcess, entities.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := c.createUC.Execute(ctx, in.UserID, itemsToProcess, in.CardToken)
	if err != nil {
		return PlaceOrderOutput{}, err
	}

	return PlaceOrderOutput{
		OrderID: order.ID,
		Total:   order.Total,
		Status:  order.Status,
	}, nil
}
