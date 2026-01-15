package users

import "context"

type UserAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
}

type Component interface {
	GetUser(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}
