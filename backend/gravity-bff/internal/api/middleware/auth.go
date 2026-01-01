package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

// ============================================================================
// WARNING: DEVELOPMENT-ONLY AUTHENTICATION MIDDLEWARE
// ============================================================================
// This authentication middleware is a PLACEHOLDER for development and testing.
// It MUST be replaced with proper JWT validation before production deployment.
//
// Current limitations (INSECURE):
// - No JWT signature verification
// - Bearer token used directly as user ID
// - Fallback to hardcoded "default-user" when no auth header present
//
// Production requirements:
// - Implement proper JWT validation (RS256 or ES256)
// - Verify token signature against public key
// - Validate token expiration (exp claim)
// - Extract user ID from validated token claims
// - Remove default-user fallback
// - Add rate limiting for failed auth attempts
// ============================================================================

// Auth returns a middleware that validates authentication.
//
// SECURITY WARNING: This is a development-only placeholder.
// DO NOT use in production without implementing proper JWT validation.
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")

		// WARNING: Development-only bypass - remove before production!
		// For development, allow requests without auth
		if authHeader == "" {
			// INSECURE: Hardcoded default user - must be removed for production
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

		// SECURITY TODO: Implement proper JWT validation here:
		// 1. Parse the JWT token
		// 2. Verify signature against public key
		// 3. Check expiration (exp) and not-before (nbf) claims
		// 4. Validate issuer (iss) and audience (aud) claims
		// 5. Extract user ID from claims (sub or custom claim)
		//
		// Example with a JWT library:
		// claims, err := jwt.ValidateToken(token, publicKey)
		// if err != nil { return unauthorized }
		// userID := claims.Subject

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.NewErrorResponse(
				model.ErrCodeUnauthorized,
				"Token is required",
			))
		}

		// INSECURE: Using token directly as user ID - replace with JWT claims
		// Set user ID in context
		c.Locals("userID", token)

		return c.Next()
	}
}
