package users

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/users")

func (c *component) GetUser(ctx context.Context, in users.GetUserInput) (users.GetUserOutput, error) {
	ctx, span := tracer.Start(ctx, "Users.Component.GetUser")
	defer span.End()

	return c.getUC.Execute(ctx, in)
}
