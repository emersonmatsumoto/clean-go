package usecases

import (
	"errors"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
)

type Repository interface {
	Save(order *entities.Order) error
}

type CreateOrderUseCase struct {
	repo     Repository
	prodComp products.Component
	payComp  payments.Component
}

func NewCreateOrderUseCase(r Repository, p products.Component, pay payments.Component) *CreateOrderUseCase {
	return &CreateOrderUseCase{repo: r, prodComp: p, payComp: pay}
}

func (uc *CreateOrderUseCase) Execute(itemsInput []entities.OrderItem, cardToken string) (*entities.Order, error) {
	var domainItems []entities.OrderItem

	for _, item := range itemsInput {
		p, err := uc.prodComp.GetProduct(products.GetProductInput{ID: item.ProductID})
		if err != nil {
			return nil, fmt.Errorf("produto %s não encontrado", item.ProductID)
		}

		domainItems = append(domainItems, entities.OrderItem{
			ProductID: p.ID,
			Price:     p.Price,
			Quantity:  item.Quantity,
		})
	}

	order := entities.NewOrder(domainItems)

	payRes, err := uc.payComp.ProcessPayment(payments.ProcessPaymentInput{
		OrderID:  order.ID,
		Amount:   order.Total,
		TokenID:  cardToken,
		Currency: "BRL",
	})

	if err != nil || payRes.Status != "SUCCESS" {
		return nil, errors.New("falha no pagamento")
	}

	order.MarkAsPaid(payRes.TransactionID)
	err = uc.repo.Save(order)

	return order, err
}
