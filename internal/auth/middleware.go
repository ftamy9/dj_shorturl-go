package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ftamy9/dj_shorturl-go/internal/config"
)

func Middleware(repo Repository, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Basic") {
				writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]
			success, userID, err := VerifyToken(token, cfg.AuthSecret, cfg.AuthTTLMinutes)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Token format invalid")
				return
			}
			if !success {
				writeError(w, http.StatusUnauthorized, "Token error "+userID)
				return
			}

			user, err := repo.GetByID(r.Context(), userID)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "No such user")
				return
			}

			r.Header.Set("X-User-ID", user.ID)
			next.ServeHTTP(w, r)
		})
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
