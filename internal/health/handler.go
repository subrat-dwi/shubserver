package health

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subrat-dwi/shubserver/internal/utils"
)

// HealthStatus represents the overall system health
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusDegraded  HealthStatus = "degraded"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status      HealthStatus           `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Checks      map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a single health check
type CheckResult struct {
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
	Duration  float64      `json:"duration_ms"`
}

// HealthHandler holds dependencies for health checks
type HealthHandler struct {
	db      *pgxpool.Pool
	version string
	env     string // "development", "staging", "production"
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *pgxpool.Pool, version, env string) *HealthHandler {
	return &HealthHandler{
		db:      db,
		version: version,
		env:     env,
	}
}

// Health performs a basic health check (quick response for Render health checks)
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"version": h.version,
		"env":     h.env,
	})
}

// Detailed performs detailed system health checks (for monitoring dashboards)
func (h *HealthHandler) Detailed(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]CheckResult)

	// Database connectivity check
	checks["supabase"] = h.checkSupabaseDatabase()

	// Determine overall status
	overallStatus := h.determineOverallStatus(checks)

	response := HealthResponse{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		Version:     h.version,
		Environment: h.env,
		Checks:      checks,
	}

	var statusCode int
	switch overallStatus {
	case StatusHealthy:
		statusCode = http.StatusOK
	case StatusDegraded:
		statusCode = http.StatusOK // Still OK, but degraded
	case StatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusInternalServerError
	}

	utils.JSON(w, statusCode, response)
}

// Status returns current system status (for status page integrations)
func (h *HealthHandler) Status(w http.ResponseWriter, r *http.Request) {
	checks := make(map[string]CheckResult)
	checks["supabase"] = h.checkSupabaseDatabase()

	overallStatus := h.determineOverallStatus(checks)

	response := HealthResponse{
		Status:      overallStatus,
		Timestamp:   time.Now(),
		Version:     h.version,
		Environment: h.env,
		Checks:      checks,
	}

	var statusCode int
	if overallStatus == StatusHealthy {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusServiceUnavailable
	}

	utils.JSON(w, statusCode, response)
}

// checkSupabaseDatabase checks Supabase PostgreSQL connectivity and responsiveness
func (h *HealthHandler) checkSupabaseDatabase() CheckResult {
	start := time.Now()
	result := CheckResult{
		Timestamp: start,
	}

	// Create a context with timeout (Supabase can be slow sometimes)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test database connectivity with a simple query
	var version string
	err := h.db.QueryRow(ctx, "SELECT version()").Scan(&version)
	duration := time.Since(start).Seconds() * 1000 // Convert to milliseconds

	result.Duration = duration

	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = "supabase connection failed: " + err.Error()
		return result
	}

	// Check if response time is acceptable
	// Supabase can have higher latency, so we're more lenient
	if duration > 500 {
		result.Status = StatusDegraded
		result.Message = "supabase response time is slow"
	} else if duration > 200 {
		result.Status = StatusDegraded
		result.Message = "supabase experiencing minor latency"
	} else {
		result.Status = StatusHealthy
		result.Message = "supabase is responsive"
	}

	return result
}

// determineOverallStatus determines overall health based on individual checks
func (h *HealthHandler) determineOverallStatus(checks map[string]CheckResult) HealthStatus {
	hasUnhealthy := false
	hasDegraded := false

	for _, check := range checks {
		switch check.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}
	if hasDegraded {
		return StatusDegraded
	}
	return StatusHealthy
}
