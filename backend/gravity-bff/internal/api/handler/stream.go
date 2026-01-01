package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/internal/service"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// StreamHandler handles stream-related HTTP requests.
type StreamHandler struct {
	service *service.StreamService
	log     *logger.Logger
}

// NewStreamHandler creates a new stream handler.
func NewStreamHandler(svc *service.StreamService, log *logger.Logger) *StreamHandler {
	return &StreamHandler{
		service: svc,
		log:     log,
	}
}

// GetStream handles GET /v2/stream requests.
// @Summary Get unified stream
// @Description Retrieves a paginated list of priority items for the authenticated user
// @Tags stream
// @Accept json
// @Produce json
// @Param filter query string false "Filter items (all, high, unread)" default(all)
// @Param limit query int false "Maximum items to return" default(20) minimum(1) maximum(100)
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} model.StreamResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v2/stream [get]
func (h *StreamHandler) GetStream(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID")
	if userID == nil {
		// For development/testing, use a default user ID
		userID = "default-user"
	}

	// Parse query parameters
	filterStr := c.Query("filter", "all")
	filter, err := service.ValidateFilter(filterStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(
			model.ErrCodeValidationFailed,
			err.Error(),
		))
	}

	limitStr := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	cursor := c.Query("cursor")
	var cursorPtr *string
	if cursor != "" {
		cursorPtr = &cursor
	}

	// Build request
	req := model.StreamRequest{
		UserID: userID.(string),
		Filter: filter,
		Limit:  limit,
		Cursor: cursorPtr,
	}

	// Call service
	response, err := h.service.GetStream(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to get stream: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewErrorResponse(
			model.ErrCodeInternalError,
			"Failed to retrieve stream",
		))
	}

	return c.JSON(response)
}

// GetStreamItem handles GET /v2/stream/:itemId requests.
// @Summary Get stream item details
// @Description Retrieves full details of a single priority item including messages
// @Tags stream
// @Accept json
// @Produce json
// @Param itemId path string true "Priority item ID"
// @Success 200 {object} model.PriorityItem
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v2/stream/{itemId} [get]
func (h *StreamHandler) GetStreamItem(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("userID")
	if userID == nil {
		// For development/testing, use a default user ID
		userID = "default-user"
	}

	// Get item ID from path
	itemID := c.Params("itemId")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(
			model.ErrCodeValidationFailed,
			"Item ID is required",
		))
	}

	// Build request
	req := model.StreamItemRequest{
		UserID: userID.(string),
		ItemID: itemID,
	}

	// Call service
	item, err := h.service.GetStreamItemDetails(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to get stream item: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewErrorResponse(
			model.ErrCodeInternalError,
			"Failed to retrieve item",
		))
	}

	if item == nil {
		return c.Status(fiber.StatusNotFound).JSON(model.NewErrorResponse(
			model.ErrCodeNotFound,
			"The requested priority item does not exist",
		))
	}

	return c.JSON(item)
}
