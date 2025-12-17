package api

import (
	"net/http"

	"github.com/Jason-Omondi/ecom/cmd/service/user"
	"github.com/Jason-Omondi/ecom/internal/config"
	"github.com/Jason-Omondi/ecom/internal/migrations"
	"github.com/Jason-Omondi/ecom/internal/repository"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type APIServer struct {
	port   string
	db     *gorm.DB
	router *mux.Router
	log    *zap.Logger
	config *config.Config
}

func NewAPIServer(port string, db *gorm.DB, cfg *config.Config, log *zap.Logger) *APIServer {
	// create a single router instance and register health on it
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add Swagger documentation endpoint
	// Serves interactive API docs at /swagger/index.html
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return &APIServer{
		port:   port,
		db:     db,
		router: router,
		log:    log,
		config: cfg,
	}
}

func (s *APIServer) Run() {
	s.log.Info("Starting API server", zap.String("port", s.port), zap.String("db_type", s.config.Database.Type))

	// Run database migrations before starting server
	// Ensures schema is up-to-date before accepting requests
	if err := migrations.MigrateDB(s.db, s.log); err != nil {
		s.log.Fatal("Failed to run migrations", zap.Error(err))
	}

	if err := s.Start(); err != nil {
		s.log.Fatal("Failed to start server", zap.Error(err))
	}

	// Get database version
	sqlDB, _ := s.db.DB()
	var version string
	if err := sqlDB.QueryRow("SELECT VERSION()").Scan(&version); err == nil {
		s.log.Info("Connected to database", zap.String("version", version))
	}
}

func (s *APIServer) Start() error {
	// Initialize repositories - data access layer
	// Repository pattern abstracts database logic, making it testable and maintainable
	userRepo := repository.NewUserRepository(s.db, s.log)

	// Initialize services - business logic layer
	// Services contain core business logic and orchestrate between repositories and handlers
	// Pass config to service if needed (e.g., for Keycloak integration)
	userService := user.NewUserService(userRepo, s.log, s.config)

	// initialize subrouter for versioned API routes (/api/v1/...)
	subrouter := s.router.PathPrefix("/api/v1").Subrouter()

	// Initialize user handler and register routes
	// Handlers receive HTTP requests and delegate to services
	userHandler := user.NewHandler(userService, s.log)
	userHandler.RegisterRoutes(subrouter)

	s.log.Info("Listening on port", zap.String("port", s.port))

	return http.ListenAndServe(s.port, s.router)
}
