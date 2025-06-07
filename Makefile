# Default target
.DEFAULT_GOAL := help

# Docker and application commands
.PHONY: app-build app-start app-stop app-logs app-restart app-shell

# Build the application Docker image
app-build:
	@echo "Building application Docker image..."
	docker-compose build app

# Start the application container
app-start:
	@echo "Starting application container..."
	docker-compose up -d app

# Stop the application container
app-stop:
	@echo "Stopping application container..."
	docker-compose stop app

# View application logs
app-logs:
	@echo "Viewing application logs..."
	docker-compose logs -f app

# Restart the application container
app-restart:
	@echo "Restarting application container..."
	docker-compose restart app

# Open a shell in the application container
app-shell:
	@echo "Opening shell in application container..."
	docker-compose exec app sh

# Full stack commands
.PHONY: stack-up stack-down stack-restart stack-logs stack-build

# Start the entire application stack
stack-up:
	@echo "Starting the entire application stack..."
	docker-compose up -d

# Stop the entire application stack
stack-down:
	@echo "Stopping the entire application stack..."
	docker-compose down
	@echo "All containers stopped. Data volumes are preserved."

# Restart the entire application stack
stack-restart:
	@echo "Restarting the entire application stack..."
	docker-compose restart

# View logs for all services
stack-logs:
	@echo "Viewing logs for all services..."
	docker-compose logs -f

# Build all services
stack-build:
	@echo "Building all services..."
	docker-compose build

# Clean up all containers and volumes
stack-clean:
	@echo "Cleaning up all containers and volumes..."
	docker-compose down -v
	@echo "All containers and volumes removed."

# Database and SQLC related commands
.PHONY: db-start db-stop db-reset db-psql sqlc-gen

# Start the database container
db-start:
	@echo "Starting PostgreSQL container..."
	docker-compose up -d postgres

# Stop the database container
db-stop:
	@echo "Stopping PostgreSQL container..."
	docker-compose stop postgres

# Reset database (delete volume and restart)
db-reset:
	@echo "Removing PostgreSQL container and volume..."
	docker-compose down -v
	@echo "Starting fresh PostgreSQL container..."
	docker-compose up -d postgres

# Connect to PostgreSQL using psql
db-psql:
	@echo "🔌 Connecting to PostgreSQL..."
	docker-compose exec postgres psql -U postgres -d lumo_db

# Database down command - removes containers and optionally volumes
.PHONY: db-down db-down-v

# Stop and remove containers but preserve the data
db-down:
	@echo "Stopping and removing PostgreSQL container..."
	docker-compose down
	@echo "Database container removed. Data volume is preserved."

# Stop and remove containers AND delete volume data
db-down-v:
	@echo "Stopping PostgreSQL container and removing all data..."
	docker-compose down -v
	@echo "Database container and volume removed completely."

# Development and testing commands
.PHONY: dev-run dev-test dev-lint

# Run the application locally (not in Docker)
dev-run:
	@echo "Running application locally..."
	go run ./go/internal/cmd/main.go

# Run tests locally
dev-test:
	@echo "Running tests locally..."
	go test ./go/...

# Run linter locally
dev-lint:
	@echo "Running linter locally..."
	go vet ./go/...

# Docker testing commands
.PHONY: test-docker test-coverage

# Run tests in Docker
test-docker:
	@echo "Running tests in Docker..."
	docker-compose run --rm app go test ./go/...

# Run tests with coverage in Docker
test-coverage:
	@echo "Running tests with coverage in Docker..."
	docker-compose run --rm app go test -coverprofile=coverage.out ./go/...
	docker-compose run --rm app go tool cover -html=coverage.out -o coverage.html

# gRPC UI commands
.PHONY: grpcui grpcui-dev grpcui-proto

# Start gRPC UI
grpcui:
	@echo "Starting grpcui..."
	grpcui -plaintext localhost:8080

# Start gRPC UI in development mode
grpcui-dev:
	@echo "Starting grpcui in development mode..."
	grpcui -plaintext -reflect-headers localhost:8080

# Start gRPC UI with proto files
grpcui-proto:
	@echo "Starting grpcui with proto files..."
	grpcui -plaintext -proto api/proto/lume/service.proto localhost:8080

# Help command
.PHONY: help

# Show help
help:
	@echo "Lumo Application Makefile Commands:"
	@echo ""
	@echo "Application Commands:"
	@echo "  app-build      - Build the application Docker image"
	@echo "  app-start      - Start the application container"
	@echo "  app-stop       - Stop the application container"
	@echo "  app-logs       - View application logs"
	@echo "  app-restart    - Restart the application container"
	@echo "  app-shell      - Open a shell in the application container"
	@echo ""
	@echo "Full Stack Commands:"
	@echo "  stack-up       - Start the entire application stack"
	@echo "  stack-down     - Stop the entire application stack"
	@echo "  stack-restart  - Restart the entire application stack"
	@echo "  stack-logs     - View logs for all services"
	@echo "  stack-build    - Build all services"
	@echo "  stack-clean    - Clean up all containers and volumes"
	@echo ""
	@echo "Database Commands:"
	@echo "  db-start       - Start the database container"
	@echo "  db-stop        - Stop the database container"
	@echo "  db-reset       - Reset database (delete volume and restart)"
	@echo "  db-psql        - Connect to PostgreSQL using psql"
	@echo "  db-down        - Stop and remove containers but preserve the data"
	@echo "  db-down-v      - Stop and remove containers AND delete volume data"
	@echo ""
	@echo "Development and Testing Commands:"
	@echo "  dev-run        - Run the application locally (not in Docker)"
	@echo "  dev-test       - Run tests locally"
	@echo "  dev-lint       - Run linter locally"
	@echo "  test-docker    - Run tests in Docker"
	@echo "  test-coverage  - Run tests with coverage in Docker"
	@echo ""
	@echo "gRPC UI Commands:"
	@echo "  grpcui         - Start gRPC UI"
	@echo "  grpcui-dev     - Start gRPC UI in development mode"
	@echo "  grpcui-proto   - Start gRPC UI with proto files"
