# Gravity V2 BFF (Backend-for-Frontend)

A high-performance, read-only Backend-for-Frontend API for the Gravity V2 personal infrastructure interface. Built with Go and Fiber to deliver sub-100ms P95 response times.

## Project Overview

Gravity V2 is a "Personal Infrastructure" interface that unifies multiple communication channels (Email, Slack, WhatsApp, Calendar, YouTube, LinkedIn, Twitter) into a single conversational workspace. This BFF serves preprocessed, read-only data to the frontend, optimized for lightning-fast performance.

**Key Features**:
- Sub-100ms P95 response time
- Unified stream of priority items from multiple sources
- Redis caching layer for sub-millisecond reads
- PostgreSQL persistence layer
- Container-first architecture
- Comprehensive test automation

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Development Environment

```bash
# Clone the repository
git clone <repository-url>
cd gravity-bff

# Copy environment variables
cp .env.example .env

# Start development environment with hot-reload
make dev

# View logs
make logs

# Stop environment
make down
```

The API will be available at `http://localhost:8080`.

### Production Environment

```bash
# Start production environment
make up

# View logs
make logs

# Stop environment
make down
```

### Running Tests

```bash
# Run unit tests
go test ./...

# Run integration tests
make test

# Run performance tests
k6 run --env API_BASE_URL=http://localhost:8080 tests/performance/stream_k6.js
```

## Project Structure

```
gravity-bff/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # API layer (handlers and routing)
│   ├── cache/                   # Redis caching layer
│   ├── config/                  # Configuration management
│   ├── domain/                  # Core domain models and interfaces
│   ├── repository/              # PostgreSQL data access layer
│   └── service/                 # Business logic layer
├── migrations/                  # Database migrations
├── plans/                       # Architecture and planning documentation
│   ├── 01-api-specification.md
│   ├── 02-tech-stack-recommendation.md
│   ├── 03-implementation-plan.md
│   ├── 04-test-automation.md
│   ├── 05-containerization.md
│   └── README.md
├── tests/
│   ├── integration/             # Integration tests
│   └── performance/             # k6 performance tests
├── .dockerignore
├── .env.example
├── docker-compose.yml           # Production configuration
├── docker-compose.dev.yml       # Development configuration
├── docker-compose.test.yml      # Testing configuration
├── Dockerfile
├── Makefile
├── tasks.md                     # Implementation task breakdown
└── README.md                    # This file
```

## Documentation

All architectural and planning documentation is located in the `plans/` directory:

- **[API Specification](plans/01-api-specification.md)**: Complete API contract with endpoints and data models
- **[Tech Stack Recommendation](plans/02-tech-stack-recommendation.md)**: Technology choices and performance analysis
- **[Implementation Plan](plans/03-implementation-plan.md)**: Step-by-step implementation guide
- **[Test Automation](plans/04-test-automation.md)**: Testing strategy and CI/CD pipeline
- **[Containerization](plans/05-containerization.md)**: Docker and deployment configuration

For a complete task breakdown, see **[tasks.md](tasks.md)**.

## API Endpoints

### `GET /v2/stream`
Retrieves the unified stream of priority items.

**Query Parameters**:
- `filter`: `all`, `high`, or `unread` (default: `all`)
- `cursor`: Pagination cursor (optional)
- `limit`: Number of items per page (default: 20, max: 100)

### `GET /v2/stream/{itemId}`
Retrieves full details of a single priority item, including complete message history.

For complete API documentation, see [plans/01-api-specification.md](plans/01-api-specification.md).

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Database Driver**: pgx
- **Redis Client**: go-redis
- **Configuration**: Viper
- **Testing**: Go testing, testify, k6
- **Containerization**: Docker, Docker Compose

## Performance

**Target**: Sub-100ms P95 response time

**Achieved Through**:
- Go's compiled performance (50,000-100,000+ RPS)
- Redis caching layer (sub-millisecond reads)
- Optimized PostgreSQL queries with proper indexes
- Container-based deployment (no cold starts)

## Contributing

1. Review the planning documentation in `plans/`
2. Check `tasks.md` for available tasks
3. Follow the implementation patterns in `plans/03-implementation-plan.md`
4. Write tests following `plans/04-test-automation.md`
5. Ensure all tests pass before submitting a PR

## License

[Your License Here]

## Contact

[Your Contact Information]
