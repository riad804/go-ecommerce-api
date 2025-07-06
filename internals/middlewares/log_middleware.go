package middlewares

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LogEntry struct {
	Timestamp string        `json:"timestamp"`
	Method    string        `json:"method"`
	Path      string        `json:"path"`
	Status    int           `json:"status"`
	Latency   time.Duration `json:"latency"`
	IP        string        `json:"ip"`
}

func LogMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Method:    c.Method(),
			Path:      c.OriginalURL(),
			Status:    c.Response().StatusCode(),
			Latency:   latency,
			IP:        c.IP(),
		}

		logJson, _ := json.Marshal(entry)
		fmt.Println(string(logJson))
		return err
	}
}
