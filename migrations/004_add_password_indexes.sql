-- +migrate Up
-- Always needed
CREATE INDEX IF NOT EXISTS idx_passwords_user_id
ON passwords(user_id);

-- Trigram index for partial search, scoped by user
CREATE INDEX IF NOT EXISTS idx_passwords_user_name_trgm
ON passwords
USING gin (user_id, name gin_trgm_ops);

-- +migrate Down

DROP INDEX IF EXISTS idx_passwords_user_name_trgm;
DROP INDEX IF EXISTS idx_passwords_user_id;
