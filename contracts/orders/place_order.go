package orders

type OrderItemInput struct {
	ProductID string
	Quantity  int
}

type PlaceOrderInput struct {
	UserID    string
	Items     []OrderItemInput
	CardToken string
}

type PlaceOrderOutput struct {
	OrderID string
	Total   float64
	Status  string
}
