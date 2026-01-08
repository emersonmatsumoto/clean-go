package db

import (
	"context"
	"time"

	"github.com/emersonmatsumoto/clean-go/orders/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	collection *mongo.Collection
}

type orderItemModel struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Price     float64            `bson:"price"`
	Quantity  int                `bson:"quantity"`
}

type orderModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Items         []orderItemModel   `bson:"items"`
	Total         float64            `bson:"total"`
	Status        string             `bson:"status"`
	TransactionID string             `bson:"transaction_id"`
	CreatedAt     time.Time          `bson:"created_at"`
}

func NewMongoRepo(client *mongo.Client) *MongoRepo {
	return &MongoRepo{
		collection: client.Database("clean_db").Collection("orders"),
	}
}

func (r *MongoRepo) Save(order *entities.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var itemsModel []orderItemModel
	for _, item := range order.Items {
		prodID, err := primitive.ObjectIDFromHex(item.ProductID)
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
		Items:         itemsModel,
		Total:         order.Total,
		Status:        order.Status,
		TransactionID: order.TransactionID,
		CreatedAt:     order.CreatedAt,
	}

	res, err := r.collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	if insertID, ok := res.InsertedID.(primitive.ObjectID); ok {
		order.ID = insertID.Hex()
	}

	return nil
}
