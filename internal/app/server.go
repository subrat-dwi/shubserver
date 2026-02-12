package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subrat-dwi/shubserver/internal/config"
)

// Server struct to hold the router and address
type Server struct {
	Router chi.Router
	Addr   string
}

func Setup(db *pgxpool.Pool, version, env string) *Server {
	// Set up the router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve the index.html file at the root path
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	// Serve static files from the ./web/static directory
	fs := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Mount API routes
	r.Mount("/api", Routes(db, version, env))

	// Return the server instance
	return &Server{
		Router: r,
		Addr:   config.Load().Port,
	}
}
