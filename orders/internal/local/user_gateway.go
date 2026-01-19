package local

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
)

type userGateway struct {
	comp users.Component
}

func NewUserGateway(comp users.Component) ports.UserGateway {
	return &userGateway{comp: comp}
}

func (a *userGateway) GetUser(ctx context.Context, in ports.GetUserInput) (ports.GetUserOutput, error) {
	return a.comp.GetUser(ctx, in)
}
