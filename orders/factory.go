package orders

import (
	"github.com/emersonmatsumoto/clean-go/orders/internal/db"
	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
	"github.com/emersonmatsumoto/clean-go/users"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewComponent(
	mongoClient *mongo.Client,
	prodComp products.Component,
	payComp payments.Component,
	userComp users.Component,
) Component {
	repo := db.NewMongoRepo(mongoClient)

	uc := usecases.NewCreateOrderUseCase(repo, prodComp, payComp, userComp)

	return &component{
		createUC: uc,
	}
}
