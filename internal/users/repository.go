package users

import "context"

type UsersRepository interface {
	CreateUser(ctx context.Context, user *User) (*UserDB, error)
	GetByEmail(ctx context.Context, email string) (*UserDB, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
