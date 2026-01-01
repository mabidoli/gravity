package middleware

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// Recovery returns a middleware that recovers from panics.
func Recovery(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Panic recovered: %v\n%s", r, debug.Stack())

				// Return 500 error
				_ = c.Status(fiber.StatusInternalServerError).JSON(model.NewErrorResponse(
					model.ErrCodeInternalError,
					"An unexpected error occurred",
				))
			}
		}()

		return c.Next()
	}
}
