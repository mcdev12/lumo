-- Enable extensions for UUIDs and PostGIS if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: lume (new schema)
CREATE TABLE IF NOT EXISTS lume (
    -- Internal database ID
    id BIGSERIAL PRIMARY KEY,
    -- Unique identifier (UUID)
    lume_id UUID NOT NULL UNIQUE,
    -- Unique identifier reference to lumo id (UUID)
    lumo_id UUID NOT NULL UNIQUE,
    -- Lume type
    type TEXT NOT NULL,
    -- Display name
    name TEXT NOT NULL,
    -- Optional scheduling dates
    date_start TIMESTAMPTZ NULL,
    date_end TIMESTAMPTZ NULL,
    -- Optional GPS coordinates
    latitude DOUBLE PRECISION NULL,
    longitude DOUBLE PRECISION NULL,
    -- Optional textual address
    address TEXT NULL,
    -- Freeform description
    description TEXT NULL,
    -- Array of image URLs
    images TEXT[] NULL DEFAULT '{}',
    -- Array of category tags
    category_tags TEXT[] NULL DEFAULT '{}',
    -- Optional booking link
    booking_link TEXT NULL,
    -- System timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_lume_lume_id ON lume (lume_id);
CREATE INDEX IF NOT EXISTS idx_lume_lumo ON lume (lumo_id);
CREATE INDEX IF NOT EXISTS idx_lume_type ON lume (type);
CREATE INDEX IF NOT EXISTS idx_lume_created_at ON lume (created_at);

-- Optional: Spatial index if you plan to do geo queries
-- CREATE INDEX IF NOT EXISTS idx_lume_location ON lume USING GIST (point(longitude, latitude));

-- Optional: GIN index for array searches
CREATE INDEX IF NOT EXISTS idx_lume_category_tags ON lume USING GIN (category_tags);