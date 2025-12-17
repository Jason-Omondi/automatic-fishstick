package main

import (
	"log"

	"github.com/Jason-Omondi/ecom/cmd/api"
	"github.com/Jason-Omondi/ecom/internal/config"
	"github.com/Jason-Omondi/ecom/internal/database"
	"github.com/Jason-Omondi/ecom/internal/logger"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	appLogger, _ := logger.NewLogger()

	// Load configuration ONCE at application startup
	// This is the single source of truth for all config throughout the app
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Validate database configuration
	if cfg.Database.Type == "" {
		log.Fatal("DB_TYPE environment variable not set")
	}
	if cfg.Database.User == "" {
		log.Fatal("DB_USER environment variable not set")
	}
	if cfg.Database.Name == "" {
		log.Fatal("DB_NAME environment variable not set")
	}
	if cfg.Database.Host == "" {
		log.Fatal("DB_HOST environment variable not set")
	}
	if cfg.Database.Password == "" {
		log.Fatal("DB_PASSWORD environment variable not set - this is required")
	}

	appLogger.Info("Configuration loaded successfully",
		zap.String("db_type", cfg.Database.Type),
		zap.String("db_user", cfg.Database.User),
		zap.String("db_host", cfg.Database.Host),
		zap.String("db_port", cfg.Database.Port),
		zap.String("db_name", cfg.Database.Name),
		zap.Bool("has_password", cfg.Database.Password != ""),
	)

	// Initialize database connection with GORM
	// Automatically uses MySQL or PostgreSQL based on DB_TYPE config
	db, err := database.InitDatabase(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("Failed to initialize database",
			zap.Error(err),
			zap.String("db_type", cfg.Database.Type),
			zap.String("db_host", cfg.Database.Host),
			zap.String("db_user", cfg.Database.User),
		)
	}

	appLogger.Info("Database connected successfully")

	// Pass config and GORM db to APIServer
	apiServer := api.NewAPIServer(":"+cfg.Server.Port, db, cfg, appLogger)
	apiServer.Run()
}
