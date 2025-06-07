-- Enable extensions if needed (for UUIDs elsewhere, not used here)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: lumes
CREATE TABLE
    IF NOT EXISTS lume (
        id BIGSERIAL PRIMARY KEY,
        lume_id UUID not NULL,
        lumo_id UUID NOT NULL,
        label TEXT NOT NULL,
        type TEXT NOT NULL,
        description TEXT NULL,
        metadata JSONB NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Index on lumo_id to speed up listing under a given Lumo
CREATE INDEX IF NOT EXISTS idx_lumes_lumo_id ON lume (lumo_id);