// Package api provides HTTP routing and middleware configuration.
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/mabidoli/gravity-bff/internal/api/handler"
	"github.com/mabidoli/gravity-bff/internal/api/middleware"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// Router holds all the handlers and configures routes.
type Router struct {
	healthHandler *handler.HealthHandler
	streamHandler *handler.StreamHandler
	log           *logger.Logger
}

// NewRouter creates a new router with the given handlers.
func NewRouter(
	healthHandler *handler.HealthHandler,
	streamHandler *handler.StreamHandler,
	log *logger.Logger,
) *Router {
	return &Router{
		healthHandler: healthHandler,
		streamHandler: streamHandler,
		log:           log,
	}
}

// Setup configures all routes and middleware on the Fiber app.
func (r *Router) Setup(app *fiber.App) {
	// Global middleware
	app.Use(requestid.New())
	app.Use(middleware.Recovery(r.log))
	app.Use(middleware.RequestLogger(r.log))
	app.Use(middleware.CORSMiddleware())

	// Health check endpoint (no auth required)
	app.Get("/health", r.healthHandler.GetHealth)

	// API v2 routes
	v2 := app.Group("/v2")

	// Stream routes (auth required)
	stream := v2.Group("/stream", middleware.Auth())
	stream.Get("/", r.streamHandler.GetStream)
	stream.Get("/:itemId", r.streamHandler.GetStreamItem)
}
