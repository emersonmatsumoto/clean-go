package ports

import "github.com/emersonmatsumoto/clean-go/orders/internal/entities"

type OrderRepository interface {
	Save(order *entities.Order) error
}
