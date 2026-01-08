package usecases

import (
	"errors"

	"github.com/emersonmatsumoto/clean-go/users/internal/entities"
	"github.com/emersonmatsumoto/clean-go/users/internal/ports"
)

type GetUserUseCase struct {
	repo ports.UserRepository
}

func NewGetUserUseCase(repo ports.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) Execute(id string) (*entities.User, error) {
	if id == "" {
		return nil, errors.New("id do utilizador é obrigatório")
	}

	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("utilizador não encontrado")
	}

	return user, nil
}
