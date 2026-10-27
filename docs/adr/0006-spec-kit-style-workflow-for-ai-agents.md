# 0006. Spec-kit style spec-driven workflow for AI agents

## Status
Accepted

## Context
This project is built largely by AI coding agents across disconnected
sessions with no shared memory. Without a structured workflow, each session
risks re-deriving requirements, missing non-goals, or losing track of why a
tradeoff was made — the same problem ADR 0001 addresses for decisions, but
for *features* before they're built.

## Decision
Adopt a spec-kit style structure (similar to GitHub's spec-kit /
constitution-driven development) under `specs/`:

- `specs/constitution.md` — repo-wide non-negotiable principles.
- `specs/templates/{spec,plan,tasks}-template.md` — templates for new
  features.
- `specs/<NNN-slug>/{spec.md,plan.md,tasks.md}` — per-feature: problem/goals
  (spec), affected layers/risks (plan), checkable steps (tasks).

Non-trivial features get this treatment before implementation starts. The
first one (`specs/001-mvp-trade-tracking/`) documents the MVP built in this
same session, written retroactively to seed the practice with a real
example rather than a hypothetical one.

## Consequences
- A new session (human or agent) can read `constitution.md` + the relevant
  `specs/<slug>/` folder and have the context an ADR alone wouldn't give
  (ADRs record decisions; specs record what a feature is supposed to do).
- Adds process overhead for genuinely small changes — use judgment; a
  one-line bug fix does not need a spec.
- Requires discipline to keep specs/plans updated when real implementation
  deviates from what was planned (see the `documentation` skill).

## Alternatives considered
- **No structured process, just issues/PR descriptions** — rejected: issues
  live outside the repo for an agent reading the checked-out tree, and PR
  descriptions aren't organized by feature over time.
- **Full GitHub spec-kit tool** — rejected for now: adopting its file
  format/conventions manually keeps this dependency-free; revisit via the
  `find-skills` skill if spec-kit tooling itself becomes valuable to install.
