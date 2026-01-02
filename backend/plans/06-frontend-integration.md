# Phase 4.1: Frontend-Backend Integration Plan

## Overview

This document outlines the plan to integrate the Next.js frontend with the Go BFF backend, replacing mock data with real API calls.

---

## 1. Type Mapping Analysis

### Field Name Differences

| Frontend Type | Frontend Field | Backend Field | Notes |
|---------------|----------------|---------------|-------|
| `User` | `avatar` | `avatarUrl` | Rename required |
| `Attachment` | `type` | `mimeType` | Rename required |
| `Attachment` | `size` | `sizeBytes` | Rename required |
| `Message` | `sender` | `senderType` | Rename required |
| `Message` | `type` | `contentType` | Rename required |
| `Message` | `fullContent` | `fullContentHtml` | Rename required |
| `PriorityItem` | `unread` | `isUnread` | Rename required |
| `SocialContent` | `authorAvatar` | `authorAvatarUrl` | Rename required |
| `SocialContent` | `thumbnail` | `thumbnailUrl` | Rename required |

### Decision: Align Frontend Types with Backend

Option A: Transform responses in frontend (adapter pattern)
Option B: Update frontend types to match backend (breaking change)
Option C: Update backend JSON tags to match frontend (cleaner API)

**Recommended: Option C** - Update backend JSON serialization to match frontend expectations. This is the cleanest approach as:
- The BFF exists specifically to serve the frontend
- No transformation logic needed in frontend
- Single source of truth

---

## 2. Implementation Tasks

### Task 4.1.1: Update Backend JSON Tags

**Files to modify:**
- `backend/gravity-bff/internal/domain/model/stream.go`

**Changes:**
```go
// User
AvatarURL *string `json:"avatar,omitempty"`  // was "avatarUrl"

// Attachment
MimeType  string `json:"type"`      // was "mimeType"
SizeBytes int64  `json:"size"`      // was "sizeBytes"

// Message
SenderType  SenderType  `json:"sender"`      // was "senderType"
ContentType ContentType `json:"type"`        // was "contentType"
FullContentHTML *string `json:"fullContent,omitempty"`  // was "fullContentHtml"

// PriorityItem
IsUnread bool `json:"unread"`  // was "isUnread"

// SocialContent
AuthorAvatarURL *string `json:"authorAvatar,omitempty"`  // was "authorAvatarUrl"
ThumbnailURL    *string `json:"thumbnail,omitempty"`     // was "thumbnailUrl"
```

### Task 4.1.2: Configure CORS in Backend

**File:** `backend/gravity-bff/internal/api/middleware/cors.go`

```go
package middleware

import "github.com/gofiber/fiber/v2/middleware/cors"

func CORSConfig() cors.Config {
    return cors.Config{
        AllowOrigins: "http://localhost:3000",  // Frontend dev server
        AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
        AllowHeaders: "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: true,
    }
}
```

**Update router.go** to use CORS middleware.

### Task 4.1.3: Create Frontend API Client

**File:** `frontend/src/lib/api.ts`

```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const api = {
  async getStream(filter?: string, limit?: number, cursor?: string) {
    const params = new URLSearchParams();
    if (filter && filter !== 'all') params.append('filter', filter);
    if (limit) params.append('limit', limit.toString());
    if (cursor) params.append('cursor', cursor);

    const response = await fetch(`${API_BASE_URL}/v2/stream?${params}`, {
      headers: { 'Authorization': `Bearer ${getUserToken()}` }
    });
    return response.json();
  },

  async getStreamItem(itemId: string) {
    const response = await fetch(`${API_BASE_URL}/v2/stream/${itemId}`, {
      headers: { 'Authorization': `Bearer ${getUserToken()}` }
    });
    return response.json();
  }
};
```

### Task 4.1.4: Add TanStack Query Provider

**File:** `frontend/src/app/providers.tsx`

```typescript
'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useState } from 'react';

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000,        // 1 minute
        refetchOnWindowFocus: false,
      },
    },
  }));

  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}
```

**Update layout.tsx** to wrap with Providers.

### Task 4.1.5: Create Custom Hooks for Data Fetching

**File:** `frontend/src/hooks/useStream.ts`

```typescript
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { FilterType, PriorityItem } from '@/types';

export function useStream(filter: FilterType = 'all') {
  return useQuery({
    queryKey: ['stream', filter],
    queryFn: () => api.getStream(filter !== 'all' ? filter : undefined),
  });
}

export function useStreamItem(itemId: string | null) {
  return useQuery({
    queryKey: ['streamItem', itemId],
    queryFn: () => api.getStreamItem(itemId!),
    enabled: !!itemId,
  });
}
```

### Task 4.1.6: Update Zustand Store

Refactor `useGravityStore.ts` to:
1. Remove mock data import
2. Initialize with empty items array
3. Add actions to set items from API
4. Keep local state management (selection, modals)

**Key changes:**

```typescript
// Remove: import { mockStreamItems } from "@/data/mockData";

interface GravityState {
  // Data from API
  items: PriorityItem[];
  isLoading: boolean;
  error: string | null;

  // ... rest of state

  // New actions
  setItems: (items: PriorityItem[]) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const useGravityStore = create<GravityState>((set, get) => ({
  items: [],  // Start empty, populated by API
  isLoading: true,
  error: null,
  // ...
}));
```

### Task 4.1.7: Update Components to Use Hooks

**StreamSidebar.tsx:**
```typescript
import { useStream } from '@/hooks/useStream';

export function StreamSidebar() {
  const { filter, setFilter, selectItem, selectedItemId } = useGravityStore();
  const { data, isLoading, error } = useStream(filter);

  // Use data?.data for items (from API response)
  const items = data?.data ?? [];
  // ...
}
```

### Task 4.1.8: Environment Configuration

**File:** `frontend/.env.local`
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

**File:** `frontend/.env.production`
```
NEXT_PUBLIC_API_URL=https://api.gravity.example.com
```

### Task 4.1.9: Docker Compose for Full Stack

**File:** `docker-compose.fullstack.yml`

```yaml
version: '3.8'

services:
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    depends_on:
      - api

  api:
    build:
      context: ./backend/gravity-bff
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=gravity
      - DB_PASSWORD=secret
      - DB_NAME=gravity_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - CORS_ORIGINS=http://localhost:3000
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: gravity
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: gravity_db
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  migrate:
    image: migrate/migrate:v4.17.0
    volumes:
      - ./backend/gravity-bff/migrations:/migrations
    command: ["-path", "/migrations", "-database", "postgres://gravity:secret@postgres:5432/gravity_db?sslmode=disable", "up"]
    depends_on:
      - postgres

volumes:
  postgres_data:
  redis_data:
```

---

## 3. Implementation Order

1. **Task 4.1.1** - Update backend JSON tags (aligns API with frontend types)
2. **Task 4.1.2** - Add CORS middleware to backend
3. **Task 4.1.3** - Create API client in frontend
4. **Task 4.1.4** - Add TanStack Query provider
5. **Task 4.1.5** - Create custom data fetching hooks
6. **Task 4.1.6** - Update Zustand store (remove mock data)
7. **Task 4.1.7** - Update components to use hooks
8. **Task 4.1.8** - Configure environment variables
9. **Task 4.1.9** - Create fullstack docker-compose

---

## 4. Testing Strategy

### Manual Testing Checklist

- [ ] Frontend loads without mock data
- [ ] Stream list fetches from API
- [ ] Filter tabs work (all, high, unread)
- [ ] Clicking item fetches details
- [ ] Messages display correctly
- [ ] Cache invalidation works on mutations

### Integration Test Updates

Update backend integration tests to verify:
- CORS headers present
- JSON field names match frontend expectations

---

## 5. Rollback Strategy

If issues arise:
1. Revert frontend to mock data by re-importing `mockStreamItems`
2. Frontend can work independently while backend issues are resolved
3. Feature flag option: `NEXT_PUBLIC_USE_MOCK_DATA=true`

---

## 6. Future Considerations

- **Phase 4.2 (Clerk Auth)**: Replace `getUserToken()` with Clerk token
- **Real-time updates**: Consider WebSocket or SSE for live stream updates
- **Optimistic updates**: For marking items as read
- **Offline support**: Service worker caching

---

## Dependencies

- TanStack Query already in package.json
- No new npm packages required
- Backend changes are non-breaking (JSON tag updates only)
