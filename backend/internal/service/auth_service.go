package service

import (
	"context"
	"fmt"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/auth"
	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
)

// AuthService glues the framework-agnostic internal/auth package (OAuth,
// sessions) to the domain/repository layer (upserting a User).
type AuthService struct {
	google   *auth.GoogleOAuth
	sessions *auth.SessionManager
	users    repository.UserRepository
}

func NewAuthService(google *auth.GoogleOAuth, sessions *auth.SessionManager, users repository.UserRepository) *AuthService {
	return &AuthService{google: google, sessions: sessions, users: users}
}

// LoginURL returns Google's consent screen URL for the given CSRF state.
func (s *AuthService) LoginURL(state string) (string, error) {
	return s.google.AuthCodeURL(state)
}

// HandleCallback exchanges the OAuth code for a verified identity, upserts
// the local User, and returns a signed session token ready to set as a cookie.
func (s *AuthService) HandleCallback(ctx context.Context, code string) (token string, expiresAt time.Time, err error) {
	identity, err := s.google.Exchange(ctx, code)
	if err != nil {
		return "", time.Time{}, err
	}

	user := &domain.User{
		GoogleID:  identity.Sub,
		Email:     identity.Email,
		Name:      identity.Name,
		AvatarURL: identity.Picture,
	}
	if err := user.Validate(); err != nil {
		return "", time.Time{}, fmt.Errorf("auth: identity from google failed validation: %w", err)
	}
	if err := s.users.UpsertByGoogleID(ctx, user); err != nil {
		return "", time.Time{}, fmt.Errorf("auth: upsert user: %w", err)
	}

	return s.sessions.IssueToken(user.ID)
}

// CurrentUser resolves the authenticated user's full profile.
func (s *AuthService) CurrentUser(ctx context.Context, userID uint) (*domain.User, error) {
	return s.users.GetByID(ctx, userID)
}
