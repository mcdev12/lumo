use serde::{Deserialize, Serialize};
use time::OffsetDateTime;
use uuid::Uuid;
use validator::Validate;

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct Label {
    pub id: i64,
    pub label_id: Uuid,
    pub label_name: String,
    pub created_at: OffsetDateTime,
    pub updated_at: OffsetDateTime,
}

#[derive(Debug, Clone, Serialize, Deserialize, Validate)]
pub struct NewLabel {
    pub label_id: Uuid,
    #[validate(length(min = 1, max = 255))]
    pub label_name: String,
}

impl Label {
    pub fn internal_id(&self) -> i64 {
        self.id
    }
}

impl NewLabel {
    pub fn new(label_name: String) -> Self {
        Self {
            label_id: Uuid::new_v4(),
            label_name,
        }
    }
}
