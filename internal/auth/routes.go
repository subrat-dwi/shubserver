package auth

import "github.com/go-chi/chi/v5"

// Routes sets up the routes for the auth handler
func Routes(h *AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.registerUser)
	r.Post("/login", h.loginUser)

	return r
}
