// Package cache provides caching functionality using Redis.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

// Cache defines the caching interface.
type Cache interface {
	// GetStream retrieves cached stream data.
	GetStream(ctx context.Context, key string) (*model.StreamResponse, error)
	// SetStream caches stream data with TTL.
	SetStream(ctx context.Context, key string, data *model.StreamResponse, ttl time.Duration) error
	// GetStreamItem retrieves a cached stream item.
	GetStreamItem(ctx context.Context, key string) (*model.PriorityItem, error)
	// SetStreamItem caches a stream item with TTL.
	SetStreamItem(ctx context.Context, key string, item *model.PriorityItem, ttl time.Duration) error
	// Delete removes a key from cache.
	Delete(ctx context.Context, keys ...string) error
	// Ping checks if Redis is reachable.
	Ping(ctx context.Context) error
}

// RedisCache implements Cache using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Key prefixes for different cache types.
const (
	streamKeyPrefix = "stream:"
	itemKeyPrefix   = "item:"
)

// StreamKey generates a cache key for stream data.
func StreamKey(userID string, filter model.StreamFilter, cursor *string) string {
	cursorPart := "none"
	if cursor != nil && *cursor != "" {
		cursorPart = *cursor
	}
	return fmt.Sprintf("%s%s:%s:%s", streamKeyPrefix, userID, filter, cursorPart)
}

// ItemKey generates a cache key for a stream item.
func ItemKey(itemID string) string {
	return fmt.Sprintf("%s%s", itemKeyPrefix, itemID)
}

// GetStream retrieves cached stream data.
func (c *RedisCache) GetStream(ctx context.Context, key string) (*model.StreamResponse, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get stream from cache: %w", err)
	}

	var response model.StreamResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stream data: %w", err)
	}

	return &response, nil
}

// SetStream caches stream data with TTL.
func (c *RedisCache) SetStream(ctx context.Context, key string, data *model.StreamResponse, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal stream data: %w", err)
	}

	if err := c.client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set stream in cache: %w", err)
	}

	return nil
}

// GetStreamItem retrieves a cached stream item.
func (c *RedisCache) GetStreamItem(ctx context.Context, key string) (*model.PriorityItem, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get item from cache: %w", err)
	}

	var item model.PriorityItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item data: %w", err)
	}

	return &item, nil
}

// SetStreamItem caches a stream item with TTL.
func (c *RedisCache) SetStreamItem(ctx context.Context, key string, item *model.PriorityItem, ttl time.Duration) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item data: %w", err)
	}

	if err := c.client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set item in cache: %w", err)
	}

	return nil
}

// Delete removes keys from cache.
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys from cache: %w", err)
	}

	return nil
}

// Ping checks if Redis is reachable.
func (c *RedisCache) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

// InvalidateUserCache invalidates all cached data for a user.
func (c *RedisCache) InvalidateUserCache(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s%s:*", streamKeyPrefix, userID)

	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		return c.Delete(ctx, keys...)
	}

	return nil
}
