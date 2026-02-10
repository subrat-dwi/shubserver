package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subrat-dwi/shubserver/internal/auth"
	"github.com/subrat-dwi/shubserver/internal/health"
	"github.com/subrat-dwi/shubserver/internal/notes"
	"github.com/subrat-dwi/shubserver/internal/users"
)

func Routes(db *pgxpool.Pool) chi.Router {

	userRepo := users.NewUsersPostgresRepository(db)
	notesRepo := notes.NewNotesPostgresRepository(db)

	authService := auth.NewAuthService(userRepo)
	// notesRepo := notes.NewMemoryRepository()

	authHandler := auth.NewAuthHandler(authService)
	notesHandler := notes.NewNotesHandler(notesRepo)

	r := chi.NewRouter()

	r.Mount("/health", health.Routes())
	r.Mount("/users", auth.Routes(authHandler))
	r.Mount("/notes", notes.Routes(notesHandler))

	return r
}
