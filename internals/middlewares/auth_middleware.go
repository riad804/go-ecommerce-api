package middlewares

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/riad804/go_ecommerce_api/internals/models"
	"github.com/riad804/go_ecommerce_api/token"
)

func AuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return models.Error(c, fiber.StatusUnauthorized, "Missing authorization")
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			return models.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyAccessToken(accessToken)
		if err != nil {
			return models.Error(c, fiber.StatusUnauthorized, "Invalid access token")
		}

		path := c.OriginalURL()
		adminPathRegex := regexp.MustCompile(`^/api/v1/admin/.*`)
		if !payload.IsAdmin && adminPathRegex.MatchString(path) {
			return models.Error(c, fiber.StatusUnauthorized, "access token is not for admin")
		}

		c.Locals("email", payload.Email)
		c.Locals("is_admin", payload.IsAdmin)

		return c.Next()
	}
}
