package migrations

import (
	"github.com/Jason-Omondi/ecom/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateDB runs all database migrations using GORM
// GORM auto-generates correct SQL for MySQL, PostgreSQL, etc.
// Migrations ensure schema is consistent across environments
// Returns: error if any migration fails
// Why here: keeps schema changes version-controlled and reversible
func MigrateDB(db *gorm.DB, log *zap.Logger) error {
	log.Info("Running database migrations")

	// Auto-migrate creates/updates table schema based on struct tags
	// GORM generates appropriate SQL for configured database type
	// This is simpler than raw SQL migrations for most use cases
	migrations := []func(*gorm.DB) error{
		migrateUsersTable,
		// Add future migrations here:
		// migrateProductsTable,
		// migrateOrdersTable,
	}

	for _, migration := range migrations {
		if err := migration(db); err != nil {
			log.Error("Migration failed", zap.Error(err))
			return err
		}
	}

	log.Info("All migrations completed successfully")
	return nil
}

// migrateUsersTable creates/updates users table
// GORM reads User struct tags and creates appropriate schema
// Works identically for MySQL and PostgreSQL
func migrateUsersTable(db *gorm.DB) error {
	// AutoMigrate creates table if not exists, adds new columns, creates indexes
	// It does NOT drop existing columns (safe for production)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return nil
}

// For complex migrations, use raw SQL that works across databases:
// func migrateComplexSchema(db *gorm.DB) error {
// 	// Raw SQL here would need to handle MySQL vs PostgreSQL syntax
// 	// Use GORM AutoMigrate when possible for simplicity
// 	return nil
// }
