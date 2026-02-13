package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subrat-dwi/shubserver/internal/auth"
	"github.com/subrat-dwi/shubserver/internal/health"
	"github.com/subrat-dwi/shubserver/internal/notes"
	passwordmanager "github.com/subrat-dwi/shubserver/internal/password-manager"

	"github.com/subrat-dwi/shubserver/internal/users"
)

func Routes(db *pgxpool.Pool, version, env string) chi.Router {

	// Initialize repositories and services
	userRepo := users.NewUsersPostgresRepository(db)
	notesRepo := notes.NewNotesPostgresRepository(db)
	passwordRepo := passwordmanager.NewPasswordsPostgresRepository(db)
	// notesRepo := notes.NewMemoryRepository() // Use in-memory repository for testing

	// Initialize services and handlers
	authService := auth.NewAuthService(userRepo)
	passwordService := passwordmanager.NewPasswordService(passwordRepo)

	// Initialize handlers
	authHandler := auth.NewAuthHandler(authService)
	notesHandler := notes.NewNotesHandler(notesRepo)
	passwordHandler := passwordmanager.NewPasswordHandler(passwordService)

	// Set up the router
	r := chi.NewRouter()

	// Mount routes
	r.Mount("/health", health.Routes(db, version, env))
	r.Mount("/users", auth.Routes(authHandler))
	r.Mount("/notes", notes.Routes(notesHandler))
	r.Mount("/passwords", passwordmanager.Routes(passwordHandler))

	return r
}
