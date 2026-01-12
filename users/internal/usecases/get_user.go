package usecases

import (
	"context"
	"errors"

	"github.com/emersonmatsumoto/clean-go/users/internal/entities"
	"github.com/emersonmatsumoto/clean-go/users/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type GetUserUseCase struct {
	repo ports.UserRepository
}

func NewGetUserUseCase(repo ports.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/users/internal/usecases")

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*entities.User, error) {
	ctx, span := tracer.Start(ctx, "Users.GetUserUseCase.Execute")
	defer span.End()

	if id == "" {
		return nil, errors.New("id do utilizador é obrigatório")
	}

	span.SetAttributes(
		attribute.String("user.id", id),
	)

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("utilizador não encontrado")
	}

	return user, nil
}
