// Package service implements the business logic layer.
package service

import (
	"context"
	"fmt"

	"github.com/mabidoli/gravity-bff/internal/cache"
	"github.com/mabidoli/gravity-bff/internal/config"
	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/internal/domain/repository"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

// StreamService provides business logic for stream operations.
type StreamService struct {
	repo   repository.StreamRepository
	cache  cache.Cache
	config *config.Config
	log    *logger.Logger
}

// NewStreamService creates a new stream service.
func NewStreamService(
	repo repository.StreamRepository,
	cache cache.Cache,
	cfg *config.Config,
	log *logger.Logger,
) *StreamService {
	return &StreamService{
		repo:   repo,
		cache:  cache,
		config: cfg,
		log:    log,
	}
}

// GetStream retrieves the priority stream for a user with caching.
func (s *StreamService) GetStream(ctx context.Context, req model.StreamRequest) (*model.StreamResponse, error) {
	// Validate and set defaults
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Filter == "" {
		req.Filter = model.FilterAll
	}

	// Generate cache key
	cacheKey := cache.StreamKey(req.UserID, req.Filter, req.Cursor)

	// Try to get from cache first
	cachedResponse, err := s.cache.GetStream(ctx, cacheKey)
	if err != nil {
		s.log.Warn("Cache get error: %v", err)
		// Continue without cache
	} else if cachedResponse != nil {
		s.log.Debug("Cache hit for stream: %s", cacheKey)
		return cachedResponse, nil
	}

	s.log.Debug("Cache miss for stream: %s", cacheKey)

	// Fetch from repository
	items, nextCursor, err := s.repo.GetStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream from repository: %w", err)
	}

	response := &model.StreamResponse{
		Data:       items,
		NextCursor: nextCursor,
	}

	// Cache the response
	if err := s.cache.SetStream(ctx, cacheKey, response, s.config.Cache.StreamTTL); err != nil {
		s.log.Warn("Failed to cache stream: %v", err)
		// Continue without caching
	}

	return response, nil
}

// GetStreamItemDetails retrieves full details of a stream item with caching.
func (s *StreamService) GetStreamItemDetails(ctx context.Context, req model.StreamItemRequest) (*model.PriorityItem, error) {
	// Generate cache key
	cacheKey := cache.ItemKey(req.ItemID)

	// Try to get from cache first
	cachedItem, err := s.cache.GetStreamItem(ctx, cacheKey)
	if err != nil {
		s.log.Warn("Cache get error: %v", err)
		// Continue without cache
	} else if cachedItem != nil {
		s.log.Debug("Cache hit for item: %s", cacheKey)
		return cachedItem, nil
	}

	s.log.Debug("Cache miss for item: %s", cacheKey)

	// Fetch from repository
	item, err := s.repo.GetStreamItemByID(ctx, req.UserID, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from repository: %w", err)
	}

	if item == nil {
		return nil, nil // Not found
	}

	// Cache the response
	if err := s.cache.SetStreamItem(ctx, cacheKey, item, s.config.Cache.ItemTTL); err != nil {
		s.log.Warn("Failed to cache item: %v", err)
		// Continue without caching
	}

	return item, nil
}

// ValidateFilter validates the filter parameter.
func ValidateFilter(filter string) (model.StreamFilter, error) {
	switch filter {
	case "", "all":
		return model.FilterAll, nil
	case "high":
		return model.FilterHigh, nil
	case "unread":
		return model.FilterUnread, nil
	default:
		return "", fmt.Errorf("invalid filter: %s. Valid values: all, high, unread", filter)
	}
}
