package ports

import "github.com/emersonmatsumoto/clean-go/users/internal/entities"

type UserRepository interface {
	FindByID(id string) (*entities.User, error)
}
