package ports

import "github.com/emersonmatsumoto/clean-go/products/internal/entities"

type ProductRepository interface {
	FindByID(id string) (*entities.Product, error)
}
