---
name: documentation
description: Use after any change to behavior, API endpoints, configuration, setup steps, or architecture in this repo — before considering the task finished. Keeps README.md, backend/README.md (if present), specs/, and docs/adr/ in sync with what the code actually does, instead of letting docs silently rot.
---

# Documentation

Docs are only trustworthy if they are updated in the same change as the code
they describe. This skill is the checklist for that, applied *after*
implementation, not as a separate pass days later.

## When to use

Run this check whenever a change:
- Adds, removes, or changes an HTTP endpoint, request/response shape, or
  CLI flag.
- Changes a setup/run step (new env var, new dependency, new `make` target).
- Changes the architecture described in `specs/constitution.md`.
- Finishes a feature whose `specs/<NNN-slug>/` spec exists.

## What to check, in order

1. **Root `README.md`** — does "Getting started" still work as written? Are
   new env vars documented? Is the feature list current?
2. **`backend/README.md` / `frontend/README.md`** (if they exist) — same
   question, scoped to that part of the stack.
3. **`specs/<NNN-slug>/spec.md`** — if the feature's actual behavior ended
   up different from what the spec described, update the spec (don't leave
   it describing a feature that no longer matches reality) and note the
   deviation in `plan.md` if it was non-trivial.
4. **`specs/constitution.md`** — only touch this for a change to a
   repo-wide, cross-feature principle (stack choice, layering rule, testing
   bar). If this is the case, it also needs an ADR (see the `adr` skill) —
   the constitution records the rule, the ADR records why it changed.
5. **Code comments** — only add one where the *why* is non-obvious (a
   constraint, an invariant, a workaround). Don't restate what the code
   already says through naming.

## Writing style for this repo

- Prefer short, concrete sections over long prose.
- Every documented endpoint/config value should be something a reader could
  verify by running the code — no aspirational docs for unbuilt features
  (those belong in `specs/`, clearly marked `Status: Draft`).
- Keep examples copy-pasteable (`curl` commands, exact env var names).

## Workflow

1. After implementation, diff what changed against what's documented.
2. Update the affected doc(s) from the list above in the same commit/PR as
   the code change.
3. If nothing needs updating, say so explicitly rather than skipping the
   check silently — it should be visible that the check happened.
