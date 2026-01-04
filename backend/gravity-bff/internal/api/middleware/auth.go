package middleware

import (
	"context"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

// clerkInitialized tracks whether Clerk has been configured.
var clerkInitialized bool

// InitClerk initializes the Clerk SDK with the secret key.
// This should be called during application startup.
func InitClerk() error {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return nil // Allow uninitialized for development mode
	}

	clerk.SetKey(secretKey)
	clerkInitialized = true
	return nil
}

// ClerkAuth returns a middleware that validates Clerk JWT tokens.
// It extracts the user ID from the validated token and sets it in the context.
func ClerkAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")

		// Development mode: allow requests without auth if CLERK_SECRET_KEY is not set
		if !clerkInitialized {
			if authHeader == "" {
				// Development fallback - allows unauthenticated access
				c.Locals("userID", "dev-user")
				return c.Next()
			}
		}

		// Require authorization header in production
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Authorization header is required",
			))
		}

		// Extract token (Bearer <token>)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Invalid authorization header format",
			))
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Token is required",
			))
		}

		// Development mode with token but Clerk not initialized
		if !clerkInitialized {
			// In development, use token as user ID for testing
			c.Locals("userID", "dev-user")
			return c.Next()
		}

		// Verify the JWT token using Clerk
		claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
			Token: token,
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Invalid or expired token",
			))
		}

		// Extract user ID from the subject claim
		userID := claims.Subject
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Token missing user identifier",
			))
		}

		// Set user ID in context for downstream handlers
		c.Locals("userID", userID)

		return c.Next()
	}
}

// Auth is an alias for ClerkAuth for backward compatibility.
// Deprecated: Use ClerkAuth() instead.
func Auth() fiber.Handler {
	return ClerkAuth()
}
