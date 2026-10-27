# 0003. GORM + SQLite for the MVP, behind a retry/backoff seam for future databases

## Status
Accepted

## Context
The MVP needs an ORM and a database. It must run with zero external
infrastructure (single-user local tool), but the domain (trade tracking,
tax export) is expected to grow — multi-user, larger datasets, and a
networked database (Postgres likely) are plausible future requirements.
SQLite only allows one writer at a time, so concurrent writes can fail with
"database is locked" even under light concurrency (e.g. two quick successive
requests).

## Decision
- Use **GORM** as the ORM — it supports SQLite, Postgres, and MySQL behind
  the same API, which keeps a future driver swap mechanical.
- Use **SQLite** (via `gorm.io/driver/sqlite`, CGO/`mattn/go-sqlite3`) as the
  only implemented driver for the MVP.
- Every repository method wraps its GORM call in `database.WithRetry`
  (`backend/internal/database/retry.go`): exponential backoff with jitter,
  retrying only errors classified as transient by `database.IsRetryable`
  (lock/busy/connection-reset style errors). Non-transient errors
  (validation, not-found, unique-constraint) fail immediately — they are
  never retried.
- `database.Open` (`backend/internal/database/connection.go`) is the single
  seam for adding a driver: a `postgres`/`mysql` case already exists and
  currently returns a clear "not implemented yet" error. Adding real support
  later means implementing that case and extending `IsRetryable` with that
  driver's transient-error strings — the repository, service, and handler
  layers never change.
- The SQLite connection pool is capped at 1 open connection
  (`SetMaxOpenConns(1)`), making the single-writer constraint explicit
  instead of letting GORM's default pool create contention the retry layer
  then has to paper over.

## Consequences
- The retry/backoff logic is shared infrastructure, not duplicated per
  repository method.
- SQLite's locking behavior is treated as a known, handled condition rather
  than an occasional unexplained 500.
- Moving to Postgres/MySQL later is additive (new driver case + new
  transient-error strings), not a rewrite of the service/handler layers.
- GORM's `AutoMigrate` is used for schema management at this stage; a real
  migration tool (e.g. `golang-migrate`) is deferred until schema changes
  need to be versioned/reviewed independently of the Go struct definitions
  (see `backend/internal/repository/sqlite/migrate.go`).

## Alternatives considered
- **`database/sql` + raw SQL / `sqlc`** — more explicit, but more
  boilerplate per entity and no built-in multi-driver abstraction; revisit
  if GORM's query builder becomes a bottleneck.
- **A separate retry library (e.g. `cenkalti/backoff`)** — rejected for now:
  the retry policy needed is small (classify error → exponential backoff
  with jitter) and a hand-rolled ~70-line implementation avoids an external
  dependency for logic this simple.
- **`modernc.org/sqlite` (pure Go, no CGO)** — would simplify
  cross-compilation and avoid needing a C toolchain in CI/Docker, but
  `mattn/go-sqlite3` is the far more battle-tested SQLite driver and GORM's
  officially documented one; CGO is not a real obstacle since both the CI
  runners and the Docker build stage already have a C toolchain. Worth
  revisiting only if cross-compilation needs change.
