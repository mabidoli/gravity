package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

// Auth returns a middleware that validates authentication.
// For now, this is a placeholder that extracts user ID from header.
// In production, this would validate JWT tokens.
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")

		// For development, allow requests without auth
		if authHeader == "" {
			// Set a default user ID for development
			c.Locals("userID", "default-user")
			return c.Next()
		}

		// Extract token (Bearer <token>)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Invalid authorization header format",
			))
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// TODO: In production, validate JWT token here
		// For now, we'll use the token as the user ID (for testing)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Token is required",
			))
		}

		// Set user ID in context
		// In production, this would be extracted from the validated JWT
		c.Locals("userID", token)

		return c.Next()
	}
}
