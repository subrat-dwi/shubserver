package notes

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userID,omitempty"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
