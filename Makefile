# Makefile

.PHONY: dev up down logs psql swagger build clean help

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev      - Start development environment"
	@echo "  make up       - Start containers in the background"
	@echo "  make down     - Stop and remove containers"
	@echo "  make logs     - Follow docker logs"
	@echo "  make psql     - Connect to PostgreSQL database"
	@echo "  make swagger  - Generate Swagger documentation"
	@echo "  make build    - Build the Go application"
	@echo "  make clean    - Clean temporary files"

# Start development environment
dev:
	@echo "Starting development environment..."
	docker compose up

# Build and run the Go application
build-dev:
	@echo "Starting development environment with force build..."
	docker compose up --build

# Start containers in the background
up:
	@echo "Starting containers in the background..."
	docker compose up -d

# Stop and remove containers
down:
	@echo "Stopping and removing containers..."
	docker compose down

# Follow docker logs
logs:
	@echo "Following logs..."
	docker compose logs -f

# Connect to PostgreSQL database
psql:
	@echo "Connecting to PostgreSQL..."
	docker compose exec postgres psql -U postgres -d trade_orders

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	docker compose exec api swag init -g main.go

# Build the Go application
build:
	@echo "Building the Go application..."
	go build -o ./bin/api ./cmd/api

# Run database migrations
migrate:
	@echo "Running database migrations..."
	docker compose --profile migrate up migrate

# Create a new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration: $(name)"
	@mkdir -p migrations
	@timestamp=$(date +%Y%m%d%H%M%S); \
	filename=migrations/${timestamp}_$(name); \
	echo "-- Up: Add your migration SQL here" > ${filename}.up.sql; \
	echo "-- Down: Add rollback SQL here" > ${filename}.down.sql; \
	echo "Created migration files: ${filename}.up.sql and ${filename}.down.sql"

# Run unit tests
test:
	@echo "Running unit tests..."
	go test -v ./pkg/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./pkg/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run integration tests with testcontainers
test-integration:
	@echo "Running integration tests with testcontainers..."
	go test -v ./test/integration/...

# Run all tests
test-all: test test-integration
	@echo "All tests completed"

# Clean temporary files
clean:
	@echo "Cleaning temporary files..."
	rm -rf ./tmp ./bin
	docker compose down -v