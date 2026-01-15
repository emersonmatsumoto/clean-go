package payments

type ProcessPaymentInput struct {
	OrderID  string
	Amount   float64
	TokenID  string
	Currency string
}

type ProcessPaymentOutput struct {
	TransactionID string
	Status        string
}
