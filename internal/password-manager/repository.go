package passwordmanager

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Password Repository Interface
type PasswordRepository interface {
	Create(ctx context.Context, password *Password) (*Password, error)
	List(ctx context.Context, userID uuid.UUID) ([]*Password, error)
	Get(ctx context.Context, userID, passwordID uuid.UUID) (*Password, error)
	Update(ctx context.Context, password *Password) (*Password, error)
	Delete(ctx context.Context, userID, passwordID uuid.UUID) error
	Search(ctx context.Context, userID uuid.UUID, searchQuery string) ([]*Password, error)
}

// Postgres Password Repo
type PasswordsPostgresRepository struct {
	db *pgxpool.Pool
}

// Postgres repo constructor
func NewPasswordsPostgresRepository(db *pgxpool.Pool) *PasswordsPostgresRepository {
	return &PasswordsPostgresRepository{db: db}
}

// CRUD implementation on DB

func (p *PasswordsPostgresRepository) Create(ctx context.Context, password *Password) (*Password, error) {
	query := `
	INSERT INTO passwords(user_id, name, username, ciphertext, nonce)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, encrypt_version, created_at, updated_at
	`
	var created Password

	err := p.db.QueryRow(ctx, query,
		password.UserID,
		password.Name,
		password.Username,
		password.Ciphertext,
		password.Nonce,
	).Scan(&created.ID, &created.EncryptVersion, &created.CreatedAt, &created.UpdatedAt)

	if err != nil {
		return nil, err
	}

	created.UserID = password.UserID
	created.Name = password.Name
	created.Username = password.Username
	created.Ciphertext = password.Ciphertext
	created.Nonce = password.Nonce

	return &created, nil
}

func (p *PasswordsPostgresRepository) List(ctx context.Context, userID uuid.UUID) ([]*Password, error) {
	query := `
	SELECT id, user_id, name, username, ciphertext, nonce, encrypt_version, created_at, updated_at
	FROM passwords
	WHERE user_id = $1
	ORDER BY created_at DESC
	`
	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []*Password
	for rows.Next() {
		var password Password
		if err := rows.Scan(
			&password.ID,
			&password.UserID,
			&password.Name,
			&password.Username,
			&password.Ciphertext,
			&password.Nonce,
			&password.EncryptVersion,
			&password.CreatedAt,
			&password.UpdatedAt,
		); err != nil {
			return nil, err
		}
		passwords = append(passwords, &password)
	}

	return passwords, nil
}

func (p *PasswordsPostgresRepository) Get(ctx context.Context, userID, passwordID uuid.UUID) (*Password, error) {
	query := `
	SELECT id, user_id, name, username, ciphertext, nonce, encrypt_version, created_at, updated_at
	FROM passwords
	WHERE id = $1 AND user_id = $2
	`
	var password Password
	err := p.db.QueryRow(ctx, query, passwordID, userID).Scan(
		&password.ID,
		&password.UserID,
		&password.Name,
		&password.Username,
		&password.Ciphertext,
		&password.Nonce,
		&password.EncryptVersion,
		&password.CreatedAt,
		&password.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &password, nil
}

func (p *PasswordsPostgresRepository) Update(ctx context.Context, password *Password) (*Password, error) {
	query := `
	UPDATE passwords
	SET name = $1, username = $2, ciphertext = $3, nonce = $4, encrypt_version = $5, updated_at = NOW()
	WHERE id = $6
	RETURNING id, user_id, name, username, ciphertext, nonce, encrypt_version, created_at, updated_at
	`
	var updated Password
	err := p.db.QueryRow(ctx, query,
		password.Name,
		password.Username,
		password.Ciphertext,
		password.Nonce,
		password.EncryptVersion,
		password.ID,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Name,
		&updated.Username,
		&updated.Ciphertext,
		&updated.Nonce,
		&updated.EncryptVersion,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (p *PasswordsPostgresRepository) Delete(ctx context.Context, userID, passwordID uuid.UUID) error {
	query := `
	DELETE FROM passwords
	WHERE id = $1 AND user_id = $2
	`
	_, err := p.db.Exec(ctx, query, passwordID, userID)
	return err
}
func (p *PasswordsPostgresRepository) Search(ctx context.Context, userID uuid.UUID, searchQuery string) ([]*Password, error) {
	query := `
	SELECT id, user_id, name, username, ciphertext, nonce, encrypt_version, created_at, updated_at
	FROM passwords
	WHERE (name ILIKE $1 OR username ILIKE $1) AND user_id = $2
	ORDER BY created_at DESC
	`
	rows, err := p.db.Query(ctx, query, "%"+searchQuery+"%", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []*Password
	for rows.Next() {
		var password Password
		if err := rows.Scan(
			&password.ID,
			&password.UserID,
			&password.Name,
			&password.Username,
			&password.Ciphertext,
			&password.Nonce,
			&password.EncryptVersion,
			&password.CreatedAt,
			&password.UpdatedAt,
		); err != nil {
			return nil, err
		}
		passwords = append(passwords, &password)
	}

	return passwords, nil
}
