package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	// Arrange
	handler := NewHealthHandler()
	app := fiber.New()
	app.Get("/health", handler.GetHealth)

	// Act
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.HealthResponse
	err = json.Unmarshal(body, &result)

	assert.NoError(t, err)
	assert.Equal(t, "ok", result.Status)
	assert.NotEmpty(t, result.Timestamp)
}

func TestHealthHandler_GetHealth_ContentType(t *testing.T) {
	// Arrange
	handler := NewHealthHandler()
	app := fiber.New()
	app.Get("/health", handler.GetHealth)

	// Act
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
}
