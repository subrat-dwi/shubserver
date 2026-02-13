package passwordmanager

import (
	"github.com/go-chi/chi/v5"
	"github.com/subrat-dwi/shubserver/internal/middleware"
)

func Routes(h *PasswordHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)

	// Mounted at /passwords in app.Routes, so use relative paths here.
	r.Get("/", h.listPasswords)
	r.Post("/", h.createPassword)
	r.Get("/{id}", h.getPassword)
	r.Put("/{id}", h.updatePassword)
	r.Delete("/{id}", h.deletePassword)

	return r
}
