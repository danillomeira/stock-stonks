package auth_test

import (
	"testing"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/auth"
)

func TestSessionManager_IssueAndVerify_RoundTrip(t *testing.T) {
	sm := auth.NewSessionManager("test-secret", time.Hour)

	token, expiresAt, err := sm.IssueToken(42)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}
	if time.Until(expiresAt) <= 0 || time.Until(expiresAt) > time.Hour {
		t.Fatalf("expiresAt = %v, want ~1h from now", expiresAt)
	}

	userID, err := sm.VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}
	if userID != 42 {
		t.Fatalf("userID = %d, want 42", userID)
	}
}

func TestSessionManager_VerifyToken_RejectsExpired(t *testing.T) {
	sm := auth.NewSessionManager("test-secret", -time.Hour) // already-expired token

	token, _, err := sm.IssueToken(1)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}

	if _, err := sm.VerifyToken(token); err != auth.ErrInvalidSession {
		t.Fatalf("VerifyToken() error = %v, want %v", err, auth.ErrInvalidSession)
	}
}

func TestSessionManager_VerifyToken_RejectsWrongSecret(t *testing.T) {
	issuer := auth.NewSessionManager("secret-a", time.Hour)
	verifier := auth.NewSessionManager("secret-b", time.Hour)

	token, _, err := issuer.IssueToken(1)
	if err != nil {
		t.Fatalf("IssueToken() error = %v", err)
	}

	if _, err := verifier.VerifyToken(token); err != auth.ErrInvalidSession {
		t.Fatalf("VerifyToken() error = %v, want %v", err, auth.ErrInvalidSession)
	}
}

func TestSessionManager_VerifyToken_RejectsGarbage(t *testing.T) {
	sm := auth.NewSessionManager("test-secret", time.Hour)

	if _, err := sm.VerifyToken("not-a-jwt"); err != auth.ErrInvalidSession {
		t.Fatalf("VerifyToken() error = %v, want %v", err, auth.ErrInvalidSession)
	}
}
