# Spec: Google login and multi-user data isolation

- **Status**: In Progress
- **Author**: Danillo Meira (with Claude Code)
- **Date**: 2026-10-27

## Problem

The MVP (spec 001) is single-user with no auth — every Asset/Operation is
global. To build anything closer to a real product (dashboard, IR export,
multiple people using their own installs), each user needs their own
isolated data, and needs to log in. Google was chosen as the login
provider (no password storage, no separate signup flow to build).

## Goals

- "Entrar com Google" login; a Google account maps 1:1 to a local `User`.
- Every `Asset` and `Operation` belongs to exactly one `User`; one user can
  never read or write another's data via the API.
- A logged-in session persists across page reloads (cookie-based), and a
  user can log out.
- Unauthenticated requests to any data endpoint are rejected (401).

## Non-goals (this spec)

- Any login method other than Google.
- Roles/permissions beyond "owns their own data" (no admin, no sharing).
- Session revocation/management UI (a logout simply drops the cookie; see
  the ADR for why this is an accepted MVP limitation).
- Real-time quotes, dashboard widgets, IR export, rating/classification —
  separate, later specs. This spec is the foundation they'll build on.

## User-facing behavior

- Frontend shows a "Entrar com Google" button when logged out.
- Clicking it does a full-page redirect through Google's consent screen
  and back; on success the user lands back on the app, logged in.
- `GET /api/v1/auth/me` returns the current user (`401` if not logged in) —
  the frontend uses this to decide whether to show the login button or the
  app.
- `POST /api/v1/auth/logout` clears the session.
- Every existing data endpoint (`/assets`, `/operations`, `/positions`) now
  requires a session and is scoped to the caller's own data.

## Domain model impact

- New entity `domain.User`: `ID`, `GoogleID` (unique, the OAuth `sub`
  claim), `Email`, `Name`, `AvatarURL`, `CreatedAt`, `UpdatedAt`.
- `domain.Asset` and `domain.Operation` gain a required `UserID` field.
- Every `AssetRepository`/`OperationRepository` method gains a `userID`
  parameter and filters/scopes by it — see `plan.md`.

## Open questions (resolved — see ADR 0009)

- Session mechanism → signed JWT in an HttpOnly cookie, stateless (no
  sessions table for the MVP).
- Google ID token verification → `google.golang.org/api/idtoken` (official
  Google library).
