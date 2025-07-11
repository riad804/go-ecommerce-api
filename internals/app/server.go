package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/middlewares"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/internals/repositories"
	"github.com/riad804/go_ecommerce_api/internals/routes"
	"github.com/riad804/go_ecommerce_api/pkg/database"
	"github.com/riad804/go_ecommerce_api/pkg/redis"
	"github.com/riad804/go_ecommerce_api/workers"
)

type Server struct {
	config *config.Config
	app    *fiber.App
	Mongo  *database.MongoConnection
}

func NewServer(config *config.Config, redisClient *redis.RedisClient, mongoConn *database.MongoConnection, distributor workers.TaskDistributor) *Server {

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})
	app.Use(middlewares.LogMiddleware())
	app.Use(cors.New())
	app.Use(cache.New())
	app.Use(compress.New())
	app.Static("/", "./public")

	if config.Server.RateLimit.Enabled {
		app.Use(middlewares.LimiterMiddleware(redisClient.Client, config.Server.RateLimit.RateLimit, config.Server.RateLimit.RateLimitWindow))
	}

	routes := routes.NewRoutes(config, app, mongoConn, distributor)
	routes.NewAuthRoutes()
	routes.NewUserRoutes()
	routes.NewAdminRoutes()

	return &Server{
		config: config,
		app:    app,
	}
}

func (s *Server) Start() {
	s.app.Listen(fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port))
}

func (s *Server) StartCron() {
	c := cron.New()
	c.AddFunc("0 0 * * *", s.DeleteOldCategories)
	c.Start()
	select {}
}

func (s *Server) DeleteOldCategories() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.Mongo.Database.Collection(repositories.CATEGORIES)
	filter := bson.M{"marked_for_delete": true}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("Error finding categories", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cat models.Category
		if err := cursor.Decode(&cat); err != nil {
			log.Println("Decode error:", err)
			continue
		}

		imagePath := "./public/" + cat.Image
		if err := os.Remove(imagePath); err != nil {
			log.Println("could not delete image:", imagePath, err)
		} else {
			fmt.Println("Deleted image:", imagePath)
		}

		_, err := collection.DeleteOne(ctx, bson.M{"_id": cat.ID})
		if err != nil {
			log.Println("error deleting category:", err)
		} else {
			fmt.Println("deleted category Id:", cat.ID)
		}
	}

	log.Println("running cron job for deleting categories")
}
