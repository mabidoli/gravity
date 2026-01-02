package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

func TestStreamKey(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		filter   model.StreamFilter
		cursor   *string
		expected string
	}{
		{
			name:     "all filter no cursor",
			userID:   "user-123",
			filter:   model.FilterAll,
			cursor:   nil,
			expected: "stream:user-123:all:none",
		},
		{
			name:     "high filter no cursor",
			userID:   "user-456",
			filter:   model.FilterHigh,
			cursor:   nil,
			expected: "stream:user-456:high:none",
		},
		{
			name:     "unread filter with cursor",
			userID:   "user-789",
			filter:   model.FilterUnread,
			cursor:   strPtr("abc123"),
			expected: "stream:user-789:unread:abc123",
		},
		{
			name:     "empty cursor string",
			userID:   "user-000",
			filter:   model.FilterAll,
			cursor:   strPtr(""),
			expected: "stream:user-000:all:none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreamKey(tt.userID, tt.filter, tt.cursor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestItemKey(t *testing.T) {
	tests := []struct {
		name     string
		itemID   string
		expected string
	}{
		{
			name:     "simple item ID",
			itemID:   "item-123",
			expected: "item:item-123",
		},
		{
			name:     "UUID item ID",
			itemID:   "550e8400-e29b-41d4-a716-446655440000",
			expected: "item:550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ItemKey(tt.itemID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
