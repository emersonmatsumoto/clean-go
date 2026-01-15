package orders

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
)

type Component interface {
	PlaceOrder(ctx context.Context, in usecases.PlaceOrderInput) (usecases.PlaceOrderOutput, error)
}

type component struct {
	createUC *usecases.CreateOrderUseCase
}
