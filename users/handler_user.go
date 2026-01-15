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

	user, err := c.getUC.Execute(ctx, in.ID)
	if err != nil {
		return users.GetUserOutput{}, err
	}

	return users.GetUserOutput{
		ID:   user.ID,
		Name: user.Name,
		Address: users.UserAddress{
			Street:  user.Address.Street,
			City:    user.Address.City,
			ZipCode: user.Address.ZipCode,
		},
	}, nil
}
