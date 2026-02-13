package users

import (
	"time"

	"github.com/google/uuid"
)

// UserDB represents the user data structure as stored in the database
type UserDB struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
	Salt         []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
