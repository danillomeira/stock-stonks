package middleware

import "net/http"

// CORS allows the Next.js frontend (running on a different origin in dev) to
// call the API. allowedOrigin is read from config so it can be locked down
// per environment instead of using a wildcard in production. Credentials
// (the session cookie) require an exact origin - the wildcard "*" is not
// valid together with Allow-Credentials per the Fetch spec.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
