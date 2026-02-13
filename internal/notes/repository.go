package notes

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository Interface
type NotesRepository interface {
	Create(ctx context.Context, note *Note) (*Note, error)
	Get(ctx context.Context, userID uuid.UUID, id string) (*Note, error)
	List(ctx context.Context, userID uuid.UUID) ([]*Note, error)
	Delete(ctx context.Context, userID uuid.UUID, id string) error
	Update(ctx context.Context, userID uuid.UUID, note *Note) error
}

// Postgres Repository
type NotesPostgresRepository struct {
	db *pgxpool.Pool
}

// Postgres Repository Constructor
func NewNotesPostgresRepository(db *pgxpool.Pool) *NotesPostgresRepository {
	return &NotesPostgresRepository{db: db}
}

// ------ CRUD Implementation on DB ------

// Create a new note in the database
func (p *NotesPostgresRepository) Create(ctx context.Context, note *Note) (*Note, error) {
	query := `
	INSERT INTO notes(user_id, title, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	if len(note.UserID) == 0 {
		return nil, errors.New("UserID not available")
	}

	err := p.db.QueryRow(ctx, query, note.UserID, note.Title, note.Content).Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	return note, err
}

// Delete a note from the database
func (p *NotesPostgresRepository) Delete(ctx context.Context, userID uuid.UUID, id string) error {
	query := `
	DELETE FROM notes
	WHERE id = $1 AND user_id = $2
	`

	cmd, err := p.db.Exec(ctx, query, id, userID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("Note not Found")
	}

	return nil
}

// Update an existing note in the database
func (p *NotesPostgresRepository) Update(ctx context.Context, userID uuid.UUID, note *Note) error {
	query := `
	UPDATE notes
	SET title = $2, content = $3
	WHERE id = $1 AND user_id = $4
	`

	cmd, err := p.db.Exec(ctx, query, note.ID, note.Title, note.Content, userID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("Note not Found")
	}

	return nil
}

// Get a specific note from the database
func (p *NotesPostgresRepository) Get(ctx context.Context, userID uuid.UUID, id string) (*Note, error) {
	query := `
    SELECT id, title, content, created_at, updated_at
    FROM notes
    WHERE id = $1 AND user_id = $2
    `
	var n Note
	err := p.db.QueryRow(ctx, query, id, userID).Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

// List all notes for a specific user from the database
func (p *NotesPostgresRepository) List(ctx context.Context, userID uuid.UUID) ([]*Note, error) {
	query := `
	SELECT id, title, content, updated_at
	FROM notes
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*Note

	for rows.Next() {
		var n Note

		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.UpdatedAt); err != nil {
			return nil, err
		}

		notes = append(notes, &n)
	}

	return notes, rows.Err()
}

// ----------------------------------------
// Memory Repository
// --------------Not Updated, Just for reference-----------------

type MemoryRepository struct {
	data map[uuid.UUID]*Note
	mu   sync.RWMutex
}

// Memory Repository Constructor
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[uuid.UUID]*Note),
	}
}

// ------ CRUD Implementation on Memory ------
func (m *MemoryRepository) Create(ctx context.Context, note *Note) (*Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[note.ID] = note
	return note, nil
}

func (m *MemoryRepository) Get(ctx context.Context, id uuid.UUID) (*Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	note, exists := m.data[id]

	if !exists {
		return nil, fmt.Errorf("not found")
	}

	return note, nil
}

func (m *MemoryRepository) List(ctx context.Context) ([]*Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list := make([]*Note, 0, len(m.data))

	for _, item := range m.data {
		list = append(list, item)
	}

	return list, nil
}

func (m *MemoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, id)
	return nil
}

func (m *MemoryRepository) Update(ctx context.Context, note *Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[note.ID]; !exists {
		return fmt.Errorf("note doesn't exist")
	}

	m.data[note.ID] = note
	return nil
}
