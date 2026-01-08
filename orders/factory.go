package orders

import (
	"github.com/emersonmatsumoto/clean-go/orders/internal/db"
	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewComponent(
	mongoClient *mongo.Client,
	prodComp products.Component,
	payComp payments.Component,
) Component {
	repo := db.NewMongoRepo(mongoClient)

	uc := usecases.NewCreateOrderUseCase(repo, prodComp, payComp)

	return &component{
		createUC: uc,
	}
}
