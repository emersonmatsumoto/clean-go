package db

import (
	"context"
	"time"

	"github.com/emersonmatsumoto/clean-go/users/internal/entities"
	"github.com/emersonmatsumoto/clean-go/users/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type addressModel struct {
	Street  string `bson:"street"`
	City    string `bson:"city"`
	ZipCode string `bson:"zip_code"`
}

type userModel struct {
	ID      bson.ObjectID `bson:"_id,omitempty"`
	Name    string        `bson:"name"`
	Email   string        `bson:"email"`
	Address addressModel  `bson:"address"`
}

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(client *mongo.Client) ports.UserRepository {
	return &mongoRepo{
		collection: client.Database("clean_db").Collection("users"),
	}
}

func (r *mongoRepo) FindByID(id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var model userModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &entities.User{
		ID:    model.ID.Hex(),
		Name:  model.Name,
		Email: model.Email,
		Address: entities.Address{
			Street:  model.Address.Street,
			City:    model.Address.City,
			ZipCode: model.Address.ZipCode,
		},
	}, nil
}
