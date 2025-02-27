package testutils

import (
	"context"
	"fmt"
	"github.com/Javlopez/go-api/pkg/database"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer represents a Postgres container for testing
type PostgresContainer struct {
	Container testcontainers.Container
	Config    database.Config
	DB        *sqlx.DB
}

// NewPostgresContainer starts a new Postgres container for testing
func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	// Default container setup
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "trade_orders_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(time.Second * 30),
	}

	// Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Get container connection details
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Create database config
	dbConfig := database.Config{
		Host:     ip,
		Port:     mappedPort.Port(),
		User:     "postgres",
		Password: "postgres",
		DBName:   "trade_orders_test",
		SSLMode:  "disable",
	}

	// Wait a moment for the database to initialize fully
	time.Sleep(time.Second * 2)

	// Create container struct
	pgContainer := &PostgresContainer{
		Container: container,
		Config:    dbConfig,
	}

	// Connect to the database
	if err := pgContainer.Connect(); err != nil {
		// Clean up container on connection failure
		pgContainer.Terminate(ctx)
		return nil, err
	}

	return pgContainer, nil
}

// Connect establishes a connection to the database
func (p *PostgresContainer) Connect() error {
	var err error
	// Try to connect multiple times (the container might need a moment to be ready)
	for i := 0; i < 5; i++ {
		p.DB, err = sqlx.Connect("postgres", p.Config.DSN())
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to test database after multiple attempts: %w", err)
	}

	return nil
}

// SetupOrdersTable creates the orders table in the test database
func (p *PostgresContainer) SetupOrdersTable() error {
	_, err := p.DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			symbol VARCHAR(20) NOT NULL,
			price DECIMAL(12, 4) NOT NULL,
			quantity INTEGER NOT NULL,
			order_type VARCHAR(10) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create test table: %w", err)
	}
	return nil
}

// CleanupData removes all data from the orders table
func (p *PostgresContainer) CleanupData() error {
	_, err := p.DB.Exec("DELETE FROM orders")
	return err
}

// Close closes the database connection
func (p *PostgresContainer) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

// Terminate stops and removes the container
func (p *PostgresContainer) Terminate(ctx context.Context) {
	if p.Container != nil {
		p.Container.Terminate(ctx)
	}
}
