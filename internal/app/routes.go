package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subrat-dwi/shubserver/internal/auth"
	"github.com/subrat-dwi/shubserver/internal/health"
	"github.com/subrat-dwi/shubserver/internal/notes"
	"github.com/subrat-dwi/shubserver/internal/users"
)

func Routes(db *pgxpool.Pool, version, env string) chi.Router {

	// Initialize repositories and services
	userRepo := users.NewUsersPostgresRepository(db)
	notesRepo := notes.NewNotesPostgresRepository(db)
	// notesRepo := notes.NewMemoryRepository() // Use in-memory repository for testing

	// Initialize services and handlers
	authService := auth.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := auth.NewAuthHandler(authService)
	notesHandler := notes.NewNotesHandler(notesRepo)

	// Set up the router
	r := chi.NewRouter()

	// Mount routes
	r.Mount("/health", health.Routes(db, version, env))
	r.Mount("/users", auth.Routes(authHandler))
	r.Mount("/notes", notes.Routes(notesHandler))

	return r
}
