# libapi

A Go library for standardized API responses with built-in support for tracing, pagination, and structured error handling.

## Features

- **Standardized API Responses**: Consistent JSON response format across your API
- **Trace ID Support**: Built-in request tracing using Go contexts
- **Pagination Metadata**: Easy pagination support with metadata
- **Structured Error Handling**: Rich error responses with validation details
- **HTTP Response Helpers**: Convenient functions for writing JSON responses

## Installation

```bash
go get github.com/nofendian17/libapi/v1
```

## Quick Start

### Basic Usage

```go
package main

import (
    "net/http"
    "github.com/nofendian17/libapi/v1"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // Create success response
    data := map[string]string{"message": "Hello, World!"}
    resp := response.NewSuccessResponse(data, nil)

    // Send JSON response
    if err := response.RespondJSON(w, http.StatusOK, resp); err != nil {
        // Handle error
    }
}
```

### With Trace ID

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Add trace ID to context
    ctx := response.WithTraceID(r.Context(), "trace-123")

    // Create response with trace ID
    apiResp := response.NewAPIResponse(ctx)
    apiResp.Status = response.StatusSuccess
    apiResp.Data = map[string]string{"result": "success"}

    response.RespondJSON(w, http.StatusOK, apiResp)
}
```

### Error Responses

```go
func errorHandler(w http.ResponseWriter, r *http.Request) {
    // Create error response
    resp := response.NewErrorResponse(
        http.StatusInternalServerError,
        "INTERNAL_ERROR",
        "Something went wrong",
    )

    response.RespondJSON(w, http.StatusInternalServerError, resp)
}
```

### Validation Errors

```go
func validationHandler(w http.ResponseWriter, r *http.Request) {
    details := []response.ValidationError{
        {Field: "email", Message: "Invalid email format"},
        {Field: "password", Message: "Password must be at least 8 characters"},
    }

    resp := response.NewValidationErrorResponse(details)
    response.RespondJSON(w, http.StatusUnprocessableEntity, resp)
}

// Or use the enhanced function with custom error code and message
func customValidationHandler(w http.ResponseWriter, r *http.Request) {
    details := []response.ValidationError{
        {Field: "email", Message: "Invalid email format"},
        {Field: "password", Message: "Password must be at least 8 characters"},
    }

    resp := response.NewValidationErrorResponseWithCodeAndMessage(
        details, 
        "CUSTOM_VALIDATION_ERROR", 
        "Please check your input fields",
    )
    response.RespondJSON(w, http.StatusUnprocessableEntity, resp)
}
```

### Pagination

```go
func paginatedHandler(w http.ResponseWriter, r *http.Request) {
    data := []map[string]string{
        {"id": "1", "name": "Item 1"},
        {"id": "2", "name": "Item 2"},
    }

    meta := &response.Metadata{
        Page:       1,
        PerPage:    10,
        TotalItems: 25,
    }

    resp := response.NewSuccessResponse(data, meta)
    response.RespondJSON(w, http.StatusOK, resp)
}
```

## API Reference

### Types

#### APIResponse

The main response structure:

```go
type APIResponse struct {
    Status  string             `json:"status"`           // "success" or "error"
    TraceID string             `json:"trace_id,omitempty"` // Request trace ID
    Data    interface{}        `json:"data,omitempty"`   // Response data
    Error   *APIError          `json:"error,omitempty"`  // Error details (for error responses)
    Meta    *Metadata          `json:"meta,omitempty"`   // Pagination metadata
}
```

#### APIError

Error response structure:

```go
type APIError struct {
    HTTPStatus int                `json:"-"`                    // HTTP status code (not serialized)
    Code       string             `json:"code"`                 // Error code
    Message    string             `json:"message"`              // Error message
    Details    []ValidationError  `json:"details,omitempty"`    // Validation errors
}
```

#### ValidationError

Individual validation error:

```go
type ValidationError struct {
    Field   string `json:"field"`   // Field name
    Message string `json:"message"` // Error message
}
```

#### Metadata

Pagination metadata:

```go
type Metadata struct {
    Page       int `json:"page,omitempty"`        // Current page
    PerPage    int `json:"per_page,omitempty"`    // Items per page
    TotalItems int `json:"total_items,omitempty"` // Total number of items
}
```

### Functions

#### Context Helpers

```go
// Add trace ID to context
func WithTraceID(ctx context.Context, traceID string) context.Context

// Create new API response with trace ID from context
func NewAPIResponse(ctx context.Context) *APIResponse
```

#### Response Builders

```go
// Create success response
func NewSuccessResponse(data interface{}, meta *Metadata) APIResponse

// Create error response
func NewErrorResponse(httpStatus int, code string, message string) APIResponse

// Create validation error response
func NewValidationErrorResponse(details []ValidationError) APIResponse

// Create validation error response with custom error code and message
func NewValidationErrorResponseWithCodeAndMessage(details []ValidationError, code string, message string) APIResponse
```

#### HTTP Helpers

```go
// Write JSON response to ResponseWriter
func RespondJSON(w http.ResponseWriter, statusCode int, payload interface{}) error
```

## Response Format

### Success Response

```json
{
  "status": "success",
  "trace_id": "trace-123",
  "data": {
    "message": "Operation successful"
  },
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 25
  }
}
```

### Error Response

```json
{
  "status": "error",
  "trace_id": "trace-123",
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Something went wrong"
  }
}
```

### Validation Error Response

```json
{
  "status": "error",
  "trace_id": "trace-123",
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "The submitted data is invalid.",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      },
      {
        "field": "password",
        "message": "Password must be at least 8 characters"
      }
    ]
  }
}
```

## Versioning

The library uses semantic versioning. The `v1` package provides the current stable API.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
