package notes

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.listNotes)
	r.Get("/{id}", h.getNote)
	r.Post("/", h.createNote)
	r.Delete("/{id}", h.deleteNote)
	r.Put("/{id}", h.updateNote)

	return r
}
