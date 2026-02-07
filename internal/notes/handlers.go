package notes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/subrat-dwi/shubserver/internal/utils"
)

// Handler struct for notes
type Handler struct {
	repo Repository
}

// Constructor for handler
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// Handler to show all notes
func (h *Handler) listNotes(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List()
	if err != nil {
		utils.Error(w, http.StatusNotFound, "can't access notes")
		return
	}

	json.NewEncoder(w).Encode(list)
}

// Handler to get a single note
func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	note, err := h.repo.Get(id)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, "note not found")
		return
	}

	json.NewEncoder(w).Encode(note)
}

// Handler to create a note
func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var note Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalide JSON")
		return
	}

	note.ID = uuid.NewString()
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	if err := h.repo.Create(&note); err != nil {
		utils.Error(w, http.StatusInternalServerError, "cannot save the note")
		return
	}

	json.NewEncoder(w).Encode(note)
}

// Handler to delete a note
func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.repo.Delete(id); err != nil {
		utils.Error(w, http.StatusBadRequest, "cannot delete note")
		return
	}
	w.Write([]byte("Deletion Successful"))
}

// Handler to update a note
func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.Error(w, http.StatusBadRequest, "missing id")
		return
	}
	existing, err := h.repo.Get(id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, "note not found")
		return
	}

	var payload Note
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	existing.Title = payload.Title
	existing.Content = payload.Content
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(existing); err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to update")
		return
	}

	json.NewEncoder(w).Encode(existing)

}
