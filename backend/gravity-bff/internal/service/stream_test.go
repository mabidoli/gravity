package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mabidoli/gravity-bff/internal/config"
	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// MockStreamRepository is a mock implementation of StreamRepository.
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

// MockCache is a mock implementation of Cache.
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

// Test helpers
func newTestConfig() *config.Config {
	return &config.Config{
		Cache: config.CacheConfig{
			StreamTTL: 2 * time.Minute,
			ItemTTL:   5 * time.Minute,
		},
	}
}

func newTestService(repo *MockStreamRepository, cache *MockCache) *StreamService {
	return NewStreamService(repo, cache, newTestConfig(), logger.New())
}

// Tests for GetStream
func TestStreamService_GetStream_CacheHit(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	expectedItems := []model.PriorityItem{
		{ID: "item-1", Title: "Test Item 1", Priority: model.PriorityHigh},
		{ID: "item-2", Title: "Test Item 2", Priority: model.PriorityMedium},
	}
	expectedResponse := &model.StreamResponse{
		Data:       expectedItems,
		NextCursor: nil,
	}

	req := model.StreamRequest{
		UserID: "user-123",
		Filter: model.FilterAll,
		Limit:  20,
	}

	// Cache returns data (hit)
	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Act
	result, err := svc.GetStream(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, "item-1", result.Data[0].ID)

	// Repository should NOT be called on cache hit
	mockRepo.AssertNotCalled(t, "GetStream")
	mockCache.AssertExpectations(t)
}

func TestStreamService_GetStream_CacheMiss(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	expectedItems := []model.PriorityItem{
		{ID: "item-1", Title: "Test Item 1", Priority: model.PriorityHigh},
	}

	req := model.StreamRequest{
		UserID: "user-123",
		Filter: model.FilterAll,
		Limit:  20,
	}

	// Cache miss
	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	// Repository returns data
	mockRepo.On("GetStream", mock.Anything, req).Return(expectedItems, (*string)(nil), nil)
	// Cache set
	mockCache.On("SetStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	result, err := svc.GetStream(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, "item-1", result.Data[0].ID)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestStreamService_GetStream_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	req := model.StreamRequest{
		UserID: "user-123",
		Filter: model.FilterAll,
		Limit:  20,
	}

	// Cache miss
	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	// Repository returns error
	mockRepo.On("GetStream", mock.Anything, req).Return(nil, (*string)(nil), errors.New("database error"))

	// Act
	result, err := svc.GetStream(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestStreamService_GetStream_DefaultLimit(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	req := model.StreamRequest{
		UserID: "user-123",
		Filter: model.FilterAll,
		Limit:  0, // Should default to 20
	}

	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("GetStream", mock.Anything, mock.MatchedBy(func(r model.StreamRequest) bool {
		return r.Limit == 20
	})).Return([]model.PriorityItem{}, (*string)(nil), nil)
	mockCache.On("SetStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	_, err := svc.GetStream(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestStreamService_GetStream_MaxLimit(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	req := model.StreamRequest{
		UserID: "user-123",
		Filter: model.FilterAll,
		Limit:  500, // Should be capped to 100
	}

	mockCache.On("GetStream", mock.Anything, mock.Anything).Return(nil, nil)
	mockRepo.On("GetStream", mock.Anything, mock.MatchedBy(func(r model.StreamRequest) bool {
		return r.Limit == 100
	})).Return([]model.PriorityItem{}, (*string)(nil), nil)
	mockCache.On("SetStream", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Act
	_, err := svc.GetStream(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Tests for GetStreamItemDetails
func TestStreamService_GetStreamItemDetails_CacheHit(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	expectedItem := &model.PriorityItem{
		ID:       "item-1",
		Title:    "Test Item",
		Priority: model.PriorityHigh,
		Messages: []model.Message{{ID: "msg-1", Content: "Hello"}},
	}

	req := model.StreamItemRequest{
		UserID: "user-123",
		ItemID: "item-1",
	}

	// Cache hit
	mockCache.On("GetStreamItem", mock.Anything, "item:item-1").Return(expectedItem, nil)

	// Act
	result, err := svc.GetStreamItemDetails(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "item-1", result.ID)
	assert.Equal(t, 1, len(result.Messages))

	mockRepo.AssertNotCalled(t, "GetStreamItemByID")
}

func TestStreamService_GetStreamItemDetails_CacheMiss(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	expectedItem := &model.PriorityItem{
		ID:       "item-1",
		Title:    "Test Item",
		Priority: model.PriorityHigh,
	}

	req := model.StreamItemRequest{
		UserID: "user-123",
		ItemID: "item-1",
	}

	// Cache miss
	mockCache.On("GetStreamItem", mock.Anything, "item:item-1").Return(nil, nil)
	// Repository returns data
	mockRepo.On("GetStreamItemByID", mock.Anything, "user-123", "item-1").Return(expectedItem, nil)
	// Cache set
	mockCache.On("SetStreamItem", mock.Anything, "item:item-1", expectedItem, mock.Anything).Return(nil)

	// Act
	result, err := svc.GetStreamItemDetails(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "item-1", result.ID)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestStreamService_GetStreamItemDetails_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache)
	svc := newTestService(mockRepo, mockCache)

	req := model.StreamItemRequest{
		UserID: "user-123",
		ItemID: "nonexistent",
	}

	// Cache miss
	mockCache.On("GetStreamItem", mock.Anything, mock.Anything).Return(nil, nil)
	// Repository returns nil (not found)
	mockRepo.On("GetStreamItemByID", mock.Anything, "user-123", "nonexistent").Return(nil, nil)

	// Act
	result, err := svc.GetStreamItemDetails(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// Tests for ValidateFilter
func TestValidateFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.StreamFilter
		hasError bool
	}{
		{"empty string defaults to all", "", model.FilterAll, false},
		{"all filter", "all", model.FilterAll, false},
		{"high filter", "high", model.FilterHigh, false},
		{"unread filter", "unread", model.FilterUnread, false},
		{"invalid filter", "invalid", "", true},
		{"uppercase invalid", "HIGH", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateFilter(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
