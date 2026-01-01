package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mabidoli/gravity-bff/internal/config"
	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/internal/service"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// MockStreamRepository for testing
type MockStreamRepository struct {
	mock.Mock
}

func (m *MockStreamRepository) GetStream(ctx context.Context, req model.StreamRequest) ([]model.PriorityItem, *string, error) {
	args := m.Called(ctx, req)
	items := args.Get(0)
	if items == nil {
		return nil, args.Get(1).(*string), args.Error(2)
	}
	return items.([]model.PriorityItem), args.Get(1).(*string), args.Error(2)
}

func (m *MockStreamRepository) GetStreamItemByID(ctx context.Context, userID, itemID string) (*model.PriorityItem, error) {
	args := m.Called(ctx, userID, itemID)
	item := args.Get(0)
	if item == nil {
		return nil, args.Error(1)
	}
	return item.(*model.PriorityItem), args.Error(1)
}

func (m *MockStreamRepository) GetParticipantsByItemID(ctx context.Context, itemID string) ([]model.User, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockStreamRepository) GetMessagesByItemID(ctx context.Context, itemID string) ([]model.Message, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).([]model.Message), args.Error(1)
}

// MockCache for testing
type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetStream(ctx context.Context, key string) (*model.StreamResponse, error) {
	args := m.Called(ctx, key)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*model.StreamResponse), args.Error(1)
}

func (m *MockCache) SetStream(ctx context.Context, key string, data *model.StreamResponse, ttl time.Duration) error {
	args := m.Called(ctx, key, data, ttl)
	return args.Error(0)
}

func (m *MockCache) GetStreamItem(ctx context.Context, key string) (*model.PriorityItem, error) {
	args := m.Called(ctx, key)
	item := args.Get(0)
	if item == nil {
		return nil, args.Error(1)
	}
	return item.(*model.PriorityItem), args.Error(1)
}

func (m *MockCache) SetStreamItem(ctx context.Context, key string, item *model.PriorityItem, ttl time.Duration) error {
	args := m.Called(ctx, key, item, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockCache) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Test setup helpers
func setupTestApp(handler *StreamHandler) *fiber.App {
	app := fiber.New()

	// Add middleware to set userID (simulating auth)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user")
		return c.Next()
	})

	app.Get("/v2/stream", handler.GetStream)
	app.Get("/v2/stream/:itemId", handler.GetStreamItem)

	return app
}

func newTestConfig() *config.Config {
	return &config.Config{
		Cache: config.CacheConfig{
			StreamTTL: 2 * time.Minute,
			ItemTTL:   5 * time.Minute,
		},
	}
}

// Tests for GetStream handler
func TestStreamHandler_GetStream_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	log := logger.New()
	cfg := newTestConfig()

	svc := service.NewStreamService(mockRepo, mockCache, cfg, log)
	handler := NewStreamHandler(svc, log)
	app := setupTestApp(handler)

	expectedItems := []model.PriorityItem{
		{
			ID:       "item-1",
			Title:    "Test Item",
			Source:   model.SourceEmail,
			Priority: model.PriorityHigh,
			IsUnread: true,
		},
	}

	// Cache miss, repo returns data
	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("GetStream", mock.Anything, mock.Anything).Return(expectedItems, (*string)(nil), nil)
	mockCache.On("SetStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	req := httptest.NewRequest("GET", "/v2/stream", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.StreamResponse
	err = json.Unmarshal(body, &result)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, "item-1", result.Data[0].ID)
}

func TestStreamHandler_GetStream_WithFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	log := logger.New()
	cfg := newTestConfig()

	svc := service.NewStreamService(mockRepo, mockCache, cfg, log)
	handler := NewStreamHandler(svc, log)
	app := setupTestApp(handler)

	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("GetStream", mock.Anything, mock.MatchedBy(func(req model.StreamRequest) bool {
		return req.Filter == model.FilterHigh
	})).Return([]model.PriorityItem{}, (*string)(nil), nil)
	mockCache.On("SetStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	req := httptest.NewRequest("GET", "/v2/stream?filter=high", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestStreamHandler_GetStream_InvalidFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	log := logger.New()
	cfg := newTestConfig()

	svc := service.NewStreamService(mockRepo, mockCache, cfg, log)
	handler := NewStreamHandler(svc, log)
	app := setupTestApp(handler)

	// Act
	req := httptest.NewRequest("GET", "/v2/stream?filter=invalid", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.ErrorResponse
	json.Unmarshal(body, &result)

	assert.Equal(t, model.ErrCodeValidationFailed, result.Error.Code)
}

func TestStreamHandler_GetStream_WithPagination(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	log := logger.New()
	cfg := newTestConfig()

	svc := service.NewStreamService(mockRepo, mockCache, cfg, log)
	handler := NewStreamHandler(svc, log)
	app := setupTestApp(handler)

	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("GetStream", mock.Anything, mock.MatchedBy(func(req model.StreamRequest) bool {
		return req.Limit == 50 && req.Cursor != nil && *req.Cursor == "abc123"
	})).Return([]model.PriorityItem{}, (*string)(nil), nil)
	mockCache.On("SetStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	req := httptest.NewRequest("GET", "/v2/stream?limit=50&cursor=abc123", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// Tests for GetStreamItem handler
func TestStreamHandler_GetStreamItem_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	log := logger.New()
	cfg := newTestConfig()

	svc := service.NewStreamService(mockRepo, mockCache, cfg, log)
	handler := NewStreamHandler(svc, log)
	app := setupTestApp(handler)

	expectedItem := &model.PriorityItem{
		ID:       "item-123",
		Title:    "Test Item",
		Source:   model.SourceEmail,
		Priority: model.PriorityHigh,
		Messages: []model.Message{
			{ID: "msg-1", Content: "Hello"},
		},
	}

	mockCache.On("GetStreamItem", mock.Anything, "item:item-123").Return(nil, nil)
	mockRepo.On("GetStreamItemByID", mock.Anything, "test-user", "item-123").Return(expectedItem, nil)
	mockCache.On("SetStreamItem", mock.Anything, "item:item-123", expectedItem, mock.Anything).Return(nil)

	// Act
	req := httptest.NewRequest("GET", "/v2/stream/item-123", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.PriorityItem
	json.Unmarshal(body, &result)

	assert.Equal(t, "item-123", result.ID)
	assert.Equal(t, 1, len(result.Messages))
}

func TestStreamHandler_GetStreamItem_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	log := logger.New()
	cfg := newTestConfig()

	svc := service.NewStreamService(mockRepo, mockCache, cfg, log)
	handler := NewStreamHandler(svc, log)
	app := setupTestApp(handler)

	mockCache.On("GetStreamItem", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("GetStreamItemByID", mock.Anything, "test-user", "nonexistent").Return(nil, nil)

	// Act
	req := httptest.NewRequest("GET", "/v2/stream/nonexistent", nil)
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.ErrorResponse
	json.Unmarshal(body, &result)

	assert.Equal(t, model.ErrCodeNotFound, result.Error.Code)
}
