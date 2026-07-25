---
name: reefops-domain-module
description: Design, implement, or review an autonomous ReefOps business-domain module using Screaming Architecture, commands, events, outbox/inbox, projections, authorization, and complete traceability. Use for new aggregates, use cases, event contracts, process managers, extraction boundaries, or changes under a business domain; do not use for platform-only Kubernetes work.
---

# ReefOps Domain Module

## Establish the boundary

1. Read `AGENTS.md`, AD-007, AD-008A, AD-010, sections 14–15 of
   `docs/arquitectura.md`, and the affected rows of
   `docs/alineacion-requisitos-arquitectura.md`.
2. Name the business capability and owning domain.
3. List invariants, commands, queries, errors and integration events.
4. Identify data owned by the module and reject foreign keys or joins to other
   domains.
5. Define event-built projections needed from external facts.

Stop and document an architectural decision if ownership remains ambiguous.

## Design the write path

Use this shape:

```text
authorized command
  → aggregate invariant
  → state + audit + outbox in one transaction
  → versioned integration event
```

Specify idempotent `command_id`, optimistic aggregate version, actor context,
authorization decision, `correlation_id` and `causation_id`.

## Design reactions

- Consume only versioned integration events.
- Use inbox deduplication and record processing outcome.
- Build local projections; never call the producer.
- Emit a new fact when the module changes its own state.
- Use a process manager for multi-event workflows without distributed
  transactions.
- Separate retry, DLQ and replay from new business facts.
- Suppress physical, notification and external side effects during projection
  rebuild.

## Verify

Test invariants, idempotency, optimistic concurrency, outbox atomicity,
consumer duplication, ordering assumptions, authorization and audit chain.
Run `task validate` after documentation and code checks.
Use `architecture_reviewer` for non-trivial boundary or multi-domain changes
and `traceability_reviewer` when the flow is sensitive.
