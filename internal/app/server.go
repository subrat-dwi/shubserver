package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subrat-dwi/shubserver/internal/config"
)

type Server struct {
	Router chi.Router
	Addr   string
}

func Setup(db *pgxpool.Pool) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	fs := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Mount("/api", Routes(db))

	return &Server{
		Router: r,
		Addr:   config.Load().Port,
	}
}
