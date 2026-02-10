package auth

import (
	"context"

	"fmt"

	"github.com/subrat-dwi/shubserver/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users *users.UsersPostgresRepository
}

func NewAuthService(users *users.UsersPostgresRepository) *AuthService {
	return &AuthService{users: users}
}

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

	// create user in DB
	user, err := a.users.CreateUser(ctx, email, string(passwordHash))
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
