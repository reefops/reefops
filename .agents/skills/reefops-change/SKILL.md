---
name: reefops-change
description: Execute general ReefOps product, documentation, application, or cross-cutting changes with the mandatory documentation-first sequence. Use when a task changes behavior, requirements, architecture, security, traceability, public/shared/private access, or more than one project area; do not use for a narrowly scoped GitOps-only or domain-module-only change when the specialized ReefOps skill covers it.
---

# ReefOps Change

## Workflow

1. Read the root `AGENTS.md` and every applicable nested `AGENTS.md`.
2. Inspect Git status and preserve unrelated work.
3. Locate the affected RF identifiers in:
   - `docs/requisitos-funcionales.md`;
   - `docs/arquitectura.md`;
   - `docs/alineacion-requisitos-arquitectura.md`.
   If the change is non-functional, identify the applicable AD or operational
   document instead; do not invent an RF identifier.
4. Document the requirement or decision before editing implementation.
5. Confirm owner, boundary, events, authorization, data, privacy and
   traceability.
6. Implement the smallest complete vertical slice.
7. Add focused tests and contract checks.
8. Run the narrowest checks, then `task validate`.
9. Report changed behavior, validation and unresolved decisions.

## Guardrails

- Never create synchronous domain-to-domain dependencies.
- Never use another domain's table, repository, entity or business type.
- Preserve correlation and causation through every asynchronous boundary.
- Treat AI output as a proposal or observation until the owning domain accepts
  it.
- Keep inference and private aquarium data local.
- Do not turn technical logs into the functional audit source.
- Do not claim Git rollback restores data, migrations or physical effects.

## Documentation outcome

Update the RF alignment matrix only when an RF identifier, owner, component or
evidence changes. For development aids, update
`docs/desarrollo-asistido.md`; for toolchains, update
`docs/entorno-desarrollo.md`. Add a stable RF identifier instead of
repurposing an unrelated one.
