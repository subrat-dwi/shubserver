package passwordmanager

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/subrat-dwi/shubserver/internal/utils"
)

// Handler structs for API requests and responses

// PasswordItem struct for API responses
type PasswordItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Request struct for creating a password item
type CreatePasswordRequest struct {
	Name       string `json:"name" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Ciphertext string `json:"password" validate:"required"` // base64 encoded ciphertext
	Nonce      string `json:"nonce" validate:"required"`    // base64 encoded nonce
}

// Response struct for creating a password item
type CreatePasswordResponse struct {
	PasswordItem
}

// Response struct for listing password items
type ListPasswordsResponse struct {
	Passwords []PasswordItem `json:"passwords"`
}

type GetPasswordResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Ciphertext string `json:"password"` // base64 encoded ciphertext
	Nonce      string `json:"nonce"`    // base64 encoded nonce
}

// Handler struct for password manager API
type PasswordHandler struct {
	passwordService *PasswordService
}

func NewPasswordHandler(service *PasswordService) *PasswordHandler {
	return &PasswordHandler{passwordService: service}
}

// Handler function to create a new password item
func (h *PasswordHandler) createPassword(w http.ResponseWriter, r *http.Request) {
	var req CreatePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	ciphertext, err := base64.RawStdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ciphertext encoding")
		return
	}

	nonce, err := base64.RawStdEncoding.DecodeString(req.Nonce)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid nonce encoding")
		return
	}

	password := &Password{
		UserID:         userID,
		Name:           req.Name,
		Username:       req.Username,
		Ciphertext:     ciphertext,
		Nonce:          nonce,
		EncryptVersion: 1,
	}
	// Call the service layer to create the password item
	created, err := h.passwordService.CreatePassword(r.Context(), password)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := CreatePasswordResponse{
		PasswordItem: PasswordItem{
			ID:        created.ID.String(),
			Name:      created.Name,
			Username:  created.Username,
			CreatedAt: created.CreatedAt.String(),
			UpdatedAt: created.UpdatedAt.String(),
		},
	}
	utils.JSON(w, http.StatusCreated, resp)
}

// Handler function to list all password items for a user
func (h *PasswordHandler) listPasswords(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	searchQuery := r.URL.Query().Get("search")

	passwords, err := h.passwordService.ListPasswords(r.Context(), userID, searchQuery)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp ListPasswordsResponse
	for _, p := range passwords {
		resp.Passwords = append(resp.Passwords, PasswordItem{
			ID:        p.ID.String(),
			Name:      p.Name,
			Username:  p.Username,
			CreatedAt: p.CreatedAt.String(),
			UpdatedAt: p.UpdatedAt.String(),
		})
	}
	utils.JSON(w, http.StatusOK, resp)
}

// Handler function to get a specific password item by ID
func (h *PasswordHandler) getPassword(w http.ResponseWriter, r *http.Request) {
	passwordIDStr := chi.URLParam(r, "id")
	passwordID, err := uuid.Parse(passwordIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid password ID")
		return
	}
	userID := r.Context().Value("userID").(uuid.UUID)
	password, err := h.passwordService.GetPassword(r.Context(), userID, passwordID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "password not found")
		return
	}

	ciphertextBase64 := base64.RawStdEncoding.EncodeToString(password.Ciphertext)
	nonceBase64 := base64.RawStdEncoding.EncodeToString(password.Nonce)

	resp := GetPasswordResponse{
		ID:         password.ID.String(),
		Name:       password.Name,
		Username:   password.Username,
		CreatedAt:  password.CreatedAt.String(),
		UpdatedAt:  password.UpdatedAt.String(),
		Ciphertext: ciphertextBase64,
		Nonce:      nonceBase64,
	}
	utils.JSON(w, http.StatusOK, resp)
}

// Handler function to update an existing password item
func (h *PasswordHandler) updatePassword(w http.ResponseWriter, r *http.Request) {
	passwordIDStr := chi.URLParam(r, "id")
	passwordID, err := uuid.Parse(passwordIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid password ID")
		return
	}

	var req CreatePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	userID := r.Context().Value("userID").(uuid.UUID)

	ciphertextBytes, err := base64.RawStdEncoding.DecodeString(req.Ciphertext)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ciphertext encoding")
		return
	}
	nonceBytes, err := base64.RawStdEncoding.DecodeString(req.Nonce)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid nonce encoding")
		return
	}
	password := &Password{
		ID:             passwordID,
		UserID:         userID,
		Name:           req.Name,
		Username:       req.Username,
		Ciphertext:     ciphertextBytes,
		Nonce:          nonceBytes,
		EncryptVersion: 1,
	}

	updated, err := h.passwordService.UpdatePassword(r.Context(), password)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := PasswordItem{
		ID:        updated.ID.String(),
		Name:      updated.Name,
		Username:  updated.Username,
		CreatedAt: updated.CreatedAt.String(),
		UpdatedAt: updated.UpdatedAt.String(),
	}
	utils.JSON(w, http.StatusOK, resp)
}

// Handler function to delete a password item
func (h *PasswordHandler) deletePassword(w http.ResponseWriter, r *http.Request) {
	passwordIDStr := chi.URLParam(r, "id")
	passwordID, err := uuid.Parse(passwordIDStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid password ID")
		return
	}
	userID := r.Context().Value("userID").(uuid.UUID)
	err = h.passwordService.DeletePassword(r.Context(), userID, passwordID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusNoContent, nil)
}
