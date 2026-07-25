# Desarrollo asistido de ReefOps

## 1. Objetivo

Las ayudas para Codex forman parte del repositorio y deben mejorar
repetibilidad, revisión y seguridad sin sustituir controles ejecutables.

Se seguirá este orden:

1. `AGENTS.md` para reglas duraderas y alcance general.
2. Skills para workflows repetibles.
3. Agentes especializados para revisiones acotadas y paralelizables.
4. MCP únicamente cuando haga falta consultar o modificar un sistema externo.
5. Hooks o CI para reglas que deban aplicarse mecánicamente.

Todas las ayudas respetarán la decisión documentation-first: documentar,
alinear, implementar y validar.

## 1.1 Estado de decisiones

| Decisión | Estado | Evidencia |
|---|---|---|
| `AGENTS.md` raíz | Adoptada e implementada | `AGENTS.md` |
| Instrucciones por subárbol | Adoptada para docs y automatización GitHub | `docs/AGENTS.md`, `.github/AGENTS.md` |
| Skills de proyecto | Adoptada e implementada | `.agents/skills/reefops-*` |
| Revisores especializados | Adoptada e implementada | `.codex/agents/*.toml` |
| Concurrencia máxima de agentes | Adoptada: 3 subagentes | `.codex/config.toml` |
| Agentes con capacidad de escritura | No adoptada | El agente principal conserva la integración |
| MCP/plugin de GitHub | Instalado; `gh` autenticado cubre operaciones no disponibles para la App | Plugin GitHub de Codex y GitHub CLI |
| Hooks Codex | Pendiente de una regla mecánica concreta | Ningún hook configurado |
| CI remoto | Implementado en GitHub Actions para el repositorio de producto | `.github/workflows/validate.yml` y validación local mediante Task |
| Plugin distribuible | No necesario mientras las skills sean propias del repositorio | Skills locales versionadas |

## 2. Instrucciones `AGENTS.md`

El archivo raíz contendrá:

- propósito y fuentes de verdad;
- secuencia documentation-first;
- prohibición de dependencias entre dominios;
- reglas de GitOps;
- trazabilidad obligatoria;
- comandos de validación;
- límites de seguridad y datos.

Los subárboles podrán tener un `AGENTS.md` más específico. Inicialmente:

- `docs/`: coherencia de requisitos, arquitectura y matriz;
- `.github/`: automatización, permisos mínimos y publicación segura;
- en el futuro, cada servicio con sus pruebas y reglas propias.

Las instrucciones cercanas concretan las generales, pero no pueden relajar
seguridad, trazabilidad ni boundaries sin una nueva decisión arquitectónica.

## 3. Skills del repositorio

Las skills se versionarán en `.agents/skills`, la ubicación de proyecto
descubierta por Codex. El conjunto inicial será deliberadamente pequeño:

| Skill | Uso |
|---|---|
| `reefops-change` | Cambios generales o transversales; cede ante una skill especializada |
| `reefops-domain-module` | Diseñar o implementar módulos autónomos, comandos, eventos y proyecciones |

Las skills serán concisas. Los detalles canónicos seguirán en `docs/`; una
skill enlazará las referencias que deba leer y no duplicará la arquitectura.
Cada skill tendrá metadata de interfaz y será validada con la herramienta
oficial de creación de skills.

## 4. Agentes especializados

Los agentes de proyecto residirán en `.codex/agents`. Inicialmente serán
revisores de solo lectura:

| Agente | Responsabilidad |
|---|---|
| `architecture_reviewer` | Detectar acoplamientos, ownership ambiguo y violaciones de Screaming Architecture |
| `traceability_reviewer` | Verificar correlación, causación, auditoría, idempotencia y reproducción |

El agente principal seguirá siendo responsable de integrar resultados. Los
revisores no editarán archivos, desplegarán infraestructura ni compartirán
estado mutable. No se delegará trabajo por defecto: se usarán cuando el cambio
sea suficientemente grande y la revisión pueda ejecutarse de forma
independiente.

Criterios de uso:

- `architecture_reviewer`: cambios de boundaries, contratos, ownership o más
  de un dominio;
- `traceability_reviewer`: comandos/eventos sensibles, IA, dispositivos,
  publicación, autorización, replay o eliminación;
- ningún revisor para ediciones triviales cuyo coste exceda el riesgo.

Las skills y revisores específicos de GitOps no se cargan desde el repositorio
de producto público. Pertenecen al workspace de plataforma y al repositorio
privado de composición, evitando que esta ayuda sugiera ownership inexistente.

Cuando un cambio afecte a varios ejes se podrán ejecutar revisiones en
paralelo, hasta el máximo configurado. Recibir el informe de un agente no
sustituye la validación ejecutable ni transfiere la decisión final.

## 5. MCP y conectores de desarrollo

No se configurará un MCP sin una necesidad externa concreta. Un MCP aumenta la
superficie de permisos y no es un sustituto de scripts locales, documentación o
APIs del producto.

Conector adoptado:

- GitHub para repositorios, issues, pull requests, CI, releases y flujo
  GitOps. La conexión no amplía por sí sola la autorización de una tarea.

Candidatos posteriores:

- un servidor de documentación técnica autorizado;
- observabilidad de un entorno profesional;
- el MCP funcional ofrecido por ReefOps, cuando se implemente, tratado como
  sistema externo desde la perspectiva de Codex.

Cada incorporación deberá documentar:

- propietario y finalidad;
- herramientas y recursos habilitados;
- operaciones de lectura y escritura;
- autenticación, scopes y secretos;
- datos que pueden abandonar la instalación;
- auditoría, revocación y comportamiento degradado.

No se conectarán a terceros imágenes, telemetría ni contexto privado del
acuario para inferencia.

Las decisiones completas de uso de GitHub se encuentran en
[Plataforma GitHub](plataforma-github.md).

## 6. Configuración de proyecto

`.codex/config.toml` habilitará un máximo conservador de agentes concurrentes.
No fijará modelo, permisos amplios, credenciales ni MCP personales. La
configuración solo se carga cuando el repositorio está marcado como confiable.

Los agentes especializados tendrán `sandbox_mode = "read-only"`. Los cambios
reales permanecerán en el agente principal y seguirán las autorizaciones de la
sesión.

## 7. Controles ejecutables

Las instrucciones se complementarán con:

- `task validate`;
- validación exacta de identificadores `RF-*`;
- validación de skills;
- pruebas de contratos y arquitectura cuando exista código;
- CI como autoridad remota y Task como puerta local reproducible.

Una regla crítica no se considerará protegida solo porque aparezca en un
prompt, skill o `AGENTS.md`.

## 8. Matriz de alineación de ayudas

| Decisión arquitectónica | Instrucción/skill | Revisor | Control ejecutable |
|---|---|---|---|
| AD-007 Screaming Architecture | `AGENTS.md`, `reefops-domain-module` | `architecture_reviewer` | Pruebas de arquitectura cuando exista código |
| AD-008A integración por eventos | `AGENTS.md`, `reefops-domain-module` | `architecture_reviewer` | Contratos, outbox/inbox y tests de consumidores |
| AD-010 trazabilidad completa | `AGENTS.md`, `reefops-change`, `reefops-domain-module` | `traceability_reviewer` | Tests E2E de correlación y auditoría |
| AD-017 operación reproducible | Repositorios plataforma/GitOps | Revisión propia de plataforma | Promoción por SHA y checks externos |
| AD-018 IA exclusivamente local | `AGENTS.md`, `reefops-change` | arquitectura y trazabilidad | Tests de red y configuración cuando exista runtime |
| AD-019 equipo vital independiente | `AGENTS.md` y plataforma externa | trazabilidad | Pruebas de degradación futuras |
| AD-020 documentation-first | Todos los `AGENTS.md` y skills | Agente principal | Matriz RF y `task validate` |
| AD-021 GitOps desde bootstrap | Repositorios plataforma/GitOps | Revisión propia de plataforma | Reconciliación Flux externa |
| AD-022 gobierno de ayudas | `docs/desarrollo-asistido.md` | Agente principal | `check-codex-aids.sh` |
| AD-023 GitHub como plataforma | `.github/AGENTS.md`, `reefops-change` | Agente principal | `actionlint` y `check-policies.sh` |

Las filas con controles futuros describen deuda explícita, no una garantía
actual. Al incorporarse código o CI deberán convertirse en comprobaciones
ejecutables antes de considerar satisfecha la decisión.
