# API Documentation

## Base URL

```
http://localhost:8085/api/v1
```

## Authentication

The API uses JWT token-based authentication. Include the token in the Authorization header:

```
Authorization: Bearer <token>
```

Tokens are returned from the `/login` and `/register` endpoints.

## Endpoints

### User Registration

**Endpoint**: `POST /register`

**Description**: Creates a new user account and returns an authentication token.

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Validation**:
- email: required, must be valid email format
- password: required, minimum 6 characters
- first_name: optional
- last_name: optional

**Success Response** (201 Created):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "expires_at": 1705325400
}
```

**Error Responses**:

```json
// 400 Bad Request - Invalid input
{
  "error": "Invalid request"
}

// 400 Bad Request - User already exists
{
  "error": "user already exists"
}

// 500 Internal Server Error
{
  "error": "Internal server error"
}
```

**Example cURL**:

```bash
curl -X POST http://localhost:8085/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

---

### User Login

**Endpoint**: `POST /login`

**Description**: Authenticates a user and returns an authentication token.

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Validation**:
- email: required, must be valid email format
- password: required, minimum 6 characters

**Success Response** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "expires_at": 1705325400
}
```

**Error Responses**:

```json
// 400 Bad Request - Invalid input
{
  "error": "Invalid request"
}

// 401 Unauthorized - Invalid credentials
{
  "error": "Invalid credentials"
}

// 500 Internal Server Error
{
  "error": "Internal server error"
}
```

**Example cURL**:

```bash
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

---

### Get User by ID

**Endpoint**: `GET /users/{id}`

**Description**: Retrieves user information by user ID.

**Path Parameters**:
- id (string, required): The unique identifier of the user

**Success Response** (200 OK):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Responses**:

```json
// 404 Not Found
{
  "error": "User not found"
}

// 500 Internal Server Error
{
  "error": "Internal server error"
}
```

**Example cURL**:

```bash
curl -X GET http://localhost:8085/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

---

### Health Check

**Endpoint**: `GET /health`

**Description**: Returns the health status of the API server.

**Success Response** (200 OK):

```
OK
```

**Example cURL**:

```bash
curl -X GET http://localhost:8085/health
```

---

## Response Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid input or client error |
| 401 | Unauthorized - Authentication failed or token invalid |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common errors:

- **Invalid request**: Malformed JSON or missing required fields
- **user already exists**: Email is already registered
- **Invalid credentials**: Wrong email or password combination
- **User not found**: User ID doesn't exist in database
- **Internal server error**: Unexpected server error

---

## Rate Limiting

Currently not implemented. Planned for production deployment.

---

## Authentication Token Format

Tokens follow JWT format:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.{payload}.{signature}
```

Include in Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## Data Types

### User Object

```json
{
  "id": "string (UUID)",
  "email": "string (email format)",
  "first_name": "string (optional)",
  "last_name": "string (optional)",
  "created_at": "string (ISO 8601 timestamp)",
  "updated_at": "string (ISO 8601 timestamp)"
}
```

Note: `password_hash` and `deleted_at` are never returned in responses.

---

## CORS

Currently not configured. Add CORS middleware for production client-side consumption.

---

## API Versioning

Current API version: `v1`

Base path: `/api/v1/`

Future versions will use `/api/v2/`, etc.

---

## Pagination

Not yet implemented. Planned for list endpoints (e.g., GET /users).

---

## Filtering and Sorting

Not yet implemented. Planned for future releases.

---

## Examples

### Complete Registration Flow

```bash
# 1. Register new user
curl -X POST http://localhost:8085/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepassword123",
    "first_name": "Jane",
    "last_name": "Smith"
  }'

# Response includes token and user data
# Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.550e8400-e29b-41d4-a716-446655440000

# 2. Use token to get user data
curl -X GET http://localhost:8085/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.550e8400-e29b-41d4-a716-446655440000"
```

### Complete Login Flow

```bash
# 1. Login
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "securepassword123"
  }'

# 2. Use returned token for authenticated requests
curl -X GET http://localhost:8085/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

---

## API Documentation Interface

Interactive Swagger documentation available at:

```
http://localhost:8085/swagger/index.html
```

This provides an interactive interface to test all endpoints.
