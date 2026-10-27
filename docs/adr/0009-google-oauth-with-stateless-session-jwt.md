# 0009. Google OAuth login with a stateless session JWT

## Status
Accepted

## Context
The app is moving from single-user to multi-user (spec
`specs/002-google-auth-multiuser/`). Google was chosen as the only login
provider to avoid building/storing passwords. Two things need deciding:
how to verify a Google identity server-side, and how to keep a user logged
in across requests without re-authenticating every time.

## Decision
- **OAuth flow**: standard server-side Authorization Code flow
  (`golang.org/x/oauth2` + `golang.org/x/oauth2/google`). The backend, not
  the frontend, holds the Google client secret and does the token exchange
  — the frontend only ever sees a link to `/api/v1/auth/google/login`.
- **ID token verification**: `google.golang.org/api/idtoken` (official
  Google library) validates the signature, issuer, audience, and expiry of
  the ID token returned by Google — no hand-rolled JWT/JWKS verification.
- **Our own session**: a JWT (`github.com/golang-jwt/jwt/v5`, HS256, signed
  with a server-held secret) containing `{sub: userID, exp}`, set as an
  HttpOnly, SameSite=Lax cookie. Stateless — there is no `sessions` table.
  `AuthMiddleware` verifies the signature and expiry on every request;
  logout just clears the cookie client-side.

## Consequences
- No password storage, no password-reset flow to build.
- No sessions table to manage/clean up; the whole mechanism is "verify a
  signature," which is simple and fast.
- **Logout does not revoke the token early** — a copied/stolen cookie
  remains valid until it expires, even after "logout." Mitigated with a
  short expiry (days, not months). This is an accepted MVP limitation for a
  personal-finance tool with no untrusted multi-device/shared-computer
  threat model yet; if that changes, the fix is a sessions table with a
  revocation check added to `AuthMiddleware`, which is the obvious upgrade
  path from here, not a redesign.
- `SameSite=Lax` requires frontend and backend to share a registrable
  domain (true for `localhost:3000`/`localhost:8080` in dev, and for any
  prod setup on one domain with subdomains). A future deployment that puts
  the frontend and backend on genuinely different domains needs
  `SameSite=None; Secure` instead — flagged here so it isn't rediscovered
  the hard way when ADR 0007's hosting decision is finally made.
- The backend now requires `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET` to be
  configured for login to work; everything else (including running the
  server at all) still works with zero config, matching the existing
  zero-setup default for local dev.

## Alternatives considered
- **Google Identity Services (client-side) handing the frontend an ID
  token directly** — rejected: would mean the frontend talks to Google
  directly and hands a raw ID token to the backend on every request (or a
  separate exchange step anyway); the server-side Authorization Code flow
  keeps all OAuth secrets and logic in one place (the backend) and the
  frontend stays a thin client.
- **Database-backed sessions (session ID in cookie, row in `sessions`
  table)** — rejected for the MVP: real revocation support, but an extra
  table, an extra query per request, and a cleanup job for expired rows.
  Worth adopting later if the stateless JWT's "can't revoke early"
  limitation becomes a real problem.
- **`golang.org/x/oauth2` only, hand-parsing the ID token's JWT ourselves**
  — rejected: `google.golang.org/api/idtoken` already does signature/aud/iss
  verification correctly against Google's rotating keys; re-implementing
  that is security-sensitive code better left to a library Google
  maintains.
