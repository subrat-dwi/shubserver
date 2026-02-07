package notes

import (
	"fmt"
	"sync"
)

type Repository interface {
	Create(note *Note) error
	Get(id string) (*Note, error)
	List() ([]*Note, error)
	Delete(id string) error
	Update(note *Note) error
}

type MemoryRepository struct {
	data map[string]*Note
	mu   sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]*Note),
	}
}

func (m *MemoryRepository) Create(note *Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[note.ID] = note
	return nil
}

func (m *MemoryRepository) Get(id string) (*Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	note, exists := m.data[id]

	if !exists {
		return nil, fmt.Errorf("not found")
	}

	return note, nil
}

func (m *MemoryRepository) List() ([]*Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	list := make([]*Note, 0, len(m.data))

	for _, item := range m.data {
		list = append(list, item)
	}

	return list, nil
}

func (m *MemoryRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, id)
	return nil
}

func (m *MemoryRepository) Update(note *Note) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[note.ID]; !exists {
		return fmt.Errorf("note doesn't exist")
	}

	m.data[note.ID] = note
	return nil
}
