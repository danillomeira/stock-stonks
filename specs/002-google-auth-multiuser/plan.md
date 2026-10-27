# Plan: Google login and multi-user data isolation

Implements `spec.md` in this folder.

## Affected layers

| Layer | Change |
|---|---|
| domain | new `User`; `UserID` added to `Asset`, `Operation` |
| repository | new `UserRepository`; `userID` param added to every `AssetRepository`/`OperationRepository` method |
| repository/sqlite | `UserRepository` impl; existing repos filter/scope every query by `user_id` |
| service | new `AuthService` (OAuth exchange + session issue/verify); `userID` threaded through `AssetService`/`OperationService`/`PositionService` |
| httpserver | new `internal/auth` package (OAuth config, session JWT); `AuthMiddleware`; new `/api/v1/auth/*` routes; every existing handler reads `userID` from request context instead of taking it as a free parameter |
| frontend | login button, `/auth/me` check gating the app shell, logout, `credentials: "include"` on all API calls |

## Data / migration impact

- New table `users` (`AutoMigrate`).
- `assets.user_id`, `operations.user_id` columns added, both `NOT NULL` +
  indexed. No existing rows to backfill in any real deployment yet (MVP has
  no production data), so no migration script beyond `AutoMigrate`.

## Retry/backoff considerations

No new transient-error class. User lookups/upserts go through the same
`database.WithRetry` as everything else.

## Session design

- `internal/auth.Config`: Google OAuth client ID/secret/redirect URL (from
  `config.Config`, empty by default — `/auth/google/login` returns a clear
  error if unset rather than crashing the whole server), session signing
  secret, cookie settings.
- Login: `GET /api/v1/auth/google/login` sets a short-lived `state` cookie
  (CSRF protection) and redirects to Google's consent screen.
- Callback: `GET /api/v1/auth/google/callback` validates `state`, exchanges
  the code, verifies the ID token via `google.golang.org/api/idtoken`,
  upserts the `User` by `GoogleID`, issues a signed session JWT in an
  HttpOnly cookie, redirects to the frontend.
- `AuthMiddleware` reads the session cookie, verifies the JWT, and injects
  `userID` into the request context; returns 401 if missing/invalid.
- Logout: clears the cookie. No server-side revocation list (see ADR 0009).

## Risks / tradeoffs

- Stateless session JWT means logout doesn't invalidate the token
  server-side before it expires — mitigated with a short expiry
  (documented in the ADR). Acceptable for a personal-finance MVP with no
  sensitive-action-without-reauth requirement; revisit if this becomes a
  concern (would need a sessions table + revocation check).
- `SameSite=Lax` works for frontend/backend on the same registrable domain
  (including different localhost ports in dev). A prod deployment on two
  *different* domains would need `SameSite=None; Secure` — noted in the ADR
  so it isn't forgotten when a real hosting target is chosen (ADR 0007).
- CORS middleware must now set `Access-Control-Allow-Credentials: true` and
  can no longer use a wildcard origin (already uses an explicit origin from
  config, so this is additive, not a breaking change).

## Rollout

Breaking change to the API (every data endpoint now requires auth) — fine
pre-release with no real users yet. CI (`backend-ci`) covers the new
domain/service unit tests; OAuth itself isn't exercised in CI (needs real
Google credentials), so it's smoke-tested manually against a real Google
OAuth client.
