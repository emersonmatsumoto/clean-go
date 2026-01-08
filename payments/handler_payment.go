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

func (c *component) ProcessPayment(in ProcessPaymentInput) (ProcessPaymentOutput, error) {
	res, err := c.payUC.Execute(in.Amount, in.TokenID, in.Currency)

	status := "SUCCESS"
	if err != nil {
		status = "FAILED"
	}

	return ProcessPaymentOutput{
		TransactionID: res.TransactionID,
		Status:        status,
	}, err
}
