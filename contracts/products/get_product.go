package products

type GetProductInput struct {
	ID string
}

type GetProductOutput struct {
	ID    string
	Name  string
	Price float64
}
