package auth

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/subrat-dwi/shubserver/internal/utils"
)

type AuthHandler struct {
	authservice *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authservice: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.authservice.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := RegisterResponse{
		ID:    user.Id,
		Email: user.Email,
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authservice.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
