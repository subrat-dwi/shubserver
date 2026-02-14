package auth

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/subrat-dwi/shubserver/internal/utils"
)

// AuthHandler struct to hold the auth service
type AuthHandler struct {
	authservice *AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authservice: authService}
}

// RegisterRequest struct to hold the registration request data
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse struct to hold the registration response data
type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Token string    `json:"token"`
	Salt  string    `json:"salt"`
}

// LoginRequest struct to hold the login request data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse struct to hold the login response data
type LoginResponse struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

// registerUser handles the user registration endpoint
func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Decode the JSON request body into the RegisterRequest struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call the Register method of the auth service to create a new user and generate a token
	user, token, err := h.authservice.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create a RegisterResponse struct with the user ID, email, and token
	resp := RegisterResponse{
		ID:    user.Id,
		Email: user.Email,
		Token: token,
		Salt:  user.Salt,
	}

	// Set the Content-Type header to application/json and encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// loginUser handles the user login endpoint
func (h *AuthHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call the Login method of the auth service to authenticate the user and generate a token
	token, salt, err := h.authservice.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := LoginResponse{
		Token: token,
		Salt:  salt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
