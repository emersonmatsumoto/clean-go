package ports

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/users"
)

type GetUserInput = users.GetUserInput
type GetUserOutput = users.GetUserOutput

type UserGateway interface {
	GetUser(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}
