package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
// Loaded from environment variables for security (no hardcoded secrets)
type Config struct {
	Database Database
	Server   Server
	Keycloak Keycloak
}

type Database struct {
	Type     string // mysql or postgres
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

type Server struct {
	Port string
}

type Keycloak struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
}

// LoadConfig reads configuration from .env file and environment variables
// Searches for .env in current directory and parent directories (up to project root)
// Returns: Config struct with all settings, or error if required vars missing
// Why here: centralizes config loading, prevents scattered getenv calls throughout app
func LoadConfig() (*Config, error) {
	// Try to load .env from current directory and parent directories
	// This allows running from any subdirectory (e.g., ./cmd)
	if err := loadDotEnvFromRoot(); err != nil {
		// Don't fail if .env not found; env vars can be used instead
		// Only log for debugging
		_ = err
	}

	cfg := &Config{
		Database: Database{
			Type:     strings.TrimSpace(getEnv("DB_TYPE", "mysql")), // Default to MySQL
			User:     strings.TrimSpace(getEnv("DB_USER", "root")),
			Password: strings.TrimSpace(getEnv("DB_PASSWORD", "")),
			Name:     strings.TrimSpace(getEnv("DB_NAME", "ecomgo")),
			Host:     strings.TrimSpace(getEnv("DB_HOST", "localhost")),
			Port:     strings.TrimSpace(getEnv("DB_PORT", "3306")),
			SSLMode:  strings.TrimSpace(getEnv("DB_SSLMODE", "disable")),
		},
		Server: Server{
			Port: strings.TrimSpace(getEnv("SERVER_PORT", "8085")),
		},
		Keycloak: Keycloak{
			URL:          strings.TrimSpace(getEnv("KEYCLOAK_URL", "http://localhost:8080")),
			Realm:        strings.TrimSpace(getEnv("KEYCLOAK_REALM", "master")),
			ClientID:     strings.TrimSpace(getEnv("KEYCLOAK_CLIENT_ID", "ecomgo")),
			ClientSecret: strings.TrimSpace(getEnv("KEYCLOAK_CLIENT_SECRET", "")),
		},
	}

	// Validate database configuration
	if cfg.Database.Type != "mysql" && cfg.Database.Type != "postgres" {
		return nil, fmt.Errorf("invalid DB_TYPE: %s (must be 'mysql' or 'postgres')", cfg.Database.Type)
	}

	return cfg, nil
}

// loadDotEnvFromRoot searches for and loads .env file from project root
// Walks up directory tree until .env is found or filesystem root is reached
// Handles running from any subdirectory (cmd/, internal/, etc.)
func loadDotEnvFromRoot() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Keep walking up until we find .env or hit filesystem root
	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// .env found, load it
			return godotenv.Load(envPath)
		}

		// Move to parent directory
		parentDir := filepath.Dir(currentDir)

		// Stop if we've reached filesystem root (parentDir == currentDir means root)
		if parentDir == currentDir {
			break
		}

		currentDir = parentDir
	}

	// .env not found, but that's okay - env vars can be used
	return nil
}

// GetDSN builds database connection string from config
// Returns: DSN string formatted for specific database type
func (db Database) GetDSN() string {
	switch db.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.Name,
		)
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			db.Host,
			db.Port,
			db.User,
			db.Password,
			db.Name,
			db.SSLMode,
		)
	default:
		return ""
	}
}

// getEnv retrieves environment variable with fallback default
// Returns: env var value if set, otherwise default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
