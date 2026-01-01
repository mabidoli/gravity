package model

import (
	"time"
)

// User represents a participant in the system.
type User struct {
	ID        string  `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Email     *string `json:"email,omitempty" db:"email"`
	AvatarURL *string `json:"avatarUrl,omitempty" db:"avatar_url"`
}

// Attachment represents a file attached to a message.
type Attachment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MimeType  string `json:"mimeType"`
	SizeBytes int64  `json:"sizeBytes"`
	URL       string `json:"url"`
}

// CalendarEvent represents calendar event details within a message.
type CalendarEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Attendees   []User    `json:"attendees"`
	Location    *string   `json:"location,omitempty"`
	MeetingLink *string   `json:"meetingLink,omitempty"`
	Description *string   `json:"description,omitempty"`
}

// SocialStats represents engagement statistics for social content.
type SocialStats struct {
	Views    *int64 `json:"views,omitempty"`
	Likes    *int64 `json:"likes,omitempty"`
	Comments *int64 `json:"comments,omitempty"`
	Shares   *int64 `json:"shares,omitempty"`
}

// SocialContent represents content from social media platforms.
type SocialContent struct {
	ID              string         `json:"id"`
	Platform        SocialPlatform `json:"platform"`
	Author          string         `json:"author"`
	AuthorAvatarURL *string        `json:"authorAvatarUrl,omitempty"`
	ThumbnailURL    *string        `json:"thumbnailUrl,omitempty"`
	Title           *string        `json:"title,omitempty"`
	Description     *string        `json:"description,omitempty"`
	Stats           SocialStats    `json:"stats"`
	URL             string         `json:"url"`
}

// AIInsight represents an AI-generated insight or suggestion.
type AIInsight struct {
	ID      string      `json:"id"`
	Type    InsightType `json:"type"`
	Label   string      `json:"label"`
	Content string      `json:"content"`
	IsDraft bool        `json:"isDraft"`
}

// Message represents a single message in a conversation thread.
type Message struct {
	ID              string         `json:"id" db:"id"`
	SenderType      SenderType     `json:"senderType" db:"sender_type"`
	SenderInfo      *User          `json:"senderInfo,omitempty"`
	Content         string         `json:"content" db:"content"`
	Timestamp       time.Time      `json:"timestamp" db:"message_timestamp"`
	ContentType     ContentType    `json:"contentType" db:"content_type"`
	EventDetails    *CalendarEvent `json:"eventDetails,omitempty"`
	SocialContent   *SocialContent `json:"socialContent,omitempty"`
	AIInsights      []AIInsight    `json:"aiInsights,omitempty"`
	Attachments     []Attachment   `json:"attachments,omitempty"`
	FullContentHTML *string        `json:"fullContentHtml,omitempty" db:"full_content_html"`
}

// PriorityItem represents a single item in the unified priority stream.
type PriorityItem struct {
	ID           string     `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Source       SourceType `json:"source" db:"source"`
	Priority     Priority   `json:"priority" db:"priority"`
	IsUnread     bool       `json:"isUnread" db:"is_unread"`
	Snippet      *string    `json:"snippet,omitempty" db:"snippet"`
	Timestamp    time.Time  `json:"timestamp" db:"item_timestamp"`
	Participants []User     `json:"participants"`
	Messages     []Message  `json:"messages,omitempty"` // Only included in detail view
}

// PriorityItemWithMessages is the full detail view of a priority item.
type PriorityItemWithMessages struct {
	PriorityItem
	Messages []Message `json:"messages"`
}
