# 0005. Position (quantity, average cost, realized profit) is always computed, never stored

## Status
Accepted

## Context
Brazilian IR rules require the "custo médio ponderado" (weighted average
cost) method: each buy updates the average cost, each sell realizes
profit/loss against the *current* average cost without changing it. This
calculation must replay a given asset's operations in date order to be
correct. If a past operation is later edited or deleted (corrections happen
— a user fixes a typo'd price), any stored/cached position snapshot
silently becomes wrong until something notices and recomputes it.

## Decision
`domain.Position` is never persisted. `service.PositionService` recomputes
it on every request by loading the relevant operations (ordered by date)
and folding them through `domain.Position.ApplyOperation`.

## Consequences
- Correctness by construction: there is no cache-invalidation bug class for
  positions, because there is no cache.
- Cost is O(n) operations per position request. Acceptable at MVP scale
  (single user, realistically hundreds to low thousands of operations).
- If this becomes a real latency problem, the fix is a materialized
  snapshot *plus* an event log to replay from a checkpoint — not caching the
  final numbers directly — to preserve the same correctness guarantee. That
  change would need its own ADR superseding this one.

## Alternatives considered
- **Store the position and update it incrementally on each write** —
  rejected: correct only if operations are always appended in
  chronological order and never edited/deleted after the fact, which this
  domain cannot guarantee (tax-relevant corrections will happen).
