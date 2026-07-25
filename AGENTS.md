# ReefOps agent guidance

## Sources of truth

- Read `docs/requisitos-funcionales.md`, `docs/arquitectura.md` and
  `docs/alineacion-requisitos-arquitectura.md` for behavior and boundaries.
- Read `docs/desarrollo-asistido.md` for Codex aids and
  `docs/entorno-desarrollo.md` for host toolchains when those areas change.
- Read `docs/plataforma-github.md` for repositories, CI, artifacts, releases
  and GitOps promotion.
- Read the nearest nested `AGENTS.md` before editing its subtree.
- Treat Git as the desired-state source for infrastructure.

## Required workflow

1. Document the requirement or architectural decision first.
2. Align ownership, contracts, security and traceability.
3. Implement only after the documentation is updated.
4. Run the narrowest relevant checks, then `task validate`.
5. Report validation and any remaining decision explicitly.

Update the RF alignment matrix only when functional requirements, ownership,
components or evidence change. Do not invent an RF identifier for tooling,
documentation governance or other non-functional work.

Never implement first and document afterward, except an emergency needed to
protect animal life or data. Reconcile emergency work immediately.

## Architecture rules

- Organize code by aquarium domain using Screaming Architecture.
- A business domain must not import, query, join, call or reference another
  domain's internals.
- Domain-to-domain reactions use versioned integration events.
- Use transactional outbox, consumer inbox and idempotent handlers.
- Cross-domain reads use event-built projections.
- The producer must not know its consumers.
- Platform services are not business domains.
- AI inference, vision, embeddings, STT and TTS run locally.

## Traceability

- Propagate `correlation_id` across requests, commands, events, jobs,
  inference, notifications and publication.
- Preserve immediate `causation_id`; never reset correlation at a boundary.
- Record actor/delegation, authorization decision, source versions, result and
  provenance for sensitive changes.
- Retrying or replaying must not repeat physical actions or notifications.
- Logs and traces support diagnosis but never replace functional audit.

## Infrastructure

- This public product repository does not own cluster desired state.
- Make platform changes in `reefops-platform` and environment composition in
  the private `reefops-gitops` repository, following their nearest guidance.
- Never copy plaintext secrets, private age keys, cluster composition or
  floating production versions into this repository.
- Remember that Git revert does not restore data, migrations or physical
  effects.

## Validation

- Run `task prerequisites` before cluster-related work.
- Run `task validate` before handing off any repository change.
- Add focused tests when changing executable behavior.
- Keep the worktree free of generated secrets and temporary artifacts.

## Code review rules

- Flag synchronous domain coupling, cross-domain database access and shared
  business models.
- Flag missing authorization, idempotency, causation or audit evidence.
- Flag replay paths that can duplicate dosing, equipment commands or alerts.
- Flag cloud inference or transmission of private aquarium context.
- Flag imperative infrastructure drift and unrecoverable migrations.

## Specialized reviewers

- Keep implementation and integration in the primary agent.
- Use `architecture_reviewer` for non-trivial boundaries, ownership, contracts
  or multi-domain changes.
- Use `traceability_reviewer` for sensitive commands/events, AI, devices,
  publication, authorization, replay or deletion.
- Use the reviewers and instructions of the platform or GitOps repository for
  IaC, persistence, secrets, migrations or recovery changes.
- Reviewers are read-only. Their findings do not replace executable checks.
- Do not delegate trivial edits. Independent applicable reviews may run in
  parallel within the configured limit.
