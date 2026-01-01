# Gravity V2 BFF: Container-First Docker System

## 1. Introduction

This document outlines a complete, container-first architecture for the Gravity V2 BFF application using Docker and Docker Compose. This setup provides a consistent, reproducible, and isolated environment for development, testing, and production, ensuring that the application runs the same way everywhere.

**Core Principles**:

-   **Container-First**: Every component of the system, including the API, database, and cache, runs in a dedicated container.
-   **Environment Parity**: Development, testing, and production environments are as similar as possible.
-   **Declarative**: Docker Compose files declaratively define the entire application stack.
-   **Scalability**: The architecture is designed to be easily scalable in a container orchestration platform like Kubernetes or AWS ECS.

### System Components

The system consists of the following services, all managed by Docker Compose:

1.  **`api`**: The Go/Fiber BFF application.
2.  **`db`**: The PostgreSQL database for data persistence.
3.  **`cache`**: The Redis instance for high-speed caching.
4.  **`migrate`**: A one-off container to run database migrations on startup.

---

## 2. Project Structure

The following directory structure should be used to organize the project:

```
gravity-bff/
├── cmd/api/main.go
├── internal/
│   └── ... (application code)
├── migrations/
│   └── 0001_initial_schema.up.sql # Database migration file
├── .dockerignore
├── .env.example
├── docker-compose.yml             # Production environment
├── docker-compose.dev.yml         # Development environment override
├── docker-compose.test.yml        # Testing environment override
├── Dockerfile                     # API container build instructions
└── Makefile
```

---

## 3. Dockerfiles & Configuration

### 3.1. API Dockerfile (`Dockerfile`)

A multi-stage `Dockerfile` is used to create a small, optimized, and secure production image for the Go application.

```dockerfile
# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application into a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/gravity-bff ./cmd/api/main.go

# ---- Production Stage ----
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gravity-bff .

# Copy environment file (optional, for reference)
COPY .env.example .

# Expose the application port
EXPOSE 8080

# Set the entrypoint
CMD ["./gravity-bff"]
```

### 3.2. Docker Ignore (`.dockerignore`)

To keep the build context small and fast, we exclude unnecessary files.

```
.git
.vscode
*.md
bin/
docker-compose*.yml
```

### 3.3. Environment Variables (`.env.example`)

This file serves as a template for the required environment variables.

```ini
# Application Configuration
API_PORT=8080

# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_USER=gravity
DB_PASSWORD=secret
DB_NAME=gravity_db
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=cache
REDIS_PORT=6379
REDIS_PASSWORD=
```

---

## 4. Docker Compose Configurations

We use multiple Compose files to tailor the environment for different use cases.

### 4.1. Production (`docker-compose.yml`)

This is the base file for production and shared services.

```yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gravity_api
    ports:
      - "${API_PORT}:${API_PORT}"
    env_file:
      - .env
    depends_on:
      db: 
        condition: service_healthy
      cache: 
        condition: service_healthy
    networks:
      - gravity_net
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 5

  db:
    image: postgres:15-alpine
    container_name: gravity_db
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    ports:
      - "${DB_PORT}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - gravity_net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 10s
      timeout: 5s
      retries: 5

  cache:
    image: redis:7-alpine
    container_name: gravity_cache
    ports:
      - "${REDIS_PORT}:6379"
    networks:
      - gravity_net
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  migrate:
    image: migrate/migrate
    container_name: gravity_migrate
    command: ["-path", "/migrations", "-database", "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}", "up"]
    volumes:
      - ./migrations:/migrations
    depends_on:
      db: 
        condition: service_healthy
    networks:
      - gravity_net
    restart: on-failure

networks:
  gravity_net:
    driver: bridge

volumes:
  postgres_data:
```

### 4.2. Development (`docker-compose.dev.yml`)

This file overrides the base configuration for development, adding hot-reloading.

```yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
      target: builder # Use the build stage with source code
    container_name: gravity_api_dev
    command: go run ./cmd/api/main.go
    volumes:
      - .:/app # Mount source code for hot-reloading
    ports:
      - "${API_PORT}:${API_PORT}"
    env_file:
      - .env
    depends_on:
      - db
      - cache
    networks:
      - gravity_net
```

### 4.3. Testing (`docker-compose.test.yml`)

This file is used for running integration tests in a clean, isolated environment.

```yaml
version: '3.8'

services:
  db_test:
    image: postgres:15-alpine
    container_name: gravity_db_test
    environment:
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
      POSTGRES_DB: test_db
    ports:
      - "5433:5432"
    networks:
      - gravity_net

  cache_test:
    image: redis:7-alpine
    container_name: gravity_cache_test
    ports:
      - "6380:6379"
    networks:
      - gravity_net
```

---

## 5. Database Migrations

We use `migrate/migrate` to manage database schema changes.

**File**: `migrations/0001_initial_schema.up.sql`

```sql
-- Your full SQL schema from the implementation plan goes here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ... (rest of the schema)
```

---

## 6. Usage

A `Makefile` simplifies common commands.

**File**: `Makefile`

```makefile
.PHONY: up down logs ps dev test

# Production environment
up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

ps:
	docker-compose ps

# Development environment
dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# Testing environment
test-up:
	docker-compose -f docker-compose.test.yml up -d

test-down:
	docker-compose -f docker-compose.test.yml down

test:
	@make test-up
	go test -v -tags=integration ./...
	@make test-down
```

### Commands

-   **Start Production Environment**: `make up`
-   **Start Development Environment**: `make dev`
-   **Stop Environment**: `make down`
-   **View Logs**: `make logs`
-   **Run Integration Tests**: `make test`

This container-first setup provides a powerful and flexible foundation for the Gravity V2 BFF, ensuring consistency and reliability across all stages of the application lifecycle.
