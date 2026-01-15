package orders

import "context"

type Component interface {
	PlaceOrder(ctx context.Context, in PlaceOrderInput) (PlaceOrderOutput, error)
}
