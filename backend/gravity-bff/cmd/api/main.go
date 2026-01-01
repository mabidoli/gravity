// Package main is the entry point for the Gravity BFF API server.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mabidoli/gravity-bff/internal/config"
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
	log.Info("Server port: %d", cfg.Server.Port)
	log.Info("Database host: %s", cfg.Database.Host)
	log.Info("Redis host: %s", cfg.Redis.Host)

	// TODO: Initialize database connection
	// TODO: Initialize Redis connection
	// TODO: Initialize repositories
	// TODO: Initialize services
	// TODO: Initialize handlers
	// TODO: Setup Fiber app and routes
	// TODO: Start server

	// Placeholder: Print startup message
	fmt.Printf("Gravity BFF API ready to start on port %d\n", cfg.Server.Port)
	fmt.Println("Phase 1 complete - waiting for Phase 2 & 3 implementation...")

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Gravity BFF API...")
}
