package auth

import (
	"context"
	"crypto/rand"

	"fmt"

	"github.com/subrat-dwi/shubserver/internal/users"
	"golang.org/x/crypto/bcrypt"
)

// AuthService struct to hold the user repository
type AuthService struct {
	users *users.UsersPostgresRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(users *users.UsersPostgresRepository) *AuthService {
	return &AuthService{users: users}
}

// GenerateSalt generates a random salt of the specified length.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16) // 16 bytes salt
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// Register registers a new user and returns the user and JWT token
func (a *AuthService) Register(ctx context.Context, email string, password string) (*users.User, string, error) {
	// check if user exists
	existing, err := a.users.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", fmt.Errorf("email already registered")
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// generate salt (not used in this implementation, BUT stored with the user for future by Clients)
	salt, _ := GenerateSalt()

	// create user in DB
	user, err := a.users.CreateUser(ctx, email, string(passwordHash), salt)
	if err != nil {
		return nil, "", err
	}

	// generate JWT
	token, err := GenerateToken(user.Id.String())
	if err != nil {
		return nil, "", err
	}

	return &users.User{
		Id:        user.Id,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, token, nil
}

// Login authenticates a user and returns a JWT token if successful
func (a *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	// check if email is registered
	user, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := GenerateToken(user.Id.String())
	if err != nil {
		return "", err
	}

	return token, nil
}
