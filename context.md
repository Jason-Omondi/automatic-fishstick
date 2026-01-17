# Project Context: EcomGo - E-Commerce API

EcomGo is a production-ready e-commerce backend API built with Go. It features a modular, layered architecture designed for maintainability, testability, and scalability.

## Key Features:
*   **User Management:** Registration and authentication with JWT tokens.
*   **Multi-database Support:** Seamless integration with both MySQL and PostgreSQL using GORM, with automatic schema migration.
*   **Structured Logging:** Utilizes Uber Zap for structured, JSON-formatted logs.
*   **API Documentation:** Self-generating Swagger/OpenAPI 2.0 documentation.
*   **Containerization:** Docker and Docker Compose support for simplified local development and deployment.
*   **Configuration:** Environment variable-based configuration with `.env` file support.

## Technology Stack:
*   **Language:** Go 1.24.3
*   **Framework:** Gorilla Mux (HTTP routing)
*   **ORM:** GORM with MySQL and PostgreSQL drivers
*   **Logging:** Uber Zap
*   **Environment Variables:** `joho/godotenv`
*   **API Documentation:** `swaggo/http-swagger` and `swaggo/swag`

## Architecture:
The application follows a layered architecture:
1.  **Handler Layer (`cmd/service/*/routes.go`):** Handles HTTP requests and responses, including JSON parsing and input validation.
2.  **Service Layer (`cmd/service/*/service.go`):** Contains the core business logic, such as authentication, password hashing, and token generation.
3.  **Repository Layer (`internal/repository/`):** Abstracts data access, providing CRUD operations and database-agnostic queries.
4.  **Model Layer (`internal/models/`):** Defines data structures (e.g., `User` model).
5.  **Database Layer (`internal/database/`):** Manages database connections using GORM.

## Design Patterns:
*   **Dependency Injection:** Components receive dependencies via constructors for loose coupling and testability.
*   **Repository Pattern:** Abstracts data access from business logic, allowing for easy database switching and testing.
*   **Service Layer:** Centralizes business logic, separating it from HTTP concerns.
*   **Configuration Management:** Centralized loading of configuration from environment variables.

## Current Functionality (API Endpoints - prefixed with `/api/v1`):
*   `POST /register`: Registers a new user.
*   `POST /login`: Authenticates a user and returns a token.
*   `GET /users/{id}`: Retrieves user details by ID.
*   `GET /health`: Checks server health.

## Future Improvements:
The project has a detailed roadmap for future improvements across several phases, including enhanced authentication (bcrypt, JWT refresh, rate limiting, email verification), Keycloak integration, comprehensive user management, product catalog, orders, advanced features (caching, API versioning), testing, CI/CD, production readiness, and further API enhancements.
