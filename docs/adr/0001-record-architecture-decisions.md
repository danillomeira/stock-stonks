# 0001. Record architecture decisions as ADRs

## Status
Accepted

## Context
This repo is built incrementally, often by AI agents across many sessions
with no shared memory of *why* earlier choices were made. Commit messages
and code show *what* changed, not the reasoning or the rejected
alternatives. Without a durable record, future work (by a human or an
agent) risks re-litigating settled questions or silently violating an
earlier constraint nobody remembers.

## Decision
Record every hard-to-reverse or multi-option decision as a numbered ADR
under `docs/adr/`, using Michael Nygard's format, per the `adr` skill in
`.claude/skills/adr/SKILL.md`.

## Consequences
- Decisions have a single, searchable source of truth.
- Adds small overhead per significant decision (a page of writing).
- Requires discipline to keep ADRs short and factual rather than letting
  them become design docs.

## Alternatives considered
- **Rely on commit messages/PR descriptions** — rejected: these get buried
  in history and aren't organized by topic.
- **A wiki page** — rejected: lives outside the repo, drifts out of sync,
  and isn't visible to an AI agent working from the checked-out tree.
