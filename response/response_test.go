package response

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAPIResponse(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected *APIResponse
	}{
		{
			name:     "nil context",
			ctx:      nil,
			expected: &APIResponse{},
		},
		{
			name:     "context without trace ID",
			ctx:      context.Background(),
			expected: &APIResponse{},
		},
		{
			name:     "context with trace ID",
			ctx:      WithTraceID(context.Background(), "test-trace-id"),
			expected: &APIResponse{TraceID: "test-trace-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAPIResponse(tt.ctx)
			if result.TraceID != tt.expected.TraceID {
				t.Errorf("NewAPIResponse() TraceID = %v, want %v", result.TraceID, tt.expected.TraceID)
			}
		})
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-123"

	newCtx := WithTraceID(ctx, traceID)

	if newCtx == ctx {
		t.Error("WithTraceID() should return a new context")
	}

	retrievedID, ok := newCtx.Value(traceIDKey{}).(string)
	if !ok {
		t.Error("WithTraceID() failed to store trace ID")
	}

	if retrievedID != traceID {
		t.Errorf("WithTraceID() stored trace ID = %v, want %v", retrievedID, traceID)
	}
}

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		payload    any
		expected   string
	}{
		{
			name:       "success response",
			statusCode: http.StatusOK,
			payload:    map[string]string{"message": "success"},
			expected:   `{"message":"success"}`,
		},
		{
			name:       "error response",
			statusCode: http.StatusBadRequest,
			payload:    map[string]string{"error": "bad request"},
			expected:   `{"error":"bad request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := RespondJSON(w, tt.statusCode, tt.payload)
			if err != nil {
				t.Errorf("RespondJSON() error = %v", err)
			}

			if w.Code != tt.statusCode {
				t.Errorf("RespondJSON() status code = %v, want %v", w.Code, tt.statusCode)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("RespondJSON() Content-Type = %v, want application/json", contentType)
			}

			var buf bytes.Buffer
			buf.ReadFrom(w.Body)
			body := buf.String()

			var expectedJSON, actualJSON interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedJSON); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(body), &actualJSON); err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !jsonEqual(expectedJSON, actualJSON) {
				t.Errorf("RespondJSON() body = %v, want %v", body, tt.expected)
			}
		})
	}
}

func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	meta := &Metadata{Page: 1, PerPage: 10, TotalItems: 100}

	response := NewSuccessResponse(data, meta)

	if response.Status != StatusSuccess {
		t.Errorf("NewSuccessResponse() Status = %v, want %v", response.Status, StatusSuccess)
	}

	if response.Data == nil {
		t.Error("NewSuccessResponse() Data should not be nil")
	}

	// Check data contents
	responseData, ok := response.Data.(map[string]string)
	if !ok {
		t.Errorf("NewSuccessResponse() Data type = %T, want map[string]string", response.Data)
	}

	if responseData["key"] != "value" {
		t.Errorf("NewSuccessResponse() Data[\"key\"] = %v, want value", responseData["key"])
	}

	if response.Meta != meta {
		t.Errorf("NewSuccessResponse() Meta = %v, want %v", response.Meta, meta)
	}

	if response.Error != nil {
		t.Error("NewSuccessResponse() Error should be nil")
	}
}

func TestNewErrorResponse(t *testing.T) {
	httpStatus := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "Internal server error"

	response := NewErrorResponse(httpStatus, code, message)

	if response.Status != StatusError {
		t.Errorf("NewErrorResponse() Status = %v, want %v", response.Status, StatusError)
	}

	if response.Error == nil {
		t.Fatal("NewErrorResponse() Error should not be nil")
	}

	if response.Error.HTTPStatus != httpStatus {
		t.Errorf("NewErrorResponse() Error.HTTPStatus = %v, want %v", response.Error.HTTPStatus, httpStatus)
	}

	if response.Error.Code != code {
		t.Errorf("NewErrorResponse() Error.Code = %v, want %v", response.Error.Code, code)
	}

	if response.Error.Message != message {
		t.Errorf("NewErrorResponse() Error.Message = %v, want %v", response.Error.Message, message)
	}

	if response.Data != nil {
		t.Error("NewErrorResponse() Data should be nil")
	}
}

func TestNewValidationErrorResponse(t *testing.T) {
	details := []ValidationError{
		{Field: "email", Message: "Invalid email format"},
		{Field: "password", Message: "Password too short"},
	}

	response := NewValidationErrorResponse(details)

	if response.Status != StatusError {
		t.Errorf("NewValidationErrorResponse() Status = %v, want %v", response.Status, StatusError)
	}

	if response.Error == nil {
		t.Fatal("NewValidationErrorResponse() Error should not be nil")
	}

	if response.Error.HTTPStatus != http.StatusUnprocessableEntity {
		t.Errorf("NewValidationErrorResponse() Error.HTTPStatus = %v, want %v", response.Error.HTTPStatus, http.StatusUnprocessableEntity)
	}

	if response.Error.Code != "VALIDATION_FAILED" {
		t.Errorf("NewValidationErrorResponse() Error.Code = %v, want VALIDATION_FAILED", response.Error.Code)
	}

	if len(response.Error.Details) != len(details) {
		t.Errorf("NewValidationErrorResponse() Error.Details length = %v, want %v", len(response.Error.Details), len(details))
	}

	for i, detail := range response.Error.Details {
		if detail != details[i] {
			t.Errorf("NewValidationErrorResponse() Error.Details[%d] = %v, want %v", i, detail, details[i])
		}
	}
}

// jsonEqual compares two JSON values for equality
func jsonEqual(a, b interface{}) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return bytes.Equal(aBytes, bBytes)
}
