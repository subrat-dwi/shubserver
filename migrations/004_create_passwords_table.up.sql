CREATE TABLE IF NOT EXISTS passwords (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    username TEXT NOT NULL,
    ciphertext BYTEA NOT NULL,
    nonce BYTEA NOT NULL,
    encrypt_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint with cascade delete
    CONSTRAINT fk_passwords_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    
    -- Data validation constraints
    CONSTRAINT name_not_empty CHECK (name != ''),
    CONSTRAINT username_not_empty CHECK (username != ''),
    CONSTRAINT name_max_length CHECK (LENGTH(name) <= 255),
    CONSTRAINT username_max_length CHECK (LENGTH(username) <= 255),
    CONSTRAINT ciphertext_not_empty CHECK (OCTET_LENGTH(ciphertext) > 0),
    CONSTRAINT ciphertext_max_size CHECK (OCTET_LENGTH(ciphertext) <= 1048576), -- 1 MB
    CONSTRAINT nonce_exact_size CHECK (OCTET_LENGTH(nonce) = 12), -- GCM standard: 96 bits
);

-- Index for user's passwords lookup (most common query)
CREATE INDEX IF NOT EXISTS idx_passwords_user_id 
ON passwords(user_id);

-- Trigram index for service name search within user's passwords
CREATE INDEX IF NOT EXISTS idx_passwords_name_username_trgm 
ON passwords 
USING gin (name, username gin_trgm_ops);

-- Comments for documentation
COMMENT ON TABLE passwords IS 'Encrypted password storage - server stores ciphertext only, never decrypts';
COMMENT ON COLUMN passwords.id IS 'Unique identifier (UUID v4)';
COMMENT ON COLUMN passwords.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN passwords.name IS 'Service/application name (max 255 characters)';
COMMENT ON COLUMN passwords.username IS 'Account username/email (max 255 characters)';
COMMENT ON COLUMN passwords.ciphertext IS 'AES-256-GCM encrypted password data (max 1 MB)';
COMMENT ON COLUMN passwords.nonce IS 'Encryption nonce - exactly 12 bytes for GCM mode';
COMMENT ON COLUMN passwords.encrypt_version IS 'Encryption algorithm version (for future migrations)';
COMMENT ON COLUMN passwords.created_at IS 'Entry creation timestamp';
COMMENT ON COLUMN passwords.updated_at IS 'Last modification timestamp';