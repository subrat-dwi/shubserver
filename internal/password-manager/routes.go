package passwordmanager

import (
	"github.com/go-chi/chi/v5"
	"github.com/subrat-dwi/shubserver/internal/middleware"
)

func Routes(h *PasswordHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware)

	r.Get("/passwords", h.listPasswords)
	r.Post("/passwords", h.createPassword)
	r.Get("/passwords/{id}", h.getPassword)
	r.Put("/passwords/{id}", h.updatePassword)
	r.Delete("/passwords/{id}", h.deletePassword)

	return r
}
