-- Initial migration for campbell‑checkbook‑go
-- Add your tables here. For now we create a simple placeholder table.
CREATE TABLE IF NOT EXISTS placeholder (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
