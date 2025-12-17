# EcomGo - E-Commerce API

A production-ready e-commerce backend API built with Go, featuring multi-database support (MySQL/PostgreSQL), structured logging, and comprehensive authentication.

## Key Features

- User registration and authentication with JWT tokens
- Multi-database support (MySQL and PostgreSQL) with automatic schema migration
- Structured logging with Zap logger
- Swagger/OpenAPI documentation
- Docker and Docker Compose support for easy local development
- Database connection pooling and performance optimization
- Modular architecture following clean code principles

## Technology Stack

- **Language**: Go 1.24.3
- **Framework**: Gorilla Mux (HTTP routing)
- **Database**: GORM (ORM) with MySQL and PostgreSQL drivers
- **Logging**: Uber Zap (structured logging)
- **API Documentation**: Swagger/OpenAPI 2.0
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.24.3 or higher
- Docker and Docker Compose (optional, for containerized setup)
- MySQL 8.0 or PostgreSQL 16 (if running outside Docker)

### Setup with Docker (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd ecomgo

# Start services (app + MySQL + Adminer)
docker-compose up --build

# Application runs on http://localhost:8085
# Database GUI available at http://localhost:8080
```

### Manual Setup

```bash
# Set up environment
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies
go mod download

# Run migrations
./scripts/migrate.sh up

# Start the application
go run cmd/main.go
```

## Environment Configuration

Copy `.env.example` to `.env` and configure:

```env
# Database Type: mysql or postgres
DB_TYPE=mysql

# Database Connection
DB_USER=root
DB_PASSWORD=your_secure_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ecomgo

# Server
SERVER_PORT=8085

# Keycloak (for future OAuth2 integration)
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=ecomgo
KEYCLOAK_CLIENT_SECRET=your_secret
```

## API Endpoints

All endpoints are prefixed with `/api/v1`

### Authentication
- `POST /register` - Register new user
- `POST /login` - Authenticate and get token

### Users
- `GET /users/{id}` - Retrieve user by ID

### Health Check
- `GET /health` - Server health status

For detailed API documentation, see [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

## Project Structure

```
ecomgo/
├── cmd/
│   ├── api/              # API server initialization
│   ├── service/user/     # User service business logic
│   └── main.go           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database initialization
│   ├── docs/             # Swagger documentation
│   ├── logger/           # Structured logging
│   ├── migrations/       # Database schema migrations
│   ├── models/           # Data models
│   └── repository/       # Data access layer
├── scripts/
│   └── migrate.sh        # Database migration script
├── docker-compose.yml    # Docker services definition
├── Dockerfile            # Application container
├── .env.example          # Example environment variables
└── README.md             # This file
```

## Architecture

The application follows a layered architecture:

1. **Handler Layer** (`cmd/service/*/routes.go`) - HTTP request/response handling
2. **Service Layer** (`cmd/service/*/service.go`) - Business logic
3. **Repository Layer** (`internal/repository/`) - Data access abstraction
4. **Model Layer** (`internal/models/`) - Data structures
5. **Database Layer** (`internal/database/`) - Connection management

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture documentation.

## Database Migrations

Migrations are automatically run on application startup. To manually run migrations:

```bash
# Run up migrations
./scripts/migrate.sh up

# Run down migrations
./scripts/migrate.sh down
```

Migrations are defined in `internal/migrations/migrations.go` and use GORM's AutoMigrate for database-agnostic schema management.

## Logging

The application uses structured logging with Zap. Logs are output as JSON for easy parsing by log aggregation systems.

```go
// Example log output
{"level":"info","ts":1234567890,"msg":"User registered successfully","email":"user@example.com"}
```

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -o ecomgo ./cmd/main.go
```

## Database Management

### Using Adminer (Web UI)

1. Start Docker Compose: `docker-compose up`
2. Open http://localhost:8080
3. Use these credentials:
   - System: MySQL
   - Server: mysql
   - Username: root
   - Password: password
   - Database: ecomgo

### Switching Databases

To switch from MySQL to PostgreSQL:

1. Update `.env`:
   ```env
   DB_TYPE=postgres
   DB_PORT=5432
   DB_USER=postgres
   ```

2. Uncomment PostgreSQL service in `docker-compose.yml`
3. Restart services: `docker-compose down && docker-compose up`

GORM automatically generates correct SQL for the selected database.

## API Documentation

Interactive Swagger documentation is available at `http://localhost:8085/swagger/index.html`

## Security Considerations

- Passwords are hashed before storage (currently SHA256, upgrade to bcrypt in production)
- Environment variables are used for all sensitive configuration
- Sensitive data (passwords) is excluded from JSON responses
- Database credentials are never logged
- Use HTTPS in production
- Implement rate limiting for authentication endpoints
- Use strong, unique database passwords

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for contribution guidelines.

## Future Improvements

See [ARCHITECTURE.md](./ARCHITECTURE.md#future-improvements) for planned enhancements.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.

You are free to use, modify, and distribute this software under the terms of the Apache License 2.0.

## Support

For issues and questions, please create an issue on the GitHub repository.
