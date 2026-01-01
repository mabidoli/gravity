# Gravity BFF Makefile
# Provides helper commands for development, testing, and production

.PHONY: all build run test test-unit test-integration test-coverage \
        lint fmt vet clean docker-build docker-up docker-down docker-logs \
        dev dev-down migrate migrate-down help

# Default target
all: lint test build

# ============================================================================
# Development
# ============================================================================

## build: Build the application binary
build:
	@echo "Building..."
	@go build -ldflags="-w -s" -o bin/gravity-bff ./cmd/api/main.go

## run: Run the application locally
run:
	@go run ./cmd/api/main.go

## fmt: Format the code
fmt:
	@echo "Formatting..."
	@gofmt -s -w .

## vet: Run go vet
vet:
	@echo "Vetting..."
	@go vet ./...

## lint: Run golangci-lint
lint:
	@echo "Linting..."
	@golangci-lint run ./... || echo "golangci-lint not installed, skipping..."

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# ============================================================================
# Testing
# ============================================================================

## test: Run all tests
test: test-unit

## test-unit: Run unit tests
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -short ./...

## test-integration: Run integration tests (requires test containers)
test-integration:
	@echo "Starting test containers..."
	@docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for containers to be healthy..."
	@sleep 5
	@echo "Running integration tests..."
	@DB_HOST=localhost DB_PORT=5433 DB_USER=test DB_PASSWORD=test DB_NAME=test_db \
		REDIS_HOST=localhost REDIS_PORT=6380 \
		go test -v -tags=integration ./...
	@echo "Stopping test containers..."
	@docker-compose -f docker-compose.test.yml down

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# ============================================================================
# Docker
# ============================================================================

## docker-build: Build the Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t gravity-bff:latest .

## docker-up: Start all containers (production)
docker-up:
	@echo "Starting production environment..."
	@docker-compose up -d

## docker-down: Stop all containers
docker-down:
	@echo "Stopping containers..."
	@docker-compose down

## docker-logs: View container logs
docker-logs:
	@docker-compose logs -f

## docker-ps: Show running containers
docker-ps:
	@docker-compose ps

## dev: Start development environment with hot-reload
dev:
	@echo "Starting development environment..."
	@docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

## dev-down: Stop development environment
dev-down:
	@docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

# ============================================================================
# Database
# ============================================================================

## migrate: Run database migrations
migrate:
	@echo "Running migrations..."
	@docker-compose run --rm migrate

## migrate-down: Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@docker run --rm -v $(PWD)/migrations:/migrations --network gravity-bff_gravity_net \
		migrate/migrate:v4.17.0 \
		-path /migrations \
		-database "postgres://gravity:secret@db:5432/gravity_db?sslmode=disable" \
		down 1

## migrate-create: Create a new migration (usage: make migrate-create name=migration_name)
migrate-create:
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir migrations -seq $(name)

# ============================================================================
# Utilities
# ============================================================================

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

## generate: Run go generate
generate:
	@echo "Running go generate..."
	@go generate ./...

## help: Show this help message
help:
	@echo "Gravity BFF - Available Commands:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
