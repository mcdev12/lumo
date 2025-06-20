use crate::label::model::{Label, NewLabel};
use sqlx::PgPool;

pub struct LabelRepository {
    pool: PgPool,
}

impl LabelRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    pub async fn create(&self, new_label: NewLabel) -> Result<Label, sqlx::Error> {
        let label = sqlx::query_as!(
            Label,
            r#"
        INSERT INTO labels (label_id, label_name)
        VALUES ($1, $2)
        RETURNING id, label_id, label_name, created_at, updated_at
        "#,
            new_label.label_id,
            new_label.label_name
        )
        .fetch_one(&self.pool)
        .await?;

        Ok(label)
    }
}
