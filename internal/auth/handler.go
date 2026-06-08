package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/ftamy9/dj_shorturl/internal/config"
)

type Handler struct {
	repo Repository
	cfg  *config.Config
}

func NewHandler(repo Repository, cfg *config.Config) *Handler {
	return &Handler{repo: repo, cfg: cfg}
}

func hashPassword(password, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodPut:
		h.login(w, r)
	default:
		w.Header().Set("Allow", "POST, PUT")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ID == "" || req.PasswordHash == "" {
		writeError(w, http.StatusBadRequest, "id and password_hash are required")
		return
	}

	if len(req.ID) > 10 {
		writeError(w, http.StatusBadRequest, "id must be at most 10 characters")
		return
	}

	user := &BaseUser{
		ID:           req.ID,
		PasswordHash: hashPassword(req.PasswordHash, h.cfg.PasswordSecret),
	}

	if err := h.repo.Create(r.Context(), user); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{
		ID:           user.ID,
		PasswordHash: "ok***hash",
	})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hashedPassword := hashPassword(req.PasswordHash, h.cfg.PasswordSecret)

	user, err := h.repo.GetByID(r.Context(), req.ID)
	if err != nil || user.PasswordHash != hashedPassword {
		writeError(w, http.StatusNotAcceptable, "Invalid credentials")
		return
	}

	token, err := CreateToken(user.ID, h.cfg.AuthSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
