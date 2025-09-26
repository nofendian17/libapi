package v1

import (
	"context"
)

// Status constants for API responses
const (
	StatusSuccess = "success" // StatusSuccess indicates a successful operation
	StatusError   = "error"   // StatusError indicates an error occurred
)

// traceIDKey is the context key for trace ID.
// It uses an unexported type to avoid key collisions.
type traceIDKey struct{}

// Metadata represents pagination metadata for API responses.
type Metadata struct {
	Page       int `json:"page,omitempty"`        // Current page number (1-based)
	PerPage    int `json:"per_page,omitempty"`    // Number of items per page
	TotalItems int `json:"total_items,omitempty"` // Total number of items across all pages
}

// APIResponse represents a standard API response structure.
// It provides a consistent format for all API responses.
type APIResponse struct {
	Status  string    `json:"status"`             // Response status: "success" or "error"
	TraceID string    `json:"trace_id,omitempty"` // Request trace ID for debugging
	Data    any       `json:"data,omitempty"`     // Response data (for success responses)
	Error   *APIError `json:"error,omitempty"`    // Error details (for error responses)
	Meta    *Metadata `json:"meta,omitempty"`     // Pagination metadata
}

// NewAPIResponse creates a new APIResponse with trace ID extracted from the context.
// If no trace ID is found in the context, TraceID will be empty.
//
// Example:
//
//	ctx := WithTraceID(context.Background(), "trace-123")
//	resp := NewAPIResponse(ctx)
//	// resp.TraceID will be "trace-123"
func NewAPIResponse(ctx context.Context) *APIResponse {
	resp := &APIResponse{}
	if ctx != nil {
		if traceID, ok := ctx.Value(traceIDKey{}).(string); ok {
			resp.TraceID = traceID
		}
	}
	return resp
}

// WithTraceID adds a trace ID to the context for request tracing.
// This allows requests to be tracked across service boundaries.
//
// Example:
//
//	ctx := WithTraceID(r.Context(), "unique-trace-id")
//	// Use ctx in subsequent operations
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}
