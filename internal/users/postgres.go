package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres Repository for Users
type UsersPostgresRepository struct {
	db *pgxpool.Pool
}

// Constructor for PostgresRepository
func NewUsersPostgresRepository(db *pgxpool.Pool) *UsersPostgresRepository {
	return &UsersPostgresRepository{db: db}
}

// ------ Implementation of UserRepository ------

// CreateUser creates a new user in the database
func (p *UsersPostgresRepository) CreateUser(ctx context.Context, email string, passwordHash string) (*UserDB, error) {
	query := `
	INSERT INTO users(email, password_hash)
	VALUES($1, $2)
	RETURNING id, email, password_hash, created_at, updated_at
	`
	var user UserDB

	err := p.db.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user from the database by their email
func (p *UsersPostgresRepository) GetByEmail(ctx context.Context, email string) (*UserDB, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	var user UserDB

	err := p.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user from the database by their ID
func (p *UsersPostgresRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
	SELECT id, email, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user User

	err := p.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
