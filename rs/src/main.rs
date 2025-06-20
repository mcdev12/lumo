mod label;

use dotenvy::dotenv;
use sqlx::{migrate::Migrator, PgPool};

static MIGRATOR: Migrator = sqlx::migrate!();

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenv().ok();

    match std::env::var("DATABASE_URL") {
        Ok(db_url) => println!("Loaded DATABASE_URL: {}", db_url),
        Err(e) => println!("❌ Could not load DATABASE_URL: {}", e),
    }

    let db_url = std::env::var("DATABASE_URL")?;
    let pool = PgPool::connect(&db_url).await?;

    MIGRATOR.run(&pool).await?;

    println!("DB Migrations Applied!");
    
    Ok(())
}

