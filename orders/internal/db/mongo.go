package db

import (
	"context"
	"time"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"github.com/emersonmatsumoto/clean-go/orders/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (r *mongoRepo) Save(order *entities.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var itemsModel []orderItemModel
	for _, item := range order.Items {
		prodID, err := bson.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return err
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
		return err
	}

	if insertID, ok := res.InsertedID.(bson.ObjectID); ok {
		order.ID = insertID.Hex()
	}

	return nil
}
