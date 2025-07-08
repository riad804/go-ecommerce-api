package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/config"
	"github.com/riad804/go_ecommerce_api/pkg/database"
	"github.com/riad804/go_ecommerce_api/token"
	"github.com/riad804/go_ecommerce_api/workers"
)

const apiPrefix = "/api/v1"

type Routes struct {
	api         fiber.Router
	tokenMaker  *token.Maker
	Mongo       *database.MongoConnection
	Validator   *validator.Validate
	cfg         *config.Config
	distributor workers.TaskDistributor
}

func NewRoutes(cfg *config.Config, app *fiber.App, mongoConn *database.MongoConnection, distributor workers.TaskDistributor) *Routes {
	maker, err := token.NewPasetoMaker(cfg.Token.AccessKey, cfg.Token.RefreshKey)
	if err != nil {
		panic("failed to create token maker: " + err.Error())
	}

	validate := validator.New()
	validate.RegisterValidation("password", config.IsValidPassword)

	api := app.Group(apiPrefix)

	return &Routes{
		api:         api,
		tokenMaker:  &maker,
		Mongo:       mongoConn,
		Validator:   validate,
		cfg:         cfg,
		distributor: distributor,
	}
}
