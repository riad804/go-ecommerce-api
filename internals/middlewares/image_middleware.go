package middlewares

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/riad804/go_ecommerce_api/internals/models"
)

func ImageMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		const maxFileSize = 2 * 1024 * 1024

		file, err := c.FormFile("image")
		if err != nil {
			return models.Error(c, fiber.StatusBadRequest, "No file uploaded")
		}

		if file.Size > maxFileSize {
			return models.Error(c, fiber.StatusRequestEntityTooLarge, "File too large (max 2MB)")
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return models.Error(c, fiber.StatusBadRequest, "Only JPG and PNG files are allowed")
		}

		uploadPath := "./public/uploads"
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			return models.Error(c, fiber.StatusInternalServerError, "Could not create upload directory")
		}

		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath := filepath.Join(uploadPath, newFileName)

		if err := c.SaveFile(file, savePath); err != nil {
			return models.Error(c, fiber.StatusInternalServerError, "Could not save file")
		}

		c.Locals("image_url", fmt.Sprintf("uploads/%s", newFileName))
		return c.Next()
	}
}
