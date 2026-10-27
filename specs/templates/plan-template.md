# Plan: <feature name>

Implements `spec.md` in this same folder. Written after the spec is approved,
before `tasks.md`.

## Affected layers

| Layer | Change |
|---|---|
| domain | |
| repository | |
| repository/sqlite | |
| service | |
| httpserver | |
| frontend | |

## Data / migration impact

Does this add/change GORM-managed tables? Any backfill needed for existing rows?

## Retry/backoff considerations

Does this introduce a new transient-error class `database.IsRetryable` must
learn (e.g. a new driver, a new external call)? If so, say what it is and how
it will be detected.

## Risks / tradeoffs

What could go wrong, and what did we decide not to do (and why)? If a choice
here is non-obvious or hard to reverse, write an ADR instead of burying the
reasoning in this plan.

## Rollout

Feature flag? Behind a migration? Needs a CI/CD change?
