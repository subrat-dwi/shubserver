package notes

import "time"

type Note struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userID,omitempty"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
