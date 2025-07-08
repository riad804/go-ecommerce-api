package database

import (
	"context"
	"fmt"
	"log"

	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(cfg *config.Config) (*MongoConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Ping failed:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")

	db := client.Database(cfg.MongoDB.Database)

	err = models.EnsureUserIndexes(db.Collection(repositories.USERS))
	if err != nil {
		log.Fatal("User indexing failed:", err)
	}

	return &MongoConnection{
		Client:   client,
		Database: db,
	}, nil
}
