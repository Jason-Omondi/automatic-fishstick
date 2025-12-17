package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Jason-Omondi/ecom/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes database connection based on config
// Supports MySQL and PostgreSQL - switch via DB_TYPE environment variable
// Returns: *gorm.DB connection, or error if connection fails
// Why here: centralizes database initialization, allows easy switching between databases
func InitDatabase(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "mysql":
		// MySQL DSN format for GORM: username:password@protocol(address)/dbname?param=value
		// IMPORTANT: Must not have spaces, all characters matter
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)

		// Log masked DSN for debugging (never log real password)
		log.Info("Connecting to MySQL",
			zap.String("user", cfg.Database.User),
			zap.String("host", cfg.Database.Host),
			zap.String("port", cfg.Database.Port),
			zap.String("database", cfg.Database.Name),
		)

		dialector = mysql.Open(dsn)

	case "postgres":
		// PostgreSQL connection string format
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)

		log.Info("Connecting to PostgreSQL",
			zap.String("user", cfg.Database.User),
			zap.String("host", cfg.Database.Host),
			zap.String("port", cfg.Database.Port),
			zap.String("database", cfg.Database.Name),
		)

		dialector = postgres.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// Connect to database with GORM
	// GORM handles connection pooling, prepared statements, etc.
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &GormLogger{log: log},
	})
	if err != nil {
		log.Error("Failed to connect to database",
			zap.Error(err),
			zap.String("db_type", cfg.Database.Type),
			zap.String("host", cfg.Database.Host),
		)
		return nil, err
	}

	log.Info("Database connection established successfully")

	// Get underlying SQL database and set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to get database instance", zap.Error(err))
		return nil, err
	}

	// Configure connection pool for better performance
	// MaxOpenConns: maximum number of open connections to the database
	// MaxIdleConns: maximum number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// maskPassword hides password in logs for security
// Returns: DSN string with password masked
// func maskPassword(dsn string) string {
// 	// Simple mask - replace password value with ***
// 	// In production, use more sophisticated masking
// 	for i, c := range dsn {
// 		if c == ':' {
// 			// Find next @ to extract password section
// 			for j := i + 1; j < len(dsn); j++ {
// 				if dsn[j] == '@' {
// 					return dsn[:i+1] + "***" + dsn[j:]
// 				}
// 			}
// 		}
// 	}
// 	return dsn
// }

// GormLogger implements GORM logger interface for structured logging
// Integrates GORM logs with Zap structured logger
type GormLogger struct {
	log *zap.Logger
}

func (gl *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return gl
}

func (gl *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	gl.log.Info(msg)
}

func (gl *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	gl.log.Warn(msg)
}

func (gl *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	gl.log.Error(msg)
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
}
