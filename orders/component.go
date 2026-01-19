package orders

import (
	"context"

	"github.com/emersonmatsumoto/clean-go/contracts/orders"
	"github.com/emersonmatsumoto/clean-go/contracts/payments"
	"github.com/emersonmatsumoto/clean-go/contracts/products"
	"github.com/emersonmatsumoto/clean-go/contracts/users"
	"github.com/emersonmatsumoto/clean-go/orders/internal/db"
	"github.com/emersonmatsumoto/clean-go/orders/internal/local"
	"github.com/emersonmatsumoto/clean-go/orders/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
)

type component struct {
	createUC *usecases.CreateOrderUseCase
}

func NewComponent(
	mongoClient *mongo.Client,
	prodComp products.Component,
	payComp payments.Component,
	userComp users.Component,
) orders.Component {
	uc := usecases.NewCreateOrderUseCase(
		db.NewMongoRepo(mongoClient),
		local.NewProductGateway(prodComp),
		local.NewPaymentGateway(payComp),
		local.NewUserGateway(userComp),
	)

	return &component{
		createUC: uc,
	}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/orders")

func (c *component) PlaceOrder(ctx context.Context, in orders.PlaceOrderInput) (orders.PlaceOrderOutput, error) {
	ctx, span := tracer.Start(ctx, "Orders.Component.PlaceOrder")
	defer span.End()

	return c.createUC.Execute(ctx, in)
}
