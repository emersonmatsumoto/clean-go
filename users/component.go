package users

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/users/internal/usecases"
)

type Component interface {
	GetUser(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}

type component struct {
	getUC *usecases.GetUserUseCase
}
