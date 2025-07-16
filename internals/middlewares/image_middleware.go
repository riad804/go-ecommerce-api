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

// i know that it's not suitable to upload image in middleware, cuase
// if there's something happen after uploading images there's no way to
// get back, so we should upload images after all the validation has been completed
func ImageMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		hasImage := false
		if file, err := c.FormFile("image"); err == nil && file != nil {
			hasImage = true
		}
		if form, err := c.MultipartForm(); err == nil {
			if files := form.File["images"]; len(files) > 0 {
				hasImage = true
			}
		}

		if hasImage {

			const maxFileSize = 2 * 1024 * 1024 // 2MB max
			const maxImageCount = 10
			allowedExts := map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
			}
			uploadPath := "./public/uploads"

			if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
				return models.Error(c, fiber.StatusInternalServerError, "Could not create upload directory")
			}

			var uploadedURLs []string

			// Handle single image upload
			if file, err := c.FormFile("image"); err == nil {
				if file.Size > maxFileSize {
					return models.Error(c, fiber.StatusRequestEntityTooLarge, "Single file too large (max 2MB)")
				}

				ext := strings.ToLower(filepath.Ext(file.Filename))
				if !allowedExts[ext] {
					return models.Error(c, fiber.StatusBadRequest, "Only JPG and PNG files are allowed")
				}

				newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
				savePath := filepath.Join(uploadPath, newFileName)

				if err := c.SaveFile(file, savePath); err != nil {
					return models.Error(c, fiber.StatusInternalServerError, "Could not save image file")
				}

				uploadedURLs = append(uploadedURLs, fmt.Sprintf("uploads/%s", newFileName))
				c.Locals("image_url", uploadedURLs[0]) // single image
			}

			// Handle multiple image uploads (images[])
			if files, err := c.MultipartForm(); err == nil {
				if imgFiles, ok := files.File["images"]; ok {
					if len(imgFiles) > maxImageCount {
						return models.Error(c, fiber.StatusBadRequest, fmt.Sprintf("Maximum %d images allowed", maxImageCount))
					}
					for _, file := range imgFiles {
						if file.Size > maxFileSize {
							return models.Error(c, fiber.StatusRequestEntityTooLarge, "One of the files is too large (max 2MB)")
						}

						ext := strings.ToLower(filepath.Ext(file.Filename))
						if !allowedExts[ext] {
							return models.Error(c, fiber.StatusBadRequest, "Only JPG and PNG files are allowed")
						}

						newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
						savePath := filepath.Join(uploadPath, newFileName)

						if err := c.SaveFile(file, savePath); err != nil {
							return models.Error(c, fiber.StatusInternalServerError, "Could not save one of the image files")
						}

						uploadedURLs = append(uploadedURLs, fmt.Sprintf("uploads/%s", newFileName))
					}

					// Set all uploaded image URLs
					c.Locals("image_urls", uploadedURLs)
				}
			}

		}

		return c.Next()
	}
}
