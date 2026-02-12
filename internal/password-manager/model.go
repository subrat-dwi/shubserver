package passwordmanager

import (
	"time"

	"github.com/google/uuid"
)

type Password struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Name           string
	Username       string
	Ciphertext     []byte
	Nonce          []byte
	EncryptVersion int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
