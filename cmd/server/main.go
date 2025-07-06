package main

import (
	"context"
	"log"

	"github.com/riad804/go_ecommerce_api/internals/app"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/pkg/database"
	"github.com/riad804/go_ecommerce_api/pkg/redis"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load configuration:", err)
	}

	mongoConn, err := database.NewMongoDB(config)
	if err != nil {
		log.Fatal("failed to connect to MongoDB:", err)
	}
	defer mongoConn.Client.Disconnect(context.Background())

	redisClient, err := redis.NewRedisClient(config)
	if err != nil {
		log.Fatal("failed to connect to Redis:", err)
	}
	defer redisClient.WithContext(context.Background()).Close()

	log.Println("✅ Successfully connected to MongoDB and Redis!")

	app.NewServer(config, redisClient, mongoConn).Start()
}
