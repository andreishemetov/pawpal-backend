CREATE INDEX idx_pets_name ON pets(name);
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_pets_name_trgm ON pets USING gin (name gin_trgm_ops);