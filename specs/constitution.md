# Stock-Stonks Constitution

Non-negotiable principles for any spec, plan, or implementation in this repository.
Changes to this document are themselves architectural decisions and require an ADR
(see [[docs/adr]]).

## 1. Purpose

Stock-Stonks tracks stock/FII/ETF/BDR trades (compra/venda) and exports data for
Brazilian income tax (Imposto de Renda) reporting, in the spirit of the Bastter
System, starting from an MVP and growing feature by feature.

## 2. Stack

- **Backend**: Go, standard `net/http` (Go 1.22+ routing), GORM as the ORM.
- **Frontend**: Next.js (App Router, TypeScript).
- **Database**: SQLite for the MVP, accessed only through the `repository`
  interfaces — never through GORM directly outside `internal/repository/sqlite`.
  Swapping to Postgres/MySQL later means adding a case in
  `internal/database.Open` plus driver-specific entries in
  `internal/database.IsRetryable`, not rewriting services or handlers.
- **Resilience**: every database call goes through `database.WithRetry`
  (exponential backoff + jitter) so transient errors (e.g. SQLite's
  "database is locked") don't surface as request failures.
- **Monetary values**: always `decimal.Decimal` (shopspring/decimal), never
  `float64` — tax calculations must reproduce exact centavo amounts.

## 3. Architecture shape

```
domain        -> entities + validation, no framework imports
repository    -> interfaces only
repository/sqlite -> GORM + SQLite implementation of those interfaces
service       -> business logic, depends on repository interfaces
httpserver    -> handlers + routing, depends on service
cmd/api       -> wiring (reads config, opens DB, builds services, starts server)
```

A layer never imports a layer below it in this list except through the
interfaces it depends on.

## 4. Process for AI agents working in this repo

- **Spec-driven**: non-trivial features get a spec under `specs/<NNN-slug>/`
  before implementation, using the templates in `specs/templates/`. See the
  [[spec-kit-workflow]] skill reference once it exists, or follow
  `specs/templates/spec-template.md` directly.
- **Decisions get recorded**: any choice that is hard to reverse, affects
  multiple features, or a future contributor would reasonably ask "why did
  we do it this way?" gets an ADR in `docs/adr/`. Use the `adr` skill.
- **Docs stay in sync**: when behavior, endpoints, or setup steps change, the
  relevant doc (README, this constitution, or an ADR) is updated in the same
  change. Use the `documentation` skill.
- **Finding new skills**: when a task would benefit from a skill this repo
  doesn't have yet, use the `find-skills` skill before building one from
  scratch.

## 5. Testing bar

- Pure domain logic (validation, position/average-cost math, retry policy)
  must have unit tests with no I/O.
- `go vet` and `gofmt -l` must be clean; `npm run lint` and `tsc --noEmit`
  must pass on the frontend.
- CI (`.github/workflows/`) enforces all of the above on every push and PR.
