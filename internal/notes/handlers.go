package notes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

// NotesHandler to show all notes
func (h *NotesHandler) listNotes(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		utils.Error(w, http.StatusNotFound, "can't access notes")
		return
	}

	json.NewEncoder(w).Encode(list)
}

// NotesHandler to get a single note
func (h *NotesHandler) getNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	note, err := h.repo.Get(r.Context(), id)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, "note not found")
		return
	}

	json.NewEncoder(w).Encode(note)
}

// NotesHandler to create a note
func (h *NotesHandler) createNote(w http.ResponseWriter, r *http.Request) {
	var note Note

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	userID := r.Context().Value("userID").(string)
	note.UserID = userID

	dbnote, err := h.repo.Create(r.Context(), &note)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(dbnote)
}

// NotesHandler to delete a note
func (h *NotesHandler) deleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.repo.Delete(r.Context(), id); err != nil {
		utils.Error(w, http.StatusBadRequest, "cannot delete note")
		return
	}
	w.Write([]byte("Deletion Successful"))
}

// NotesHandler to update a note
func (h *NotesHandler) updateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.Error(w, http.StatusBadRequest, "missing id")
		return
	}
	existing, err := h.repo.Get(r.Context(), id)
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

	if err := h.repo.Update(r.Context(), existing); err != nil {
		utils.Error(w, http.StatusInternalServerError, "failed to update")
		return
	}

	json.NewEncoder(w).Encode(existing)
}
