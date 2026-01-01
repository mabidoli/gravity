// Package handler contains HTTP request handlers for the API.
package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

// HealthHandler handles health check requests.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// GetHealth returns the health status of the API.
// @Summary Health check
// @Description Returns the health status of the API
// @Tags health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /health [get]
func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	response := model.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.JSON(response)
}
