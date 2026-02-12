package health

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes sets up health check routes
func Routes(db *pgxpool.Pool, version, env string) chi.Router {
	r := chi.NewRouter()

	handler := NewHealthHandler(db, version, env)

	// Basic health check (lightweight - used by Render for auto-restart detection)
	r.Get("/", handler.Health)

	// Detailed health check (for monitoring tools like Uptime Robot, Datadog, etc.)
	r.Get("/detailed", handler.Detailed)

	// Status endpoint (for status page integrations)
	r.Get("/status", handler.Status)

	return r
}
