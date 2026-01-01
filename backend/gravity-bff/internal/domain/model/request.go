package model

// StreamFilter represents filter options for the stream endpoint.
type StreamFilter string

const (
	FilterAll    StreamFilter = "all"
	FilterHigh   StreamFilter = "high"
	FilterUnread StreamFilter = "unread"
)

// StreamRequest represents the query parameters for fetching the stream.
type StreamRequest struct {
	UserID string       `json:"-"`              // Extracted from auth token
	Filter StreamFilter `json:"filter"`         // all, high, unread
	Limit  int          `json:"limit"`          // Max items to return (default: 20, max: 100)
	Cursor *string      `json:"cursor"`         // Pagination cursor
}

// StreamResponse represents the paginated response for the stream endpoint.
type StreamResponse struct {
	Data       []PriorityItem `json:"data"`
	NextCursor *string        `json:"nextCursor"`
}

// StreamItemRequest represents the request for a single stream item.
type StreamItemRequest struct {
	UserID string `json:"-"`      // Extracted from auth token
	ItemID string `json:"itemId"` // The item ID from URL path
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response.
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// Common error codes
const (
	ErrCodeBadRequest       = "bad_request"
	ErrCodeUnauthorized     = "unauthorized"
	ErrCodeForbidden        = "forbidden"
	ErrCodeNotFound         = "resource_not_found"
	ErrCodeInternalError    = "internal_error"
	ErrCodeValidationFailed = "validation_failed"
)
