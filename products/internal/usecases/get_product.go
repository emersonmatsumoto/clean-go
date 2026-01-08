package usecases

import (
	"errors"
	"github.com/emersonmatsumoto/clean-go/products/internal/entities"
	"github.com/emersonmatsumoto/clean-go/products/internal/ports"
)

type GetProductUseCase struct {
	repo ports.ProductRepository
}

func NewGetProductUseCase(repo ports.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{
		repo: repo,
	}
}

func (uc *GetProductUseCase) Execute(id string) (*entities.Product, error) {
	if id == "" {
		return nil, errors.New("id do produto é obrigatório")
	}

	product, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("produto não encontrado")
	}

	return product, nil
}
