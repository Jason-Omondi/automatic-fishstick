# Development Setup Guide

## Prerequisites

- Go 1.24.3 or higher
- Git
- Docker and Docker Compose (optional, for containerized setup)
- MySQL 8.0 or PostgreSQL 16 (if not using Docker)

## Environment Setup

### Option 1: Docker Setup (Recommended)

Fastest way to get started with all services running.

```bash
# Clone repository
git clone https://github.com/your-org/ecomgo.git
cd ecomgo

# Start all services
docker-compose up --build

# In another terminal, run migrations (optional - auto-runs on app start)
docker-compose exec app ./scripts/migrate.sh up

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

**Services running**:
- Application: http://localhost:8085
- MySQL: localhost:3306
- Database GUI (Adminer): http://localhost:8080
- Swagger Docs: http://localhost:8085/swagger/index.html

### Option 2: Local Setup

Install dependencies locally for development.

#### Install Go

```bash
# macOS with Homebrew
brew install go

# Or download from https://golang.org/dl/
```

#### Install MySQL

```bash
# macOS with Homebrew
brew install mysql
brew services start mysql

# Or use Docker for just MySQL
docker run -d \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=password \
  -e MYSQL_DATABASE=ecomgo \
  mysql:8.0
```

#### Setup Environment

```bash
# Copy example env file
cp .env.example .env

# Edit .env with your settings
# For local MySQL on same machine:
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=password
DB_NAME=ecomgo
```

#### Install Dependencies

```bash
# Download Go modules
go mod download

# Verify modules
go mod verify
```

#### Run Application

```bash
# Run migrations
./scripts/migrate.sh up

# Start application
go run cmd/main.go

# Application running on http://localhost:8085
```

## Project Structure Overview

```
ecomgo/
├── cmd/                              # Application entry points
│   ├── api/
│   │   └── api.go                   # Server initialization
│   ├── service/user/
│   │   ├── service.go               # User business logic
│   │   └── routes.go                # HTTP handlers
│   └── main.go                      # Application entry
│
├── internal/                         # Internal packages (not exported)
│   ├── config/
│   │   └── config.go                # Configuration loading
│   ├── database/
│   │   └── database.go              # Database initialization
│   ├── docs/
│   │   └── docs.go                  # Swagger documentation
│   ├── logger/
│   │   └── logger.go                # Logging setup
│   ├── migrations/
│   │   └── migrations.go            # Database schema
│   ├── models/
│   │   └── user.go                  # Data models
│   └── repository/
│       └── user.go                  # Data access layer
│
├── scripts/
│   └── migrate.sh                   # Database migration CLI
│
├── docker-compose.yml               # Services definition
├── Dockerfile                       # Application container
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── .env.example                    # Example configuration
├── .gitignore                      # Git ignore rules
│
└── docs/
    ├── README.md                   # Project overview
    ├── ARCHITECTURE.md             # System design
    ├── API_DOCUMENTATION.md        # API reference
    ├── CONTRIBUTING.md             # Contribution guidelines
    └── SETUP.md                    # This file
```

## Common Development Tasks

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (install first: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
golangci-lint run

# Run vet
go vet ./...
```

### Database Operations

```bash
# Run migrations up
./scripts/migrate.sh up

# Run migrations down
./scripts/migrate.sh down

# Create new migration (if using golang-migrate)
migrate create -ext sql -dir internal/migrations/sql add_products_table
```

### API Testing

#### Using cURL

```bash
# Register user
curl -X POST http://localhost:8085/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login
curl -X POST http://localhost:8085/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Get user (replace token and id)
curl -X GET http://localhost:8085/api/v1/users/{user-id} \
  -H "Authorization: Bearer {token}"
```

#### Using Swagger UI

1. Open http://localhost:8085/swagger/index.html
2. Click on endpoint
3. Click "Try it out"
4. Enter parameters
5. Click "Execute"

#### Using Postman

1. Import API endpoints to Postman
2. Set variables: `base_url=http://localhost:8085`, `token={from login}`
3. Test endpoints

### Building for Production

```bash
# Build binary
go build -o ecomgo ./cmd/main.go

# Build with version info
VERSION=$(git describe --tags)
go build -ldflags "-X main.Version=$VERSION" -o ecomgo ./cmd/main.go

# Build optimized for production
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ecomgo ./cmd/main.go
```

### Docker Operations

```bash
# Build application image
docker build -t ecomgo:latest .

# Build and tag for registry
docker build -t myregistry.azurecr.io/ecomgo:v1.0.0 .

# Push to registry
docker push myregistry.azurecr.io/ecomgo:v1.0.0

# Run container
docker run -p 8085:8085 --env-file .env ecomgo:latest

# Run container with docker-compose
docker-compose up
docker-compose down
docker-compose logs -f
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8085
lsof -i :8085

# Kill process
kill -9 <PID>

# Or use different port
SERVER_PORT=8086 go run cmd/main.go
```

### Database Connection Failed

```bash
# Check MySQL is running
brew services list  # macOS
systemctl status mysql  # Linux

# Check credentials in .env
cat .env | grep DB_

# Test connection
mysql -u root -p -h localhost ecomgo
```

### Module Download Issues

```bash
# Clear module cache
go clean -modcache

# Download modules again
go mod download

# Update modules
go get -u ./...
```

### Docker Build Issues

```bash
# Clean docker cache
docker system prune

# Rebuild without cache
docker-compose build --no-cache

# Check Docker daemon
docker ps  # Should not error
```

## Development Workflow

1. Create feature branch: `git checkout -b feature/my-feature`
2. Make changes following code standards
3. Run tests: `go test ./...`
4. Format code: `go fmt ./...`
5. Run linter: `golangci-lint run`
6. Commit with clear message: `git commit -m "feat: add new feature"`
7. Push to fork: `git push origin feature/my-feature`
8. Open pull request with description
9. Address review comments
10. Merge to main branch

## IDE Setup

### VS Code

Install extensions:
- Go (golang.go)
- Docker (ms-azuretools.vscode-docker)
- REST Client (humao.rest-client)

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to Container",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${fileDirname}",
      "cwd": "${workspaceFolder}"
    }
  ]
}
```

### GoLand / IntelliJ IDEA

- File > Open > Select ecomgo directory
- Go > Go Modules > Enable Go Modules
- Run > Edit Configurations > Add Go Application

## Performance Testing

```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:8085/health

# Using wrk
wrk -t4 -c100 -d30s http://localhost:8085/health

# Using hey
hey -z 30s -c 100 http://localhost:8085/health
```

## Security Checklist

- [ ] No hardcoded credentials in code
- [ ] .env file in .gitignore
- [ ] Environment variables for all secrets
- [ ] No passwords in logs or error messages
- [ ] Input validation on all endpoints
- [ ] SQL injection protection (GORM handles this)
- [ ] HTTPS in production (TODO)
- [ ] Rate limiting on auth endpoints (TODO)
- [ ] CORS configured properly (TODO)

## Documentation

- Update README.md for major changes
- Update API_DOCUMENTATION.md for new endpoints
- Update ARCHITECTURE.md for design changes
- Add comments to complex code
- Document configuration options

## Resources

- Go Documentation: https://golang.org/doc/
- GORM Guides: https://gorm.io/docs/
- Docker Documentation: https://docs.docker.com/
- Swagger/OpenAPI: https://swagger.io/
- MySQL: https://dev.mysql.com/doc/
- PostgreSQL: https://www.postgresql.org/docs/

## Getting Help

- Check existing issues
- Search discussions
- Ask in issue with "question" label
- Contact maintainers

Happy coding!
