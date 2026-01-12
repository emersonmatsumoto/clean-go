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

func (ctrl *OrderController) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "OrderController.PlaceOrder")
	defer span.End()

	var input orders.PlaceOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid_payload", http.StatusBadRequest)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", input.UserID),
	)

	output, err := ctrl.orderComp.PlaceOrder(ctx, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
