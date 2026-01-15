package products

import (
	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"github.com/emersonmatsumoto/clean-go/products/internal/db"
	"github.com/emersonmatsumoto/clean-go/products/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
