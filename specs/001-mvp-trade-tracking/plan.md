# Plan: MVP trade tracking

Implements `spec.md` in this folder.

## Affected layers

| Layer | Change |
|---|---|
| domain | `Asset`, `Operation`, `Position` + validation + average-cost math |
| repository | `AssetRepository`, `OperationRepository` interfaces |
| repository/sqlite | GORM implementations, both wrapped in `database.WithRetry` |
| service | `AssetService`, `OperationService`, `PositionService` |
| httpserver | handlers + routes under `/api/v1`, `/health` |
| frontend | not yet implemented (tracked as task 6 below) |

## Data / migration impact

Two new tables via `AutoMigrate`: `assets`, `operations` (FK `asset_id`).
No existing data to backfill — this is the first migration.

## Retry/backoff considerations

SQLite with a single writer connection surfaces "database is locked" under
concurrent writes. `database.IsRetryable` matches that message (plus a few
other transient classes anticipated for a future networked DB) and
`database.WithRetry` retries with exponential backoff + jitter, capped by
`DB_RETRY_MAX_ATTEMPTS`/`DB_RETRY_MAX_DELAY`.

## Risks / tradeoffs

- Position is always recomputed from full operation history (O(n) per
  request). Fine at MVP scale; revisit with a materialized snapshot + event
  log only if a real user has thousands of operations and feels the latency.
- `SetMaxOpenConns(1)` on the SQLite connection trades write concurrency for
  correctness; acceptable for a single-user local tool.

## Rollout

No feature flag needed (first feature). CI runs `go test ./... -race -cover`
on every push (see `.github/workflows/backend-ci.yml`).
