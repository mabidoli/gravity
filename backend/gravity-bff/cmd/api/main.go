// Package main is the entry point for the Gravity BFF API server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/mabidoli/gravity-bff/internal/api"
	"github.com/mabidoli/gravity-bff/internal/api/handler"
	"github.com/mabidoli/gravity-bff/internal/api/middleware"
	"github.com/mabidoli/gravity-bff/internal/cache"
	"github.com/mabidoli/gravity-bff/internal/config"
	"github.com/mabidoli/gravity-bff/internal/repository"
	"github.com/mabidoli/gravity-bff/internal/service"
	"github.com/mabidoli/gravity-bff/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.New()
	log.Info("Starting Gravity BFF API...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	log.Info("Configuration loaded successfully")

	// Initialize Clerk authentication
	if err := middleware.InitClerk(); err != nil {
		log.Fatal("Failed to initialize Clerk: %v", err)
	}
	log.Info("Clerk authentication initialized")

	// Initialize database connection
	db, err := initDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Info("Database connection established")

	// Initialize Redis connection
	redisClient, err := initRedis(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize Redis: %v", err)
	}
	defer redisClient.Close()
	log.Info("Redis connection established")

	// Initialize cache
	redisCache := cache.NewRedisCache(redisClient)

	// Initialize repositories
	streamRepo := repository.NewPgStreamRepository(db)

	// Initialize services
	streamService := service.NewStreamService(streamRepo, redisCache, cfg, log)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	streamHandler := handler.NewStreamHandler(streamService, log)

	// Initialize router
	router := api.NewRouter(healthHandler, streamHandler, log)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Gravity BFF",
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           cfg.Server.IdleTimeout,
		DisableStartupMessage: true,
	})

	// Setup routes
	router.Setup(app)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("Server listening on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatal("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Gravity BFF API...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}

	log.Info("Server shutdown complete")
}

// initDatabase initializes the PostgreSQL connection pool.
func initDatabase(cfg *config.Config, log *logger.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.Database.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = cfg.Database.MaxConns
	poolConfig.MinConns = cfg.Database.MinConns
	poolConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database pool created: max=%d, min=%d", cfg.Database.MaxConns, cfg.Database.MinConns)
	return pool, nil
}

// initRedis initializes the Redis client.
func initRedis(cfg *config.Config, log *logger.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Info("Redis connected: %s, pool=%d", cfg.Redis.Address(), cfg.Redis.PoolSize)
	return client, nil
}
