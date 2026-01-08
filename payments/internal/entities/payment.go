package entities

import "time"

type Payment struct {
	ID            string
	TransactionID string
	Amount        float64
	Currency      string
	Status        string
	CreatedAt     time.Time
}
