-- Enable extensions for UUIDs if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: link
CREATE TABLE IF NOT EXISTS link (
    -- Internal database ID
    id BIGSERIAL PRIMARY KEY,
    -- External UUID for API clients
    link_id UUID NOT NULL UNIQUE,
    -- Which two Lumés this edge connects
    from_lume_id UUID NOT NULL,
    to_lume_id UUID NOT NULL,
    -- The high-level type of relation
    link_type TEXT NOT NULL,
    -- Travel details (JSON for flexibility)
    travel_details JSONB NULL,
    -- Freeform notes about this relationship
    notes TEXT NULL,
    -- Optional hint for rendering order in lists/timelines
    sequence_index INTEGER NULL,
    -- Audit timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Ensure we don't have duplicate connections
    UNIQUE(from_lume_id, to_lume_id, link_type)
);

-- Foreign key constraints
ALTER TABLE link
    ADD CONSTRAINT fk_link_from_lume
    FOREIGN KEY (from_lume_id)
    REFERENCES lume(lume_id)
    ON DELETE CASCADE;

ALTER TABLE link
    ADD CONSTRAINT fk_link_to_lume
    FOREIGN KEY (to_lume_id)
    REFERENCES lume(lume_id)
    ON DELETE CASCADE;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_link_link_id ON link (link_id);
CREATE INDEX IF NOT EXISTS idx_link_from_lume_id ON link (from_lume_id);
CREATE INDEX IF NOT EXISTS idx_link_to_lume_id ON link (to_lume_id);
CREATE INDEX IF NOT EXISTS idx_link_type ON link (link_type);
CREATE INDEX IF NOT EXISTS idx_link_created_at ON link (created_at);