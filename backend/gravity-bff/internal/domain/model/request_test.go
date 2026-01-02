package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	// Act
	resp := NewErrorResponse("test_code", "Test message")

	// Assert
	assert.Equal(t, "test_code", resp.Error.Code)
	assert.Equal(t, "Test message", resp.Error.Message)
}

func TestErrorResponse_JSON(t *testing.T) {
	// Arrange
	resp := NewErrorResponse(ErrCodeNotFound, "Resource not found")

	// Act
	data, err := json.Marshal(resp)

	// Assert
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	errorObj := result["error"].(map[string]interface{})
	assert.Equal(t, "resource_not_found", errorObj["code"])
	assert.Equal(t, "Resource not found", errorObj["message"])
}

func TestStreamFilter_Constants(t *testing.T) {
	assert.Equal(t, StreamFilter("all"), FilterAll)
	assert.Equal(t, StreamFilter("high"), FilterHigh)
	assert.Equal(t, StreamFilter("unread"), FilterUnread)
}

func TestErrorCode_Constants(t *testing.T) {
	assert.Equal(t, "bad_request", ErrCodeBadRequest)
	assert.Equal(t, "unauthorized", ErrCodeUnauthorized)
	assert.Equal(t, "forbidden", ErrCodeForbidden)
	assert.Equal(t, "resource_not_found", ErrCodeNotFound)
	assert.Equal(t, "internal_error", ErrCodeInternalError)
	assert.Equal(t, "validation_failed", ErrCodeValidationFailed)
}

func TestStreamRequest_Defaults(t *testing.T) {
	// Default values
	req := StreamRequest{}

	assert.Equal(t, "", req.UserID)
	assert.Equal(t, StreamFilter(""), req.Filter)
	assert.Equal(t, 0, req.Limit)
	assert.Nil(t, req.Cursor)
}

func TestHealthResponse_JSON(t *testing.T) {
	// Arrange
	resp := HealthResponse{
		Status:    "ok",
		Timestamp: "2026-01-01T00:00:00Z",
	}

	// Act
	data, err := json.Marshal(resp)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"status":"ok"`)
	assert.Contains(t, string(data), `"timestamp":"2026-01-01T00:00:00Z"`)
}
