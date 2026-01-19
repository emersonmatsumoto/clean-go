package users

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"github.com/emersonmatsumoto/clean-go/users/internal/db"
	"github.com/emersonmatsumoto/clean-go/users/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
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

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/users")

func (c *component) GetUser(ctx context.Context, in users.GetUserInput) (users.GetUserOutput, error) {
	ctx, span := tracer.Start(ctx, "Users.Component.GetUser")
	defer span.End()

	return c.getUC.Execute(ctx, in)
}
