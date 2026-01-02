import { useQuery, useQueryClient, useInfiniteQuery } from "@tanstack/react-query";
import { api, StreamResponse } from "@/lib/api";
import { FilterType, PriorityItem } from "@/types";

// Query keys for cache management
export const streamKeys = {
  all: ["stream"] as const,
  lists: () => [...streamKeys.all, "list"] as const,
  list: (filter: FilterType) => [...streamKeys.lists(), filter] as const,
  details: () => [...streamKeys.all, "detail"] as const,
  detail: (id: string) => [...streamKeys.details(), id] as const,
};

/**
 * Hook to fetch the priority stream with filtering.
 */
export function useStream(filter: FilterType = "all") {
  return useQuery({
    queryKey: streamKeys.list(filter),
    queryFn: () => api.getStream(filter),
    select: (data) => data.data,
  });
}

/**
 * Hook to fetch the priority stream with infinite scrolling.
 */
export function useInfiniteStream(filter: FilterType = "all", limit = 20) {
  return useInfiniteQuery({
    queryKey: [...streamKeys.list(filter), "infinite"],
    queryFn: ({ pageParam }) => api.getStream(filter, limit, pageParam),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) =>
      lastPage.meta.hasMore ? lastPage.meta.cursor : undefined,
    select: (data) => ({
      items: data.pages.flatMap((page) => page.data),
      hasMore: data.pages[data.pages.length - 1]?.meta.hasMore ?? false,
    }),
  });
}

/**
 * Hook to fetch a single stream item with all messages.
 */
export function useStreamItem(itemId: string | null) {
  return useQuery({
    queryKey: streamKeys.detail(itemId ?? ""),
    queryFn: () => api.getStreamItem(itemId!),
    enabled: !!itemId,
    select: (data) => data.data,
  });
}

/**
 * Hook to prefetch a stream item (for hover preview).
 */
export function usePrefetchStreamItem() {
  const queryClient = useQueryClient();

  return (itemId: string) => {
    queryClient.prefetchQuery({
      queryKey: streamKeys.detail(itemId),
      queryFn: () => api.getStreamItem(itemId),
      staleTime: 60 * 1000, // 1 minute
    });
  };
}

/**
 * Hook to invalidate stream queries (call after mutations).
 */
export function useInvalidateStream() {
  const queryClient = useQueryClient();

  return {
    invalidateList: (filter?: FilterType) => {
      if (filter) {
        queryClient.invalidateQueries({ queryKey: streamKeys.list(filter) });
      } else {
        queryClient.invalidateQueries({ queryKey: streamKeys.lists() });
      }
    },
    invalidateItem: (itemId: string) => {
      queryClient.invalidateQueries({ queryKey: streamKeys.detail(itemId) });
    },
    invalidateAll: () => {
      queryClient.invalidateQueries({ queryKey: streamKeys.all });
    },
  };
}

/**
 * Hook to update cached item data optimistically.
 */
export function useOptimisticUpdate() {
  const queryClient = useQueryClient();

  return {
    markAsRead: (itemId: string) => {
      // Update all list caches
      queryClient.setQueriesData<StreamResponse>(
        { queryKey: streamKeys.lists() },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((item) =>
              item.id === itemId ? { ...item, unread: false } : item
            ),
          };
        }
      );

      // Update detail cache
      queryClient.setQueryData<{ data: PriorityItem }>(
        streamKeys.detail(itemId),
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: { ...old.data, unread: false },
          };
        }
      );
    },
  };
}
