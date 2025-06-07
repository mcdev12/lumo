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
	docker-compose exec postgres psql -U postgres -d moss_db

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

# Add this to your Makefile

.PHONY: grpcui
grpcui:
	@echo "Starting grpcui..."
	grpcui -plaintext localhost:8080

# Alternative with more options
.PHONY: grpcui-dev
grpcui-dev:
	@echo "Starting grpcui in development mode..."
	grpcui -plaintext -reflect-headers localhost:8080

# If you need to specify proto files directly
.PHONY: grpcui-proto
grpcui-proto:
	@echo "Starting grpcui with proto files..."
	grpcui -plaintext -proto api/proto/lume/service.proto localhost:8080