# Agent instructions for Stock-Stonks

This file follows the [agents.md](https://agents.md) convention: instructions
any coding agent should read before working in this repo.

## What this project is

Stock-Stonks tracks stock/FII/ETF/BDR trades and will export data for
Brazilian Imposto de Renda (income tax) reporting — in the spirit of the
Bastter System, built incrementally.

## Read these first

1. `specs/constitution.md` — non-negotiable stack/architecture rules.
2. `docs/adr/README.md` — why past decisions were made the way they were.
3. The relevant `specs/<NNN-slug>/spec.md` if one exists for the feature
   you're touching.

## Skills

This repo ships project-level skills in `.claude/skills/`:

- **find-skills** — search skills.sh for an existing skill before building
  ad-hoc tooling for a recurring kind of task.
- **adr** — when and how to record an Architecture Decision Record.
- **documentation** — keep README/specs/ADRs in sync after a change.

## Workflow for non-trivial changes

1. Check whether a spec already exists under `specs/`. If not and the
   change is non-trivial, write one using `specs/templates/spec-template.md`
   → `plan-template.md` → `tasks-template.md`.
2. Implement, respecting the layering in `specs/constitution.md` §3
   (domain → repository → repository/sqlite → service → httpserver → cmd).
3. Write/run tests: `cd backend && go test ./... -race -cover`, and
   `cd frontend && npm run lint && npx tsc --noEmit && npm run build`.
4. Update docs (the `documentation` skill) and write an ADR if the change
   made a hard-to-reverse decision (the `adr` skill).

## Running it locally

```sh
# Backend
cd backend && go run ./cmd/api        # http://localhost:8080

# Frontend (separate shell)
cd frontend && npm install && npm run dev   # http://localhost:3000

# Or both, containerized (Docker or Podman — the compose file is
# runtime-agnostic):
podman-compose up --build
# (docker compose up --build works identically)
```

## Do not

- Use `float64` for money or share quantities — see
  `docs/adr/0004-decimal-for-monetary-values.md`.
- Call a GORM/`*sql.DB` method directly outside `backend/internal/repository/sqlite/`
  — go through the `repository` interfaces.
- Add a new database driver without reading
  `docs/adr/0003-gorm-sqlite-with-retry-backoff-seam.md` first.
