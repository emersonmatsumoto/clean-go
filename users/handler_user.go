package users

type UserAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zip_code"`
}

type GetUserInput struct {
	ID string `json:"id"`
}

type GetUserOutput struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Address UserAddress `json:"address"`
}

func (c *component) GetUser(in GetUserInput) (GetUserOutput, error) {
	user, err := c.getUC.Execute(in.ID)
	if err != nil {
		return GetUserOutput{}, err
	}

	return GetUserOutput{
		ID:   user.ID,
		Name: user.Name,
		Address: UserAddress{
			Street:  user.Address.Street,
			City:    user.Address.City,
			ZipCode: user.Address.ZipCode,
		},
	}, nil
}
