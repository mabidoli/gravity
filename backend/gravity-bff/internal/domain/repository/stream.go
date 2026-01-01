// Package repository defines the data access interfaces for the domain.
// These interfaces abstract the underlying data store implementation.
package repository

import (
	"context"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
)

// StreamRepository defines the interface for priority stream data access.
type StreamRepository interface {
	// GetStream retrieves a paginated list of priority items for a user.
	// It supports filtering by priority or unread status.
	// Returns items, next cursor (nil if no more items), and any error.
	GetStream(ctx context.Context, req model.StreamRequest) ([]model.PriorityItem, *string, error)

	// GetStreamItemByID retrieves a single priority item with all its messages.
	// Returns the full item details including messages, or an error if not found.
	GetStreamItemByID(ctx context.Context, userID, itemID string) (*model.PriorityItem, error)

	// GetParticipantsByItemID retrieves all participants for a priority item.
	GetParticipantsByItemID(ctx context.Context, itemID string) ([]model.User, error)

	// GetMessagesByItemID retrieves all messages for a priority item.
	GetMessagesByItemID(ctx context.Context, itemID string) ([]model.Message, error)
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	// GetUserByID retrieves a user by their ID.
	GetUserByID(ctx context.Context, userID string) (*model.User, error)

	// GetUsersByIDs retrieves multiple users by their IDs.
	GetUsersByIDs(ctx context.Context, userIDs []string) ([]model.User, error)
}
