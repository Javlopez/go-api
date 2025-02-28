# Trade Orders API

A RESTful API for managing trade orders built with Go, Gin framework, and PostgreSQL. The application follows clean architecture principles with repository pattern, includes Swagger documentation, and is containerized with Docker.

## Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Development Setup](#development-setup)
    - [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Database Migrations](#database-migrations)
- [Testing](#testing)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
- [Deployment](#deployment)
- [GitHub Workflows](#github-workflows)
- [Configuration](#configuration)

## Features

- Create and retrieve trade orders
- Validate incoming requests
- Structured error responses
- PostgreSQL database with migrations
- Swagger API documentation
- Containerized with Docker
- Hot reload for development
- Comprehensive test suite

## Technology Stack

- **Go**: Main programming language
- **Gin**: Web framework
- **SQLx**: Database access layer
- **PostgreSQL**: Relational database
- **Swagger**: API documentation
- **Docker/Docker Compose**: Containerization
- **Air**: Hot reloading for development
- **Golang-migrate**: Database migrations
- **Testcontainers**: For integration testing

## Getting Started

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Git

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Javlopez/go-api.git
   cd go-api
   ```

2. Start the development environment:
   ```bash
   make dev
   ```
   This will:
    - Start a PostgreSQL container
    - Start the API with hot reloading
    - Set up the environment for development

3. Run database migrations:
   ```bash
   make migrate
   ```

### Running the Application

- **Development mode** (with hot reloading):
  ```bash
  make dev
  ```

- **Run in background**:
  ```bash
  make up
  ```

- **Stop the application**:
  ```bash
  make down
  ```

- **View logs**:
  ```bash
  make logs
  ```

## Project Structure

```
├── cmd
│   ├── api                 # Main API application
│   └── migrate             # Database migration tool
├── docs                    # Swagger documentation
├── infra
│   ├── terraform           # Terraform configuration
├── pkg
│   ├── api                 # API handlers and routes
│   ├── config              # Configuration management
│   ├── models              # Data models
│   ├── repository          # Database repository layer
│   ├── testutils           # Testing utilities
│   └── utils               # Helper utilities
├── test
│   └── integration         # Integration tests
├── .air.toml               # Air configuration for hot reload
├── .env                    # Environment variables
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Main application Dockerfile
├── Dockerfile.migrate      # Migrations Dockerfile
├── go.mod                  # Go modules
└── Makefile                # Makefile with common commands
```

## API Endpoints

### Documentation

```
GET /docs/index.html
```

### Create Order

```
POST /api/v1/orders
```

Example request body:
```json
{
  "symbol": "AAPL",
  "price": 150.50,
  "quantity": 10,
  "order_type": "BUY"
}
```

### Get Orders

```
GET /api/v1/orders
```

## Database Migrations

The project uses golang-migrate for database migrations. The migrations are stored in the `migrations` directory.

- **Create a new migration**:
  ```bash
  make migrate-create name=add_status_column
  ```

- **Run migrations**:
  ```bash
  make migrate
  ```

## Testing

### Unit Tests

Run unit tests with:
```bash
make test
```

Generate test coverage report:
```bash
make test-coverage
```

### Integration Tests

Run integration tests (requires Docker):
```bash
make test-integration
```

The integration tests use Testcontainers to spin up a PostgreSQL instance for each test run.

## Deployment

The application can be deployed using Docker Compose:

1. Set up environment variables in `.env` file
2. Run migrations:
   ```bash
   docker-compose --profile migrate up migrate
   ```
3. Start the application:
   ```bash
   docker-compose up -d
   ```

For production deployment, a specialized `docker-compose.production.yml` can be used with the provided GitHub Actions workflows.

## GitHub Workflows

The project includes several GitHub workflows:

1. **Unit Tests**: Runs on every push and pull request
2. **Integration Tests**: Runs on every push and pull request using Testcontainers
3. **Build**: Builds Docker images and pushes them to Docker Hub
4. **Deploy**: Deploys the application to a production server via SSH

See the [GitHub Workflows documentation](/.github/workflows/README.md) for more details.

## Configuration

The application is configured through environment variables, which can be set in the `.env` file or directly in the environment.

Key environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Port for the API server | 8080 |
| DB_HOST | PostgreSQL hostname | postgres |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL username | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | PostgreSQL database name | trade_orders |
| DB_SSLMODE | PostgreSQL SSL mode | disable |
| GIN_MODE | Gin framework mode (debug/release) | debug |
