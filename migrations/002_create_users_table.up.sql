CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    salt BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT email_not_empty CHECK (email != ''),
    CONSTRAINT password_hash_not_empty CHECK (password_hash != '')
);

-- Index for email lookups (login queries)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Comment for documentation
COMMENT ON TABLE users IS 'Core user accounts table';
COMMENT ON COLUMN users.id IS 'Unique identifier (UUID v4)';
COMMENT ON COLUMN users.email IS 'Unique email address for login';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password (never store plaintext)';
COMMENT ON COLUMN users.salt IS 'Random salt for password hashing (stored for use by clients)';
COMMENT ON COLUMN users.created_at IS 'Account creation timestamp';
COMMENT ON COLUMN users.updated_at IS 'Last update timestamp';
