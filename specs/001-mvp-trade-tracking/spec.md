# Spec: MVP trade tracking (assets, operations, positions)

- **Status**: Done
- **Author**: Danillo Meira (with Claude Code)
- **Date**: 2026-10-26

## Problem

Before any tax-export feature can exist, we need somewhere to register
assets (tickers) and the buy/sell operations on them, and to know the
resulting position (quantity held, average cost, realized profit) at any
point in time.

## Goals

- CRUD for assets (ticker, name, type).
- CRUD for operations (buy/sell, date, quantity, unit price, fees, broker).
- Compute current position per asset (and for all assets) from the full
  operation history, using the weighted average cost method required by
  Brazilian IR rules.
- Database access resilient to SQLite's single-writer lock contention.

## Non-goals

- Tax export / DARF calculation (future spec).
- Day-trade vs. swing-trade classification (future spec — changes applicable
  IR rate, not handled yet).
- Dividends/proventos tracking (future spec).
- Authentication/multi-user (single-user local tool for now — superseded by
  `specs/002-google-auth-multiuser/`).
- Frontend UI beyond a basic list/registration page.

## User-facing behavior

REST API under `/api/v1`:

- `POST/GET /assets`, `GET/PUT/DELETE /assets/{id}` — asset CRUD.
- `POST/GET /operations`, `GET/PUT/DELETE /operations/{id}` — operation CRUD.
  `GET /operations` accepts `asset_id`, `from`, `to` (YYYY-MM-DD) filters.
- `GET /positions` — current position for every asset with operations.
- `GET /positions/{assetID}` — current position for one asset.
- `GET /health` — liveness check.

## Domain model impact

New entities: `domain.Asset`, `domain.Operation`, `domain.Position` (the last
one computed, never persisted — see ADR [[0005-position-is-computed-not-stored]]).
All monetary/quantity fields on `Operation` are `decimal.Decimal`.

## Open questions (resolved)

- ORM/DB choice → GORM + SQLite, see
  [[0003-go-gorm-sqlite-with-retry-backoff-seam]].
- Retry policy → exponential backoff with jitter, see same ADR.
