package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/subrat-dwi/shubserver/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// Expecting header format: "Bearer <token>"
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
			return
		}

		// Verify the token and extract claims
		claims, err := auth.VerifyToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		// Ensure the UserID claim is present
		if claims.UserID == "" {
			http.Error(w, "UserID is empty", http.StatusUnauthorized)
			return
		}

		// Parse userID to UUID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			http.Error(w, "Invalid UserID format", http.StatusUnauthorized)
			return
		}

		// Add userID to the request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
