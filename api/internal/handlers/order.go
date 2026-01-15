package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/emersonmatsumoto/clean-go/orders"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type OrderController struct {
	orderComp orders.Component
}

func NewOrderController(oc orders.Component) *OrderController {
	return &OrderController{
		orderComp: oc,
	}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/api/internal/handlers")

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type PlaceOrderRequest struct {
	UserID    string             `json:"user_id"`
	Items     []OrderItemRequest `json:"items"`
	CardToken string             `json:"card_token"`
}

type PlaceOrderResponse struct {
	OrderID string  `json:"order_id"`
	Total   float64 `json:"total"`
	Status  string  `json:"status"`
}

func (ctrl *OrderController) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "OrderController.PlaceOrder")
	defer span.End()

	var input PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid_payload", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", input.UserID),
	)

	output, err := ctrl.orderComp.PlaceOrder(ctx, orders.PlaceOrderInput{
		UserID:    input.UserID,
		CardToken: input.CardToken,
		Items: func() []orders.OrderItemInput {
			var items []orders.OrderItemInput
			for _, item := range input.Items {
				items = append(items, orders.OrderItemInput{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
				})
			}
			return items
		}(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
