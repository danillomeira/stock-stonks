# Architecture Decision Records

Numbered, append-only log of decisions that are hard to reverse or involved
a real tradeoff. See the `adr` skill (`.claude/skills/adr/SKILL.md`) for
when and how to add one.

| # | Title |
|---|---|
| [0001](0001-record-architecture-decisions.md) | Record architecture decisions as ADRs |
| [0002](0002-monorepo-with-go-backend-and-nextjs-frontend.md) | Monorepo with a Go backend and a Next.js frontend |
| [0003](0003-gorm-sqlite-with-retry-backoff-seam.md) | GORM + SQLite for the MVP, behind a retry/backoff seam for future databases |
| [0004](0004-decimal-for-monetary-values.md) | Use shopspring/decimal for all monetary and quantity fields |
| [0005](0005-position-is-computed-not-stored.md) | Position is always computed, never stored |
| [0006](0006-spec-kit-style-workflow-for-ai-agents.md) | Spec-kit style spec-driven workflow for AI agents |
| [0007](0007-ci-only-pipeline-no-deploy-target-yet.md) | CI-only pipeline for now; no live deploy target yet |
| [0008](0008-reuse-catalog-for-components-and-helpers.md) | Maintain a reuse catalog for generic components/helpers |
