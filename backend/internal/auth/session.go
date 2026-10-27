package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SessionCookieName is the HttpOnly cookie carrying the session JWT. Shared
// between the handler that sets it (login callback / clears it on logout)
// and the middleware that reads it, so both sides can never drift apart.
const SessionCookieName = "session"

// SessionManager issues and verifies the app's own session token (a JWT
// signed with a server-held secret, carried in an HttpOnly cookie) - kept
// separate from Google's ID token, which is only used once at login time.
// Stateless by design: no sessions table, see ADR 0009 for the tradeoff.
type SessionManager struct {
	secret []byte
	ttl    time.Duration
}

func NewSessionManager(secret string, ttl time.Duration) *SessionManager {
	return &SessionManager{secret: []byte(secret), ttl: ttl}
}

// ErrInvalidSession covers every way a session token can fail to verify:
// bad signature, expired, malformed. Callers don't need to distinguish why.
var ErrInvalidSession = fmt.Errorf("auth: invalid or expired session")

// IssueToken returns a signed session token for userID and its expiry time.
func (s *SessionManager) IssueToken(userID uint) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(s.ttl)
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(userID), 10),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("auth: sign session token: %w", err)
	}
	return signed, expiresAt, nil
}

// VerifyToken validates the signature and expiry and returns the userID it carries.
func (s *SessionManager) VerifyToken(tokenString string) (userID uint, err error) {
	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return 0, ErrInvalidSession
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, ErrInvalidSession
	}
	return uint(id), nil
}
