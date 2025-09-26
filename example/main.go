package main

import (
	"log"
	"net/http"

	v1 "github.com/nofendian17/libapi/v1"
)

func main() {
	// Example HTTP server demonstrating libapi usage
	http.HandleFunc("/api/success", successHandler)
	http.HandleFunc("/api/error", errorHandler)
	http.HandleFunc("/api/validation", validationHandler)
	http.HandleFunc("/api/paginated", paginatedHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// successHandler demonstrates a basic success response
func successHandler(w http.ResponseWriter, r *http.Request) {
	// Add trace ID to context
	ctx := v1.WithTraceID(r.Context(), "trace-123")

	// Create response data
	data := map[string]any{
		"message": "Operation successful",
		"user": map[string]string{
			"id":   "123",
			"name": "John Doe",
		},
	}

	// Create API response with trace ID
	apiResp := v1.NewAPIResponse(ctx)
	apiResp.Status = v1.StatusSuccess
	apiResp.Data = data

	// Send response
	if err := v1.RespondJSON(w, http.StatusOK, apiResp); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// errorHandler demonstrates error response handling
func errorHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an error condition
	err := simulateError()

	if err != nil {
		// Create error response
		resp := v1.NewErrorResponse(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Database connection failed",
		)

		// Add trace ID if available
		if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
			ctx := v1.WithTraceID(r.Context(), traceID)
			resp.TraceID = v1.NewAPIResponse(ctx).TraceID
		}

		v1.RespondJSON(w, http.StatusInternalServerError, resp)
		return
	}

	// Success case
	resp := v1.NewSuccessResponse(map[string]string{"status": "ok"}, nil)
	v1.RespondJSON(w, http.StatusOK, resp)
}

// validationHandler demonstrates validation error responses
func validationHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate validation errors
	details := []v1.ValidationError{
		{
			Field:   "email",
			Message: "Must be a valid email address",
		},
		{
			Field:   "password",
			Message: "Must be at least 8 characters long",
		},
		{
			Field:   "age",
			Message: "Must be a positive number",
		},
	}

	resp := v1.NewValidationErrorResponse(details)

	// Add trace ID
	ctx := v1.WithTraceID(r.Context(), "validation-trace-456")
	resp.TraceID = v1.NewAPIResponse(ctx).TraceID

	v1.RespondJSON(w, http.StatusUnprocessableEntity, resp)
}

// paginatedHandler demonstrates paginated responses
func paginatedHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate paginated data
	items := []map[string]any{
		{"id": 1, "name": "Item 1", "description": "First item"},
		{"id": 2, "name": "Item 2", "description": "Second item"},
		{"id": 3, "name": "Item 3", "description": "Third item"},
	}

	// Pagination metadata
	meta := &v1.Metadata{
		Page:       1,
		PerPage:    10,
		TotalItems: 25,
	}

	resp := v1.NewSuccessResponse(items, meta)

	// Add trace ID
	ctx := v1.WithTraceID(r.Context(), "pagination-trace-789")
	resp.TraceID = v1.NewAPIResponse(ctx).TraceID

	v1.RespondJSON(w, http.StatusOK, resp)
}

// simulateError simulates an error condition
func simulateError() error {
	// Simulate some error logic
	return nil // Change to return an error to test error handling
}
