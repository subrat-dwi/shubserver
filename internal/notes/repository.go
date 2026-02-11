package notes

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository Interface
type NotesRepository interface {
	Create(ctx context.Context, note *Note) (*Note, error)
	Get(ctx context.Context, id string) (*Note, error)
	List(ctx context.Context) ([]*Note, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, note *Note) error
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
func (p *NotesPostgresRepository) Create(ctx context.Context, note *Note) (*Note, error) {
	query := `
	INSERT INTO notes(user_id, title, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`

	if len(note.UserID) == 0 {
		return nil, errors.New("UserID not available")
	}
	// userID, err := uuid.Parse(note.UserID)
	// if err != nil {
	// 	return err
	// }
	// userID := ctx.Value("userID")
	err := p.db.QueryRow(ctx, query, note.UserID, note.Title, note.Content).Scan(note.ID, note.CreatedAt, note.UpdatedAt)
	return note, err
}

func (p *NotesPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM notes
	WHERE id = $1
	`

	cmd, err := p.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("Note not Found")
	}

	return nil
}

func (p *NotesPostgresRepository) Update(ctx context.Context, note *Note) error {
	query := `
	UPDATE notes
	SET title = $2, content = $3
	WHERE id = $1
	`

	cmd, err := p.db.Exec(ctx, query, note.ID, note.Title, note.Content)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("Note not Found")
	}

	return nil
}

func (p *NotesPostgresRepository) Get(ctx context.Context, id string) (*Note, error) {
	query := `
	SELECT id, title, content
	FROM notes
	WHERE id = $1
	`
	var n Note
	err := p.db.QueryRow(ctx, query, id).Scan(&n.ID, &n.Title, &n.Content)

	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (p *NotesPostgresRepository) List(ctx context.Context) ([]*Note, error) {
	query := `
	SELECT id, title, content, updated_at
	FROM notes
	ORDER BY created_at DESC
	`

	rows, err := p.db.Query(ctx, query)
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
type MemoryRepository struct {
	data map[string]*Note
	mu   sync.RWMutex
}

// Memory Repository Constructor
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]*Note),
	}
}

// ------ CRUD Implementation on Memory ------
func (m *MemoryRepository) Create(ctx context.Context, note *Note) (*Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[note.ID] = note
	return note, nil
}

func (m *MemoryRepository) Get(ctx context.Context, id string) (*Note, error) {
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

func (m *MemoryRepository) Delete(ctx context.Context, id string) error {
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
