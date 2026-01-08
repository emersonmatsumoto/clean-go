package orders

import (
	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
)

type Component interface {
	PlaceOrder(in PlaceOrderInput) (PlaceOrderOutput, error)
}

type component struct {
	createUC *usecases.CreateOrderUseCase
}
