package usecases

import (
	"errors"
	"fmt"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
	"github.com/emersonmatsumoto/clean-go/users"
)

type CreateOrderUseCase struct {
	repo     ports.OrderRepository
	prodComp products.Component
	payComp  payments.Component
	userComp users.Component
}

func NewCreateOrderUseCase(r ports.OrderRepository, p products.Component, pay payments.Component, user users.Component) *CreateOrderUseCase {
	return &CreateOrderUseCase{repo: r, prodComp: p, payComp: pay, userComp: user}
}

func (uc *CreateOrderUseCase) Execute(userID string, itemsInput []entities.OrderItem, cardToken string) (*entities.Order, error) {
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

	userData, err := uc.userComp.GetUser(users.GetUserInput{ID: userID})
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	addressStr := fmt.Sprintf("%s, %s - %s", userData.Address.Street, userData.Address.City, userData.Address.ZipCode)
	order := entities.NewOrder(userID, domainItems, addressStr)

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
