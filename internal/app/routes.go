package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/subrat-dwi/shubserver/internal/health"
	"github.com/subrat-dwi/shubserver/internal/notes"
)

func Routes() chi.Router {

	notesRepo := notes.NewMemoryRepository()
	notesHandler := notes.NewHandler(notesRepo)

	r := chi.NewRouter()

	r.Mount("/health", health.Routes())
	r.Mount("/notes", notes.Routes(notesHandler))

	return r
}
