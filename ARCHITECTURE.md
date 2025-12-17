# Architecture Documentation

## System Architecture

EcomGo follows a layered architecture pattern for maintainability, testability, and scalability.

```
┌─────────────────────────────────────────────────────────┐
│                     HTTP Clients                         │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              HTTP Handler Layer                          │
│  (routes.go - Request/Response marshaling)              │
│  - Parse JSON requests                                  │
│  - Validate input                                       │
│  - Return HTTP responses                                │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Service Layer                              │
│  (service.go - Business Logic)                          │
│  - Authentication logic                                 │
│  - Password hashing/verification                        │
│  - Token generation                                     │
│  - Business rules validation                            │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Repository Layer                           │
│  (repository/*.go - Data Access)                        │
│  - Database queries abstraction                         │
│  - CRUD operations                                      │
│  - Database-agnostic queries                            │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│              Database Layer                             │
│  (GORM ORM - Query Generation)                          │
│  - SQL generation for MySQL/PostgreSQL                  │
│  - Connection pooling                                   │
│  - Transaction management                              │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│          Database (MySQL or PostgreSQL)                 │
└─────────────────────────────────────────────────────────┘
```

## Design Patterns

### 1. Dependency Injection

All components receive dependencies through constructors:

```go
// Handler receives service and logger
func NewHandler(service *UserService, log *zap.Logger) *Handler {
    return &Handler{
        service: service,
        log:     log,
    }
}
```

**Benefits**: Loose coupling, easy testing, flexible configuration

### 2. Repository Pattern

Data access is abstracted through repositories:

```go
// Service calls repository, not database directly
user, err := r.userRepo.GetUserByEmail(ctx, email)
```

**Benefits**: Database-agnostic, testable with mock repositories, easy migration between databases

### 3. Service Layer

Business logic is centralized in services:

```go
// Handler delegates to service
authResp, err := h.service.Register(context.Background(), &req)
```

**Benefits**: Reusable business logic, separated from HTTP concerns, easier testing

### 4. Configuration Management

Configuration is loaded once at startup and passed to all components:

```go
// Single config instance used throughout application
cfg, err := config.LoadConfig()
apiServer := api.NewAPIServer(":"+cfg.Server.Port, db, cfg, appLogger)
```

**Benefits**: Consistent configuration, environment-specific settings, no magic strings

## Database Design

### User Table

```sql
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

**Features**:
- UUID primary key (CHAR(36) for compatibility)
- Unique email constraint
- Soft deletes (deleted_at)
- Automatic timestamp management

### Multi-Database Support

GORM handles database differences automatically:

**MySQL**:
```go
dialector = mysql.Open("user:pass@tcp(host:3306)/dbname")
```

**PostgreSQL**:
```go
dialector = postgres.Open("host=... user=... password=...")
```

Same code works for both - GORM generates appropriate SQL.

## Configuration Flow

```
1. LoadConfig() reads .env file
2. Validates required environment variables
3. Creates Config struct with nested types
4. Config passed to main components:
   - Database initialization
   - API server setup
   - Service instantiation
5. Components access config values as needed
```

## Data Flow - User Registration

```
POST /api/v1/register
       ↓
Handler.handleRegister()
       ↓
Validates JSON request
       ↓
Service.Register()
       ↓
Check user doesn't exist
       ↓
Hash password
       ↓
Repository.CreateUser()
       ↓
Database INSERT
       ↓
Generate token
       ↓
Return AuthResponse (200)
```

## Logging Strategy

Structured logging with contextual information:

```go
h.log.Info("User registered successfully",
    zap.String("email", user.Email),
    zap.String("id", user.ID),
)
```

**Output**: JSON format for log aggregation systems

```json
{
  "level": "info",
  "ts": 1234567890.123,
  "msg": "User registered successfully",
  "email": "user@example.com",
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Error Handling

### HTTP Errors

```go
// Return appropriate HTTP status codes
http.Error(w, "Invalid request", http.StatusBadRequest)    // 400
http.Error(w, "User not found", http.StatusNotFound)       // 404
http.Error(w, "Invalid credentials", http.StatusUnauthorized) // 401
```

### Error Logging

```go
// Log context with errors
h.log.Error("Registration failed",
    zap.String("email", req.Email),
    zap.Error(err),
)
```

## Security Architecture

### Password Storage

Current: SHA256 hashing
Planned: bcrypt for production

```go
// Current (demo)
hash := sha256.Sum256([]byte(password))

// Production
bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

### Token Management

Current: Simple token format
Planned: Real JWT or Keycloak OAuth2

```go
// Current (demo)
token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." + userID

// Planned
token := generateJWT(userID, config.JWT.Secret, config.JWT.ExpiresIn)
```

### Configuration Security

- No hardcoded credentials
- Environment variables for all secrets
- .env file in .gitignore
- .env.example with placeholder values

## Scalability Considerations

### Database Connection Pooling

```go
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

Configurable pool size for different workloads.

### Stateless Design

- No session state in application
- Token-based authentication
- Can scale horizontally with load balancer

### Context Usage

All database operations use context for timeout/cancellation:

```go
db.WithContext(ctx).Create(user)
```

## Testing Strategy

### Unit Tests

Test services with mock repositories:

```go
mockRepo := &MockUserRepository{}
service := user.NewUserService(mockRepo, log, cfg)
resp, err := service.Register(context.Background(), req)
```

### Integration Tests

Test with real database in Docker:

```bash
docker-compose -f docker-compose.test.yml up
go test ./...
```

## Deployment

### Docker Deployment

```bash
docker build -t ecomgo:latest .
docker run -p 8085:8085 --env-file .env ecomgo:latest
```

### Environment-Specific Configs

```bash
# Development
.env.development

# Production
.env.production

# Load with
source .env.production
go run cmd/main.go
```

## Monitoring and Observability

### Health Check Endpoint

```
GET /health -> 200 OK
```

### Structured Logs

All logs output as JSON for easy parsing.

### Database Monitoring

Connection pool metrics available via GORM:

```go
sqlDB, _ := db.DB()
sqlDB.Stats()  // Get connection pool stats
```

---

## Future Improvements

### Phase 1: Authentication Enhancement

- [ ] Implement bcrypt for password hashing
- [ ] Add JWT token generation and validation
- [ ] Implement token refresh mechanism
- [ ] Add rate limiting on auth endpoints
- [ ] Add email verification on registration

### Phase 2: Keycloak Integration

- [ ] Integrate Keycloak OAuth2/OpenID Connect
- [ ] Replace simple tokens with Keycloak tokens
- [ ] Use config.Keycloak settings throughout
- [ ] Implement single sign-on (SSO)
- [ ] Add social login providers

### Phase 3: User Management

- [ ] Add user profile endpoints (PUT /users/{id})
- [ ] Add user deletion (soft delete)
- [ ] Add password reset flow
- [ ] Add user role management
- [ ] Add user preferences storage

### Phase 4: Product Catalog

- [ ] Create Product model and migration
- [ ] Add ProductRepository and ProductService
- [ ] Implement CRUD endpoints for products
- [ ] Add product search and filtering
- [ ] Add product categories

### Phase 5: Orders

- [ ] Create Order and OrderItem models
- [ ] Implement shopping cart
- [ ] Add order creation and management
- [ ] Add order history endpoints
- [ ] Implement payment integration

### Phase 6: Advanced Features

- [ ] Add database migrations versioning (golang-migrate)
- [ ] Implement caching layer (Redis)
- [ ] Add request/response validation middleware
- [ ] Add API versioning strategy
- [ ] Implement GraphQL API alongside REST

### Phase 7: Testing & CI/CD

- [ ] Add comprehensive unit tests
- [ ] Add integration tests with Docker
- [ ] Setup GitHub Actions CI/CD pipeline
- [ ] Add code coverage reporting
- [ ] Add automated security scanning

### Phase 8: Production Readiness

- [ ] Add request tracing (Jaeger)
- [ ] Add performance monitoring (Prometheus)
- [ ] Setup application metrics collection
- [ ] Add graceful shutdown handling
- [ ] Implement request/response compression
- [ ] Add CORS configuration
- [ ] Setup SSL/TLS certificates
- [ ] Add database backup strategy

### Phase 9: API Enhancements

- [ ] Add pagination to list endpoints
- [ ] Add filtering and sorting capabilities
- [ ] Implement partial response selection
- [ ] Add webhook support for events
- [ ] Add audit logging for critical operations

### Phase 10: Documentation & DevOps

- [ ] Add OpenAPI 3.0 specification
- [ ] Create development setup guide
- [ ] Add troubleshooting documentation
- [ ] Create Kubernetes deployment configs
- [ ] Add Helm charts for K8s deployment

## Code Quality

### Standards

- Follow Go style guide: https://golang.org/doc/effective_go
- Use golangci-lint for code analysis
- Maintain >80% test coverage
- Document all exported functions

### Review Checklist

- No hardcoded credentials
- Error handling for all operations
- Logs include context (IDs, emails, etc.)
- Thread-safe database operations
- Input validation on all endpoints

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
