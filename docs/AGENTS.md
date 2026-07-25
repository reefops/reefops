# Documentation guidance

- Update documentation before implementation.
- Preserve stable RF identifiers; add a new identifier instead of silently
  repurposing an existing requirement.
- Keep `requisitos-funcionales.md`, `arquitectura.md` and
  `alineacion-requisitos-arquitectura.md` consistent.
- Keep `desarrollo-asistido.md` aligned with `AGENTS.md`, `.agents/skills`,
  `.codex/agents` and `.codex/config.toml`.
- Keep `entorno-desarrollo.md`, `Brewfile` and version files aligned.
- Every RF identifier must appear exactly in the alignment matrix.
- Record resolved decisions as decisions, not as pending questions.
- Distinguish functional audit from technical observability.
- Use precise aquarium terminology and state uncertainty where domain evidence
  is incomplete.
- Run `task validate` after changing requirements or the alignment matrix.
