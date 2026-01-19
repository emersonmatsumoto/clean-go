package products

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"github.com/emersonmatsumoto/clean-go/products/internal/db"
	"github.com/emersonmatsumoto/clean-go/products/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
)

type component struct {
	getUC *usecases.GetProductUseCase
}

func NewComponent(mongoClient *mongo.Client) products.Component {
	repo := db.NewMongoRepo(mongoClient)

	return &component{
		getUC: usecases.NewGetProductUseCase(repo),
	}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/products")

func (c *component) GetProduct(ctx context.Context, in products.GetProductInput) (products.GetProductOutput, error) {
	ctx, span := tracer.Start(ctx, "Products.Component.GetProduct")
	defer span.End()

	return c.getUC.Execute(ctx, in)
}
