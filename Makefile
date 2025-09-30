# WhatsApp Meow Service Makefile

.PHONY: build run test clean docker-build docker-run help

# Default target
help:
	@echo "Available targets:"
	@echo "  build         - Build the Go application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker"
	@echo "  setup         - Setup development environment"
	@echo "  deps          - Download dependencies"

# Build the application
build:
	@echo "Building WhatsApp Meow service..."
	go build -o bin/whatsmeow-service main.go

# Run the application
run:
	@echo "Running WhatsApp Meow service..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Setup development environment
setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp env.example .env; fi
	@echo "Please edit .env file with your configuration"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t whatsmeow-service:latest .

# Run with Docker
docker-run:
	@echo "Running with Docker..."
	docker run -d \
		--name whatsmeow-service \
		-p 8081:8081 \
		--env-file .env \
		whatsmeow-service:latest

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop whatsmeow-service || true
	docker rm whatsmeow-service || true

# View Docker logs
docker-logs:
	@echo "Viewing Docker logs..."
	docker logs -f whatsmeow-service

# Database setup
db-setup:
	@echo "Setting up database schema..."
	@echo "Please run the SQL commands from database/schema.sql in your PostgreSQL database"

# Health check
health:
	@echo "Checking service health..."
	@curl -s http://localhost:8081/health | jq . || echo "Service not running or jq not installed"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Generate code
generate:
	@echo "Generating code..."
	go generate ./...

# Development with hot reload
dev:
	@echo "Starting development server with hot reload..."
	air

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/whatsmeow-service main.go

# Create release
release: clean prod-build
	@echo "Creating release..."
	tar -czf whatsmeow-service-$(shell date +%Y%m%d-%H%M%S).tar.gz bin/ README.md database/ integration/

# Full setup for development
dev-setup: setup deps install-tools
	@echo "Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Edit .env file with your database configuration"
	@echo "2. Run 'make db-setup' to create database schema"
	@echo "3. Run 'make run' to start the service"
