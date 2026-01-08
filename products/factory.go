package products

import (
	"github.com/emersonmatsumoto/clean-go/products/internal/db"
	"github.com/emersonmatsumoto/clean-go/products/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewComponent(mongoClient *mongo.Client) Component {
	repo := db.NewMongoRepo(mongoClient)

	return &component{
		getUC: usecases.NewGetProductUseCase(repo),
	}
}
