package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"github.com/emersonmatsumoto/clean-go/users/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type GetUserUseCase struct {
	repo ports.UserRepository
}

func NewGetUserUseCase(repo ports.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/users/internal/usecases")

func (uc *GetUserUseCase) Execute(ctx context.Context, in users.GetUserInput) (users.GetUserOutput, error) {
	ctx, span := tracer.Start(ctx, "Users.GetUserUseCase.Execute")
	defer span.End()

	if in.ID == "" {
		err := fmt.Errorf("id do usuário é obrigatório")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return users.GetUserOutput{}, errors.New("id do usuário é obrigatório")
	}

	span.SetAttributes(
		attribute.String("user.id", in.ID),
	)

	user, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "erro ao buscar usuário")
		return users.GetUserOutput{}, err
	}

	if user == nil {
		err := fmt.Errorf("usuário não encontrado")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return users.GetUserOutput{}, errors.New("usuário não encontrado")
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
