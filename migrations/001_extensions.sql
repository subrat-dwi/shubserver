-- Needed for uuid generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Needed for ILIKE %search%
CREATE EXTENSION IF NOT EXISTS pg_trgm;