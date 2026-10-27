# 0008. Maintain a reuse catalog for generic components/helpers

## Status
Accepted

## Context
Like specs (ADR 0006) and decisions (ADR 0001), discoverability of
*reusable code* has the same disconnected-session problem: an AI agent with
no memory of past sessions has to re-read the codebase to find out whether
a generic helper or component already exists before deciding whether to
write a new one. Without an index, that search costs tokens and turns every
time, and the usual failure mode is a near-duplicate getting written
because the existing one wasn't found.

## Decision
Add the `reuse` skill (`.claude/skills/reuse/SKILL.md`) and a flat index,
`docs/catalog.md`, listing generic/reusable components and helpers (name,
one-line purpose, file path) for both `backend/internal/` and `frontend/`.
The rule: check the catalog before writing new code; if what you build is
generic enough for another feature to plausibly need, extract it to a
shared location and add a catalog line in the same change.

## Consequences
- Checking "does this already exist" becomes one grep of a short file
  instead of a codebase search.
- The catalog is only as good as the discipline to update it — it will
  drift if extraction happens without the catalog line being added in the
  same change. The `reuse` skill states this explicitly to keep it from
  becoming a separate, forgettable step.
- Small overhead per reusable piece (one line), justified by the token/time
  cost it saves every future session that would otherwise re-derive the
  same search.

## Alternatives considered
- **Rely on grepping the codebase each time** — rejected: works, but is
  exactly the repeated-cost problem this ADR exists to avoid; a short index
  is cheaper to read than the code it points to.
- **Enforce reuse via a linter/dependency-graph tool** — rejected for now:
  overkill for the current codebase size; revisit if real duplication
  starts accumulating despite the catalog (check `docs/catalog.md`'s "Known
  gaps" section — it already flags one case to extract later).
