# SQLx Repository Pattern for Labels - Detailed Breakdown

This guide provides a comprehensive explanation of implementing a repository pattern using SQLx in Rust. We'll break down each component and explain how it works.

## Overview

The repository pattern provides an abstraction layer between your business logic and data access. In Rust with SQLx, this pattern offers:
- Type-safe SQL queries verified at compile time
- Clean separation of concerns
- Testable database operations
- Connection pooling out of the box

## Core Components Explained

### 1. The Label Model

```rust
use sqlx::{PgPool, postgres::PgRow};
use sqlx::FromRow;
use uuid::Uuid;
use chrono::{DateTime, Utc};

// The FromRow derive macro is crucial - it tells SQLx how to map database rows to this struct
#[derive(Debug, Clone, FromRow)]
pub struct Label {
    pub label_id: Uuid,
    pub label_name: String,
    pub created_at: Option<DateTime<Utc>>,
    pub updated_at: Option<DateTime<Utc>>,
}
```

**Key Points:**
- `#[derive(FromRow)]` - This macro automatically generates code to convert database rows into your struct
- Field names must match column names exactly (or use `#[sqlx(rename = "column_name")]`)
- Types must be compatible with PostgreSQL types:
  - `Uuid` ↔ PostgreSQL `UUID`
  - `String` ↔ PostgreSQL `TEXT/VARCHAR`
  - `DateTime<Utc>` ↔ PostgreSQL `TIMESTAMPTZ`
  - `Option<T>` for nullable columns

### 2. The Repository Structure

```rust
pub struct LabelRepository {
    pool: PgPool,  // PgPool is a connection pool - it manages multiple connections efficiently
}

impl LabelRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }
```

**Why PgPool?**
- Connection pooling prevents the overhead of creating new connections for each query
- Automatically manages connection lifecycle
- Thread-safe - can be cloned and shared across async tasks
- Configurable limits (min/max connections, idle timeout, etc.)

## Repository Methods Deep Dive

### 3. Create Method - Inserting Data

```rust
pub async fn create(&self, label_name: String) -> Result<Label, sqlx::Error> {
    let label_id = Uuid::new_v4();  // Generate ID in Rust rather than using database defaults
    
    let label = sqlx::query_as!(
        Label,                      // Type to map the result to
        r#"                        // r#"..."# is a raw string literal - no escaping needed
        INSERT INTO labels (label_id, label_name)
        VALUES ($1, $2)            // $1, $2 are parameter placeholders (prevent SQL injection)
        RETURNING label_id, label_name, created_at, updated_at  // Get the inserted row back
        "#,
        label_id,                  // Binds to $1
        label_name                 // Binds to $2
    )
    .fetch_one(&self.pool)         // Execute and expect exactly one row back
    .await?;                       // ? propagates any errors
    
    Ok(label)
}
```

**Important Concepts:**
- **`query_as!` macro**: Validates SQL at compile time against your actual database
- **Parameter binding**: `$1, $2` prevents SQL injection automatically
- **RETURNING clause**: PostgreSQL feature to get the inserted row without a second query
- **Error handling**: `Result<T, sqlx::Error>` follows Rust's explicit error handling pattern

### 4. Get by ID - Fetching Optional Data

```rust
pub async fn get_by_id(&self, label_id: Uuid) -> Result<Option<Label>, sqlx::Error> {
    let label = sqlx::query_as!(
        Label,
        r#"
        SELECT label_id, label_name, created_at, updated_at
        FROM labels
        WHERE label_id = $1
        "#,
        label_id
    )
    .fetch_optional(&self.pool)    // Returns Option<Label> - None if no row found
    .await?;
    
    Ok(label)
}
```

**`fetch_optional` vs `fetch_one`:**
- `fetch_one`: Expects exactly one row - errors if none or multiple
- `fetch_optional`: Returns `Option<T>` - gracefully handles missing data
- Choose based on your business logic requirements

### 5. List All - Fetching Multiple Rows

```rust
pub async fn list(&self) -> Result<Vec<Label>, sqlx::Error> {
    let labels = sqlx::query_as!(
        Label,
        r#"
        SELECT label_id, label_name, created_at, updated_at
        FROM labels
        ORDER BY created_at DESC   // Most recent first
        "#
    )
    .fetch_all(&self.pool)         // Returns Vec<Label> - empty vec if no rows
    .await?;
    
    Ok(labels)
}
```

**Performance Considerations:**
- `fetch_all` loads all results into memory at once
- For large datasets, consider pagination or streaming with `fetch()`
- Add `LIMIT` and `OFFSET` for pagination

### 6. Update - Modifying Existing Data

```rust
pub async fn update(&self, label_id: Uuid, label_name: String) -> Result<Label, sqlx::Error> {
    let label = sqlx::query_as!(
        Label,
        r#"
        UPDATE labels
        SET label_name = $2,
            updated_at = now()     // Use database function for consistency
        WHERE label_id = $1
        RETURNING label_id, label_name, created_at, updated_at
        "#,
        label_id,
        label_name
    )
    .fetch_one(&self.pool)         // Expects the row to exist
    .await?;
    
    Ok(label)
}
```

**Update Pattern:**
- Always use RETURNING to get the updated data
- `now()` ensures consistent timestamps from the database
- Consider optimistic locking with a version field for concurrent updates

### 7. Delete - Removing Data

```rust
pub async fn delete(&self, label_id: Uuid) -> Result<bool, sqlx::Error> {
    let result = sqlx::query!(     // Note: query! not query_as! (no data returned)
        r#"
        DELETE FROM labels
        WHERE label_id = $1
        "#,
        label_id
    )
    .execute(&self.pool)           // Returns QueryResult with metadata
    .await?;
    
    Ok(result.rows_affected() > 0) // Convert to boolean: true if deleted
}
```

**Delete Patterns:**
- `execute()` for queries that don't return data
- `rows_affected()` tells you how many rows were modified
- Consider soft deletes (adding a `deleted_at` column) for audit trails

### 8. Search - Pattern Matching

```rust
pub async fn search(&self, search_term: &str) -> Result<Vec<Label>, sqlx::Error> {
    let pattern = format!("%{}%", search_term);  // SQL LIKE pattern
    
    let labels = sqlx::query_as!(
        Label,
        r#"
        SELECT label_id, label_name, created_at, updated_at
        FROM labels
        WHERE label_name ILIKE $1  // ILIKE = case-insensitive LIKE in PostgreSQL
        ORDER BY label_name
        "#,
        pattern
    )
    .fetch_all(&self.pool)
    .await?;
    
    Ok(labels)
}
```

**Search Best Practices:**
- `ILIKE` for case-insensitive search (PostgreSQL specific)
- Consider full-text search for better performance on large datasets
- Be careful with user input - SQLx handles escaping automatically
- Add indexes on searched columns for performance

```rust
}  // End of impl LabelRepository
```

## Deep Dive: SQLx Macros and Methods

### Query Macros Comparison

| Macro | Purpose | Returns | Use Case |
|-------|---------|---------|----------|
| `query!` | Basic query | Raw results | When you need custom processing |
| `query_as!` | Typed query | Your struct type | Most common - automatic mapping |
| `query_scalar!` | Single value | Primitive type | COUNT, SUM, single column |
| `query_unchecked!` | Skip compile checks | Any | Rarely used - dynamic SQL |

### Fetch Methods Explained

```rust
// 1. fetch_one() - Expects exactly one row
let user = query_as!(User, "SELECT * FROM users WHERE id = $1", id)
    .fetch_one(&pool)    // ERROR if 0 or 2+ rows
    .await?;

// 2. fetch_optional() - Maybe one row
let user = query_as!(User, "SELECT * FROM users WHERE email = $1", email)
    .fetch_optional(&pool)  // Ok(None) if no rows
    .await?;

// 3. fetch_all() - All rows as Vec
let users = query_as!(User, "SELECT * FROM users")
    .fetch_all(&pool)      // Empty Vec if no rows
    .await?;

// 4. fetch() - Stream of rows (for large datasets)
let mut stream = query_as!(User, "SELECT * FROM users")
    .fetch(&pool);
    
while let Some(user) = stream.try_next().await? {
    // Process each row without loading all into memory
}

// 5. execute() - No data returned
let result = query!("DELETE FROM users WHERE id = $1", id)
    .execute(&pool)
    .await?;
println!("Deleted {} rows", result.rows_affected());
```

## Practical Usage Example

```rust
use sqlx::postgres::PgPoolOptions;
use std::time::Duration;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // 1. Configure and create connection pool
    let pool = PgPoolOptions::new()
        .max_connections(5)                  // Maximum connections in pool
        .min_connections(1)                  // Keep at least 1 connection alive
        .connect_timeout(Duration::from_secs(10))
        .idle_timeout(Duration::from_secs(600))
        .connect("postgres://user:password@localhost/dbname")
        .await?;
    
    // 2. Run migrations (if using sqlx migrate)
    // sqlx::migrate!("./migrations").run(&pool).await?;
    
    // 3. Create repository instance
    let repo = LabelRepository::new(pool.clone());
    
    // 4. Example operations with error handling
    
    // Create a label
    match repo.create("Important".to_string()).await {
        Ok(label) => println!("Created: {:?}", label),
        Err(e) => eprintln!("Failed to create label: {}", e),
    }
    
    // Search for labels
    let search_results = repo.search("port").await?;
    println!("Found {} labels containing 'port'", search_results.len());
    
    // Update a label (if exists)
    if let Some(label) = repo.get_by_id(some_uuid).await? {
        let updated = repo.update(label.label_id, "Very Important".to_string()).await?;
        println!("Updated: {:?}", updated);
    }
    
    // Transaction example
    let mut tx = pool.begin().await?;
    // Multiple operations in a transaction
    sqlx::query!("INSERT INTO labels (label_id, label_name) VALUES ($1, $2)", id1, "Label 1")
        .execute(&mut *tx).await?;
    sqlx::query!("INSERT INTO labels (label_id, label_name) VALUES ($1, $2)", id2, "Label 2")
        .execute(&mut *tx).await?;
    tx.commit().await?;  // or tx.rollback().await?
    
    Ok(())
}
```

## Dependencies Deep Dive

```toml
[dependencies]
# SQLx with specific features
sqlx = { 
    version = "0.7", 
    features = [
        "runtime-tokio-native-tls",  # Async runtime (alternatives: runtime-async-std)
        "postgres",                  # PostgreSQL driver
        "uuid",                      # UUID type support
        "chrono",                    # DateTime support
        "migrate",                   # Migration support (optional)
        "macros",                    # query! macros (included by default)
        "json",                      # JSON/JSONB support (optional)
        "bigdecimal",               # Decimal type support (optional)
    ]
}

# Async runtime
tokio = { version = "1", features = ["full"] }  # Can use ["rt-multi-thread", "macros"] for minimal

# Types
uuid = { version = "1.6", features = ["v4", "serde"] }
chrono = { version = "0.4", features = ["serde"] }
serde = { version = "1.0", features = ["derive"] }  # For JSON serialization
```

## Environment Setup and Compile-Time Checking

### 1. Development Setup

```bash
# .env file (git-ignored)
DATABASE_URL=postgres://user:password@localhost/dbname

# For SQLx CLI
cargo install sqlx-cli --features postgres

# Create database
sqlx database create

# Run migrations
sqlx migrate run
```

### 2. Compile-Time Query Verification

SQLx validates queries at compile time. This requires:

```bash
# Option 1: Live database (default)
export DATABASE_URL=postgres://user:password@localhost/dbname
cargo build

# Option 2: Offline mode (for CI/CD)
# First, prepare query data:
cargo sqlx prepare

# This creates .sqlx/ directory with query metadata
# Check .sqlx/ into version control
# Then build without database:
SQLX_OFFLINE=true cargo build
```

### 3. Migration Management

```bash
# Create a new migration
sqlx migrate add create_labels_table

# Edit the generated file in migrations/
# Then run:
sqlx migrate run

# Revert last migration
sqlx migrate revert
```

## Common Patterns and Best Practices

### 1. Error Handling

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum LabelError {
    #[error("Label not found")]
    NotFound,
    
    #[error("Database error: {0}")]
    Database(#[from] sqlx::Error),
    
    #[error("Invalid label name")]
    InvalidName,
}

impl LabelRepository {
    pub async fn get_by_id(&self, id: Uuid) -> Result<Label, LabelError> {
        self.pool
            .get_by_id(id)
            .await?
            .ok_or(LabelError::NotFound)
    }
}
```

### 2. Testing with SQLx

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use sqlx::PgPool;
    
    #[sqlx::test]  // Automatically sets up test database
    async fn test_create_label(pool: PgPool) {
        let repo = LabelRepository::new(pool);
        
        let label = repo.create("Test".to_string()).await.unwrap();
        assert_eq!(label.label_name, "Test");
    }
}
```

### 3. Connection Pool Best Practices

```rust
// Production configuration
let pool = PgPoolOptions::new()
    .max_connections(32)
    .min_connections(5)
    .max_lifetime(Duration::from_secs(30 * 60))  // 30 minutes
    .idle_timeout(Duration::from_secs(10 * 60))  // 10 minutes
    .connect_lazy(&database_url)?;  // Lazy = don't connect until first use
```

## Key Differences from Go's Approach

| Aspect | Go (with sqlc) | Rust (with SQLx) |
|--------|----------------|------------------|
| Query validation | Build step | Compile time |
| Error handling | Explicit returns | Result<T, E> with ? |
| Null handling | sql.NullString | Option<T> |
| Connection pooling | Manual or library | Built into SQLx |
| Type mapping | Generated structs | Derive macros |
| Transactions | Begin/Commit/Rollback | RAII with Drop |

## Troubleshooting Tips

1. **"No rows returned" on INSERT**: Add `RETURNING` clause
2. **Type mismatch errors**: Check nullable columns use `Option<T>`
3. **Compile fails without database**: Use offline mode or Docker for CI
4. **Performance issues**: Enable query logging with `RUST_LOG=sqlx=debug`
5. **Connection pool exhausted**: Increase max_connections or check for leaks