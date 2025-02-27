# Dockerfile
FROM golang:1.23-alpine

# Install git and Air for hot reloading
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Install swag for Swagger generation
RUN go install github.com/swaggo/swag/cmd/swag

# Generate Swagger docs
RUN swag init -g main.go

# Expose port
EXPOSE 8080

# Air will watch for file changes and rebuild the app
CMD ["air", "-c", ".air.toml"]