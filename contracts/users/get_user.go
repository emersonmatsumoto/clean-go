package users

type GetUserInput struct {
	ID string `json:"id"`
}

type GetUserOutput struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Address UserAddress `json:"address"`
}
