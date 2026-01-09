package db

import (
	"context"
	"errors"
	"time"

	"github.com/emersonmatsumoto/clean-go/products/internal/entities"
	"github.com/emersonmatsumoto/clean-go/products/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoRepo struct {
	collection *mongo.Collection
}

type productModel struct {
	ID    bson.ObjectID `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Price float64       `bson:"price"`
}

func NewMongoRepo(client *mongo.Client) ports.ProductRepository {
	return &mongoRepo{
		collection: client.Database("clean_db").Collection("products"),
	}
}

func (r *mongoRepo) FindByID(id string) (*entities.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("formato de ID inválido")
	}

	var model productModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &entities.Product{
		ID:    model.ID.Hex(),
		Name:  model.Name,
		Price: model.Price,
	}, nil
}
