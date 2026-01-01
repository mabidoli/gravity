# Gravity V2 - Backend-for-Frontend (BFF) API Specification

## 1. Overview

This document specifies the design for the Gravity V2 Backend-for-Frontend (BFF) API. The primary purpose of this API is to serve pre-processed, aggregated, and read-only data to the Gravity V2 frontend application. It acts as a tailored data source, abstracting away the complexities of underlying data sources and services.

This API assumes that another service is responsible for connecting to third-party platforms (e.g., Gmail, Slack, Calendar), ingesting data, processing it, and storing it in a database accessible to this BFF.

### Guiding Principles

- **Client-Centric**: The API is designed specifically for the needs of the Gravity V2 frontend, providing data in the exact format required.
- **Read-Only**: This API only supports data retrieval (`GET` operations). All data creation, updates, and deletions are handled by upstream services.
- **Stateless**: Each API request is independent and contains all necessary information. The server does not maintain client session state.
- **Framework-Agnostic**: The design is based on universal RESTful principles and can be implemented using any modern backend framework (e.g., Express.js, FastAPI, Spring Boot).

---

## 2. General Architecture

### 2.1. Protocol and Base URL

The API will be served over HTTPS. A versioned base URL structure will be used to ensure future compatibility.

- **Base URL**: `https://api.gravity.com/v2/`

### 2.2. Authentication

All endpoints must be protected and require authentication. While the specific implementation is out of scope for this design, the recommended approach is to use **OAuth 2.0** with **JSON Web Tokens (JWT)**.

- The client would include a JWT in the `Authorization` header of each request:
  `Authorization: Bearer <jwt_token>`

### 2.3. Data Format

All data will be exchanged in **JSON** format. Request bodies (where applicable in other services) and response bodies will use `Content-Type: application/json`.

### 2.4. Error Handling

The API will use standard HTTP status codes to indicate the outcome of a request. Error responses will include a structured JSON body with a machine-readable error code and a human-readable message.

**Example Error Response:**
```json
{
  "error": {
    "code": "resource_not_found",
    "message": "The requested priority item with ID item-999 does not exist."
  }
}
```

| Status Code | Meaning |
| :--- | :--- |
| `200 OK` | The request was successful. |
| `400 Bad Request` | The request was malformed (e.g., invalid query parameters). |
| `401 Unauthorized` | The request lacks valid authentication credentials. |
| `403 Forbidden` | The authenticated user does not have permission to access the resource. |
| `404 Not Found` | The requested resource does not exist. |
| `500 Internal Server Error` | An unexpected error occurred on the server. |

### 2.5. Pagination

For endpoints that return a list of resources (e.g., the unified stream), cursor-based pagination will be used for efficient and stable navigation.

- **Query Parameters**: `limit` (integer, default: 20, max: 100) and `cursor` (string).
- **Response Body**: The response will include a `nextCursor` field. If `nextCursor` is `null`, the client has reached the end of the list.

---

## 3. Data Models

These models define the structure of the resources returned by the API, directly corresponding to the frontend's TypeScript types.

<details>
<summary><strong>User Object</strong></summary>

```json
{
  "id": "user-1",
  "name": "Sarah Chen",
  "email": "sarah.chen@company.com",
  "avatarUrl": "https://cdn.gravity.com/avatars/sarah.jpg"
}
```
</details>

<details>
<summary><strong>Attachment Object</strong></summary>

```json
{
  "id": "att-1",
  "name": "Q4_Revenue_Report.pdf",
  "mimeType": "application/pdf",
  "sizeBytes": 2456000,
  "url": "https://cdn.gravity.com/attachments/q4-report.pdf"
}
```
</details>

<details>
<summary><strong>CalendarEvent Object</strong></summary>

```json
{
  "id": "event-1",
  "title": "Product Sync",
  "startTime": "2026-01-15T14:00:00Z",
  "endTime": "2026-01-15T15:00:00Z",
  "attendees": [ { "id": "user-2", "name": "Mike Johnson", ... } ],
  "location": "Conference Room A / Zoom",
  "meetingLink": "https://zoom.us/j/123456789",
  "description": "Weekly sync to discuss product roadmap..."
}
```
</details>

<details>
<summary><strong>SocialContent Object</strong></summary>

```json
{
  "id": "social-yt-1",
  "platform": "youtube",
  "author": "Tech Explained",
  "authorAvatarUrl": "https://cdn.gravity.com/avatars/tech-explained.jpg",
  "thumbnailUrl": "https://img.youtube.com/vi/xyz123/0.jpg",
  "title": "The Future of AI",
  "description": "A deep dive into the next generation of AI models...",
  "stats": {
    "views": 1500000,
    "likes": 82000,
    "comments": 4500
  },
  "url": "https://www.youtube.com/watch?v=xyz123"
}
```
</details>

<details>
<summary><strong>AIInsight Object</strong></summary>

```json
{
  "id": "insight-1",
  "type": "draft",
  "label": "✨ Draft Available",
  "content": "Hi Sarah, Thank you for the report...",
  "isDraft": true
}
```
</details>

<details>
<summary><strong>Message Object</strong></summary>

```json
{
  "id": "msg-1-1",
  "senderType": "other",
  "senderInfo": { "id": "user-1", "name": "Sarah Chen", ... },
  "content": "Hi, Please review the attached Q4 projections...",
  "timestamp": "2026-01-01T10:00:00Z",
  "contentType": "text",
  "eventDetails": null,
  "socialContent": null,
  "aiInsights": [ { "id": "insight-1", ... } ],
  "attachments": [ { "id": "att-1", ... } ],
  "fullContentHtml": "<!DOCTYPE html>..."
}
```
</details>

<details>
<summary><strong>PriorityItem Object (Core Resource)</strong></summary>

```json
{
  "id": "item-1",
  "title": "Q4 Revenue Report - Action Required",
  "source": "email",
  "priority": "high",
  "isUnread": true,
  "snippet": "Please review the attached Q4 projections...",
  "timestamp": "2026-01-01T10:00:00Z",
  "participants": [ { "id": "user-1", "name": "Sarah Chen", ... } ]
}
```
*Note: The `messages` array is only included when fetching a single item, not in the list view.*
</details>

---

## 4. Endpoint Definitions

### 4.1. Health Check

Provides a simple health check to verify that the API is running and accessible.

- **Endpoint**: `GET /health`
- **Authentication**: Not required.
- **Success Response (200 OK)**:
  ```json
  {
    "status": "ok",
    "timestamp": "2026-01-01T07:45:00Z"
  }
  ```

### 4.2. Get Unified Stream

Retrieves a paginated list of `PriorityItem` resources for the authenticated user, representing their unified stream.

- **Endpoint**: `GET /stream`
- **Query Parameters**:
  - `filter` (string, optional): Filters the items. Allowed values: `all` (default), `high` (for high priority), `unread`.
  - `limit` (integer, optional): The maximum number of items to return. Default: `20`, Max: `100`.
  - `cursor` (string, optional): The cursor from the previous response to fetch the next page.
- **Success Response (200 OK)**:
  ```json
  {
    "data": [
      {
        "id": "item-1",
        "title": "Q4 Revenue Report - Action Required",
        "source": "email",
        "priority": "high",
        "isUnread": true,
        "snippet": "Please review the attached Q4 projections...",
        "timestamp": "2026-01-01T10:00:00Z",
        "participants": [ { "id": "user-1", "name": "Sarah Chen", ... } ]
      },
      {
        "id": "item-2",
        "title": "Product Sync - Starting in 15 minutes",
        "source": "calendar",
        "priority": "high",
        "isUnread": true,
        "snippet": "Weekly product team sync with engineering leads",
        "timestamp": "2026-01-01T10:30:00Z",
        "participants": [ { "id": "user-2", ... }, { "id": "user-3", ... } ]
      }
    ],
    "nextCursor": "c_aXRlbS0y"
  }
  ```

### 4.3. Get Stream Item Details

Retrieves the full details of a single `PriorityItem`, including its complete message history.

- **Endpoint**: `GET /stream/{itemId}`
- **URL Parameters**:
  - `itemId` (string, required): The unique identifier of the priority item.
- **Success Response (200 OK)**:
  ```json
  {
    "id": "item-1",
    "title": "Q4 Revenue Report - Action Required",
    "source": "email",
    "priority": "high",
    "isUnread": true,
    "snippet": "Please review the attached Q4 projections...",
    "timestamp": "2026-01-01T10:00:00Z",
    "participants": [ { "id": "user-1", "name": "Sarah Chen", ... } ],
    "messages": [
      {
        "id": "msg-1-1",
        "senderType": "other",
        "senderInfo": { "id": "user-1", "name": "Sarah Chen", ... },
        "content": "Hi, Please review the attached Q4 projections...",
        "timestamp": "2026-01-01T10:00:00Z",
        "contentType": "text",
        "eventDetails": null,
        "socialContent": null,
        "aiInsights": [ { "id": "insight-1", ... } ],
        "attachments": [ { "id": "att-1", ... } ],
        "fullContentHtml": "<!DOCTYPE html>..."
      }
    ]
  }
  ```
- **Error Response (404 Not Found)**: Returned if `itemId` does not exist or the user does not have access.

---

## 5. Frontend Interaction Flow

This section illustrates a typical sequence of API calls from the frontend.

1.  **Initial Load**: The frontend loads and calls `GET /stream` (with default filter `all`) to populate the unified stream sidebar.
2.  **User Selects an Item**: The user clicks on a `PriorityItem` in the sidebar.
3.  **Fetch Details**: The frontend calls `GET /stream/{itemId}` using the ID of the selected item.
4.  **Display Conversation**: The response from the details endpoint is used to render the full conversation in the main chat interface, including all messages, attachments, and AI insights.
5.  **Filtering**: If the user clicks a filter tab (e.g., "High Priority"), the frontend calls `GET /stream?filter=high` to get a new list of items.
6.  **Infinite Scroll**: As the user scrolls down the stream sidebar, the frontend uses the `nextCursor` from the previous `/stream` response to call `GET /stream?cursor={nextCursor}` to load the next page of items.
