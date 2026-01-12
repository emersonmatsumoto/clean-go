package ports

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/products/internal/entities"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*entities.Product, error)
}
