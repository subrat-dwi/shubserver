package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/subrat-dwi/shubserver/internal/config"
)

type Server struct {
	Router chi.Router
	Addr   string
}

func Setup() *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", Routes())

	return &Server{
		Router: r,
		Addr:   config.Load().Port,
	}
}
