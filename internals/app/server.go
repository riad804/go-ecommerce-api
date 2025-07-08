package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/internals/middlewares"
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

	if config.Server.RateLimit.Enabled {
		app.Use(middlewares.LimiterMiddleware(redisClient.Client, config.Server.RateLimit.RateLimit, config.Server.RateLimit.RateLimitWindow))
	}

	routes := routes.NewRoutes(config, app, mongoConn, distributor)
	routes.NewAuthRoutes()
	routes.NewUserRoutes()

	return &Server{
		config: config,
		app:    app,
	}
}

func (s *Server) Start() {
	s.app.Listen(fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port))
}
