// cmd/migrate/main.go
package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Javlopez/go-api/pkg/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Initialize DB
	db := database.New(database.NewConfig())
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Run migrations
	if err := migrateDB(dbConn); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	os.Exit(0)
}

func migrateDB(dbConn *sqlx.DB) error {

	// Get raw database connection
	sqlDB := dbConn.DB

	// Create postgres driver for migrations
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	migrationsPath := "file://" + filepath.Join(currentDir, "migrations")

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
