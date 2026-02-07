package health

import "github.com/go-chi/chi/v5"

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", Health)

	return r
}
