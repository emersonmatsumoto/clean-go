package orders

import (
	"github.com/emersonmatsumoto/clean-go/contracts/orders"
	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"github.com/emersonmatsumoto/clean-go/orders/internal/db"
	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type component struct {
	createUC *usecases.CreateOrderUseCase
}

func NewComponent(
	mongoClient *mongo.Client,
	prodComp products.Component,
	payComp payments.Component,
	userComp users.Component,
) orders.Component {
	repo := db.NewMongoRepo(mongoClient)

	uc := usecases.NewCreateOrderUseCase(repo, prodComp, payComp, userComp)

	return &component{
		createUC: uc,
	}
}
