package shortener

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	if len(req.URL) > 32779 {
		writeError(w, http.StatusBadRequest, "url too long (max 32779 characters)")
		return
	}

	addr := &Address{URL: req.URL}

	existing, err := h.repo.GetByURL(r.Context(), req.URL)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateAddressResponse{
			ID:  existing.ID,
			URL: existing.URL,
		})
		return
	}

	if err := h.repo.Create(r.Context(), addr); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateAddressResponse{
		ID:  addr.ID,
		URL: addr.URL,
	})
}

func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/shorter/url/")
	if id == "" || id == r.URL.Path {
		writeError(w, http.StatusBadRequest, "Missing address ID")
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		writeError(w, http.StatusNotAcceptable, "No item found matching the url")
		return
	}

	addr, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotAcceptable, "No item found matching the url")
		return
	}

	http.Redirect(w, r, addr.URL, http.StatusFound)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
