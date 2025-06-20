CREATE TABLE
    labels (
               id BIGSERIAL PRIMARY KEY,
               label_id UUID UNIQUE NOT NULL,
               label_name TEXT NOT NULL,
               created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
               updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
--      typ TEXT CHECK (typ IN ('TEXT', 'NUMBER', 'BOOLEAN', 'DROPDOWN')) NOT NULL,
--      config JSONB, -- e.g. {"options": ["Food", "Museum"]} for dropdowns
--      user_id UUID NOT NULL,
--      UNIQUE (name, user_id)
);