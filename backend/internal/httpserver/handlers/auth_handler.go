package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/auth"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
	"github.com/danillomeira/stock-stonks/backend/internal/service"
)

const oauthStateCookieName = "oauth_state"

type AuthHandler struct {
	auth         *service.AuthService
	frontendURL  string
	cookieSecure bool
}

func NewAuthHandler(authSvc *service.AuthService, frontendURL string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{auth: authSvc, frontendURL: frontendURL, cookieSecure: cookieSecure}
}

// Login starts the Google OAuth flow: stash a CSRF state value in a
// short-lived cookie, then redirect the browser to Google's consent screen.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := auth.GenerateState()
	if err != nil {
		writeError(w, err)
		return
	}

	url, err := h.auth.LoginURL(state)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   int((5 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, url, http.StatusFound)
}

// Callback completes the flow: validate the CSRF state, exchange the code,
// set the session cookie, and send the browser back to the frontend.
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie(oauthStateCookieName)
	if err != nil || r.URL.Query().Get("state") != stateCookie.Value {
		http.Redirect(w, r, h.frontendURL+"?auth_error=invalid_state", http.StatusFound)
		return
	}
	clearCookie(w, oauthStateCookieName, h.cookieSecure)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.frontendURL+"?auth_error=missing_code", http.StatusFound)
		return
	}

	token, expiresAt, err := h.auth.HandleCallback(r.Context(), code)
	if err != nil {
		slog.Error("auth: callback failed", "error", err)
		http.Redirect(w, r, h.frontendURL+"?auth_error=login_failed", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, h.frontendURL, http.StatusFound)
}

// Logout clears the session cookie. No server-side revocation - see ADR 0009.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, auth.SessionCookieName, h.cookieSecure)
	writeJSON(w, http.StatusNoContent, nil)
}

// Me returns the authenticated user's profile.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	user, err := h.auth.CurrentUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusUnauthorized, errorBody{Error: "not authenticated"})
			return
		}
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func clearCookie(w http.ResponseWriter, name string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
