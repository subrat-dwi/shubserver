package passwordmanager

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Validation constants
const (
	MinNameLength     = 1
	MaxNameLength     = 255
	MinUsernameLength = 1
	MaxUsernameLength = 255
	MaxCiphertextSize = 1024 * 1024 // 1 MB
	ExpectedNonceSize = 12          // GCM standard: 96 bits = 12 bytes
	MinCiphertextSize = 1           // At least 1 byte
	MaxEncryptVersion = 1
)

// PasswordService struct
type PasswordService struct {
	repo PasswordRepository
}

// NewPasswordService creates a new password manager service
func NewPasswordService(repo PasswordRepository) *PasswordService {
	return &PasswordService{repo: repo}
}

// validatePasswordInput validates all password fields according to security standards
func (s *PasswordService) validatePasswordInput(password *Password) error {
	// Validate UserID
	if password.UserID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Validate Name
	if err := s.validateName(password.Name); err != nil {
		return err
	}

	// Validate Username
	if err := s.validateUsername(password.Username); err != nil {
		return err
	}

	// Validate Ciphertext
	if err := s.validateCiphertext(password.Ciphertext); err != nil {
		return err
	}

	// Validate Nonce
	if err := s.validateNonce(password.Nonce); err != nil {
		return err
	}

	// Validate EncryptVersion
	if err := s.validateEncryptVersion(password.EncryptVersion); err != nil {
		return err
	}

	return nil
}

// validateName validates the password name/service name
func (s *PasswordService) validateName(name string) error {
	trimmedName := strings.TrimSpace(name)

	if trimmedName == "" {
		return fmt.Errorf("password name cannot be empty")
	}

	if len(trimmedName) < MinNameLength {
		return fmt.Errorf("password name must be at least %d character(s), got %d", MinNameLength, len(trimmedName))
	}

	if len(trimmedName) > MaxNameLength {
		return fmt.Errorf("password name must not exceed %d characters, got %d", MaxNameLength, len(trimmedName))
	}

	return nil
}

// validateUsername validates the username/account identifier
func (s *PasswordService) validateUsername(username string) error {
	trimmedUsername := strings.TrimSpace(username)

	if trimmedUsername == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(trimmedUsername) < MinUsernameLength {
		return fmt.Errorf("username must be at least %d character(s), got %d", MinUsernameLength, len(trimmedUsername))
	}

	if len(trimmedUsername) > MaxUsernameLength {
		return fmt.Errorf("username must not exceed %d characters, got %d", MaxUsernameLength, len(trimmedUsername))
	}

	return nil
}

// validateCiphertext validates the encrypted password data
func (s *PasswordService) validateCiphertext(ciphertext []byte) error {
	if len(ciphertext) == 0 {
		return fmt.Errorf("ciphertext cannot be empty")
	}

	if len(ciphertext) < MinCiphertextSize {
		return fmt.Errorf("ciphertext must be at least %d byte(s), got %d", MinCiphertextSize, len(ciphertext))
	}

	if len(ciphertext) > MaxCiphertextSize {
		return fmt.Errorf("ciphertext must not exceed %d bytes (1 MB), got %d bytes", MaxCiphertextSize, len(ciphertext))
	}

	return nil
}

// validateNonce validates the encryption nonce (GCM expects exactly 12 bytes)
func (s *PasswordService) validateNonce(nonce []byte) error {
	if len(nonce) == 0 {
		return fmt.Errorf("nonce cannot be empty")
	}

	if len(nonce) != ExpectedNonceSize {
		return fmt.Errorf("nonce must be exactly %d bytes for AES-GCM, got %d bytes", ExpectedNonceSize, len(nonce))
	}

	return nil
}

// validateEncryptVersion validates the encryption version for future compatibility
func (s *PasswordService) validateEncryptVersion(version int) error {
	if version <= 0 {
		return fmt.Errorf("encryption version must be greater than 0, got %d", version)
	}

	if version > MaxEncryptVersion {
		return fmt.Errorf("unsupported encryption version %d (current supported: %d)", version, MaxEncryptVersion)
	}

	return nil
}

// CreatePassword creates a new password entry
func (s *PasswordService) CreatePassword(ctx context.Context, password *Password) (*Password, error) {
	// Validate all input fields
	if err := s.validatePasswordInput(password); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, password)
}

// ListPasswords lists all password entries for a user
func (s *PasswordService) ListPasswords(ctx context.Context, userID uuid.UUID, searchQuery string) ([]*Password, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	return s.repo.List(ctx, userID, searchQuery)
}

// GetPassword retrieves a specific password entry by ID
func (s *PasswordService) GetPassword(ctx context.Context, userID, passwordID uuid.UUID) (*Password, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	if passwordID == uuid.Nil {
		return nil, fmt.Errorf("password ID cannot be empty")
	}
	return s.repo.Get(ctx, userID, passwordID)
}

// UpdatePassword updates an existing password entry
func (s *PasswordService) UpdatePassword(ctx context.Context, password *Password) (*Password, error) {
	// Validate ID fields
	if password.UserID == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	if password.ID == uuid.Nil {
		return nil, fmt.Errorf("password ID cannot be empty")
	}

	// Validate all input fields
	if err := s.validatePasswordInput(password); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, password)
}

// DeletePassword deletes a password entry by ID
func (s *PasswordService) DeletePassword(ctx context.Context, userID, passwordID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf("user ID cannot be empty")
	}
	if passwordID == uuid.Nil {
		return fmt.Errorf("password ID cannot be empty")
	}
	return s.repo.Delete(ctx, userID, passwordID)
}
