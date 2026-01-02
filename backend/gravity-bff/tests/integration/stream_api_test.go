//go:build integration

// Package integration contains integration tests for the Gravity BFF API.
// These tests require running PostgreSQL and Redis containers.
//
// Run with: go test -tags=integration ./tests/integration/...
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mabidoli/gravity-bff/internal/api"
	"github.com/mabidoli/gravity-bff/internal/api/handler"
	"github.com/mabidoli/gravity-bff/internal/cache"
	"github.com/mabidoli/gravity-bff/internal/config"
	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/internal/repository"
	"github.com/mabidoli/gravity-bff/internal/service"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

var (
	testDB    *pgxpool.Pool
	testRedis *redis.Client
	testApp   *fiber.App
)

func TestMain(m *testing.M) {
	// Setup
	if err := setup(); err != nil {
		fmt.Printf("Failed to setup integration tests: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown
	teardown()

	os.Exit(code)
}

func setup() error {
	ctx := context.Background()

	// Get database connection from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "test")
	dbPass := getEnv("DB_PASSWORD", "test")
	dbName := getEnv("DB_NAME", "test_db")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	testDB, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := testDB.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Get Redis connection from environment
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6380")

	testRedis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	if err := testRedis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	// Setup test data
	if err := seedTestData(ctx); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	// Initialize application
	testApp = setupTestApp()

	return nil
}

func teardown() {
	ctx := context.Background()

	// Clean up test data
	if testDB != nil {
		cleanupTestData(ctx)
		testDB.Close()
	}

	if testRedis != nil {
		testRedis.FlushDB(ctx)
		testRedis.Close()
	}
}

func setupTestApp() *fiber.App {
	log := logger.New()
	cfg := &config.Config{
		Cache: config.CacheConfig{
			StreamTTL: 2 * time.Minute,
			ItemTTL:   5 * time.Minute,
		},
	}

	// Initialize layers
	redisCache := cache.NewRedisCache(testRedis)
	streamRepo := repository.NewPgStreamRepository(testDB)
	streamService := service.NewStreamService(streamRepo, redisCache, cfg, log)

	healthHandler := handler.NewHealthHandler()
	streamHandler := handler.NewStreamHandler(streamService, log)

	router := api.NewRouter(healthHandler, streamHandler, log)

	app := fiber.New()
	router.Setup(app)

	return app
}

func seedTestData(ctx context.Context) error {
	// Create test user
	_, err := testDB.Exec(ctx, `
		INSERT INTO users (id, name, email) VALUES
		('test-user-1', 'Test User', 'test@example.com')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return err
	}

	// Create test priority items
	_, err = testDB.Exec(ctx, `
		INSERT INTO priority_items (id, user_id, title, source, priority, is_unread, snippet, item_timestamp) VALUES
		('item-1', 'test-user-1', 'High Priority Email', 'email', 'high', true, 'Important message...', NOW()),
		('item-2', 'test-user-1', 'Medium Priority Task', 'task', 'medium', false, 'Task description...', NOW() - INTERVAL '1 hour'),
		('item-3', 'test-user-1', 'Low Priority Update', 'slack', 'low', true, 'Update message...', NOW() - INTERVAL '2 hours')
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		return err
	}

	// Add participants
	_, err = testDB.Exec(ctx, `
		INSERT INTO priority_item_participants (item_id, user_id) VALUES
		('item-1', 'test-user-1')
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return err
	}

	// Add messages
	_, err = testDB.Exec(ctx, `
		INSERT INTO messages (id, item_id, sender_id, sender_type, content_type, content, message_timestamp) VALUES
		('msg-1', 'item-1', 'test-user-1', 'other', 'text', 'Hello, this is a test message.', NOW())
		ON CONFLICT (id) DO NOTHING
	`)

	return err
}

func cleanupTestData(ctx context.Context) {
	testDB.Exec(ctx, "DELETE FROM messages WHERE id LIKE 'msg-%'")
	testDB.Exec(ctx, "DELETE FROM priority_item_participants WHERE item_id LIKE 'item-%'")
	testDB.Exec(ctx, "DELETE FROM priority_items WHERE id LIKE 'item-%'")
	testDB.Exec(ctx, "DELETE FROM users WHERE id LIKE 'test-user-%'")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Integration Tests

func TestHealthEndpoint_Integration(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := testApp.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.HealthResponse
	err = json.Unmarshal(body, &result)

	assert.NoError(t, err)
	assert.Equal(t, "ok", result.Status)
}

func TestGetStream_Integration(t *testing.T) {
	// Clear Redis cache first
	testRedis.FlushDB(context.Background())

	req := httptest.NewRequest("GET", "/v2/stream", nil)
	req.Header.Set("Authorization", "Bearer test-user-1")

	resp, err := testApp.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.StreamResponse
	err = json.Unmarshal(body, &result)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result.Data), 1)
}

func TestGetStream_FilterHigh_Integration(t *testing.T) {
	testRedis.FlushDB(context.Background())

	req := httptest.NewRequest("GET", "/v2/stream?filter=high", nil)
	req.Header.Set("Authorization", "Bearer test-user-1")

	resp, err := testApp.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.StreamResponse
	json.Unmarshal(body, &result)

	// All returned items should be high priority
	for _, item := range result.Data {
		assert.Equal(t, model.PriorityHigh, item.Priority)
	}
}

func TestGetStream_FilterUnread_Integration(t *testing.T) {
	testRedis.FlushDB(context.Background())

	req := httptest.NewRequest("GET", "/v2/stream?filter=unread", nil)
	req.Header.Set("Authorization", "Bearer test-user-1")

	resp, err := testApp.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.StreamResponse
	json.Unmarshal(body, &result)

	// All returned items should be unread
	for _, item := range result.Data {
		assert.True(t, item.IsUnread)
	}
}

func TestGetStreamItem_Integration(t *testing.T) {
	testRedis.FlushDB(context.Background())

	req := httptest.NewRequest("GET", "/v2/stream/item-1", nil)
	req.Header.Set("Authorization", "Bearer test-user-1")

	resp, err := testApp.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result model.PriorityItem
	err = json.Unmarshal(body, &result)

	assert.NoError(t, err)
	assert.Equal(t, "item-1", result.ID)
	assert.Equal(t, "High Priority Email", result.Title)
}

func TestGetStreamItem_NotFound_Integration(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/stream/nonexistent-item", nil)
	req.Header.Set("Authorization", "Bearer test-user-1")

	resp, err := testApp.Test(req, -1)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGetStream_CacheHit_Integration(t *testing.T) {
	testRedis.FlushDB(context.Background())

	// First request - cache miss
	req1 := httptest.NewRequest("GET", "/v2/stream", nil)
	req1.Header.Set("Authorization", "Bearer test-user-1")
	resp1, _ := testApp.Test(req1, -1)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Second request - should be cache hit (faster)
	start := time.Now()
	req2 := httptest.NewRequest("GET", "/v2/stream", nil)
	req2.Header.Set("Authorization", "Bearer test-user-1")
	resp2, _ := testApp.Test(req2, -1)
	duration := time.Since(start)

	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
	// Cache hit should be fast (under 50ms typically)
	assert.Less(t, duration, 100*time.Millisecond)
}
