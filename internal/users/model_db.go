package users

import (
	"time"

	"github.com/google/uuid"
)

type UserDB struct {
	Id           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
