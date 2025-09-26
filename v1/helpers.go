package v1

import (
	"encoding/json"
	"net/http"
)

// RespondJSON writes a JSON response to the provided ResponseWriter.
// It sets the Content-Type header to "application/json" and writes the status code.
// Returns an error if JSON encoding fails.
//
// Example:
//
//	data := map[string]string{"message": "success"}
//	if err := RespondJSON(w, http.StatusOK, data); err != nil {
//	    log.Printf("Failed to write response: %v", err)
//	}
func RespondJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(payload)
}

// NewSuccessResponse creates a standard success response with the provided data and metadata.
// The response will have status "success" and no error field.
//
// Example:
//
//	data := []User{user1, user2}
//	meta := &Metadata{Page: 1, PerPage: 10, TotalItems: 25}
//	resp := NewSuccessResponse(data, meta)
func NewSuccessResponse(data any, meta *Metadata) APIResponse {
	return APIResponse{
		Status: StatusSuccess,
		Data:   data,
		Meta:   meta,
	}
}

// NewErrorResponse creates a standard error response with the specified HTTP status, error code, and message.
// The response will have status "error" and no data field.
//
// Example:
//
//	resp := NewErrorResponse(http.StatusInternalServerError, "DB_ERROR", "Database connection failed")
//	RespondJSON(w, http.StatusInternalServerError, resp)
func NewErrorResponse(httpStatus int, code string, message string) APIResponse {
	return APIResponse{
		Status: StatusError,
		Error: &APIError{
			HTTPStatus: httpStatus,
			Code:       code,
			Message:    message,
		},
	}
}

// NewValidationErrorResponse creates a validation error response with detailed field-level errors.
// It automatically sets HTTP status to 422 Unprocessable Entity and error code to "VALIDATION_FAILED".
//
// Example:
//
//	details := []ValidationError{
//	    {Field: "email", Message: "Invalid email format"},
//	    {Field: "password", Message: "Too short"},
//	}
//	resp := NewValidationErrorResponse(details)
//	RespondJSON(w, http.StatusUnprocessableEntity, resp)
func NewValidationErrorResponse(details []ValidationError) APIResponse {
	return APIResponse{
		Status: StatusError,
		Error: &APIError{
			HTTPStatus: http.StatusUnprocessableEntity,
			Code:       "VALIDATION_FAILED",
			Message:    "Data yang dikirim tidak valid.",
			Details:    details,
		},
	}
}
