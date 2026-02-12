CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key constraint with cascade delete
    CONSTRAINT fk_notes_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE, 
    
    -- Data validation constraints
    CONSTRAINT title_not_empty CHECK (title != ''),
    CONSTRAINT content_not_empty CHECK (content != ''),
    CONSTRAINT title_max_length CHECK (LENGTH(title) <= 255),
    CONSTRAINT content_max_length CHECK (LENGTH(content) <= 5000)
);

-- Index for user's notes lookup (most common query)
CREATE INDEX IF NOT EXISTS idx_notes_user_id 
ON notes(user_id);

-- Composite index for user-specific sorting
CREATE INDEX IF NOT EXISTS idx_notes_user_created 
ON notes(user_id, created_at DESC);

-- Trigram index for full-text search within user's notes
CREATE INDEX IF NOT EXISTS idx_notes_title_trgm 
ON notes 
USING gin (title gin_trgm_ops);

-- Comments for documentation
COMMENT ON TABLE notes IS 'User notes storage with encryption-ready design';
COMMENT ON COLUMN notes.id IS 'Unique identifier (UUID v4)';
COMMENT ON COLUMN notes.user_id IS 'Foreign key to users table';
COMMENT ON COLUMN notes.title IS 'Note title (max 255 characters)';
COMMENT ON COLUMN notes.content IS 'Note content (max 5000 characters)';
COMMENT ON COLUMN notes.created_at IS 'Creation timestamp';
COMMENT ON COLUMN notes.updated_at IS 'Last modification timestamp';