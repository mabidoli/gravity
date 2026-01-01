# Gravity V2 BFF: Test Automation Strategy

## 1. Introduction & Philosophy

This document outlines the comprehensive test automation strategy for the Gravity V2 BFF Go project. The primary goal is to ensure the application is reliable, performant, and maintainable by embedding quality into every stage of the development lifecycle. Our philosophy is based on the **Testing Pyramid**, emphasizing a strong foundation of fast, isolated unit tests, supported by more comprehensive integration tests, and validated by targeted end-to-end and performance tests.

**Key Objectives**:

-   **Confidence**: Every change pushed to production should be backed by a suite of passing tests.
-   **Performance Validation**: Continuously verify that the API meets the **sub-100ms P95 response time** requirement.
-   **Fast Feedback**: Developers should be able to run tests quickly and get immediate feedback.
-   **Automation**: All tests will be automated and integrated into a CI/CD pipeline to act as a quality gate.

### Testing Tools

-   **Unit & Integration Testing**: Go's built-in `testing` package.
-   **Assertions**: `testify/assert` and `testify/require` for readable assertions.
-   **Mocking**: `stretchr/testify/mock` for mocking dependencies in unit tests.
-   **HTTP Testing**: Go's built-in `net/http/httptest` for API handler tests.
-   **Performance Testing**: `k6` (by Grafana) for load and stress testing.
-   **Integration Test Environment**: Docker Compose to orchestrate the API, PostgreSQL, and Redis containers.

---

## 2. Test Structure

Tests will be co-located with the code they are testing, following Go conventions. This makes them easy to discover and maintain.

-   Test files will be named `_test.go` (e.g., `stream_test.go`).
-   Tests will be placed in the same package as the code under test.
-   Integration tests will be distinguished by a build tag `//go:build integration`.

**Directory Structure Example**:

```
gravity-bff/
├── internal/
│   ├── service/
│   │   ├── stream.go
│   │   └── stream_test.go      # Unit tests for StreamService
│   └── api/
│       └── handler/
│           ├── stream.go
│           └── stream_test.go      # Unit tests for StreamHandler
├── tests/                    # Integration & Performance tests
│   ├── integration/
│   │   └── stream_api_test.go  # Integration tests for the Stream API
│   └── performance/
│       └── stream_k6.js        # k6 script for performance testing
└── ...
```

---

## 3. Unit Tests

Unit tests form the base of the pyramid. They test a single unit (a function or a method) in isolation. Dependencies are mocked to ensure the test is fast and not reliant on external systems.

**Scope**: `service`, `handler`, `cache`, and `repository` layers.

### Example: Testing a Service Method

We will test the `StreamService` by mocking the repository and cache interfaces.

**File**: `internal/service/stream_test.go`

```go
package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-username/gravity-bff/internal/domain/model"
)

// Mock Repository
type MockStreamRepository struct {
	mock.Mock
}

func (m *MockStreamRepository) GetStreamItemByID(ctx context.Context, itemID string) (*model.PriorityItem, error) {
	args := m.Called(ctx, itemID)
	return args.Get(0).(*model.PriorityItem), args.Error(1)
}

// ... other mocked methods

func TestStreamService_GetStreamItemDetails(t *testing.T) {
	// Arrange
	mockRepo := new(MockStreamRepository)
	mockCache := new(MockCache) // Assuming a similar mock for the cache

	service := NewStreamService(mockRepo, mockCache)

	expectedItem := &model.PriorityItem{ID: "item-1", Title: "Test Item"}

	// Setup expectations
	mockCache.On("Get", mock.Anything, "item:item-1").Return(nil, redis.Nil) // Cache miss
	mockRepo.On("GetStreamItemByID", mock.Anything, "item-1").Return(expectedItem, nil)
	mockCache.On("Set", mock.Anything, "item:item-1", expectedItem, mock.Anything).Return(nil)

	// Act
	item, err := service.GetStreamItemDetails(context.Background(), "item-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "item-1", item.ID)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
```

---

## 4. Integration Tests

Integration tests verify the interaction between different components of the system. They are slower than unit tests but provide higher confidence that the system works as a whole.

**Scope**: Testing the API endpoints against real (but containerized) database and cache instances.

### Setup

1.  **Docker Compose**: A `docker-compose.test.yml` file will define the services needed for testing: `api`, `postgres-test`, and `redis-test`.
2.  **Test Main**: A `TestMain` function will be used to set up the test environment (start containers, run migrations) before any tests are run, and tear it down afterward.
3.  **Build Tag**: Integration tests will use the `//go:build integration` tag to separate them from unit tests.

**File**: `tests/integration/stream_api_test.go`

```go
//go:build integration

package integration

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/your-username/gravity-bff/internal/domain/model"
)

// TestMain would handle docker-compose up/down

func TestGetStreamItemDetails_Integration(t *testing.T) {
	// Arrange: This would involve setting up the app and seeding the test DB
	app := setupTestApp() // This function initializes the app with test DB connections

	// Seed the database with a test item
	// seedDB(db, &model.PriorityItem{ID: "seeded-item-1", ...})

	// Act
	req := httptest.NewRequest("GET", "/v2/stream/seeded-item-1", nil)
	resp, err := app.Test(req, -1) // -1 disables timeout

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var item model.PriorityItem
	json.NewDecoder(resp.Body).Decode(&item)

	assert.Equal(t, "seeded-item-1", item.ID)
}
```

### Running Tests

-   **Run Unit Tests**: `go test ./...`
-   **Run Integration Tests**: `go test -tags=integration ./...`

---

## 5. Performance Tests

Performance testing is critical to validate the sub-100ms requirement under load. We will use `k6` for this purpose.

**Scope**: Load testing the deployed API endpoints in a staging environment that mirrors production.

### k6 Script

The script will define scenarios, thresholds, and checks.

**File**: `tests/performance/stream_k6.js`

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

// A custom metric to track response times
const streamItemResponseTrend = new Trend('stream_item_response_time');

export const options = {
  stages: [
    { duration: '30s', target: 100 }, // Ramp up to 100 virtual users
    { duration: '1m', target: 100 }, // Stay at 100 for 1 minute
    { duration: '10s', target: 0 },   // Ramp down
  ],
  thresholds: {
    'http_req_failed': ['rate<0.01'], // < 1% error rate
    'http_req_duration': ['p(95)<100'], // P95 response time < 100ms
    'stream_item_response_time': ['p(95)<100'],
  },
};

export default function () {
  // Fetch the list of items
  const res = http.get(`${__ENV.API_BASE_URL}/v2/stream`);
  check(res, { 'status was 200': (r) => r.status === 200 });

  // Assuming the first item is the one we want to test in detail
  const items = res.json('data');
  if (items && items.length > 0) {
    const itemId = items[0].id;
    const itemRes = http.get(`${__ENV.API_BASE_URL}/v2/stream/${itemId}`);
    check(itemRes, { 'item details status was 200': (r) => r.status === 200 });
    streamItemResponseTrend.add(itemRes.timings.duration);
  }

  sleep(1);
}
```

### Running Performance Tests

This test would be run against a deployed staging environment.

```bash
docker-compose -f docker-compose.test.yml up -d
k6 run --env API_BASE_URL=http://localhost:8080 tests/performance/stream_k6.js
```

---

## 6. CI/CD Pipeline Integration

All tests will be integrated into a CI/CD pipeline (e.g., GitHub Actions) to automate the quality assurance process.

**File**: `.github/workflows/ci.yml`

```yaml
name: Gravity BFF CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.21'

    - name: Run Unit Tests
      run: go test -v -race ./...

    - name: Build
      run: go build -v ./...

  integration-test:
    runs-on: ubuntu-latest
    services:
      postgres: 
        image: postgres:15-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test
        ports: ["5432:5432"]
      redis:
        image: redis:7-alpine
        ports: ["6379:6379"]
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.21'

    - name: Run Integration Tests
      run: go test -v -tags=integration ./...
      env:
        DB_HOST: localhost
        REDIS_HOST: localhost
        # ... other env vars

  performance-test:
    runs-on: ubuntu-latest
    needs: integration-test
    if: github.ref == 'refs/heads/main' # Only run on merge to main
    steps:
      # ... steps to deploy to a staging environment
      - name: Run k6 Performance Test
        uses: grafana/k6-action@v0.2.0
        with:
          filename: tests/performance/stream_k6.js
          flags: --env API_BASE_URL=${{ secrets.STAGING_API_URL }}
```

This pipeline ensures that every pull request is validated with unit and integration tests, and every merge to `main` is validated against our performance targets in a staging environment.
