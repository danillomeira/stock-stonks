---
name: reuse
description: Use before writing any new component, hook, helper, or utility (frontend or backend) — check docs/catalog.md and the shared locations (frontend/components/, frontend/lib/, backend/internal/*) for something that already does this, instead of re-deriving it from scratch. Use again immediately after creating something generic enough that another feature could plausibly need it — extract it to a shared location and add one line to docs/catalog.md in the same change, so a future session (with no memory of this one) can grep the catalog instead of re-reading the codebase.
---

# Reuse before building, catalog after building

## Why this exists

An AI session has no memory of the last one. Without a catalog, "is there
already something that does this?" costs a full codebase search every time
— tokens and turns spent re-discovering what a past session already built.
A short, accurate catalog turns that search into one grep of
`docs/catalog.md`.

## Before writing new code

1. Check `docs/catalog.md` for an entry matching what you're about to build.
2. If nothing is cataloged, check the shared locations directly before
   writing something new:
   - Frontend: `frontend/components/` (UI), `frontend/lib/` (framework-agnostic logic).
   - Backend: cross-cutting packages under `backend/internal/` (e.g.
     `internal/database`, `internal/httpserver/middleware`) — not buried
     inside the service/handler that happened to need it first.
3. Only write new code once you've confirmed nothing reusable already
   covers this. Prefer extending an existing reusable piece over writing a
   near-duplicate.

## When you create something reusable

Ask: would a different page/feature/service plausibly need this exact
logic? Examples: a generic table/form field, a validation helper, a
formatting/parsing function, a retry policy, an HTTP client wrapper, a
loading/error-state pattern. If yes:

1. **Put it in a shared location**, not inline in the feature that
   happened to need it first:
   - Frontend: generic UI → `frontend/components/`; framework-agnostic
     logic (API clients, formatters, validators) → `frontend/lib/`.
   - Backend: a dedicated package under `backend/internal/`, scoped by what
     it does (not by which feature first needed it).
2. **Give it a short doc comment** stating what it does and when to use
   it — specific enough that a future agent can decide from the comment
   alone whether it fits, without reading the implementation.
3. **Add one line to `docs/catalog.md`** in the same change: name,
   one-line purpose, file path. Keep entries short — the catalog is an
   index, not documentation (full docs, if needed, live in the file itself
   or in `specs/`).

## What does NOT belong in the catalog

- One-off, page/feature-specific logic with no plausible second caller.
- Anything already obvious from the layering in `specs/constitution.md`
  (e.g. don't catalog "AssetService" — that's the architecture, not a
  reuse decision).
- Anything that duplicates what an existing catalog entry already covers —
  extend that entry's file instead of adding a near-identical one.

## Workflow

1. Before implementing: check `docs/catalog.md`, then the shared
   locations above.
2. Implement.
3. After implementing: if what you built is reusable, extract it to the
   right shared location, document it, and add the catalog line — in the
   same change, not as a later cleanup pass.
