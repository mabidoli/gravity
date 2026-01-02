package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPriorityItem_JSON(t *testing.T) {
	// Arrange
	timestamp := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	snippet := "Test snippet"

	item := PriorityItem{
		ID:        "item-1",
		Title:     "Test Item",
		Source:    SourceEmail,
		Priority:  PriorityHigh,
		IsUnread:  true,
		Snippet:   &snippet,
		Timestamp: timestamp,
		Participants: []User{
			{ID: "user-1", Name: "John Doe"},
		},
	}

	// Act
	data, err := json.Marshal(item)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, "item-1", result["id"])
	assert.Equal(t, "Test Item", result["title"])
	assert.Equal(t, "email", result["source"])
	assert.Equal(t, "high", result["priority"])
	assert.Equal(t, true, result["isUnread"])
	assert.Equal(t, "Test snippet", result["snippet"])
}

func TestUser_JSON(t *testing.T) {
	// Arrange
	email := "test@example.com"
	avatar := "https://example.com/avatar.jpg"

	user := User{
		ID:        "user-1",
		Name:      "John Doe",
		Email:     &email,
		AvatarURL: &avatar,
	}

	// Act
	data, err := json.Marshal(user)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, "user-1", result["id"])
	assert.Equal(t, "John Doe", result["name"])
	assert.Equal(t, "test@example.com", result["email"])
	assert.Equal(t, "https://example.com/avatar.jpg", result["avatarUrl"])
}

func TestUser_JSON_OptionalFieldsOmitted(t *testing.T) {
	// Arrange - user without optional fields
	user := User{
		ID:   "user-1",
		Name: "John Doe",
	}

	// Act
	data, err := json.Marshal(user)

	// Assert
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "email")
	assert.NotContains(t, string(data), "avatarUrl")
}

func TestMessage_JSON(t *testing.T) {
	// Arrange
	timestamp := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	content := "Hello, World!"

	msg := Message{
		ID:          "msg-1",
		SenderType:  SenderOther,
		Content:     content,
		Timestamp:   timestamp,
		ContentType: ContentText,
	}

	// Act
	data, err := json.Marshal(msg)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, "msg-1", result["id"])
	assert.Equal(t, "other", result["senderType"])
	assert.Equal(t, "Hello, World!", result["content"])
	assert.Equal(t, "text", result["contentType"])
}

func TestAttachment_JSON(t *testing.T) {
	// Arrange
	att := Attachment{
		ID:        "att-1",
		Name:      "document.pdf",
		MimeType:  "application/pdf",
		SizeBytes: 1024000,
		URL:       "https://example.com/doc.pdf",
	}

	// Act
	data, err := json.Marshal(att)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, "att-1", result["id"])
	assert.Equal(t, "document.pdf", result["name"])
	assert.Equal(t, "application/pdf", result["mimeType"])
	assert.Equal(t, float64(1024000), result["sizeBytes"])
}

func TestAIInsight_JSON(t *testing.T) {
	// Arrange
	insight := AIInsight{
		ID:      "insight-1",
		Type:    InsightDraft,
		Label:   "Draft Available",
		Content: "Here's a suggested response...",
		IsDraft: true,
	}

	// Act
	data, err := json.Marshal(insight)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, "insight-1", result["id"])
	assert.Equal(t, "draft", result["type"])
	assert.Equal(t, true, result["isDraft"])
}

func TestStreamResponse_JSON(t *testing.T) {
	// Arrange
	nextCursor := "cursor123"
	resp := StreamResponse{
		Data: []PriorityItem{
			{ID: "item-1", Title: "Test"},
		},
		NextCursor: &nextCursor,
	}

	// Act
	data, err := json.Marshal(resp)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	assert.Equal(t, "cursor123", result["nextCursor"])
	dataArr := result["data"].([]interface{})
	assert.Equal(t, 1, len(dataArr))
}

func TestStreamResponse_NullCursor(t *testing.T) {
	// Arrange
	resp := StreamResponse{
		Data:       []PriorityItem{},
		NextCursor: nil,
	}

	// Act
	data, err := json.Marshal(resp)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"nextCursor":null`)
}
