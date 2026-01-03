import { PriorityItem, FilterType } from "@/types";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface StreamResponse {
  data: PriorityItem[];
  meta: {
    cursor: string | null;
    hasMore: boolean;
    totalCount: number;
  };
}

interface StreamItemResponse {
  data: PriorityItem;
}

class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// Token getter function - will be set by the auth hook
let getTokenFn: (() => Promise<string | null>) | null = null;

/**
 * Sets the token getter function. Called by useAuthenticatedApi hook.
 */
export function setTokenGetter(fn: () => Promise<string | null>) {
  getTokenFn = fn;
}

async function fetchWithAuth<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...options.headers,
  };

  // Get Clerk session token
  if (getTokenFn) {
    const token = await getTokenFn();
    if (token) {
      (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
    }
  }

  const response = await fetch(url, {
    ...options,
    headers,
    credentials: "include",
  });

  if (!response.ok) {
    if (response.status === 401) {
      throw new ApiError(401, "Unauthorized - please sign in");
    }
    const errorMessage = await response.text().catch(() => "Unknown error");
    throw new ApiError(response.status, errorMessage);
  }

  return response.json();
}

export const api = {
  /**
   * Fetches the priority stream with optional filtering and pagination.
   */
  async getStream(
    filter?: FilterType,
    limit?: number,
    cursor?: string
  ): Promise<StreamResponse> {
    const params = new URLSearchParams();

    if (filter && filter !== "all") {
      params.append("filter", filter);
    }
    if (limit) {
      params.append("limit", limit.toString());
    }
    if (cursor) {
      params.append("cursor", cursor);
    }

    const queryString = params.toString();
    const endpoint = `/v2/stream${queryString ? `?${queryString}` : ""}`;

    return fetchWithAuth<StreamResponse>(endpoint);
  },

  /**
   * Fetches a single stream item by ID, including all messages.
   */
  async getStreamItem(itemId: string): Promise<StreamItemResponse> {
    return fetchWithAuth<StreamItemResponse>(`/v2/stream/${itemId}`);
  },

  /**
   * Marks an item as read (placeholder for future implementation).
   */
  async markAsRead(itemId: string): Promise<void> {
    return fetchWithAuth(`/v2/stream/${itemId}/read`, {
      method: "POST",
    });
  },
};

export { ApiError };
export type { StreamResponse, StreamItemResponse };
