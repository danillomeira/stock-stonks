// Package auth holds framework-agnostic OAuth2/OIDC and session-token logic:
// no HTTP handlers, no domain/repository knowledge. The HTTP glue (cookies,
// redirects) lives in internal/httpserver/handlers; the domain glue (upsert
// a User) lives in internal/service.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Identity is what Google's ID token told us about the person who logged in.
type Identity struct {
	Sub     string // stable Google user id
	Email   string
	Name    string
	Picture string
}

// GoogleOAuth wraps the OAuth2 Authorization Code flow against Google plus
// ID token verification. ClientID/ClientSecret empty means Google login is
// not configured - callers should surface a clear error rather than let a
// nil-ish oauth2.Config fail confusingly deep in the token exchange.
type GoogleOAuth struct {
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
}

// ErrNotConfigured is returned by every GoogleOAuth method when no client
// ID/secret was provided - the server can still run (other endpoints work),
// it just can't log anyone in yet.
var ErrNotConfigured = fmt.Errorf("auth: Google OAuth is not configured (set GOOGLE_CLIENT_ID/GOOGLE_CLIENT_SECRET)")

// NewGoogleOAuth builds a GoogleOAuth client. If clientID/clientSecret are
// empty, the returned client always reports ErrNotConfigured - this lets
// the rest of the server start normally in dev without Google credentials.
func NewGoogleOAuth(ctx context.Context, clientID, clientSecret, redirectURL string) (*GoogleOAuth, error) {
	if clientID == "" || clientSecret == "" {
		return &GoogleOAuth{}, nil
	}

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("auth: discover google oidc provider: %w", err)
	}

	return &GoogleOAuth{
		config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (g *GoogleOAuth) configured() bool {
	return g.config.ClientID != ""
}

// AuthCodeURL builds the URL to redirect the browser to for Google's consent screen.
func (g *GoogleOAuth) AuthCodeURL(state string) (string, error) {
	if !g.configured() {
		return "", ErrNotConfigured
	}
	return g.config.AuthCodeURL(state), nil
}

// Exchange trades an authorization code for a verified Identity.
func (g *GoogleOAuth) Exchange(ctx context.Context, code string) (*Identity, error) {
	if !g.configured() {
		return nil, ErrNotConfigured
	}

	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("auth: exchange code: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return nil, fmt.Errorf("auth: token response had no id_token")
	}

	idToken, err := g.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("auth: verify id_token: %w", err)
	}

	var claims struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("auth: parse id_token claims: %w", err)
	}

	return &Identity{
		Sub:     idToken.Subject,
		Email:   claims.Email,
		Name:    claims.Name,
		Picture: claims.Picture,
	}, nil
}

// GenerateState returns a random, URL-safe token for CSRF protection on the
// OAuth redirect (compared against a short-lived cookie set alongside it).
func GenerateState() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("auth: generate state: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
