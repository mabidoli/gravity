# Gravity V2 BFF: Go & Fiber Implementation Plan

## 1. Introduction

This document provides a comprehensive, step-by-step implementation plan for building the Gravity V2 Backend-for-Frontend (BFF) API using Go and the Fiber web framework. The plan is designed to meet the primary requirement of a **sub-100ms P95 response time** by leveraging Go's performance, a clean architecture, and a robust caching strategy.

This guide follows professional Go development patterns, emphasizing modularity, clear separation of concerns, and testability.

### Technology Stack

-   **Language**: Go (version 1.21+)
-   **Web Framework**: Fiber (v2)
-   **Database (Persistence)**: PostgreSQL
-   **Database (Caching)**: Redis
-   **Database Driver**: pgx (for PostgreSQL)
-   **Redis Client**: go-redis
-   **Configuration**: Viper
-   **Containerization**: Docker

---

## 2. Project Structure

A standard, layered architecture will be used to ensure separation of concerns and maintainability. This structure separates API handling, business logic, and data access.

```
gravity-bff/
├── cmd/
│   └── api/
│       └── main.go         # Application entry point
├── internal/
│   ├── api/                # API layer (handlers and routing)
│   │   ├── handler/
│   │   │   └── stream.go   # HTTP handlers for the /stream endpoint
│   │   └── router.go       # API route definitions
│   ├── cache/              # Caching layer (Redis implementation)
│   │   └── redis_cache.go
│   ├── config/             # Configuration management (Viper)
│   │   └── config.go
│   ├── domain/             # Core domain models and interfaces
│   │   ├── model/
│   │   │   └── stream.go   # Go structs for data models
│   │   └── repository/
│   │       └── stream.go   # Interfaces for the repository layer
│   ├── repository/         # Data access layer (PostgreSQL)
│   │   └── pg_stream.go
│   └── service/            # Business logic layer
│       └── stream.go
├── pkg/                    # Shared utility packages (e.g., logger)
│   └── logger/
│       └── logger.go
├── .env                    # Environment variables
├── Dockerfile              # Container build instructions
├── go.mod                  # Go module dependencies
├── go.sum
└── Makefile                # Helper commands for development
```

---

## 3. Database Schema (PostgreSQL)

The following SQL schema should be applied to the PostgreSQL database. It is designed to be normalized and efficient for the read-heavy patterns of the BFF.

```sql
-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table to store participant information
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Priority items, the core of the unified stream
CREATE TABLE priority_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- The owner of the stream
    title VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL, -- e.g., 'email', 'calendar'
    priority VARCHAR(50) NOT NULL, -- e.g., 'high', 'medium', 'low'
    is_unread BOOLEAN NOT NULL DEFAULT TRUE,
    snippet TEXT,
    item_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Junction table for participants in a priority item
CREATE TABLE priority_item_participants (
    item_id UUID NOT NULL REFERENCES priority_items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, user_id)
);

-- Messages within each priority item's conversation
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES priority_items(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id), -- Can be null for system messages
    sender_type VARCHAR(50) NOT NULL, -- e.g., 'user', 'other', 'system'
    content_type VARCHAR(50) NOT NULL, -- e.g., 'text', 'event', 'social'
    content TEXT,
    full_content_html TEXT,
    message_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Optional: Separate tables for complex message types if needed
-- e.g., calendar_events, social_content, attachments, ai_insights
-- For simplicity, these can be stored as JSONB in the messages table initially.
ALTER TABLE messages ADD COLUMN event_details JSONB;
ALTER TABLE messages ADD COLUMN social_details JSONB;
ALTER TABLE messages ADD COLUMN attachments JSONB; -- Array of attachment objects
ALTER TABLE messages ADD COLUMN ai_insights JSONB; -- Array of insight objects

-- Indexes for performance
CREATE INDEX idx_priority_items_user_id_timestamp ON priority_items (user_id, item_timestamp DESC);
CREATE INDEX idx_priority_items_user_id_priority ON priority_items (user_id, priority);
CREATE INDEX idx_messages_item_id_timestamp ON messages (item_id, message_timestamp DESC);
```

---

## 4. Implementation Steps

This section provides a phased, step-by-step guide to building the application.

### Phase 1: Project Setup & Configuration

1.  **Initialize Go Module**:
    ```bash
    go mod init github.com/your-username/gravity-bff
    ```

2.  **Install Dependencies**:
    ```bash
    go get github.com/gofiber/fiber/v2
    go get github.com/jackc/pgx/v5/pgxpool
    go get github.com/redis/go-redis/v9
    go get github.com/spf13/viper
    ```

3.  **Implement Configuration (`internal/config/config.go`)**:
    -   Use Viper to load environment variables from a `.env` file and system variables.
    -   Define a `Config` struct to hold all configuration values (DB connection strings, Redis address, API port, etc.).

4.  **Setup Main Application (`cmd/api/main.go`)**:
    -   Load configuration.
    -   Initialize database and Redis connections.
    -   Initialize Fiber app.
    -   Instantiate repositories, services, and handlers (dependency injection).
    -   Setup API routes.
    -   Start the server.

### Phase 2: Domain and Data Layers

1.  **Define Data Models (`internal/domain/model/stream.go`)**:
    -   Create Go `structs` that map to the API specification and database tables.
    -   Use `json` and `db` struct tags for serialization and database mapping.

    ```go
    // Example PriorityItem struct
    type PriorityItem struct {
        ID         string    `json:"id" db:"id"`
        Title      string    `json:"title" db:"title"`
        Source     string    `json:"source" db:"source"`
        Priority   string    `json:"priority" db:"priority"`
        IsUnread   bool      `json:"isUnread" db:"is_unread"`
        Snippet    string    `json:"snippet" db:"snippet"`
        Timestamp  time.Time `json:"timestamp" db:"item_timestamp"`
        Participants []User  `json:"participants"`
    }
    ```

2.  **Define Repository Interfaces (`internal/domain/repository/stream.go`)**:
    -   Define the contracts for data access operations (e.g., `GetStream`, `GetStreamItemByID`). This decouples the business logic from the database implementation.

3.  **Implement PostgreSQL Repository (`internal/repository/pg_stream.go`)**:
    -   Implement the repository interfaces using `pgx` to query the PostgreSQL database.
    -   Write SQL queries to fetch the required data.

### Phase 3: Caching Layer

1.  **Implement Redis Cache (`internal/cache/redis_cache.go`)**:
    -   Create a struct that wraps the `go-redis` client.
    -   Implement methods for `Get` and `Set` operations.
    -   Handle JSON serialization/deserialization for storing structs in Redis.

2.  **Caching Strategy**:
    -   **`GET /stream`**: Cache the list of `PriorityItem` IDs for each filter combination (e.g., `user:123:filter:all:cursor:xyz`). Cache individual `PriorityItem` objects by their ID. The service layer will fetch the list of IDs from the cache, then fetch each item individually (checking the cache first for each).
    -   **`GET /stream/{itemId}`**: Cache the full `PriorityItem` object with its `messages` array using a key like `item:abc-123`.
    -   **Cache Invalidation**: Since this is a read-only BFF, cache invalidation will be handled by the upstream data ingestion service. This can be done by publishing invalidation events (e.g., via a message queue) or simply by using a Time-To-Live (TTL) on cache keys (e.g., 5 minutes).

### Phase 4: Business Logic and API Layers

1.  **Implement Service Layer (`internal/service/stream.go`)**:
    -   This layer contains the core business logic.
    -   The `StreamService` will orchestrate calls to the repository and the cache.
    -   **Logic for `GetStream`**: 
        1.  Generate a cache key based on user ID and filter parameters.
        2.  Try to fetch the result from Redis.
        3.  If cache miss, fetch data from the PostgreSQL repository.
        4.  Store the result in Redis with a TTL.
        5.  Return the data.

2.  **Implement API Handlers (`internal/api/handler/stream.go`)**:
    -   These are the Fiber handlers that receive HTTP requests.
    -   They parse request parameters (query params, path params).
    -   Call the appropriate service layer methods.
    -   Format the service response into a JSON response and send it to the client.

3.  **Define API Routes (`internal/api/router.go`)**:
    -   Map the API endpoints to their respective handlers.
    -   Group routes under `/v2` and apply any middleware (e.g., logging, auth).

    ```go
    func SetupRoutes(app *fiber.App, streamHandler *handler.StreamHandler) {
        v2 := app.Group("/v2")
        stream := v2.Group("/stream")
        stream.Get("/", streamHandler.GetStream)
        stream.Get("/:itemId", streamHandler.GetStreamItemDetails)
    }
    ```

### Phase 5: Deployment

1.  **Create `Dockerfile`**:
    -   Use a multi-stage build to create a small, optimized production image.
    -   Stage 1: Build the Go binary using the official `golang` image.
    -   Stage 2: Copy the compiled binary into a minimal `distroless` or `alpine` image.

    ```dockerfile
    # Build stage
    FROM golang:1.21-alpine AS builder
    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/gravity-bff ./cmd/api/main.go

    # Final stage
    FROM alpine:latest
    WORKDIR /app
    COPY --from=builder /app/gravity-bff .
    COPY .env .
    EXPOSE 8080
    CMD ["./gravity-bff"]
    ```

2.  **Create `Makefile`**:
    -   Add commands for common tasks like `run`, `build`, `test`, and `docker-build`.

    ```makefile
    .PHONY: run build test docker-build

    run:
    	go run ./cmd/api/main.go

    build:
    	go build -o bin/gravity-bff ./cmd/api/main.go

    test:
    	go test ./...

    docker-build:
    	docker build -t gravity-bff:latest .
    ```

---

## 5. Conclusion

This implementation plan provides a robust foundation for building a high-performance, scalable, and maintainable BFF API for Gravity V2. By following the layered architecture, leveraging Go and Fiber's performance, and implementing a smart caching strategy with Redis, the application will be well-equipped to meet the stringent sub-100ms response time requirement.
