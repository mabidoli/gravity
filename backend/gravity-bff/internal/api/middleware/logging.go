// Package middleware provides HTTP middleware for the API.
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// RequestLogger returns a middleware that logs HTTP requests.
func RequestLogger(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)
		log.Info("%s %s %d %v",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
