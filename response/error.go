package response

// ValidationError represents a validation error for a specific field.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIError represents an API error with structured information.
type APIError struct {
	HTTPStatus int               `json:"-"`                 // HTTP status code (not included in JSON)
	Code       string            `json:"code"`              // Error code for programmatic handling
	Message    string            `json:"message"`           // Human-readable error message
	Details    []ValidationError `json:"details,omitempty"` // Detailed validation errors
}
