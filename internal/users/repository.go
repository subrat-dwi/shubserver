package users

import "context"

// UsersRepository defines the interface for user-related database operations
type UsersRepository interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (*UserDB, error)
	GetByEmail(ctx context.Context, email string) (*UserDB, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
