// Package repository implements the data access layer using PostgreSQL.
package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mabidoli/gravity-bff/internal/domain/model"
	"github.com/mabidoli/gravity-bff/internal/domain/repository"
)

// Ensure PgStreamRepository implements the StreamRepository interface.
var _ repository.StreamRepository = (*PgStreamRepository)(nil)

// PgStreamRepository implements StreamRepository using PostgreSQL.
type PgStreamRepository struct {
	db *pgxpool.Pool
}

// NewPgStreamRepository creates a new PostgreSQL stream repository.
func NewPgStreamRepository(db *pgxpool.Pool) *PgStreamRepository {
	return &PgStreamRepository{db: db}
}

// cursor represents pagination cursor data.
type cursor struct {
	Timestamp time.Time `json:"t"`
	ID        string    `json:"id"`
}

// encodeCursor encodes a cursor for pagination.
func encodeCursor(timestamp time.Time, id string) string {
	c := cursor{Timestamp: timestamp, ID: id}
	data, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(data)
}

// decodeCursor decodes a pagination cursor.
func decodeCursor(encoded string) (*cursor, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}
	var c cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("invalid cursor data: %w", err)
	}
	return &c, nil
}

// GetStream retrieves a paginated list of priority items for a user.
func (r *PgStreamRepository) GetStream(ctx context.Context, req model.StreamRequest) ([]model.PriorityItem, *string, error) {
	// Build the query based on filter
	baseQuery := `
		SELECT id, title, source, priority, is_unread, snippet, item_timestamp
		FROM priority_items
		WHERE user_id = $1
	`

	args := []interface{}{req.UserID}
	argPos := 2

	// Apply filter
	switch req.Filter {
	case model.FilterHigh:
		baseQuery += fmt.Sprintf(" AND priority = $%d", argPos)
		args = append(args, string(model.PriorityHigh))
		argPos++
	case model.FilterUnread:
		baseQuery += " AND is_unread = TRUE"
	}

	// Apply cursor-based pagination
	if req.Cursor != nil && *req.Cursor != "" {
		c, err := decodeCursor(*req.Cursor)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid cursor: %w", err)
		}
		baseQuery += fmt.Sprintf(" AND (item_timestamp, id) < ($%d, $%d)", argPos, argPos+1)
		args = append(args, c.Timestamp, c.ID)
		argPos += 2
	}

	// Order by timestamp descending, then by ID for consistent ordering
	baseQuery += " ORDER BY item_timestamp DESC, id DESC"

	// Fetch one extra to determine if there are more items
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	baseQuery += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, limit+1)

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query stream: %w", err)
	}
	defer rows.Close()

	items := make([]model.PriorityItem, 0, limit)
	for rows.Next() {
		var item model.PriorityItem
		var source, priority string
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&source,
			&priority,
			&item.IsUnread,
			&item.Snippet,
			&item.Timestamp,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}
		item.Source = model.SourceType(source)
		item.Priority = model.Priority(priority)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("row iteration error: %w", err)
	}

	// Check if there are more items
	var nextCursor *string
	if len(items) > limit {
		items = items[:limit]
		lastItem := items[len(items)-1]
		encoded := encodeCursor(lastItem.Timestamp, lastItem.ID)
		nextCursor = &encoded
	}

	// Fetch participants for each item
	for i := range items {
		participants, err := r.GetParticipantsByItemID(ctx, items[i].ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get participants: %w", err)
		}
		items[i].Participants = participants
	}

	return items, nextCursor, nil
}

// GetStreamItemByID retrieves a single priority item with all its messages.
func (r *PgStreamRepository) GetStreamItemByID(ctx context.Context, userID, itemID string) (*model.PriorityItem, error) {
	query := `
		SELECT id, title, source, priority, is_unread, snippet, item_timestamp
		FROM priority_items
		WHERE id = $1 AND user_id = $2
	`

	var item model.PriorityItem
	var source, priority string

	err := r.db.QueryRow(ctx, query, itemID, userID).Scan(
		&item.ID,
		&item.Title,
		&source,
		&priority,
		&item.IsUnread,
		&item.Snippet,
		&item.Timestamp,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Item not found
		}
		return nil, fmt.Errorf("failed to get stream item: %w", err)
	}

	item.Source = model.SourceType(source)
	item.Priority = model.Priority(priority)

	// Fetch participants
	participants, err := r.GetParticipantsByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	item.Participants = participants

	// Fetch messages
	messages, err := r.GetMessagesByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	item.Messages = messages

	return &item, nil
}

// GetParticipantsByItemID retrieves all participants for a priority item.
func (r *PgStreamRepository) GetParticipantsByItemID(ctx context.Context, itemID string) ([]model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.avatar_url
		FROM users u
		JOIN priority_item_participants pip ON u.id = pip.user_id
		WHERE pip.item_id = $1
	`

	rows, err := r.db.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query participants: %w", err)
	}
	defer rows.Close()

	var participants []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.AvatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return participants, nil
}

// GetMessagesByItemID retrieves all messages for a priority item.
func (r *PgStreamRepository) GetMessagesByItemID(ctx context.Context, itemID string) ([]model.Message, error) {
	query := `
		SELECT m.id, m.sender_id, m.sender_type, m.content_type, m.content,
			   m.full_content_html, m.message_timestamp,
			   m.event_details, m.social_details, m.attachments, m.ai_insights,
			   u.id, u.name, u.email, u.avatar_url
		FROM messages m
		LEFT JOIN users u ON m.sender_id = u.id
		WHERE m.item_id = $1
		ORDER BY m.message_timestamp ASC
	`

	rows, err := r.db.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		var senderType, contentType string
		var senderID, userName, userEmail, userAvatar *string
		var eventDetails, socialDetails, attachments, aiInsights []byte

		err := rows.Scan(
			&msg.ID,
			&senderID,
			&senderType,
			&contentType,
			&msg.Content,
			&msg.FullContentHTML,
			&msg.Timestamp,
			&eventDetails,
			&socialDetails,
			&attachments,
			&aiInsights,
			&senderID,
			&userName,
			&userEmail,
			&userAvatar,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		msg.SenderType = model.SenderType(senderType)
		msg.ContentType = model.ContentType(contentType)

		// Parse sender info if available
		if senderID != nil && userName != nil {
			msg.SenderInfo = &model.User{
				ID:        *senderID,
				Name:      *userName,
				Email:     userEmail,
				AvatarURL: userAvatar,
			}
		}

		// Parse JSONB fields
		if len(eventDetails) > 0 {
			var event model.CalendarEvent
			if err := json.Unmarshal(eventDetails, &event); err == nil {
				msg.EventDetails = &event
			}
		}

		if len(socialDetails) > 0 {
			var social model.SocialContent
			if err := json.Unmarshal(socialDetails, &social); err == nil {
				msg.SocialContent = &social
			}
		}

		if len(attachments) > 0 {
			var atts []model.Attachment
			if err := json.Unmarshal(attachments, &atts); err == nil {
				msg.Attachments = atts
			}
		}

		if len(aiInsights) > 0 {
			var insights []model.AIInsight
			if err := json.Unmarshal(aiInsights, &insights); err == nil {
				msg.AIInsights = insights
			}
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return messages, nil
}
