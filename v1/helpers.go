// Package v1 provides version 1 of the libapi response helpers.
// This package re-exports the response package functions for API stability.
package v1

import (
	"net/http"

	"github.com/nofendian17/libapi/response"
)

// RespondJSON writes a JSON response to the ResponseWriter.
// This is a convenience wrapper around response.RespondJSON.
func RespondJSON(w http.ResponseWriter, statusCode int, payload interface{}) error {
	return response.RespondJSON(w, statusCode, payload)
}

// NewSuccessResponse creates a standard success response.
// This is a convenience wrapper around response.NewSuccessResponse.
func NewSuccessResponse(data any, meta *response.Metadata) response.APIResponse {
	return response.NewSuccessResponse(data, meta)
}

// NewErrorResponse creates a standard error response.
// This is a convenience wrapper around response.NewErrorResponse.
func NewErrorResponse(httpStatus int, code string, message string) response.APIResponse {
	return response.NewErrorResponse(httpStatus, code, message)
}

// NewValidationErrorResponse creates a validation error response.
// This is a convenience wrapper around response.NewValidationErrorResponse.
func NewValidationErrorResponse(details []response.ValidationError) response.APIResponse {
	return response.NewValidationErrorResponse(details)
}
