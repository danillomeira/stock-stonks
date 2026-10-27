# 0007. CI-only pipeline for now; no live deploy target yet

## Status
Accepted

## Context
The project needs a complete GitHub Actions pipeline as the very first
piece of infrastructure. No hosting target (VPS, Fly.io, Railway, etc.) has
been chosen yet, and committing to one now would mean managing secrets and
infra for a target that might change once real usage patterns are known.

## Decision
`.github/workflows/` runs on every push/PR:
- `backend-ci.yml` — `go vet`, `gofmt -l`, `go test ./... -race -cover`,
  `go build`, path-filtered to `backend/**`.
- `frontend-ci.yml` — `npm ci`, `npm run lint`, `tsc --noEmit`,
  `npm run build`, path-filtered to `frontend/**`.
- `docker-publish.yml` — on push to `main`, builds `backend/Dockerfile` and
  `frontend/Dockerfile` and pushes both images to GitHub Container Registry
  (`ghcr.io`), tagged with the commit SHA and `latest`. No deployment step.

Actually running these images somewhere is deferred to a follow-up ADR once
a hosting target is chosen.

## Consequences
- The pipeline is genuinely "complete" for its current scope (lint, test,
  build, package) without taking on infra decisions prematurely.
- Deploying today means manually pulling the published image — acceptable
  for an MVP with no external users yet.
- Adding real deployment later is additive: a new workflow (or job) that
  pulls the already-published image and ships it, without touching the
  CI jobs above.

## Alternatives considered
- **Deploy directly to Fly.io/Railway now** — rejected: no hosting decision
  has been made, and premature infra commitment would need unwinding if the
  target changes.
- **Skip Docker, just build binaries as artifacts** — rejected: container
  images are the more portable handoff to whatever hosting target gets
  chosen later, and GHCR is free and already tied to the GitHub repo.
