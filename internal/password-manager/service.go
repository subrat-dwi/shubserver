package passwordmanager

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// PasswordService struct
type PasswordService struct {
	repo PasswordRepository
}

// NewPasswordService creates a new password manager service
func NewPasswordService(repo PasswordRepository) *PasswordService {
	return &PasswordService{repo: repo}
}

// CreatePassword creates a new password entry
func (s *PasswordService) CreatePassword(ctx context.Context, password *Password) (*Password, error) {

	if password.UserID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	if password.Name == "" {
		return nil, errors.New("password name cannot be empty")
	}
	if password.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if len(password.Ciphertext) == 0 {
		return nil, errors.New("ciphertext cannot be empty")
	}
	if len(password.Nonce) == 0 {
		return nil, errors.New("nonce cannot be empty")
	}

	return s.repo.Create(ctx, password)
}

// ListPasswords lists all password entries for a user
func (s *PasswordService) ListPasswords(ctx context.Context, userID uuid.UUID) ([]*Password, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	return s.repo.List(ctx, userID)
}

// GetPassword retrieves a specific password entry by ID
func (s *PasswordService) GetPassword(ctx context.Context, userID, passwordID uuid.UUID) (*Password, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	if passwordID == uuid.Nil {
		return nil, errors.New("password ID cannot be empty")
	}
	return s.repo.Get(ctx, userID, passwordID)
}

// UpdatePassword updates an existing password entry
func (s *PasswordService) UpdatePassword(ctx context.Context, password *Password) (*Password, error) {
	if password.UserID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}
	if password.ID == uuid.Nil {
		return nil, errors.New("password ID cannot be empty")
	}
	if password.Name == "" {
		return nil, errors.New("password name cannot be empty")
	}
	if password.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if len(password.Ciphertext) == 0 {
		return nil, errors.New("ciphertext cannot be empty")
	}
	if len(password.Nonce) == 0 {
		return nil, errors.New("nonce cannot be empty")
	}
	return s.repo.Update(ctx, password)
}

// DeletePassword deletes a password entry by ID
func (s *PasswordService) DeletePassword(ctx context.Context, userID, passwordID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("user ID cannot be empty")
	}
	if passwordID == uuid.Nil {
		return errors.New("password ID cannot be empty")
	}
	return s.repo.Delete(ctx, userID, passwordID)
}
