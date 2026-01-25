package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(connectionString string) error {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// HealthCheck checks if the database connection is alive
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.Ping()
}
