package notes

import (
	"github.com/go-chi/chi/v5"
	"github.com/subrat-dwi/shubserver/internal/middleware"
)

// Routes sets up the routes for the notes module
func Routes(h *NotesHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)

	r.Get("/", h.listNotes)
	r.Get("/{id}", h.getNote)
	r.Post("/", h.createNote)
	r.Delete("/{id}", h.deleteNote)
	r.Put("/{id}", h.updateNote)

	return r
}
