package users

import (
	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"github.com/emersonmatsumoto/clean-go/users/internal/db"
	"github.com/emersonmatsumoto/clean-go/users/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type component struct {
	getUC *usecases.GetUserUseCase
}

func NewComponent(mongoClient *mongo.Client) users.Component {
	repo := db.NewMongoRepo(mongoClient)

	return &component{
		getUC: usecases.NewGetUserUseCase(repo),
	}
}
