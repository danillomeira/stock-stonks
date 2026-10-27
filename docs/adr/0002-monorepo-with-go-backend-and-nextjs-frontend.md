# 0002. Monorepo with a Go backend and a Next.js frontend

## Status
Accepted

## Context
The project needs a backend (Go, for the trade/tax domain logic and a
SQLite-backed API) and a frontend (Next.js, for the UI). These could live in
one repository or two.

## Decision
Use a single monorepo (`backend/`, `frontend/`) with one GitHub Actions
pipeline that runs separate jobs per subproject. Each subproject keeps its
own module/package manifest (`backend/go.mod`, `frontend/package.json`).

## Consequences
- One PR can span backend + frontend when a feature touches both (e.g. a
  new endpoint and the page that calls it) without cross-repo coordination.
- CI must be structured with path filters so a frontend-only change doesn't
  re-run Go tests and vice versa (see `.github/workflows/`).
- If the two ever need independent release cadences or separate access
  control, splitting later means extracting with `git subtree`/`git
  filter-repo`, not a trivial move — acceptable risk at this project's
  current size (single maintainer, early MVP).

## Alternatives considered
- **Separate `stock-stonks-backend` / `stock-stonks-frontend` repos** —
  rejected for now: adds cross-repo versioning overhead with no current
  need for independent deploy cadences.
