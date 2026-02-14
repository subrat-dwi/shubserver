package users

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user data structure used in the application
type User struct {
	Id        uuid.UUID
	Email     string
	Salt      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
