// Package config provides configuration management using Viper.
// It loads configuration from environment variables and .env files.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration values for the application.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Cache    CacheConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds PostgreSQL connection configuration.
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// CacheConfig holds caching configuration.
type CacheConfig struct {
	DefaultTTL time.Duration
	StreamTTL  time.Duration
	ItemTTL    time.Duration
}

// ConnectionString returns the PostgreSQL connection string.
func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// Address returns the Redis address in host:port format.
func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// Load reads configuration from environment variables and .env file.
// Environment variables take precedence over .env file values.
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read from .env file if it exists
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		// It's okay if .env doesn't exist; we'll use environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only log, don't fail - env vars might be set directly
		}
	}

	// Automatically read environment variables
	v.AutomaticEnv()

	// Build the config struct
	cfg := &Config{
		Server: ServerConfig{
			Port:         v.GetInt("API_PORT"),
			ReadTimeout:  v.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  v.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetInt("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSL_MODE"),
			MaxConns:        int32(v.GetInt("DB_MAX_CONNS")),
			MinConns:        int32(v.GetInt("DB_MIN_CONNS")),
			MaxConnLifetime: v.GetDuration("DB_MAX_CONN_LIFETIME"),
			MaxConnIdleTime: v.GetDuration("DB_MAX_CONN_IDLE_TIME"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetInt("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
			PoolSize: v.GetInt("REDIS_POOL_SIZE"),
		},
		Cache: CacheConfig{
			DefaultTTL: v.GetDuration("CACHE_DEFAULT_TTL"),
			StreamTTL:  v.GetDuration("CACHE_STREAM_TTL"),
			ItemTTL:    v.GetDuration("CACHE_ITEM_TTL"),
		},
	}

	return cfg, nil
}

// setDefaults sets sensible default values for all configuration options.
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("API_PORT", 8080)
	v.SetDefault("SERVER_READ_TIMEOUT", "5s")
	v.SetDefault("SERVER_WRITE_TIMEOUT", "10s")
	v.SetDefault("SERVER_IDLE_TIMEOUT", "120s")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "gravity")
	v.SetDefault("DB_PASSWORD", "secret")
	v.SetDefault("DB_NAME", "gravity_db")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_CONNS", 25)
	v.SetDefault("DB_MIN_CONNS", 5)
	v.SetDefault("DB_MAX_CONN_LIFETIME", "1h")
	v.SetDefault("DB_MAX_CONN_IDLE_TIME", "30m")

	// Redis defaults
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_POOL_SIZE", 10)

	// Cache defaults - short TTLs for BFF pattern
	v.SetDefault("CACHE_DEFAULT_TTL", "5m")
	v.SetDefault("CACHE_STREAM_TTL", "2m")
	v.SetDefault("CACHE_ITEM_TTL", "5m")
}
