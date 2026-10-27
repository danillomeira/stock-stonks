package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/danillomeira/stock-stonks/backend/internal/auth"
)

// RequireAuth rejects any request without a valid session cookie (401) and,
// for everything else, injects the authenticated user's ID into the request
// context via auth.WithUserID - handlers read it with auth.UserIDFromContext
// instead of trusting any user id the client might put in the request.
func RequireAuth(sessions *auth.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(auth.SessionCookieName)
			if err != nil {
				unauthorized(w)
				return
			}

			userID, err := sessions.VerifyToken(cookie.Value)
			if err != nil {
				unauthorized(w)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.WithUserID(r.Context(), userID)))
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "not authenticated"})
}
