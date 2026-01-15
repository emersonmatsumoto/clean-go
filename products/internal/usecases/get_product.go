package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"github.com/emersonmatsumoto/clean-go/products/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (uc *GetProductUseCase) Execute(ctx context.Context, in products.GetProductInput) (products.GetProductOutput, error) {
	ctx, span := tracer.Start(ctx, "Products.GetProductUseCase.Execute")
	defer span.End()

	if in.ID == "" {
		err := fmt.Errorf("id do produto é obrigatório")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return products.GetProductOutput{}, errors.New("id do produto é obrigatório")
	}

	span.SetAttributes(
		attribute.String("product.id", in.ID),
	)

	product, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "erro ao buscar produto")
		return products.GetProductOutput{}, err
	}

	if product == nil {
		err := fmt.Errorf("produto não encontrado")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return products.GetProductOutput{}, errors.New("produto não encontrado")
	}

	return products.GetProductOutput{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}
