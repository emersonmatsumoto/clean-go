package users

import (
	"context"

	"go.opentelemetry.io/otel"
)

type UserAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
}

type GetUserInput struct {
	ID string `json:"id"`
}

type GetUserOutput struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Address UserAddress `json:"address"`
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/users")

func (c *component) GetUser(ctx context.Context, in GetUserInput) (GetUserOutput, error) {
	ctx, span := tracer.Start(ctx, "Users.Component.GetUser")
	defer span.End()

	user, err := c.getUC.Execute(ctx, in.ID)
	if err != nil {
		return GetUserOutput{}, err
	}

	return GetUserOutput{
		ID:   user.ID,
		Name: user.Name,
		Address: UserAddress{
			Street:  user.Address.Street,
			City:    user.Address.City,
			ZipCode: user.Address.ZipCode,
		},
	}, nil
}
