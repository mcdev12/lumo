-- 001_lumo_schema.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    IF NOT EXISTS lumo (
        id BIGSERIAL PRIMARY KEY,
        lumo_id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4 (),
        user_id UUID NOT NULL,
        title TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- trigger to keep updated_at in sync omitted for brevity
CREATE INDEX IF NOT EXISTS idx_lumo_user_id ON lumo (user_id);

CREATE INDEX IF NOT EXISTS idx_lumo_created_at ON lumo (created_at);