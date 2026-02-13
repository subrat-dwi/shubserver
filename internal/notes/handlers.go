package notes

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/subrat-dwi/shubserver/internal/utils"
)

// Handler struct for notes
type NotesHandler struct {
	repo NotesRepository
}

// Constructor for handler
func NewNotesHandler(repo NotesRepository) *NotesHandler {
	return &NotesHandler{repo: repo}
}

// Validation helper
func (h *NotesHandler) validateNoteInput(title, content string) error {
	if strings.TrimSpace(title) == "" {
		return utils.NewValidationError("title is required")
	}
	if strings.TrimSpace(content) == "" {
		return utils.NewValidationError("content is required")
	}
	if len(title) > 255 {
		return utils.NewValidationError("title must be less than 255 characters")
	}
	if len(content) > 5000 {
		return utils.NewValidationError("content must be less than 5000 characters")
	}
	return nil
}

// NotesHandler to show all notes
func (h *NotesHandler) listNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	list, err := h.repo.List(r.Context(), userID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "can't access notes")
		return
	}

	utils.JSON(w, http.StatusOK, list)
}

// NotesHandler to get a single note
func (h *NotesHandler) getNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(uuid.UUID)
	note, err := h.repo.Get(r.Context(), userID, id)

	if err != nil {
		utils.Error(w, http.StatusNotFound, "note not found")
		return
	}

	utils.JSON(w, http.StatusOK, note)
}

// NotesHandler to create a note
func (h *NotesHandler) createNote(w http.ResponseWriter, r *http.Request) {
	var note Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Validate input
	if err := h.validateNoteInput(note.Title, note.Content); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)
	note.UserID = userID

	dbnote, err := h.repo.Create(r.Context(), &note)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, dbnote)
}

// NotesHandler to delete a note
func (h *NotesHandler) deleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(uuid.UUID)

	if err := h.repo.Delete(r.Context(), userID, id); err != nil {
		utils.Error(w, http.StatusNotFound, "cannot delete note")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{
		"message": "note deleted successfully",
	})
}

// NotesHandler to update a note
func (h *NotesHandler) updateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID := r.Context().Value("userID").(uuid.UUID)

	if id == "" {
		utils.Error(w, http.StatusBadRequest, "missing id")
		return
	}

	existing, err := h.repo.Get(r.Context(), userID, id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "note not found")
		return
	}

	var payload Note
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Validate input
	if err := h.validateNoteInput(payload.Title, payload.Content); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	existing.Title = payload.Title
	existing.Content = payload.Content
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(r.Context(), userID, existing); err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to update")
		return
	}

	utils.JSON(w, http.StatusOK, existing)
}
