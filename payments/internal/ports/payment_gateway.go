package ports

type PaymentGateway interface {
	Charge(amount float64, token string, currency string) (string, error)
}
