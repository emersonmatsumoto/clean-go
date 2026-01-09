package db

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func startMongoContainer(ctx context.Context, t *testing.T) *mongo.Client {
	mongodbContainer, err := mongodb.Run(ctx, "mongo:8.2")
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %s", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("failed to ping MongoDB: %s", err)
	}

	t.Cleanup(func() {
		// Desconecta o cliente MongoDB com um contexto em branco para garantir a desconexão caso o contexto original tenha expirado
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			t.Errorf("warning: failed to disconnect mongo client: %v", err)
		}

		if err := testcontainers.TerminateContainer(mongodbContainer); err != nil {
			// Para evitar falhas silenciosas (t.Fatalf executa runtime.Goexit()), reportamos o erro de término
			t.Errorf("failed to terminate container: %s", err)
		}
	})

	return mongoClient
}

func TestProductRepository_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := startMongoContainer(ctx, t)
	repo := NewMongoRepo(client)
	db := client.Database("clean_db")
	coll := db.Collection("products")

	setupTest := func(t *testing.T) {
		t.Helper()
		if err := coll.Drop(ctx); err != nil {
			t.Fatalf("failed to drop collection: %v", err)
		}
	}

	t.Run("Success_FindByID", func(t *testing.T) {
		setupTest(t)

		ins, err := coll.InsertOne(ctx, bson.M{"name": "Integration Product", "price": 42.5})
		if err != nil {
			t.Fatalf("insert failed: %v", err)
		}

		oid, ok := ins.InsertedID.(bson.ObjectID)
		if !ok {
			t.Fatalf("inserted id is not an ObjectID: %T", ins.InsertedID)
		}

		prod, err := repo.FindByID(oid.Hex())
		if err != nil {
			t.Fatalf("FindByID returned error: %v", err)
		}
		if prod == nil {
			t.Fatalf("expected product, got nil")
		}
		if prod.ID != oid.Hex() || prod.Name != "Integration Product" || prod.Price != 42.5 {
			t.Fatalf("unexpected product: %+v", prod)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		setupTest(t)

		id := bson.NewObjectID().Hex()
		prod, err := repo.FindByID(id)
		if err != nil {
			t.Fatalf("expected no error for not found, got: %v", err)
		}
		if prod != nil {
			t.Fatalf("expected nil product for not found, got: %+v", prod)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		setupTest(t)

		prod, err := repo.FindByID("bad-id")
		if prod != nil {
			t.Fatalf("expected nil product for invalid id, got: %+v", prod)
		}
		if err == nil {
			t.Fatalf("expected error for invalid id, got nil")
		}
	})
}
