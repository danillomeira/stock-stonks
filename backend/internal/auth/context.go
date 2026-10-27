package auth

import "context"

type contextKey int

const userIDKey contextKey = iota

// WithUserID returns a context carrying the authenticated user's ID.
func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext returns the authenticated user's ID. ok is false if the
// request never went through the auth middleware - callers must check it
// rather than treat a missing ID as user 0.
func UserIDFromContext(ctx context.Context) (userID uint, ok bool) {
	userID, ok = ctx.Value(userIDKey).(uint)
	return userID, ok
}
