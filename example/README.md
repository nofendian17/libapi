# Example Usage

This directory contains example code demonstrating how to use the libapi library.

## Running the Example

```bash
cd example
go run main.go
```

The server will start on `http://localhost:8080` with the following endpoints:

- `GET /api/success` - Basic success response
- `GET /api/error` - Error response example
- `GET /api/validation` - Validation error response
- `GET /api/paginated` - Paginated response with metadata

## Example Requests

```bash
# Success response
curl http://localhost:8080/api/success

# Error response
curl http://localhost:8080/api/error

# Validation errors
curl http://localhost:8080/api/validation

# Paginated response
curl http://localhost:8080/api/paginated
```

## Response Examples

### Success Response
```json
{
  "status": "success",
  "trace_id": "trace-123",
  "data": {
    "message": "Operation successful",
    "user": {
      "id": "123",
      "name": "John Doe"
    }
  }
}
```

### Error Response
```json
{
  "status": "error",
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "Database connection failed"
  }
}
```

### Validation Error Response
```json
{
  "status": "error",
  "trace_id": "validation-trace-456",
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Data yang dikirim tidak valid.",
    "details": [
      {
        "field": "email",
        "message": "Must be a valid email address"
      }
    ]
  }
}
```

### Paginated Response
```json
{
  "status": "success",
  "trace_id": "pagination-trace-789",
  "data": [
    {"id": 1, "name": "Item 1", "description": "First item"}
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total_items": 25
  }
}
```