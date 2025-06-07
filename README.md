# Lumo

Lumo is a Go-based microservice application that provides a RESTful API using Connect RPC protocol. It's designed to be containerized and easily deployable.

## Technologies

- **Go**: Core programming language
- **Connect RPC**: API protocol for service communication
- **PostgreSQL**: Database for persistent storage
- **Docker**: Containerization for consistent deployment
- **Docker Compose**: Multi-container orchestration
- **gRPC UI**: Interactive API exploration

## Project Structure

```
lumo/
├── Dockerfile              # Docker configuration for the Go application
├── Makefile                # Commands for development, testing, and deployment
├── README.md               # This file
├── buf.gen.yaml            # Buf configuration for code generation
├── buf.yaml                # Buf configuration
├── docker-compose.yaml     # Docker Compose configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── go/                     # Go source code
│   └── internal/           # Internal application code
│       ├── app/            # Application business logic
│       ├── cmd/            # Command-line entry points
│       ├── genproto/       # Generated protocol buffer code
│       ├── models/         # Domain models
│       ├── repository/     # Data access layer
│       └── service/        # Service implementations
└── protobuf/               # Protocol buffer definitions
```

## Prerequisites

- Go 1.24 or later
- Docker and Docker Compose
- PostgreSQL (if running locally)
- grpcui (for API exploration)

## Setup and Installation

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/lumo.git
   cd lumo
   ```

2. Start the entire stack:
   ```bash
   make stack-up
   ```

This will build and start both the application and PostgreSQL database containers.

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/lumo.git
   cd lumo
   ```

2. Start only the database:
   ```bash
   make db-start
   ```

3. Run the application locally:
   ```bash
   make dev-run
   ```

## Running the Application

The application exposes a Connect RPC service on port 8080. You can interact with it using any Connect RPC client or the included gRPC UI.

### Using gRPC UI

To explore the API using gRPC UI:

```bash
make grpcui
```

This will open a web interface at http://localhost:8080 where you can interactively test the API endpoints.

## Development Workflow

### Running Tests

Run tests locally:
```bash
make dev-test
```

Run tests in Docker:
```bash
make test-docker
```

Run tests with coverage:
```bash
make test-coverage
```

### Linting

Run the linter:
```bash
make dev-lint
```

## Database Management

### Connect to PostgreSQL

```bash
make db-psql
```

### Reset the Database

```bash
make db-reset
```

## Available Commands

Run `make help` to see all available commands:

### Application Commands
- `app-build` - Build the application Docker image
- `app-start` - Start the application container
- `app-stop` - Stop the application container
- `app-logs` - View application logs
- `app-restart` - Restart the application container
- `app-shell` - Open a shell in the application container

### Full Stack Commands
- `stack-up` - Start the entire application stack
- `stack-down` - Stop the entire application stack
- `stack-restart` - Restart the entire application stack
- `stack-logs` - View logs for all services
- `stack-build` - Build all services
- `stack-clean` - Clean up all containers and volumes

### Database Commands
- `db-start` - Start the database container
- `db-stop` - Stop the database container
- `db-reset` - Reset database (delete volume and restart)
- `db-psql` - Connect to PostgreSQL using psql
- `db-down` - Stop and remove containers but preserve the data
- `db-down-v` - Stop and remove containers AND delete volume data

### Development and Testing Commands
- `dev-run` - Run the application locally (not in Docker)
- `dev-test` - Run tests locally
- `dev-lint` - Run linter locally
- `test-docker` - Run tests in Docker
- `test-coverage` - Run tests with coverage in Docker

### gRPC UI Commands
- `grpcui` - Start gRPC UI
- `grpcui-dev` - Start gRPC UI in development mode
- `grpcui-proto` - Start gRPC UI with proto files

## Docker Configuration

The application is containerized using Docker with a multi-stage build process for efficiency:

1. The first stage builds the Go application
2. The second stage creates a minimal runtime image

The Docker Compose configuration sets up:
- The Go application container
- A PostgreSQL database container
- Proper networking between containers
- Volume persistence for the database

## Environment Variables

The application uses the following environment variables (with defaults):

- `DB_HOST` (default: "localhost")
- `DB_PORT` (default: 5432)
- `DB_USER` (default: "postgres")
- `DB_PASSWORD` (default: "postgres")
- `DB_NAME` (default: "lumo_db")
- `DB_SSLMODE` (default: "disable")

These can be configured in the docker-compose.yaml file or set directly in your environment.