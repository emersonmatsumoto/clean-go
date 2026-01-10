package handlers

import "net/http"

func RegisterRoutes(mux *http.ServeMux, orderCtrl *OrderController) {
	mux.HandleFunc("POST /orders", orderCtrl.PlaceOrder)
}
