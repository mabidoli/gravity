/**
 * Gravity BFF Performance Test Script
 *
 * This k6 script tests the performance of the Gravity BFF API endpoints.
 * Primary goal: Validate sub-100ms P95 response time requirement.
 *
 * Usage:
 *   k6 run --env API_BASE_URL=http://localhost:8080 tests/performance/stream_k6.js
 *
 * With reporting:
 *   k6 run --env API_BASE_URL=http://localhost:8080 --out json=results.json tests/performance/stream_k6.js
 */

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

// Custom metrics
const streamListDuration = new Trend('stream_list_duration', true);
const streamItemDuration = new Trend('stream_item_duration', true);
const healthCheckDuration = new Trend('health_check_duration', true);
const errorRate = new Rate('errors');
const requestCount = new Counter('requests');

// Test configuration
export const options = {
  // Staged load test
  stages: [
    { duration: '30s', target: 50 },   // Ramp up to 50 users
    { duration: '1m', target: 100 },   // Ramp up to 100 users
    { duration: '2m', target: 100 },   // Stay at 100 users for 2 minutes
    { duration: '30s', target: 200 },  // Spike to 200 users
    { duration: '1m', target: 200 },   // Stay at 200 users
    { duration: '30s', target: 0 },    // Ramp down
  ],

  // Performance thresholds - PRIMARY REQUIREMENT: sub-100ms P95
  thresholds: {
    // Overall HTTP metrics
    'http_req_duration': ['p(95)<100', 'p(99)<200'], // P95 < 100ms, P99 < 200ms
    'http_req_failed': ['rate<0.01'],                // < 1% error rate

    // Endpoint-specific metrics
    'stream_list_duration': ['p(95)<100', 'avg<50'],
    'stream_item_duration': ['p(95)<100', 'avg<50'],
    'health_check_duration': ['p(95)<50', 'avg<10'],

    // Custom error rate
    'errors': ['rate<0.01'],
  },

  // Tags for better organization
  tags: {
    testType: 'performance',
    service: 'gravity-bff',
  },
};

// Configuration
const BASE_URL = __ENV.API_BASE_URL || 'http://localhost:8080';
const AUTH_TOKEN = __ENV.AUTH_TOKEN || 'test-user-1';

// Headers
const headers = {
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${AUTH_TOKEN}`,
};

// Test setup
export function setup() {
  // Verify API is accessible
  const healthRes = http.get(`${BASE_URL}/health`);
  if (healthRes.status !== 200) {
    throw new Error(`API health check failed: ${healthRes.status}`);
  }

  console.log(`Starting performance test against ${BASE_URL}`);
  return { baseUrl: BASE_URL };
}

// Main test scenario
export default function (data) {
  const baseUrl = data.baseUrl;

  group('Health Check', () => {
    const res = http.get(`${baseUrl}/health`);
    healthCheckDuration.add(res.timings.duration);
    requestCount.add(1);

    const success = check(res, {
      'health status is 200': (r) => r.status === 200,
      'health response is ok': (r) => {
        const body = JSON.parse(r.body);
        return body.status === 'ok';
      },
    });

    errorRate.add(!success);
  });

  sleep(0.1); // Small pause between groups

  group('Get Stream List', () => {
    // Test different filter combinations
    const filters = ['all', 'high', 'unread'];
    const filter = filters[Math.floor(Math.random() * filters.length)];
    const limit = Math.floor(Math.random() * 50) + 10; // Random limit 10-60

    const res = http.get(
      `${baseUrl}/v2/stream?filter=${filter}&limit=${limit}`,
      { headers }
    );

    streamListDuration.add(res.timings.duration);
    requestCount.add(1);

    const success = check(res, {
      'stream list status is 200': (r) => r.status === 200,
      'stream list has data array': (r) => {
        const body = JSON.parse(r.body);
        return Array.isArray(body.data);
      },
      'stream list response time < 100ms': (r) => r.timings.duration < 100,
    });

    errorRate.add(!success);

    // If we got items, test fetching item details
    if (res.status === 200) {
      const body = JSON.parse(res.body);
      if (body.data && body.data.length > 0) {
        const itemId = body.data[0].id;
        testGetStreamItem(baseUrl, itemId);
      }
    }
  });

  sleep(Math.random() * 2 + 0.5); // Random sleep 0.5-2.5 seconds
}

// Helper function to test stream item endpoint
function testGetStreamItem(baseUrl, itemId) {
  group('Get Stream Item Details', () => {
    const res = http.get(`${baseUrl}/v2/stream/${itemId}`, { headers });

    streamItemDuration.add(res.timings.duration);
    requestCount.add(1);

    const success = check(res, {
      'stream item status is 200': (r) => r.status === 200,
      'stream item has id': (r) => {
        const body = JSON.parse(r.body);
        return body.id === itemId;
      },
      'stream item response time < 100ms': (r) => r.timings.duration < 100,
    });

    errorRate.add(!success);
  });
}

// Test with pagination
export function testPagination() {
  group('Pagination Test', () => {
    let cursor = null;
    let pageCount = 0;
    const maxPages = 5;

    while (pageCount < maxPages) {
      let url = `${BASE_URL}/v2/stream?limit=20`;
      if (cursor) {
        url += `&cursor=${cursor}`;
      }

      const res = http.get(url, { headers });
      requestCount.add(1);

      check(res, {
        'pagination status is 200': (r) => r.status === 200,
      });

      if (res.status !== 200) break;

      const body = JSON.parse(res.body);
      cursor = body.nextCursor;
      pageCount++;

      if (!cursor) break; // No more pages

      sleep(0.1);
    }
  });
}

// Teardown
export function teardown(data) {
  console.log('Performance test completed');
  console.log(`Tested against: ${data.baseUrl}`);
}

// Summary handler for custom reporting
export function handleSummary(data) {
  const p95 = data.metrics.http_req_duration.values['p(95)'];
  const p99 = data.metrics.http_req_duration.values['p(99)'];
  const avgDuration = data.metrics.http_req_duration.values['avg'];
  const errorRate = data.metrics.http_req_failed.values.rate;

  console.log('\n=== Performance Test Results ===');
  console.log(`P95 Response Time: ${p95.toFixed(2)}ms (target: <100ms)`);
  console.log(`P99 Response Time: ${p99.toFixed(2)}ms (target: <200ms)`);
  console.log(`Average Response Time: ${avgDuration.toFixed(2)}ms`);
  console.log(`Error Rate: ${(errorRate * 100).toFixed(2)}% (target: <1%)`);
  console.log('================================\n');

  // Return summary for k6 to output
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
  };
}

// Text summary helper
function textSummary(data, options) {
  const indent = options.indent || '';
  let output = '';

  output += `${indent}Performance Test Summary\n`;
  output += `${indent}========================\n\n`;

  // Key metrics
  const metrics = [
    ['http_req_duration', 'HTTP Request Duration'],
    ['stream_list_duration', 'Stream List Duration'],
    ['stream_item_duration', 'Stream Item Duration'],
    ['health_check_duration', 'Health Check Duration'],
  ];

  for (const [key, name] of metrics) {
    if (data.metrics[key]) {
      const m = data.metrics[key].values;
      output += `${indent}${name}:\n`;
      output += `${indent}  avg: ${m.avg?.toFixed(2) || 'N/A'}ms\n`;
      output += `${indent}  p95: ${m['p(95)']?.toFixed(2) || 'N/A'}ms\n`;
      output += `${indent}  p99: ${m['p(99)']?.toFixed(2) || 'N/A'}ms\n`;
      output += `${indent}  max: ${m.max?.toFixed(2) || 'N/A'}ms\n\n`;
    }
  }

  // Request counts
  if (data.metrics.requests) {
    output += `${indent}Total Requests: ${data.metrics.requests.values.count}\n`;
  }

  // Error rate
  if (data.metrics.errors) {
    output += `${indent}Error Rate: ${(data.metrics.errors.values.rate * 100).toFixed(2)}%\n`;
  }

  return output;
}
