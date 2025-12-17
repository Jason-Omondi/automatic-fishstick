package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// NewMySQLConfig creates a MySQL DSN configuration
// Returns: *mysql.Config with connection parameters
// eg "user:password@tcp(host:port)/dbname?param=value"
func NewMySQLStorage(cf *mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
