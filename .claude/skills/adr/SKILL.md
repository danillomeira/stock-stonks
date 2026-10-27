---
name: adr
description: Use whenever a decision is made in this repo that is hard to reverse, affects multiple features or layers, chooses between real tradeoffs (not a mechanical fix), or that a future contributor would reasonably ask "why did we do it this way?" — e.g. picking a library/driver, a data model shape, an API contract, an auth strategy, or deviating from the constitution. Records it as an ADR under docs/adr/ instead of letting the reasoning live only in a commit message or chat history.
---

# Architecture Decision Records (ADR)

Records *why*, not just *what*. Code and commit history already show what
changed; an ADR is the only durable record of the reasoning, the
alternatives considered, and why they were rejected.

## When to write one

Write an ADR when the decision:
- Is hard to reverse once other code depends on it (a DB driver, a public
  API shape, an auth scheme, a monorepo-vs-polyrepo split).
- Involves a real tradeoff — multiple viable options existed and one was
  picked for a reason worth preserving.
- Deviates from, or clarifies, `specs/constitution.md`.
- A future contributor (human or agent) would plausibly ask "why is it like
  this?" and the answer isn't obvious from reading the code.

Do **not** write one for:
- Routine bug fixes.
- Naming/style choices with no real alternative.
- Anything already fully explained by the spec in `specs/<feature>/spec.md`
  or `plan.md` — link to those instead of duplicating.

## Format

Use `docs/adr/NNNN-short-title.md` (4-digit, zero-padded, sequential — check
the highest existing number in `docs/adr/` first). Michael Nygard's template:

```markdown
# NNNN. Title in imperative mood (e.g. "Use GORM with a retry/backoff wrapper")

## Status
Proposed | Accepted | Deprecated | Superseded by [NNNN](NNNN-slug.md)

## Context
What forces are at play? What problem needs a decision? Keep this factual —
describe the situation, not the solution.

## Decision
What we decided, stated plainly. One or two sentences.

## Consequences
What becomes easier or harder as a result. Include the honest downsides, not
just the upsides — a one-sided ADR is a sales pitch, not a decision record.

## Alternatives considered
What else was on the table and why it was not chosen.
```

## Workflow

1. Check `docs/adr/` for the next sequential number and for any existing ADR
   this might supersede.
2. Write the ADR using the format above — keep it short (a page or less);
   this is a decision record, not a design doc.
3. If it supersedes a prior ADR, update that ADR's `Status` line to point
   forward, and this one's to point back.
4. Cross-link from `specs/constitution.md` or the relevant `specs/<feature>/`
   folder when the decision affects how future specs should be written.
5. Commit the ADR in the same change as the decision it records — never as
   an afterthought in a later, unrelated commit.
