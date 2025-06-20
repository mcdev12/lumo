# Lumo Go to Rust Migration Guide

## Table of Contents
1. [Project Overview](#project-overview)
2. [Current Architecture Analysis](#current-architecture-analysis)
3. [Rust Architecture Design](#rust-architecture-design)
4. [Migration Strategy](#migration-strategy)
5. [Step-by-Step Migration Plan](#step-by-step-migration-plan)
6. [Technology Choices](#technology-choices)
7. [Code Structure](#code-structure)
8. [Database Migration](#database-migration)
9. [API Design](#api-design)
10. [Testing Strategy](#testing-strategy)
11. [Deployment Considerations](#deployment-considerations)

## Project Overview

Lumo is a travel planning application that helps users organize their trips through a graph-based structure:

- **Lumo**: A travel plan/journey (the container)
- **Lume**: Individual travel nodes (cities, attractions, restaurants, etc.)
- **Link**: Connections between Lumes (travel routes, recommendations, etc.)

The application provides a RESTful API for managing these entities with full CRUD operations, supporting features like:
- User-specific travel plans
- Location-based nodes with GPS coordinates
- Travel connections with duration, cost, and distance
- Image uploads and categorization
- Scheduling with date/time ranges

## Current Architecture Analysis

### Technology Stack
- **Language**: Go 1.24
- **API Protocol**: Connect RPC (gRPC-compatible)
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Code Generation**: Protocol Buffers with Buf
- **Architecture Pattern**: Clean Architecture with layered approach

### Layer Structure
1. **Service Layer**: HTTP handlers using Connect RPC
2. **Application Layer**: Business logic and validation
3. **Repository Layer**: Data access with SQLC
4. **Model Layer**: Domain models
5. **Transport Layer**: Protocol Buffer definitions

## Rust Architecture Design

### Proposed Technology Stack
- **Language**: Rust (latest stable)
- **Web Framework**: Axum (modern, performant, Tower-based)
- **Database**: PostgreSQL with SQLx (compile-time checked queries)
- **API Design**: RESTful JSON API with OpenAPI documentation
- **Serialization**: Serde for JSON
- **Async Runtime**: Tokio
- **Error Handling**: thiserror + anyhow
- **Validation**: validator crate
- **Configuration**: config-rs
- **Logging**: tracing + tracing-subscriber
- **Testing**: Built-in test framework + mockall for mocking

### Why These Choices?

1. **Axum over gRPC**: 
   - Simpler client integration (standard HTTP/JSON)
   - Better browser compatibility
   - Excellent performance with Tower middleware
   - Native async/await support

2. **SQLx over Diesel**:
   - Compile-time query verification
   - Pure async support
   - Simpler API
   - Better migration tooling

3. **RESTful API**:
   - Universal client support
   - Better debugging tools
   - Simpler deployment
   - OpenAPI for automatic documentation

## Migration Strategy

### Phase 1: Foundation (Week 1-2)
- Set up Rust project structure
- Implement configuration management
- Database connection pool
- Basic error handling framework
- Logging infrastructure

### Phase 2: Core Domain (Week 2-3)
- Implement domain models
- Database schema migrations
- Repository trait definitions
- SQLx query implementations

### Phase 3: Business Logic (Week 3-4)
- Port application layer logic
- Implement validation rules
- Business rule enforcement
- Transaction management

### Phase 4: API Layer (Week 4-5)
- Axum router setup
- Request/response DTOs
- Middleware (CORS, logging, auth)
- OpenAPI documentation

### Phase 5: Testing & Polish (Week 5-6)
- Comprehensive test suite
- Performance optimization
- Docker containerization
- Documentation

## Step-by-Step Migration Plan

### Step 1: Initialize Rust Project

```bash
# Create new Rust project
cargo new lumo-rust --bin
cd lumo-rust

# Initialize git
git init
git add .
git commit -m "Initial Rust project setup"
```

### Step 2: Project Structure

Create the following directory structure:
```
lumo-rust/
├── Cargo.toml
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── migrations/
│   ├── 001_create_lumo.sql
│   ├── 002_create_lume.sql
│   └── 003_create_link.sql
├── src/
│   ├── main.rs
│   ├── lib.rs
│   ├── config/
│   │   └── mod.rs
│   ├── domain/
│   │   ├── mod.rs
│   │   ├── lumo.rs
│   │   ├── lume.rs
│   │   └── link.rs
│   ├── infrastructure/
│   │   ├── mod.rs
│   │   ├── database/
│   │   │   ├── mod.rs
│   │   │   └── connection.rs
│   │   └── repositories/
│   │       ├── mod.rs
│   │       ├── lumo_repository.rs
│   │       ├── lume_repository.rs
│   │       └── link_repository.rs
│   ├── application/
│   │   ├── mod.rs
│   │   ├── services/
│   │   │   ├── mod.rs
│   │   │   ├── lumo_service.rs
│   │   │   ├── lume_service.rs
│   │   │   └── link_service.rs
│   │   └── dto/
│   │       ├── mod.rs
│   │       ├── lumo_dto.rs
│   │       ├── lume_dto.rs
│   │       └── link_dto.rs
│   ├── api/
│   │   ├── mod.rs
│   │   ├── routes/
│   │   │   ├── mod.rs
│   │   │   ├── lumo_routes.rs
│   │   │   ├── lume_routes.rs
│   │   │   └── link_routes.rs
│   │   ├── middleware/
│   │   │   ├── mod.rs
│   │   │   ├── cors.rs
│   │   │   └── logging.rs
│   │   └── error.rs
│   └── utils/
│       ├── mod.rs
│       └── validation.rs
├── tests/
│   ├── integration/
│   │   ├── lumo_tests.rs
│   │   ├── lume_tests.rs
│   │   └── link_tests.rs
│   └── common/
│       └── mod.rs
└── benches/
    └── api_bench.rs
```

### Step 3: Dependencies (Cargo.toml)

```toml
[package]
name = "lumo-rust"
version = "0.1.0"
edition = "2021"

[dependencies]
# Web framework
axum = { version = "0.7", features = ["macros"] }
tower = "0.4"
tower-http = { version = "0.5", features = ["cors", "trace"] }

# Async runtime
tokio = { version = "1", features = ["full"] }

# Database
sqlx = { version = "0.8", features = ["runtime-tokio-rustls", "postgres", "uuid", "time", "json"] }

# Serialization
serde = { version = "1", features = ["derive"] }
serde_json = "1"

# UUID support
uuid = { version = "1", features = ["v4", "serde"] }

# Time handling
time = { version = "0.3", features = ["serde"] }

# Error handling
thiserror = "1"
anyhow = "1"

# Validation
validator = { version = "0.18", features = ["derive"] }

# Configuration
config = "0.14"

# Logging
tracing = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter"] }

# Environment variables
dotenvy = "0.15"

# API documentation
utoipa = { version = "4", features = ["axum_extras"] }
utoipa-swagger-ui = { version = "6", features = ["axum"] }

[dev-dependencies]
# Testing
mockall = "0.12"
reqwest = { version = "0.12", features = ["json"] }
tower = { version = "0.4", features = ["util"] }

[profile.release]
lto = true
opt-level = 3
codegen-units = 1
```

### Step 4: Domain Models

#### src/domain/lumo.rs
```rust
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use time::OffsetDateTime;
use uuid::Uuid;
use validator::Validate;

#[derive(Debug, Clone, Serialize, Deserialize, FromRow, Validate)]
pub struct Lumo {
    pub id: i64,
    pub lumo_id: Uuid,
    pub user_id: Uuid,
    #[validate(length(min = 1, max = 255))]
    pub title: String,
    #[serde(with = "time::serde::rfc3339")]
    pub created_at: OffsetDateTime,
    #[serde(with = "time::serde::rfc3339")]
    pub updated_at: OffsetDateTime,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct CreateLumoRequest {
    pub user_id: Uuid,
    #[validate(length(min = 1, max = 255))]
    pub title: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct UpdateLumoRequest {
    #[validate(length(min = 1, max = 255))]
    pub title: String,
}
```

### Step 5: Repository Implementation

#### src/infrastructure/repositories/lumo_repository.rs
```rust
use async_trait::async_trait;
use sqlx::{PgPool, Result};
use uuid::Uuid;

use crate::domain::lumo::{CreateLumoRequest, Lumo, UpdateLumoRequest};

#[async_trait]
pub trait LumoRepository: Send + Sync {
    async fn create(&self, req: &CreateLumoRequest) -> Result<Lumo>;
    async fn find_by_id(&self, id: i64) -> Result<Option<Lumo>>;
    async fn find_by_lumo_id(&self, lumo_id: Uuid) -> Result<Option<Lumo>>;
    async fn find_by_user_id(&self, user_id: Uuid, limit: i32, offset: i32) -> Result<Vec<Lumo>>;
    async fn update(&self, id: i64, req: &UpdateLumoRequest) -> Result<Lumo>;
    async fn delete(&self, id: i64) -> Result<()>;
    async fn count_by_user_id(&self, user_id: Uuid) -> Result<i64>;
}

pub struct PostgresLumoRepository {
    pool: PgPool,
}

impl PostgresLumoRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl LumoRepository for PostgresLumoRepository {
    async fn create(&self, req: &CreateLumoRequest) -> Result<Lumo> {
        let lumo = sqlx::query_as!(
            Lumo,
            r#"
            INSERT INTO lumo (lumo_id, user_id, title)
            VALUES ($1, $2, $3)
            RETURNING id, lumo_id, user_id, title, created_at, updated_at
            "#,
            Uuid::new_v4(),
            req.user_id,
            req.title
        )
        .fetch_one(&self.pool)
        .await?;

        Ok(lumo)
    }

    async fn find_by_id(&self, id: i64) -> Result<Option<Lumo>> {
        let lumo = sqlx::query_as!(
            Lumo,
            r#"
            SELECT id, lumo_id, user_id, title, created_at, updated_at
            FROM lumo
            WHERE id = $1
            "#,
            id
        )
        .fetch_optional(&self.pool)
        .await?;

        Ok(lumo)
    }

    // ... implement other methods
}
```

### Step 6: Service Layer

#### src/application/services/lumo_service.rs
```rust
use std::sync::Arc;
use uuid::Uuid;

use crate::{
    domain::lumo::{CreateLumoRequest, Lumo, UpdateLumoRequest},
    infrastructure::repositories::lumo_repository::LumoRepository,
};

pub struct LumoService {
    repository: Arc<dyn LumoRepository>,
}

impl LumoService {
    pub fn new(repository: Arc<dyn LumoRepository>) -> Self {
        Self { repository }
    }

    pub async fn create_lumo(&self, req: CreateLumoRequest) -> anyhow::Result<Lumo> {
        // Validate request
        req.validate()?;
        
        // Create lumo
        let lumo = self.repository.create(&req).await?;
        
        Ok(lumo)
    }

    pub async fn get_lumo(&self, lumo_id: Uuid) -> anyhow::Result<Option<Lumo>> {
        self.repository.find_by_lumo_id(lumo_id).await
            .map_err(Into::into)
    }

    pub async fn list_lumos_by_user(
        &self,
        user_id: Uuid,
        limit: i32,
        offset: i32,
    ) -> anyhow::Result<Vec<Lumo>> {
        let limit = limit.clamp(1, 100);
        let offset = offset.max(0);
        
        self.repository.find_by_user_id(user_id, limit, offset).await
            .map_err(Into::into)
    }

    // ... implement other methods
}
```

### Step 7: API Routes

#### src/api/routes/lumo_routes.rs
```rust
use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::IntoResponse,
    Json,
};
use serde::Deserialize;
use std::sync::Arc;
use uuid::Uuid;

use crate::{
    api::error::ApiError,
    application::{dto::lumo_dto::*, services::lumo_service::LumoService},
    domain::lumo::{CreateLumoRequest, UpdateLumoRequest},
};

#[derive(Deserialize)]
pub struct ListQuery {
    #[serde(default = "default_limit")]
    limit: i32,
    #[serde(default)]
    offset: i32,
}

fn default_limit() -> i32 {
    50
}

pub async fn create_lumo(
    State(service): State<Arc<LumoService>>,
    Json(req): Json<CreateLumoRequest>,
) -> Result<impl IntoResponse, ApiError> {
    let lumo = service.create_lumo(req).await?;
    Ok((StatusCode::CREATED, Json(lumo)))
}

pub async fn get_lumo(
    State(service): State<Arc<LumoService>>,
    Path(lumo_id): Path<Uuid>,
) -> Result<Json<Lumo>, ApiError> {
    let lumo = service
        .get_lumo(lumo_id)
        .await?
        .ok_or(ApiError::NotFound)?;
    
    Ok(Json(lumo))
}

pub async fn list_lumos(
    State(service): State<Arc<LumoService>>,
    Path(user_id): Path<Uuid>,
    Query(query): Query<ListQuery>,
) -> Result<Json<Vec<Lumo>>, ApiError> {
    let lumos = service
        .list_lumos_by_user(user_id, query.limit, query.offset)
        .await?;
    
    Ok(Json(lumos))
}

// ... implement other handlers
```

### Step 8: Main Application Setup

#### src/main.rs
```rust
use axum::{
    routing::{get, post, put, delete},
    Router,
};
use sqlx::postgres::PgPoolOptions;
use std::{net::SocketAddr, sync::Arc};
use tower_http::cors::CorsLayer;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod api;
mod application;
mod config;
mod domain;
mod infrastructure;
mod utils;

use crate::{
    api::routes::{lumo_routes, lume_routes, link_routes},
    application::services::{lumo_service::LumoService, lume_service::LumeService, link_service::LinkService},
    config::Config,
    infrastructure::repositories::{
        lumo_repository::PostgresLumoRepository,
        lume_repository::PostgresLumeRepository,
        link_repository::PostgresLinkRepository,
    },
};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Initialize tracing
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "lumo_rust=debug,tower_http=debug".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    // Load configuration
    let config = Config::from_env()?;

    // Create database pool
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&config.database_url)
        .await?;

    // Run migrations
    sqlx::migrate!("./migrations").run(&pool).await?;

    // Initialize repositories
    let lumo_repo = Arc::new(PostgresLumoRepository::new(pool.clone()));
    let lume_repo = Arc::new(PostgresLumeRepository::new(pool.clone()));
    let link_repo = Arc::new(PostgresLinkRepository::new(pool.clone()));

    // Initialize services
    let lumo_service = Arc::new(LumoService::new(lumo_repo));
    let lume_service = Arc::new(LumeService::new(lume_repo));
    let link_service = Arc::new(LinkService::new(link_repo));

    // Build router
    let app = Router::new()
        // Lumo routes
        .route("/api/v1/lumos", post(lumo_routes::create_lumo))
        .route("/api/v1/lumos/:lumo_id", get(lumo_routes::get_lumo))
        .route("/api/v1/lumos/:lumo_id", put(lumo_routes::update_lumo))
        .route("/api/v1/lumos/:lumo_id", delete(lumo_routes::delete_lumo))
        .route("/api/v1/users/:user_id/lumos", get(lumo_routes::list_lumos))
        // Lume routes
        .route("/api/v1/lumes", post(lume_routes::create_lume))
        .route("/api/v1/lumes/:lume_id", get(lume_routes::get_lume))
        .route("/api/v1/lumes/:lume_id", put(lume_routes::update_lume))
        .route("/api/v1/lumes/:lume_id", delete(lume_routes::delete_lume))
        .route("/api/v1/lumos/:lumo_id/lumes", get(lume_routes::list_lumes))
        // Link routes
        .route("/api/v1/links", post(link_routes::create_link))
        .route("/api/v1/links/:link_id", get(link_routes::get_link))
        .route("/api/v1/links/:link_id", put(link_routes::update_link))
        .route("/api/v1/links/:link_id", delete(link_routes::delete_link))
        .route("/api/v1/lumes/:lume_id/links", get(link_routes::list_links))
        // Health check
        .route("/health", get(|| async { "OK" }))
        // Add state
        .with_state(lumo_service.clone())
        .with_state(lume_service.clone())
        .with_state(link_service.clone())
        // Add middleware
        .layer(CorsLayer::permissive())
        .layer(tower_http::trace::TraceLayer::new_for_http());

    // Start server
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::info!("Server listening on {}", addr);
    
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await?;

    Ok(())
}
```

## Technology Choices

### Core Framework Decisions

1. **Axum Web Framework**
   - Built on Tower service architecture
   - Excellent performance
   - Type-safe routing
   - Great middleware ecosystem

2. **SQLx for Database**
   - Compile-time SQL verification
   - Pure async implementation
   - Migration support built-in
   - No ORM overhead

3. **Error Handling Strategy**
   - `thiserror` for library errors
   - `anyhow` for application errors
   - Custom API error types with proper HTTP status codes

4. **Validation Approach**
   - Domain-level validation with `validator` crate
   - Request validation at API boundary
   - Database constraints as final safety net

## Code Structure

### Clean Architecture Principles

1. **Domain Layer** (`src/domain/`)
   - Pure business entities
   - No framework dependencies
   - Validation rules

2. **Infrastructure Layer** (`src/infrastructure/`)
   - Database implementations
   - External service integrations
   - Repository implementations

3. **Application Layer** (`src/application/`)
   - Business logic orchestration
   - Use case implementations
   - DTOs for API contracts

4. **API Layer** (`src/api/`)
   - HTTP routing
   - Request/response handling
   - Middleware configuration

## Database Migration

### Migration Files

The existing PostgreSQL schema can be reused with minor adjustments:

#### migrations/001_create_lumo.sql
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS lumo (
    id BIGSERIAL PRIMARY KEY,
    lumo_id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lumo_user_id ON lumo(user_id);
CREATE INDEX idx_lumo_created_at ON lumo(created_at);

-- Update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_lumo_updated_at BEFORE UPDATE
    ON lumo FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## API Design

### RESTful Endpoints

#### Lumo Endpoints
- `POST /api/v1/lumos` - Create a new lumo
- `GET /api/v1/lumos/:lumo_id` - Get a specific lumo
- `PUT /api/v1/lumos/:lumo_id` - Update a lumo
- `DELETE /api/v1/lumos/:lumo_id` - Delete a lumo
- `GET /api/v1/users/:user_id/lumos` - List lumos for a user

#### Lume Endpoints
- `POST /api/v1/lumes` - Create a new lume
- `GET /api/v1/lumes/:lume_id` - Get a specific lume
- `PUT /api/v1/lumes/:lume_id` - Update a lume
- `DELETE /api/v1/lumes/:lume_id` - Delete a lume
- `GET /api/v1/lumos/:lumo_id/lumes` - List lumes for a lumo

#### Link Endpoints
- `POST /api/v1/links` - Create a new link
- `GET /api/v1/links/:link_id` - Get a specific link
- `PUT /api/v1/links/:link_id` - Update a link
- `DELETE /api/v1/links/:link_id` - Delete a link
- `GET /api/v1/lumes/:lume_id/links` - List links for a lume

### OpenAPI Documentation

Use `utoipa` for automatic OpenAPI generation:

```rust
#[derive(OpenApi)]
#[openapi(
    paths(
        lumo_routes::create_lumo,
        lumo_routes::get_lumo,
        // ... other routes
    ),
    components(
        schemas(Lumo, CreateLumoRequest, UpdateLumoRequest)
    ),
    tags(
        (name = "lumo", description = "Lumo management endpoints")
    )
)]
pub struct ApiDoc;
```

## Testing Strategy

### Unit Tests
- Domain logic validation
- Service layer business rules
- Repository mocks for isolation

### Integration Tests
- Full API endpoint testing
- Database interaction verification
- Transaction rollback for test isolation

### Example Integration Test
```rust
#[cfg(test)]
mod tests {
    use super::*;
    use axum::http::StatusCode;
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_create_lumo() {
        let app = create_test_app().await;
        
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/v1/lumos")
                    .header("content-type", "application/json")
                    .body(Body::from(
                        serde_json::to_string(&CreateLumoRequest {
                            user_id: Uuid::new_v4(),
                            title: "Test Journey".to_string(),
                        })
                        .unwrap(),
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);
    }
}
```

## Deployment Considerations

### Docker Configuration

#### Dockerfile
```dockerfile
# Build stage
FROM rust:1.75 as builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src
COPY migrations ./migrations

RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/target/release/lumo-rust /app/
COPY migrations ./migrations

ENV RUST_LOG=info

EXPOSE 8080

CMD ["./lumo-rust"]
```

### Performance Optimizations

1. **Connection Pooling**
   - Configure appropriate pool sizes
   - Use PgBouncer for additional pooling if needed

2. **Caching Strategy**
   - Redis for session/frequently accessed data
   - HTTP caching headers
   - Query result caching

3. **Monitoring**
   - Prometheus metrics endpoint
   - Structured logging with tracing
   - Health check endpoints

### Migration Checklist

- [ ] Set up Rust project structure
- [ ] Configure dependencies
- [ ] Implement domain models
- [ ] Create database migrations
- [ ] Implement repositories
- [ ] Port business logic to services
- [ ] Create API routes
- [ ] Set up middleware
- [ ] Add comprehensive tests
- [ ] Configure Docker builds
- [ ] Set up CI/CD pipeline
- [ ] Performance testing
- [ ] Load testing
- [ ] Documentation
- [ ] Deployment scripts

## Conclusion

This migration guide provides a comprehensive roadmap for converting the Lumo Go backend to Rust. The proposed architecture leverages Rust's strengths:

- **Memory safety** without garbage collection
- **Fearless concurrency** with async/await
- **Zero-cost abstractions** for clean architecture
- **Compile-time guarantees** for correctness

The migration maintains the core business logic while modernizing the technology stack for better performance, safety, and developer experience.