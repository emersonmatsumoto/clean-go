package users

import "github.com/emersonmatsumoto/clean-go/users/internal/usecases"

type Component interface {
	GetUser(input GetUserInput) (GetUserOutput, error)
}

type component struct {
	getUC *usecases.GetUserUseCase
}
