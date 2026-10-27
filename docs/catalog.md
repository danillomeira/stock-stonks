# Reuse catalog

Index of generic, reusable components/helpers in this repo. Maintained by
the `reuse` skill (`.claude/skills/reuse/SKILL.md`): check here before
writing new code; add a line here when you create something reusable.

## Backend (`backend/internal/`)

| Name | Purpose | File |
|---|---|---|
| `database.WithRetry` | Exponential backoff + jitter wrapper for any DB call; retries only transient errors | `internal/database/retry.go` |
| `database.IsRetryable` | Classifies an error as transient (lock/busy/connection) vs. permanent | `internal/database/retry.go` |
| `database.Open` | Driver factory seam — add a case here to support a new DB without touching repository/service/handler code | `internal/database/connection.go` |
| `middleware.Logging` | Generic `net/http` request logger (method/path/status/duration) | `internal/httpserver/middleware/logging.go` |
| `middleware.CORS` | Generic `net/http` CORS middleware, origin read from config | `internal/httpserver/middleware/cors.go` |
| `handlers.writeJSON` / `writeError` / `parseIDParam` | Shared HTTP response/error helpers for any new handler in that package | `internal/httpserver/handlers/respond.go` |

## Frontend (`frontend/`)

| Name | Purpose | File |
|---|---|---|
| `api` client (`request`, `ApiError`) | Typed fetch wrapper for the backend API — base URL, JSON headers, error mapping | `frontend/lib/api.ts` |

## Known gaps (not yet extracted)

The `assets`, `operations`, and `positions` pages currently duplicate their
own loading/error-message and table markup inline instead of sharing a
component. Flagged here rather than silently left as debt — extract into
`frontend/components/` (e.g. a data table, a labeled-field wrapper) the
next time a fourth page needs the same pattern, per this skill's rule.
