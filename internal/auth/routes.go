package auth

import "github.com/go-chi/chi/v5"

func Routes(h *AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.registerUser)
	r.Post("/login", h.loginUser)

	return r
}
