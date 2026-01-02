package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSMiddleware returns a configured CORS middleware.
// It reads allowed origins from the CORS_ORIGINS environment variable.
// Defaults to localhost:3000 for development.
func CORSMiddleware() fiber.Handler {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000"
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     strings.Join([]string{fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete, fiber.MethodOptions}, ","),
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: true,
		ExposeHeaders:    "X-Request-ID,X-Cursor",
		MaxAge:           86400, // 24 hours
	})
}
