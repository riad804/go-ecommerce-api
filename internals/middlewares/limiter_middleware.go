package middlewares

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func LimiterMiddleware(client *redis.Client, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		print(window.Minutes())
		ctx := context.Background() // Redis context
		ip := c.IP()
		key := "rate_limit:" + ip

		current, err := client.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		if current >= limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}

		// Increment the request count
		_, err = client.Incr(ctx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// Set expiry only when the key is new
		if current == 0 {
			client.Expire(ctx, key, window)
		}
		return c.Next()
	}
}
