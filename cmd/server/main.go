package main

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/riad804/go_ecommerce_api/internals/app"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/mail"
	"github.com/riad804/go_ecommerce_api/pkg/database"
	"github.com/riad804/go_ecommerce_api/pkg/redis"
	"github.com/riad804/go_ecommerce_api/workers"
)

func main() {
	config, err := config.LoadConfig(".")
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

	redisOpt := asynq.RedisClientOpt{
		Addr: redisClient.Options().Addr,
	}
	distributor := workers.NewRedisTaskDistributor(redisOpt)
	go RunTaskProcessor(config, redisOpt)

	server := app.NewServer(config, redisClient, mongoConn, distributor)
	server.Start()
	server.StartCron()
}

func RunTaskProcessor(config *config.Config, redisOpt asynq.RedisClientOpt) {
	mailer := mail.NewGmailSender(config.Email.Name, config.Email.Address, config.Email.Password)
	taskProcessor := workers.NewRedisTaskProcessor(redisOpt, mailer, config)
	log.Default().Println("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Default().Println("failed to start task processor")
	}
}
