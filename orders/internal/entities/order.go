package entities

import "time"

type OrderItem struct {
	ProductID string
	Price     float64
	Quantity  int
}

type Order struct {
	ID              string
	UserID          string
	Items           []OrderItem
	Total           float64
	ShippingAddress string
	Status          string
	TransactionID   string
	CreatedAt       time.Time
}

func NewOrder(userID string, items []OrderItem, shippingAddress string) *Order {
	order := &Order{
		UserID:          userID,
		Items:           items,
		Status:          "PENDING",
		ShippingAddress: shippingAddress,
		CreatedAt:       time.Now(),
	}
	order.calculateTotal()
	return order
}

func (o *Order) calculateTotal() {
	var total float64
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	o.Total = total
}

func (o *Order) SetID(id string) {
	o.ID = id
}

func (o *Order) MarkAsPaid(txID string) {
	o.Status = "PAID"
	o.TransactionID = txID
}
