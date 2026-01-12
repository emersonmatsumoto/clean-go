package usecases

import (
	"context"
	"errors"

	"github.com/emersonmatsumoto/clean-go/products/internal/entities"
	"github.com/emersonmatsumoto/clean-go/products/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type GetProductUseCase struct {
	repo ports.ProductRepository
}

func NewGetProductUseCase(repo ports.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{
		repo: repo,
	}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/products/internal/usecases")

func (uc *GetProductUseCase) Execute(ctx context.Context, id string) (*entities.Product, error) {
	ctx, span := tracer.Start(ctx, "Products.GetProductUseCase.Execute")
	defer span.End()

	if id == "" {
		return nil, errors.New("id do produto é obrigatório")
	}

	span.SetAttributes(
		attribute.String("product.id", id),
	)

	product, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("produto não encontrado")
	}

	return product, nil
}
