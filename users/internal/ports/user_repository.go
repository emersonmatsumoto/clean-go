package ports

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/users/internal/entities"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entities.User, error)
}
