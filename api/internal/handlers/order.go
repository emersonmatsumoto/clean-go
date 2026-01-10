package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/emersonmatsumoto/clean-go/orders"
)

type OrderController struct {
	orderComp orders.Component
}

func NewOrderController(oc orders.Component) *OrderController {
	return &OrderController{
		orderComp: oc,
	}
}

func (ctrl *OrderController) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var input orders.PlaceOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid_payload", http.StatusBadRequest)
		return
	}

	output, err := ctrl.orderComp.PlaceOrder(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
