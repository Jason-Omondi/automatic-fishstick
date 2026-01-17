package api

import (
	"net/http"
	"os"

	"github.com/Jason-Omondi/ecomgo/cmd/service/user"
	"github.com/Jason-Omondi/ecomgo/internal/config"
	"github.com/Jason-Omondi/ecomgo/internal/migrations"
	"github.com/Jason-Omondi/ecomgo/internal/repository"
	"github.com/gorilla/mux"
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

	// Serve swagger.json file
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

	// Read swagger.json from docs directory
	content, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		// Try from different path if running from subdirectory
		content, err = os.ReadFile("../../../docs/swagger.json")
		if err != nil {
			log.Error("Failed to read swagger.json", zap.Error(err))
			http.Error(w, `{"error":"Failed to load API definition"}`, http.StatusInternalServerError)
			return
		}
	}
	w.Write(content)
	})

	router.HandleFunc("/swagger/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>EcomGo API</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css" rel="stylesheet">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "BaseLayout"
            })
        }
    </script>
</body>
</html>
		`
		w.Write([]byte(html))
	})

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
