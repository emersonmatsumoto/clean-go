package db

import (
	"context"
	"errors"
	"time"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel"
)

type mongoRepo struct {
	collection *mongo.Collection
}

type orderItemModel struct {
	ProductID bson.ObjectID `bson:"product_id"`
	Price     float64       `bson:"price"`
	Quantity  int           `bson:"quantity"`
}

type orderModel struct {
	ID              bson.ObjectID    `bson:"_id,omitempty"`
	Items           []orderItemModel `bson:"items"`
	Total           float64          `bson:"total"`
	Status          string           `bson:"status"`
	ShippingAddress string           `bson:"shipping_address"`
	TransactionID   string           `bson:"transaction_id"`
	UserID          string           `bson:"user_id"`
	CreatedAt       time.Time        `bson:"created_at"`
}

func NewMongoRepo(client *mongo.Client) ports.OrderRepository {
	return &mongoRepo{
		collection: client.Database("clean_db").Collection("orders"),
	}
}

var tracer = otel.Tracer("github.com/emersonmatsumoto/clean-go/orders/internal/db")

func (r *mongoRepo) Save(ctx context.Context, order *entities.Order) (string, error) {
	ctx, span := tracer.Start(ctx, "Orders.MongoRepo.Save")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var itemsModel []orderItemModel
	for _, item := range order.Items {
		prodID, err := bson.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return "", err
		}
		itemsModel = append(itemsModel, orderItemModel{
			ProductID: prodID,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
	}

	model := orderModel{
		Items:           itemsModel,
		Total:           order.Total,
		Status:          order.Status,
		TransactionID:   order.TransactionID,
		UserID:          order.UserID,
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       order.CreatedAt,
	}

	res, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return "", err
	}

	insertID, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return "", errors.New("failed to convert inserted id to objectID")
	}

	return insertID.Hex(), nil
}
