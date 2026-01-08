package products

import "github.com/emersonmatsumoto/clean-go/products/internal/usecases"

type Component interface {
	GetProduct(input GetProductInput) (GetProductOutput, error)
}

type component struct {
	getUC *usecases.GetProductUseCase
}
