package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/subrat-dwi/shubserver/internal/health"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Mount("/health", health.Routes())

	return r
}
