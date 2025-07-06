package database

import (
	"context"
	"fmt"
	"log"

	"github.com/riad804/go_ecommerce_api/internals/config"
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

	return &MongoConnection{
		Client:   client,
		Database: client.Database(cfg.MongoDB.Database),
	}, nil
}
