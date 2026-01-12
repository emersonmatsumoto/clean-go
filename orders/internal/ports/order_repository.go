package ports

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entities.Order) error
}
