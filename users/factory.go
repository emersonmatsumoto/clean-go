package users

import (
	"github.com/emersonmatsumoto/clean-go/users/internal/db"
	"github.com/emersonmatsumoto/clean-go/users/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewComponent(mongoClient *mongo.Client) Component {
	repo := db.NewMongoRepo(mongoClient)

	return &component{
		getUC: usecases.NewGetUserUseCase(repo),
	}
}
