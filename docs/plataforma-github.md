# Plataforma GitHub

## 1. Decisión

GitHub será la plataforma de colaboración y entrega de ReefOps:

- Git y pull requests;
- issues, milestones, Projects y Discussions cuando aporten valor;
- GitHub Actions para CI y automatización;
- GHCR para imágenes y artefactos OCI;
- GitHub Releases para binarios y paquetes destinados a personas;
- artefactos temporales de Actions para resultados de jobs;
- Dependabot y funciones de seguridad disponibles en el plan;
- attestations y SBOM cuando el plan y la visibilidad lo permitan;
- GitHub CLI y MCP/plugin de GitHub para operación asistida.

GitHub no será fuente de verdad de datos del acuario, secretos runtime ni
backups. La pérdida de acceso a GitHub no impedirá operar localmente un sistema
ya desplegado.

## 2. Organización y topología de repositorios

ReefOps se alojará en una organización de GitHub. La organización proporciona
propiedad independiente de una cuenta personal, instalación centralizada de
GitHub Apps, equipos y permisos separados por responsabilidad. El permiso base
será `none`; pertenecer a la organización no concederá escritura.

Se adoptan tres repositorios:

| Repositorio | Propiedad | Visibilidad inicial | Escritura automatizada |
|---|---|---|---|
| `reefops` | Producto, dominios, Angular, servicios, contratos, docs y CI | Público, con historial raíz sanitizado | Actions publica artefactos; no despliega |
| `reefops-platform` | Componentes Kubernetes reutilizables, bootstrap, políticas y toolchain CI | Público | Actions publica charts, bundles e imágenes; no despliega |
| `reefops-gitops` | Composición privada por clúster, configuración no secreta, versiones y digests | Privado siempre | Cambios mediante PR; Flux solo lectura |

La visibilidad pública se limita a `reefops` y `reefops-platform`. Antes de
publicar un repositorio se escanearán tanto su árbol como todo su historial Git.
No contendrá datos de acuarios, nombres de hosts o redes privadas,
identificadores de instalaciones, credenciales, material SOPS, claves,
destinatarios de recuperación ni configuración concreta de un clúster.

`reefops` se publica desde un commit raíz nuevo que excluye las copias
transitorias de infraestructura, `.sops.yaml` y ayudas operativas propias de
plataforma. El historial privado anterior se conserva únicamente como bundle
local cifrado de recuperación, fuera del workspace y sin ninguna referencia
remota.

`reefops-gitops` permanece privado incluso si todos sus secretos están cifrados:
la composición revela topología, versiones desplegadas, ventanas de cambio y
metadatos operativos. Hacer público producto o plataforma no autoriza a copiar
allí configuración específica para eludir esta separación.

La publicación permite aprovechar gratuitamente protección de ramas, análisis
de dependencias y controles de seguridad disponibles para repositorios
públicos. Si un futuro cambio introduce información no publicable, debe
rechazarse o moverse al repositorio privado antes del merge; no se confía en
cambiar posteriormente la visibilidad para reparar una filtración.

En `reefops-platform` público se habilitarán dependency graph, alertas y
actualizaciones de seguridad de Dependabot, secret scanning, comprobación de
validez y push protection siempre que GitHub los ofrezca gratuitamente. Estos
controles complementan Gitleaks y Trivy; no sustituyen el gate previo porque la
publicación del historial ocurre antes de que un control posterior pueda
remediarla.

Un cuarto repositorio `reefops-actions` solo se creará cuando haya workflows
reutilizables usados por más de un repositorio. Firmware, entrenamiento de
modelos o un dominio solo se separarán cuando tengan ciclo de vida, permisos o
responsables realmente independientes. No se creará un repositorio por
microservicio.

`reefops` será un monorepo modular. Los límites entre dominios se impondrán con
contratos de eventos, reglas de imports, ownership y tests; un repositorio no
sustituye un límite arquitectónico.

Los antiguos directorios `infra/` del historial privado se usaron únicamente
como fuente transitoria para la extracción inicial:

- `infrastructure/`, `platform/` y `bootstrap/` se sembrarán en
  `reefops-platform`;
- `clusters/` y la selección concreta de aplicaciones se sembrarán en
  `reefops-gitops`;
- `reefops` conserva únicamente charts y descriptores que pertenezcan al
  producto; actualmente no mantiene estado deseado de clúster.

La creación inicial seguirá una transición controlada:

1. crear la organización desde GitHub y ejecutar `task github-provision` para
   crear los tres repositorios inicialmente privados;
2. publicar el historial local con `task product-publish`;
3. sembrar `reefops-platform` desde la infraestructura reutilizable;
4. publicar sus componentes versionados;
5. sembrar `reefops-gitops` desde la composición de entorno;
6. revisar los tres historiales iniciales;
7. verificar que existe `clusters/local/kustomization.yaml`;
8. ejecutar el bootstrap excepcional, que añade la credencial de solo lectura
   para GitOps; la plataforma pública se obtiene por HTTPS sin credenciales;
9. construir y auditar el historial raíz publicable de `reefops`;
10. ejecutar los gates explícitos de publicación de producto y plataforma;
11. ejecutar `task github-protect` inmediatamente después.

El bootstrap fallará antes de tocar el clúster si el repositorio GitOps no
contiene el estado esperado.

El bootstrap es la única ventana de escritura directa sobre `main`. Se ejecuta
una vez, queda asociado al actor autenticado y se cierra aplicando protección
pull-request-only. Flux runtime conserva acceso de solo lectura tanto a GitOps
como a plataforma.

Una repetición del bootstrap con Flux ya instalado no vuelve a escribir Git. Si
la protección está activa pero los controladores no existen, el procedimiento
falla antes de modificar el clúster y exige una ventana de bootstrap explícita.

El bootstrap fija Flux Toolkit en `v2.9.3`. Cambiar esa versión requiere un PR,
verificación de compatibilidad de CRD/API y prueba de reconciliación; la versión
instalada por el estado accidental de Homebrew no decide los controladores.

`reefops-gitops` fija `reefops-platform` por commit completo. Un merge en
plataforma no modifica un clúster hasta que otro PR promociona explícitamente
ese commit; el rollback restaura el commit anterior, no una rama flotante.
La composición GitOps no contiene copias ni paths relativos hacia los
directorios internos de plataforma.

## 3. Flujo de ramas y releases

Se utilizará GitHub Flow con desarrollo basado en trunk, no GitFlow:

- `main` será la única rama permanente y estará siempre integrable;
- cada cambio usará una rama corta y un pull request;
- el merge será squash por defecto y eliminará la rama;
- no habrá ramas permanentes `develop`, `staging` o `production`;
- los entornos se representan en `reefops-gitops`, no mediante ramas;
- una versión publicable se identifica con tag firmado y artefactos por digest;
- solo se abrirá una rama `release/x.y` si hay que mantener simultáneamente una
  versión anterior.

El clúster Docker Desktop solo materializa `development`. La composición
privada no reservará recursos de `production`; añadirá una raíz de entorno y su
destino cuando sean necesarios. Los componentes reutilizables permanecerán
independientes del clúster y la separación de datos, autoridades y credenciales
se realizará por entorno.

Cuando exista production, la promoción modificará mediante PR el digest
seleccionado después de validarlo en development. El workflow construirá una
vez: production no consumirá un tag flotante ni un artefacto recompilado.

## 4. Flujo GitOps

```text
Pull request en reefops
  → checks obligatorios
  → merge a main
  → build multi-arquitectura
  → SBOM + escaneo
  → imagen GHCR por digest
  → firma/attestation
  → PR en reefops-gitops con el digest
  → checks y aprobación
  → merge
  → Flux detecta Git
  → reconciliación y health checks
```

La primera imagen de producto es
`ghcr.io/reefops/reefops-authorizer-migrator`. El workflow sólo la publica desde
`main` después de la validación del repositorio y una aceptación sobre
PostgreSQL real; usa una etiqueta derivada del commit y entrega SBOM, provenance
y attestation asociados al digest multi-arquitectura. GitOps consumirá el
digest, no la etiqueta.

No habrá despliegue directo desde GitHub Actions al Kubernetes local. El
clúster iniciará todas las conexiones hacia GitHub/GHCR mediante polling y
descarga. No habrá webhook, túnel ni runner hospedado que acepte conexiones
entrantes hacia la red local. Flux utilizará una deploy key de solo lectura
solo para `reefops-gitops`; leerá `reefops-platform` públicamente por HTTPS sin
un Secret Kubernetes.
La automatización de imágenes no tendrá permiso directo sobre `main`; propondrá
un pull request.

El bootstrap utilizará `flux bootstrap github` sobre un repositorio ya creado
y sembrado.
Tomará temporalmente la identidad de `gh auth token` para registrar la deploy
key y no persistirá ese token en Git ni Kubernetes. La clave de Flux será de
solo lectura; los controladores de image automation con escritura no se
instalarán inicialmente.

El bootstrap no creará una deploy key para la plataforma pública. Si la
plataforma volviera a ser privada, esa decisión requeriría restaurar
explícitamente una credencial independiente y actualizar la fuente GitOps; no
se reutilizará la clave de GitOps.

Si el owner coincide con el usuario autenticado, el bootstrap seleccionará el
modo personal de Flux; en otro caso lo tratará como organización.

## 5. Artefactos

| Tipo | Destino |
|---|---|
| Imágenes OCI multi-arquitectura | GHCR |
| Charts Helm | GHCR como OCI |
| Binarios CLI o bundles instalables | GitHub Releases |
| SBOM y procedencia | Attestation asociada y asset cuando proceda |
| Informes de tests y análisis | Artefacto temporal de Actions |
| Datos, backups y secretos | Fuera de GitHub |

Las imágenes se referenciarán por digest en GitOps. Los tags facilitarán
descubrimiento, pero no serán autoridad de despliegue.

Las attestations de GitHub se habilitarán únicamente si están disponibles para
la visibilidad y plan elegidos. Cosign seguirá permitiendo firma verificable
mediante identidad OIDC de Actions. Producir una attestation sin verificarla en
el consumo no satisface el control de cadena de suministro.

## 6. Secretos, seguridad y gobierno

OpenBao será la autoridad local de secretos runtime. Los workloads se
autenticarán con identidades de Kubernetes y recibirán secretos de corta
duración o montados en memoria sin incorporar SDK de OpenBao al dominio. SOPS
y `age` se limitarán a material de bootstrap cifrado; las claves privadas y de
recuperación permanecerán fuera de Git y del clúster respaldado.

GitHub Secrets podrá contener únicamente secretos necesarios para CI y nunca
será la fuente primaria:

- se preferirán `GITHUB_TOKEN`, GitHub Apps y OIDC frente a secretos duraderos;
- un secreto CI nacerá y rotará en OpenBao;
- una tarea local autenticada leerá una versión concreta y la enviará cifrada
  a GitHub mediante entrada estándar;
- se registrarán nombre, versión, destino, actor y fecha, nunca el valor;
- la sincronización será solo saliente desde la instalación local;
- GitHub Actions no accederá a OpenBao ni a Kubernetes local;
- los secretos runtime, datos del acuario y credenciales de dispositivos no se
  copiarán a GitHub Secrets.

Un GitHub Secret es una réplica de entrega revocable. Si pierde su relación con
una versión activa de OpenBao, el workflow deberá fallar de forma cerrada hasta
rotarlo o retirarlo.

- `main` y tags de release estarán protegidos mediante rulesets.
- Pull request obligatorio, historial lineal y bloqueo de force-push/delete.
- Checks requeridos después de que el primer workflow haya ejecutado y exista
  su nombre estable.
- Acciones externas fijadas por commit SHA completo.
- `permissions: {}` por defecto y permisos mínimos por job.
- `GITHUB_TOKEN` en lugar de PAT para publicar artefactos del mismo repo.
- OIDC para firma y accesos federados; no almacenar credenciales cloud
  permanentes.
- Environments para operaciones que requieran aprobación.
- CODEOWNERS para infraestructura, seguridad, contratos y migraciones cuando
  exista más de un mantenedor.
- Secret scanning, Dependabot, dependency review y CodeQL según disponibilidad.
- Renovate no se añadirá mientras Dependabot cubra los ecosistemas necesarios.

La primera CI utiliza un `Brewfile.ci` mínimo y nombres de fórmulas versionadas
cuando Homebrew los ofrece. Esto fija intención pero no el contenido exacto de
todas las botellas. Se registra como deuda hasta publicar una imagen de
validación en GHCR por digest; desde ese momento los jobs utilizarán ese digest
para obtener una toolchain reproducible.

Los PAT se limitarán a bootstrap o integraciones que no admitan GitHub App,
deploy key, `GITHUB_TOKEN` u OIDC. Tendrán scopes mínimos, caducidad y
revocación documentada.

## 7. Trazabilidad

Cada build y propuesta de despliegue conservará:

- repositorio, commit, PR y workflow run;
- actor que aprobó y actor técnico;
- digest, SBOM, firma o attestation;
- resultado de tests y escaneos;
- PR y commit del repositorio GitOps;
- reconciliación Flux y `deployment_id`;
- rollback o compensación posterior.

La relación mínima será:

```text
requisito/issue
  → PR de código
  → commit
  → workflow
  → digest GHCR
  → PR GitOps
  → commit GitOps
  → reconciliación Flux
  → deployment_id
```

No se copiarán tokens, secretos ni payloads privados a logs, issues o
artefactos de Actions.

La sincronización de un secreto CI añadirá evidencia no sensible:

```text
secret_id + versión OpenBao
  → actor local + autorización
  → repositorio/entorno GitHub destino
  → instante y resultado
  → workflow que consumió la réplica
```

## 8. Licencia

El código que se publique se distribuirá bajo Apache License 2.0. Su concesión
explícita de patentes ofrece una base más adecuada que MIT para una plataforma
con automatización, dispositivos e IA, sin impedir uso comercial ni derivados
privados.

La licencia no alcanza:

- `reefops-gitops`;
- datos o medios del acuario;
- backups;
- secretos y configuraciones privadas;
- pesos o datasets cuya licencia sea distinta;
- marcas y material de terceros.

Cada dependencia, modelo, dataset y asset conservará su licencia y procedencia.
Publicar `reefops` o `reefops-platform` requerirá antes una revisión automática
de compatibilidad de licencias.

## 9. Herramientas locales

Se gestionarán mediante Homebrew:

- `gh`;
- `actionlint`;
- `cosign`;
- `syft`;
- `trivy`;
- `oras`.

El plugin/MCP de GitHub se conectará mediante Codex y respetará la autorización
del usuario. Las operaciones mutables externas seguirán requiriendo alcance
explícito; disponer del MCP no autoriza crear repositorios, rulesets, releases
o modificar pull requests sin una tarea que lo solicite.

## 10. Decisiones pendientes

- Plan de GitHub y disponibilidad de Advanced Security/attestations privadas.
- Política de aprobación y CODEOWNERS cuando exista más de un mantenedor.
- Retención y límites de artefactos de Actions.
- Digest inicial de la imagen GHCR usada para validación reproducible.
