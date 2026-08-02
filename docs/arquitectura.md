# ReefOps — Arquitectura

La cobertura respecto a los requisitos se mantiene en la
[matriz de alineación](alineacion-requisitos-arquitectura.md).

## 1. Objetivos

ReefOps se desarrollará inicialmente en un entorno local mediante contenedores,
pero su arquitectura permitirá ejecutar la misma aplicación:

1. en un equipo local del usuario;
2. en un servidor local, NAS o mini-PC;
3. en un servidor privado accesible por Internet;
4. en infraestructura cloud;
5. en un modelo híbrido con núcleo local y publicación externa.

El sistema debe continuar cuidando y registrando el acuario aunque se pierda la
conexión a Internet. Publicar o compartir información no deberá implicar exponer
directamente la red local.

## 2. Decisiones fundamentales

### AD-001. Local-first

El núcleo de ReefOps y la fuente de verdad se ejecutarán localmente por defecto.

- Los datos pertenecen a la instalación local.
- La interfaz funciona en la red local sin servicios externos obligatorios.
- Mediciones, tareas, bitácora, sensores y reglas básicas siguen disponibles
  sin Internet.
- La nube será una opción de despliegue o extensión, no una dependencia
  estructural.
- Las integraciones externas degradarán de forma controlada.

Local-first no significa limitarse a un único dispositivo. Una instancia local
podrá atender teléfonos, tabletas y ordenadores de la misma red.

### AD-002. Aplicación modular monolítica

La primera implementación del núcleo será un monolito modular en Go, acompañado
por servicios especializados cuando el tipo de carga justifique otro lenguaje o
ciclo de vida.

Se evitará comenzar con microservicios independientes porque aumentarían la
complejidad de despliegue, observabilidad, actualizaciones y copias de seguridad
sin aportar valor inicial. Los módulos tendrán límites explícitos para poder
extraer servicios cuando exista una necesidad demostrada.

### AD-003. Contenedores como unidad de despliegue

- Desarrollo integrado y despliegue inicial mediante Kubernetes de Docker
  Desktop.
- Helm como unidad de empaquetado y configuración.
- Compose limitado a pruebas aisladas de un servicio cuando resulte útil.
- Imágenes sin estado para aplicación, trabajadores y servicios auxiliares.
- Datos persistentes en volúmenes claramente identificados.
- Configuración mediante archivos y variables de entorno.
- Comprobaciones de salud y dependencias declaradas.
- Versiones de imagen fijadas, no etiquetas flotantes en producción.
- Arquitectura compatible con `amd64` y `arm64` cuando las dependencias lo
  permitan.

### AD-004. Publicación por proyección

La parte pública y compartida se generará mediante una **proyección de
publicación** separada del modelo privado.

Una proyección contiene únicamente:

- campos seleccionados;
- recursos expresamente autorizados;
- imágenes ya procesadas;
- redacciones de datos sensibles;
- permisos, caducidad y versión;
- información necesaria para renderizar la vista.

El servicio público nunca recibirá una copia completa de la base de datos
privada.

### AD-005. Salidas iniciadas desde el núcleo

En el modelo híbrido, la instancia local establecerá una conexión saliente
cifrada con el servicio de publicación. No se abrirán puertos entrantes en la
red doméstica por defecto.

Esto permite:

- publicar páginas;
- actualizar vistas compartidas;
- recibir comentarios y documentos autorizados;
- revocar accesos;
- funcionar detrás de NAT sin exponer el núcleo.

### AD-006. Arquitectura políglota por tipo de carga

No se impondrá un único lenguaje a toda la plataforma. Cada unidad desplegable
utilizará la tecnología que mejor se adapte a su carga:

- Go para dominio transaccional, dispositivos, publicación, MCP y servicios de
  red;
- Python para visión artificial, aprendizaje automático, RAG y cálculo
  científico;
- TypeScript para Angular;
- SQL y PostGIS para consultas transaccionales, temporales y espaciales;
- Rust o C++ solo si una necesidad medida no queda razonablemente cubierta por
  Go, Python o librerías existentes.

La elección de lenguaje no será por sí sola motivo suficiente para crear un
servicio. Deberá existir al menos una frontera de dominio, seguridad, hardware,
escalado, dependencias o ciclo de vida.

### AD-007. Screaming Architecture

La estructura del repositorio y del código expresará primero el negocio:

- sistemas acuáticos;
- topología y espacio;
- habitantes y bienestar;
- bioseguridad y salud;
- mediciones;
- operación y bitácora;
- nutrición y dosificación;
- equipos;
- alertas;
- compartición y publicación.

No se organizará la aplicación principalmente en carpetas globales como
`controllers`, `services`, `models`, `repositories` o `utils`. Esos conceptos
podrán existir dentro de un dominio, pero no ocultarán el propósito del sistema.

### AD-008. Monolito distribuible

Cada módulo del núcleo será una unidad lógica extraíble:

- propietario de sus reglas e invariantes;
- propietario de su esquema lógico de datos;
- sin llamadas directas desde otros dominios;
- emisor y consumidor de eventos versionados;
- sin acceso directo desde otros módulos a sus tablas internas;
- sin ciclos de dependencias;
- con pruebas de arquitectura automatizadas.

Inicialmente varios dominios podrán compilarse y desplegarse juntos por
comodidad operativa, pero se comunicarán a través del mismo bus y los mismos
contratos de eventos que usarían estando separados. Compartir proceso no
concederá acceso directo al código o los datos de otro dominio.

### AD-008A. Integración exclusivamente mediante eventos

Entre dominios:

- no habrá imports de paquetes de negocio ajenos;
- no habrá llamadas síncronas;
- no habrá consultas a APIs internas;
- no habrá repositorios compartidos;
- no habrá joins ni claves foráneas entre esquemas;
- no habrá transacciones distribuidas;
- no se compartirán entidades de dominio.

Un dominio publicará hechos que ya han ocurrido y reaccionará a hechos
publicados por otros. Ejemplo:

```text
Identity registra un usuario
        │
        └──▶ identity.user-created.v1
                     │
                     ├──▶ Notifications prepara preferencias
                     ├──▶ Sharing crea su proyección
                     └──▶ Audit conserva el hecho
```

Identity no conoce esos consumidores ni modifica su comportamiento cuando se
añade uno nuevo.

### AD-009. Datos según naturaleza y escala

- PostgreSQL será la fuente transaccional.
- PostGIS gestionará geometrías y relaciones espaciales.
- Series temporales comenzarán en PostgreSQL particionado.
- pgvector almacenará embeddings iniciales.
- SeaweedFS Community, mediante su API S3 compatible, almacenará inicialmente
  los objetos.
- SQLite podrá utilizarse como buffer en gateways desconectados.
- TimescaleDB y ClickHouse solo se incorporarán al superar criterios medibles.

No se introducirá una base diferente únicamente porque un módulo pueda
beneficiarse marginalmente de ella.

### AD-010. Trazabilidad completa por diseño

Toda operación relevante será rastreable desde su origen hasta sus efectos:

```text
actor/dispositivo/agente
        │
        ▼
solicitud o comando
        │
        ▼
decisión del dominio + transacción
        │
        ▼
evento de integración
        │
        ├──▶ consumidor ──▶ proyección
        ├──▶ consumidor ──▶ alerta
        └──▶ consumidor ──▶ nuevo evento
```

Cada eslabón conservará identidad, tiempo, correlación, causación, versión,
resultado y procedencia. La trazabilidad funcional será independiente de los
logs técnicos y no dependerá de que estos se conserven indefinidamente.

Trazabilidad completa no implica event sourcing obligatorio. Cada dominio
decidirá su persistencia interna, pero deberá conservar un historial auditable
de cambios y publicar sus hechos de integración mediante outbox.

### AD-011. Kubernetes como plataforma desde el inicio

El entorno inicial será Kubernetes de Docker Desktop sobre un Mac mini M4. Se
utilizará el provisionador `kind` de Docker Desktop porque permite seleccionar
versión de Kubernetes, trabajar con varios nodos virtuales y migrar manifiestos
estándar con menos particularidades del provisionador.

La aplicación no dependerá de APIs exclusivas de Docker Desktop. Todo recurso
deberá poder instalarse posteriormente en otro Kubernetes conforme mediante
Helm y recursos estándar.

Docker Desktop se considerará una instalación inicial de un único host, no un
entorno de alta disponibilidad. La pérdida o reinicio del host afecta a todo el
clúster.

### AD-012. Malla de servicio mínima con Linkerd

Linkerd proporcionará:

- identidad de workload basada en ServiceAccount;
- mTLS entre pods incorporados a la malla;
- políticas de comunicación entre workloads;
- métricas de tráfico;
- timeouts y retries únicamente donde sean seguros.

No se instalarán extensiones de visualización o multiclúster que no sean
necesarias. Prometheus y Grafana podrán consumir métricas sin obligar a mantener
todo el paquete de visualización de Linkerd.

`cert-manager` gestionará certificados de entrada y la rotación del issuer de
Linkerd. El trust anchor tendrá backup, monitorización de caducidad y un
procedimiento de rotación probado.
Sus charts OCI se fijan por digest y seleccionan explícitamente la capa
`application/vnd.cncf.helm.chart.content.v1.tar+gzip` con operación `copy`;
no se acepta la selección implícita de la primera capa del artefacto.

La malla no sustituye:

- autenticación de usuarios;
- autorización sobre recursos del negocio;
- Gateway/API Management;
- trazabilidad funcional;
- idempotencia de eventos.

Los retries de la malla estarán desactivados para operaciones no idempotentes y
no se aplicarán ciegamente sobre publicación de comandos.

### AD-013. Envoy Gateway como APIM y punto de aplicación de políticas

Envoy Gateway constituirá la entrada norte-sur y el Policy Enforcement Point.

- Configuración declarativa y GitOps.
- Kubernetes Gateway API como interfaz principal.
- `HTTPRoute` y `GRPCRoute` cuando proceda.
- Terminación TLS.
- Routing y versionado externo.
- Validación de tokens.
- CORS, límites de tamaño y rate limiting.
- Correlation/request ID.
- Métricas y logs de acceso.
- Protección de APIs públicas, compartidas y administrativas.
- Delegación de AuthN en un proveedor OIDC.
- Delegación de AuthZ en un Policy Decision Point externo.

Se utilizarán Gateway API y las políticas declarativas de Envoy Gateway, sin
base de datos propia. Su superficie administrativa no será pública.

El APIM no almacenará ni implementará las relaciones detalladas de acuarios,
organismos o vistas. Consultará el servicio gestionado de AuthZ.

Su despliegue tendrá dos puertas independientes. La primera instala mediante
GitOps las CRD y el controlador, pero no declara `Gateway`, listeners, rutas ni
plano de datos. La segunda crea la entrada únicamente con TLS, AuthN y AuthZ
preparadas y verificadas. Instalar el controlador no constituye autorización
para exponer una interfaz. El contrato completo está en
[Entrada norte-sur y acceso](entrada-y-acceso.md).

### AD-013A. AuthN y AuthZ delegadas en el APIM

Para cada petición protegida, el APIM:

1. valida firma, emisor, audiencia, expiración y scopes del token;
2. extrae la identidad autenticada;
3. determina acción y tipo de recurso desde la ruta y su metadata;
4. consulta un servicio de autorización externo;
5. deniega o enruta;
6. acepta del ReefOps Authorizer únicamente un `ActorContext` firmado y
   cabeceras incluidas en una allowlist;
7. registra `authz_decision_id`, sujeto, acción, recurso y resultado.

```text
Cliente
   │ bearer token
   ▼
APIM ── OIDC/JWKS ──▶ Identity Provider
   │
   ├── ext_authz ───▶ ReefOps Authorizer ──▶ AuthZ PDP
   │
   ▼
Dominio recibe ActorContext firmado
```

Los dominios no validarán tokens OIDC ni consultarán directamente el PDP. El
adaptador verificará únicamente la autenticidad, audiencia y antigüedad del
`ActorContext` interno.

El APIM eliminará cualquier cabecera de identidad proporcionada por el cliente.
ReefOps Authorizer generará y firmará el contexto interno después de recibir la
decisión de OpenFGA. Envoy solo propagará la respuesta allowlisted del
Authorizer. El contexto se transmitirá mediante un JWT de vida muy corta o una
firma equivalente, no mediante cabeceras confiadas sin protección.

### AD-014. Autenticación gestionada con ZITADEL

ZITADEL será responsable de:

- OpenID Connect y OAuth 2.x;
- usuarios humanos;
- organizaciones;
- clientes web;
- service accounts;
- MFA y passkeys cuando se habiliten;
- sesiones, recuperación y políticas de acceso;
- emisión y rotación de tokens.

ReefOps no almacenará contraseñas. Conservará únicamente el identificador
externo necesario y su proyección funcional.

Identity en ReefOps no será propietario de credenciales. Consumirá eventos de
provisión y mantendrá la relación entre sujeto autenticado y actor de negocio.

### AD-015. Autorización gestionada con OpenFGA

OpenFGA resolverá la autorización fina:

- propietario de instalación;
- miembro o cuidador;
- veterinario invitado;
- soporte de tienda;
- enlace compartido;
- agente que actúa en nombre de una persona;
- permisos sobre sistemas, organismos, documentos y operaciones.

Los modelos de autorización serán inmutables y versionados. Las relaciones o
hechos de autorización se actualizarán mediante consumidores de eventos, no
mediante escrituras ocultas desde cualquier dominio.

El PDP responde si una acción está permitida; el dominio sigue siendo
responsable de sus invariantes. Una decisión permitida no obliga a aceptar un
comando inválido.

### AD-016. Plataforma no equivale a dominio

AuthN, AuthZ, APIM, mesh, bus, almacenamiento y observabilidad son servicios de
plataforma. No participan como dominios de negocio.

La regla de integración exclusivamente mediante eventos se mantiene entre
dominios. Se permiten dependencias técnicas controladas:

- el APIM valida la identidad con el IdP;
- el APIM consulta el autorizador/PDP antes de enrutar;
- el dominio recibe un `ActorContext` ya autenticado y autorizado;
- el dominio no importa SDK del IdP, PDP, APIM o mesh.

Las consultas de autorización son síncronas, fallan de forma cerrada y quedan
auditadas en el APIM y el autorizador. Esta excepción no permite que un dominio
consulte datos funcionales de otro.

### AD-017. Operación declarativa, reproducible y verificable

Helm será el formato canónico de empaquetado y Flux CD reconciliará el estado
del clúster desde Git desde la primera fase. `Taskfile` ofrecerá comandos
repetibles para bootstrap, validación y operación; Helmfile no será una segunda
fuente de verdad.

OpenBao será la autoridad local de secretos runtime. Los workloads se
autenticarán con identidades de Kubernetes y consumirán secretos sin acoplar el
dominio al gestor. SOPS y `age` se limitarán al bootstrap cifrado; nunca habrá
secretos en claro en Git, charts, imágenes ni copias de configuración. Las
claves de recuperación se guardarán fuera del Mac mini y tendrán un
procedimiento de rotación probado.

Cada despliegue tendrá un `deployment_id` y conservará commit, digest inmutable
del chart OCI y de las imágenes, migraciones, operador, instante y resultado. Las
operaciones manuales de emergencia deberán dejar la misma evidencia.

### AD-018. IA local desacoplada del runtime de Kubernetes

Los pods de Docker Desktop se ejecutan en una VM Linux y no se asumirá acceso
directo a Metal ni al Apple Neural Engine. Intelligence consumirá una API de
inferencia intercambiable, pero dicha API pertenecerá siempre a la instalación
local.

En el Mac mini, los LLM locales podrán ejecutarse en macOS mediante un runtime
con aceleración Metal —inicialmente Docker Model Runner o equivalente—. Visión
comenzará con CPU si satisface sus SLO; de no hacerlo podrá ejecutarse en otro
nodo GPU de la red privada administrada por el propietario, sin cambiar
contratos de dominio.

Todo runtime declarará capacidades, versión de modelo, digest de pesos,
dispositivo de ejecución y política de tratamiento de datos.

No se admitirán proveedores de inferencia cloud ni APIs de modelos de terceros.
Imágenes, vídeo, embeddings, prompts, telemetría, documentos y contexto del
acuario no saldrán de la instalación para ejecutar IA. La descarga de modelos y
actualizaciones será una operación separada, verificable y no incluirá datos
del usuario.

### AD-019. Equipo vital independiente del clúster

Kubernetes podrá supervisar y solicitar acciones, pero la seguridad básica de
bombas, calentadores, refrigeración, oxigenación y dosificación no dependerá de
que el Mac mini, la red o el clúster estén disponibles. Los controladores
tendrán límites físicos, programas seguros y modo degradado autónomo.

El host contará, cuando el riesgo lo justifique, con SAI, apagado ordenado,
watchdog, monitorización externa y un canal de alerta independiente.

### AD-020. Documentación antes de implementación

Todo cambio seguirá este orden:

1. documentar requisito, decisión, alcance, riesgos y criterio de aceptación;
2. revisar su alineación con dominios y plataforma;
3. implementar IaC, aplicación o migración;
4. validar el resultado y actualizar la evidencia.

No se instalará ni desarrollará un componente cuya función, propietario,
dependencias y recuperación no estén documentados. Una corrección urgente podrá
aplicarse para proteger datos o animales, pero deberá registrarse y
reconciliarse inmediatamente con documentación y Git.

### AD-021. GitOps desde el bootstrap

Git será la fuente de verdad del estado deseado. El único procedimiento
imperativo normal será un bootstrap idempotente que:

- comprueba host, herramientas, contexto y versión de Kubernetes;
- instala los controladores de Flux;
- configura la fuente Git y la reconciliación raíz;
- instala la identidad de descifrado de SOPS sin almacenarla en Git;
- espera y verifica el estado resultante.

Después del bootstrap, namespaces, CRD, charts, políticas, configuración y
aplicaciones se modificarán mediante commits. No se utilizarán `kubectl apply`
o `helm install` manuales salvo diagnóstico o emergencia auditada.

Un revert de Git revierte estado declarativo, no datos ni efectos consumados.
PostgreSQL, objetos, volúmenes, migraciones, claves rotadas y órdenes físicas
requieren sus propios procedimientos de backup, restauración o compensación.

### AD-022. Gobierno de ayudas de desarrollo

`AGENTS.md`, skills, agentes especializados, configuración Codex, MCP y hooks
son controles auxiliares versionados. Su jerarquía será:

1. documentos de requisitos y arquitectura como fuente de verdad;
2. `AGENTS.md` como reglas duraderas por alcance;
3. skills como procedimientos reutilizables;
4. agentes especializados como revisión de solo lectura;
5. MCP únicamente para sistemas externos necesarios;
6. scripts, tests y CI como enforcement verificable.

Las ayudas no introducirán nuevas decisiones por sí mismas ni duplicarán
normativa canónica extensa. Una modificación de su comportamiento deberá
actualizar primero [Desarrollo asistido](desarrollo-asistido.md) y su matriz.
El agente principal conserva la integración, la autoridad y la responsabilidad
de validar. No habrá agentes de escritura especializados en la fase inicial.

### AD-023. GitHub como plataforma de ingeniería

GitHub será la plataforma preferente para Git, colaboración, CI, seguridad de
dependencias, releases y artefactos. La propiedad se alojará en una
organización y se separará en:

- `reefops`: producto, documentación, contratos y construcción, público desde
  un historial raíz sanitizado;
- `reefops-platform`: componentes Kubernetes reutilizables y bootstrap,
  público;
- `reefops-gitops`: composición privada, versiones y digests desplegados.

La publicación exige escaneo del árbol y del historial y excluye datos de
usuario, secretos y detalles de instalaciones. Producto y plataforma cumplen
esa condición. GitOps conserva privada la topología y la configuración
específica aunque los secretos estén cifrados.

GitHub Actions construirá y publicará; no tendrá acceso directo al Kubernetes
local. Un cambio de versión abrirá un PR en GitOps y Flux reconciliará desde Git
con credenciales de solo lectura. GHCR almacenará imágenes y charts OCI. Los
detalles de seguridad, trazabilidad y gobierno se definen en
[Plataforma GitHub](plataforma-github.md).

Se utilizará GitHub Flow sobre `main`, ramas cortas y pull requests. No se
adopta GitFlow ni se representarán entornos mediante ramas. El código publicado
usará Apache License 2.0; datos, despliegues privados, secretos, modelos y
material de terceros conservarán su régimen propio.

GitHub Secrets podrá ser una réplica revocable para CI, provisionada desde
OpenBao por una conexión saliente local. No será autoridad de secretos ni
contendrá credenciales runtime o datos privados.

### AD-024. Development local y capacidad multi-entorno

La primera etapa tendrá un único entorno operativo, `development`, dentro del
Kubernetes `docker-desktop` del Mac Mini. No se reservarán namespaces ni se
desplegarán autoridades, datos o workloads de `production` en este clúster
mientras no exista una necesidad operativa.

La arquitectura seguirá siendo multi-entorno. Cada entorno futuro tendrá una
composición GitOps y, según su destino, namespaces o clúster propios, además de
ServiceAccounts, NetworkPolicies, quotas, OpenBao, PKI, bases de datos,
almacenamiento de objetos, NATS/MQTT, stores de autorización, proyectos de
identidad, auditoría y backups independientes. Subjects, topics, buckets, URLs
y recursos llevarán identidad de entorno. Un evento o evidencia funcional
conservará `environment_id` además de correlación y causación para impedir
mezclar cadenas entre entornos.

El OpenBao actualmente operativo pertenece a `development`. No se renombra ni
se mueve de forma imperativa: su futura adopción por una raíz explícita de
development será una migración GitOps que preserve PVC, identidad Raft, PKI y
evidencia. Un futuro `production` recibirá una instalación y una ceremonia
independientes; nunca se clonarán material Shamir, token inicial, identidad
`age`, CA, snapshot o secretos desde development.

Los componentes de plataforma serán reutilizables y no codificarán el destino
local. La composición privada decidirá qué entornos existen y en qué clúster.
La separación podrá ser lógica dentro de un clúster o física entre clústeres,
pero no se dará por existente hasta materializar y verificar sus controles.

Los entornos no se representarán mediante ramas. CI construirá una sola vez y
development consumirá el artefacto inmutable por digest. Cuando exista
production, su promoción será otro PR GitOps que seleccione exactamente ese
digest después de las verificaciones; no se reconstruirá el artefacto.

### AD-025. Entrega declarativa de secretos con ESO

External Secrets Operator —ESO— será el adaptador inicial entre Kubernetes y
OpenBao. OpenBao seguirá siendo la autoridad; ESO solo materializará en un
Secret Kubernetes el subconjunto que un operador no pueda consumir como
volumen efímero. Los dominios no importarán SDK de ESO u OpenBao.

La integración inicial residirá en `reefops-secret-delivery`, aislada de las
claves privadas de OpenBao, utilizará un
`SecretStore` namespaced, TLS validado contra la CA de OpenBao y autenticación
Kubernetes con una ServiceAccount exclusiva sin token permanente. La política
OpenBao se limitará a una ruta sintética. No habrá token estático,
`ClusterSecretStore` global, comodines de namespace ni `skipVerify`.

Cada entorno futuro tendrá stores, identidades, políticas, CA y rutas propias.
ESO se desplegará por frontera de confianza y namespace, no por producto: un
controller scoped de `reefops-data` podrá reconciliar los `ExternalSecret` de
SeaweedFS, PostgreSQL y futuros servicios de datos sin que ninguno dependa de
la raíz GitOps de otro. Cada consumidor conservará ServiceAccount, TokenRequest,
`SecretStore` y política OpenBao propios. No se ampliará el controller sintético
a todo el clúster ni se introducirá un controller por cada componente.

Antes de entregar secretos a otro namespace se distribuirá únicamente el trust
root público y se creará una identidad por consumidor. El ciclo de aceptación
comprobará autenticación, lectura, refresco, revocación y auditoría sin revelar
el valor. La retirada tendrá un procedimiento separado que preserve o elimine
cada Secret destino según su contrato.

### AD-026. Home Assistant como frontera IoT preferente

MQTT no forma parte del camino crítico inicial. ReefOps integrará
preferentemente enchufes, sensores y automatizaciones mediante un adaptador
local de Home Assistant. Home Assistant será propietario del descubrimiento,
protocolos y entidad física; ReefOps conservará la intención acuícola, reglas
de seguridad, autorización, trazabilidad y valoración del resultado.

Una aceptación de la API de Home Assistant no demostrará por sí sola que la
acción física ocurrió. El adaptador correlacionará orden, respuesta y estado
posteriormente observado, y será idempotente para evitar repetir dosificación,
alimentación o conmutaciones. MQTT podrá existir detrás de Home Assistant o
adoptarse más adelante si aparece un caso no cubierto, sin entrar en los
dominios.

## 3. Vista general

```text
                             INTERNET

              ┌──────────────────────────────────┐
              │ Servicio opcional de publicación│
              │                                  │
 Visitantes ─▶│ Público · Compartido · Comentarios│
              └───────────────▲──────────────────┘
                              │ conexión saliente
                              │ cifrada y autenticada
┌─────────────────────────────┴───────────────────────────────┐
│                    INSTALACIÓN REEFOPS                       │
│                                                             │
│  Navegador / PWA ──▶ ReefOps Core (Go) ──▶ PostgreSQL       │
│                         │             Almacenamiento objetos │
│                         │                                   │
│                         ├── Trabajadores y planificador      │
│                         ├── Adaptador Home Assistant (Go)    │
│                         ├── Motor de alertas                 │
│                         ├── MCP Server (Go)                  │
│                         ├── Intelligence (Python)            │
│                         └── Vision (Go + Python)             │
│                                                             │
│  Cámaras · Sensores · Controladores · Equipos                │
└─────────────────────────────────────────────────────────────┘
```

En un despliegue cloud o de servidor privado, estos mismos componentes podrán
ejecutarse dentro de una única infraestructura y el módulo de publicación podrá
conectarse internamente.

## 4. Componentes

### 4.1 Interfaz web

Aplicación web adaptable y PWA:

- escritorio para modelado espacial, análisis e informes;
- móvil para registro rápido, tareas, voz y fotografías;
- tableta para mantenimiento y edición espacial;
- caché local de operaciones básicas;
- cola local de cambios cuando no exista conexión con el servidor;
- comunicación exclusivamente mediante la API.

No contendrá reglas de negocio que no estén también protegidas en el servidor.

### 4.2 API y aplicación principal

El núcleo se implementará en Go y será responsable de:

- autenticación y autorización;
- sistemas, componentes y topología;
- habitantes, bienestar y bioseguridad;
- parámetros y mediciones;
- tareas, bitácora e incidencias;
- inventario, aditivos y dosificación;
- modelos espaciales y máscaras;
- vistas privadas, compartidas y públicas;
- coordinación de agentes e integraciones.

La aplicación se dividirá internamente en módulos de dominio. Un módulo no
accederá directamente a las tablas privadas de otro sin una interfaz definida.

Stack inicial:

- Go 1.26 o última revisión compatible fijada por el proyecto;
- `net/http` y un router pequeño como `chi`;
- `pgx` para PostgreSQL;
- `sqlc` para generar acceso tipado desde SQL;
- OpenAPI para el contrato externo;
- `oapi-codegen` para servidor, modelos y cliente TypeScript;
- `slog` y OpenTelemetry;
- Atlas o Goose para migraciones.

No se utilizará un ORM generalista como centro del modelo. Las reglas residirán
en el dominio y las consultas complejas se expresarán explícitamente.

### 4.3 Base de datos

PostgreSQL será la fuente de verdad transaccional.

Motivos:

- relaciones y restricciones sólidas;
- datos temporales y espaciales;
- JSON cuando se necesite flexibilidad;
- índices y extensiones maduras;
- ejecución local o como servicio administrado;
- herramientas fiables de copia y restauración.

Se podrá utilizar PostGIS para geometría, máscaras, regiones y consultas
espaciales. Las series temporales comenzarán en PostgreSQL particionado. Solo se
añadirá una extensión o base especializada si el volumen real de sensores lo
justifica.

Los módulos del Core podrán compartir un clúster físico, pero cada dominio
tendrá base, owner, rol de migración y rol runtime propios. No existirán grants,
vistas, funciones, claves foráneas ni joins entre bases de dominios. Una
migración pertenecerá a un único módulo. Las necesidades de lectura cruzada se
resolverán mediante proyecciones alimentadas por eventos, nunca mediante
consultas a datos ajenos.

Development usará CloudNativePG con una sola instancia PostgreSQL: otra réplica
sobre el mismo Mac no se presentará como HA. Backup físico y WAL se gestionarán
mediante el plugin CNPG-I Barman Cloud hacia SeaweedFS y se exportarán cifrados
fuera de la VM. El contrato está en
[PostgreSQL y CloudNativePG](postgresql.md).

### 4.4 Almacenamiento de objetos

Fotografías, vídeo, documentos, modelos y evidencias se guardarán fuera de la
base relacional mediante una interfaz compatible con objetos:

- SeaweedFS Community en el despliegue local;
- S3 compatible en servidor o nube;
- referencias, hash, propietario y permisos en PostgreSQL.

SeaweedFS se adopta por su licencia Apache-2.0, soporte de Kubernetes y ARM64,
API S3 y capacidad de evolucionar desde una instalación ligera hacia una
topología distribuida. MinIO no se adoptará para instalaciones nuevas: su
repositorio comunitario está archivado, la distribución comunitaria mantenida
ha dejado de ofrecer binarios oficiales y su modelo actual diferencia el
código AGPLv3 de la oferta comercial.

Los dominios y trabajadores dependerán de un `ObjectStoragePort`, no de APIs,
tipos o SDK propios de SeaweedFS. El adaptador inicial utilizará un SDK S3 y
recibirá endpoint, región, credenciales y opciones por configuración. La
sustitución por Ceph RGW, un servicio S3 administrado u otro proveedor
compatible no modificará el modelo de dominio ni los contratos de eventos.

La compatibilidad se verificará con pruebas contractuales para el subconjunto
S3 usado por ReefOps: `PUT`, `GET`, `HEAD`, `DELETE`, multipart, rangos, URLs
prefirmadas, checksums y semántica de `ETag`. Versionado, lifecycle,
notificaciones, retención y Object Lock no se usarán hasta tener una prueba
contractual específica en cada proveedor soportado.

Los objetos privados, compartidos y públicos utilizarán espacios lógicos
separados. Una imagen pública será una derivación procesada, no un enlace al
original privado.

PostgreSQL conservará el identificador lógico, clave de objeto, bucket o
espacio, propietario, clasificación, hash, tamaño, tipo de medio, versión,
estado, retención y procedencia. Cada escritura, lectura sensible, publicación,
derivación y eliminación mantendrá `correlation_id`, `causation_id`, actor o
servicio, decisión de autorización y resultado. Las URLs prefirmadas serán de
corta duración y solo se emitirán después de autorizar la acción y el recurso.

SeaweedFS en un único Mac mini aporta portabilidad y una interfaz S3, pero no
alta disponibilidad. Sus volúmenes no constituyen un backup y deberán incluirse
en copias cifradas, verificadas y almacenadas fuera del clúster y de la VM de
Docker Desktop.

El contrato operativo, el perfil de un nodo, la entrega de credenciales, las
capacidades S3 admitidas y la recuperación se detallan en
[Almacenamiento de objetos](almacenamiento-objetos.md). En development se usará
una sola réplica de cada rol y colocación `000`: crear copias en el mismo Mac no
reduce el dominio de fallo. Master, filer y volume tendrán persistencia
independiente, acceso exclusivamente interno y backup lógico cifrado con
restauración aislada.

### 4.5 Trabajadores en segundo plano

Los trabajos de negocio escritos en Go podrán ejecutarse desde el mismo binario
con un comando o rol distinto. Procesos separados ejecutarán:

- miniaturas y transformación de medios;
- importaciones y exportaciones;
- informes;
- análisis de imágenes y vídeo;
- cálculos espaciales pesados;
- notificaciones;
- publicación y sincronización;
- resúmenes y tareas del agente.

La API creará trabajos duraderos. Los trabajadores podrán reintentarlos sin
duplicar resultados.

Los trabajos de inferencia y cálculo científico serán consumidos por servicios
Python. Los objetos grandes se intercambiarán mediante referencias S3,
no dentro del mensaje del bus.

### 4.6 Planificador

Responsable de:

- tareas recurrentes;
- análisis visual cada intervalo configurable;
- comprobación de sensores ausentes;
- resúmenes;
- caducidad de enlaces;
- mantenimiento y copias de seguridad;
- reglas temporales.

Los trabajos planificados se persistirán en la base de datos para sobrevivir a
reinicios.

### 4.7 Bus de eventos e integraciones

La primera versión utilizará eventos internos persistentes mediante patrón
outbox. NATS JetStream será el bus entre unidades desplegables. La integración
física se delegará inicialmente en Home Assistant mediante un adaptador local;
MQTT queda como detalle opcional detrás de esa frontera.

Ejemplos de eventos:

- `measurement.recorded`;
- `parameter.threshold_crossed`;
- `dose.confirmed`;
- `organism.moved`;
- `vision.anomaly_detected`;
- `incident.opened`;
- `share.revoked`.

No se introducirá una plataforma de streaming más pesada hasta que la escala lo
requiera. Kafka no formará parte del despliegue local inicial.

Los eventos usarán esquemas versionados y consumidores idempotentes. Cada
mensaje contendrá identificador, tipo, versión, instante, productor,
correlación, causación y ámbito del propietario.

### 4.8 Motor de reglas y alertas

- Evalúa umbrales, persistencia, tendencias y combinación de señales.
- Separa observaciones, advertencias y alarmas críticas.
- Deduplica eventos.
- Mantiene estado, confirmación y escalado.
- Puede ejecutar acciones locales previamente autorizadas.
- No dependerá del agente generativo para acciones de seguridad.

Las reglas críticas serán deterministas y auditables.

Cada evaluación conservará regla y versión, datos de entrada, resultado,
umbrales, acciones generadas y relación con la alerta o incidencia resultante.

### 4.9 IA y visión

Se separarán tres funciones:

1. **Orquestación:** permisos, contexto, herramientas y auditoría dentro de la
   aplicación principal.
2. **Inferencia:** modelos ejecutados dentro de la instalación.
3. **Procesamiento visual:** captura, segmentación, seguimiento y detección.

La aplicación podrá seleccionar por función:

- runtime en el host local;
- servicio en Kubernetes;
- nodo de cómputo de la red privada;
- función desactivada.

Cada resultado conservará modelo, versión, entrada, configuración, confianza y
revisión humana. La falta de IA no impedirá el uso manual del producto.

Los adaptadores abstraerán runtimes locales, no proveedores cloud. Se
verificará que los contenedores y procesos de inferencia no tengan salida a
Internet durante el tratamiento de datos.

La captura y preparación de medios podrá implementarse en Go sobre
FFmpeg/GStreamer. La inferencia, segmentación, entrenamiento y análisis
científico se implementarán en Python. Vision no escribirá en tablas de otros
dominios ni les devolverá resultados mediante llamadas síncronas. Publicará
observaciones versionadas; el dominio propietario decidirá, al consumirlas, si
crea una propuesta, alerta o versión pendiente de revisión.

Cada inferencia mantendrá el linaje desde los objetos originales hasta el
resultado, incluyendo transformaciones, modelo, versión, prompt o configuración,
herramientas, fuentes recuperadas, confianza y revisión humana.

### 4.10 Servidor MCP

El servidor MCP se implementará en Go y utilizará la misma capa de autorización
y casos de uso que la API. No accederá directamente a la base de datos.

- Desactivado por defecto.
- Escucha local por defecto.
- Herramientas de lectura como configuración inicial.
- Escrituras con permisos y confirmación.
- Tokens con alcance y caducidad.
- Auditoría completa.

Podrá ejecutarse dentro del proceso principal al inicio y separarse más
adelante sin cambiar sus contratos.

### 4.11 Device Gateway

Servicio Go responsable de la frontera física:

- API y eventos locales de Home Assistant;
- MQTT o protocolos de fabricantes solo si un caso futuro lo exige;
- descubrimiento local;
- normalización de telemetría;
- estado y diagnóstico de dispositivos;
- buffer local con SQLite durante desconexiones;
- entrega idempotente al Core;
- ejecución local de reglas de seguridad preautorizadas.

El gateway no utilizará NATS como sustituto de protocolos de dispositivo ni la
API de Home Assistant o MQTT como sustitutos del bus de dominio.

### 4.12 Intelligence Service

Servicio Python responsable de:

- recuperación de conocimiento;
- agentes especializados;
- RAG;
- evaluación de compatibilidad;
- generación de explicaciones y planes;
- adaptadores para runtimes y modelos locales.

No tendrá acceso general a PostgreSQL. Obtendrá contexto autorizado mediante
proyecciones propias alimentadas exclusivamente por eventos sanitizados. No
consultará APIs ni bases de otros dominios. Sus recomendaciones serán hechos
propios; cualquier cambio funcional posterior será iniciado por el usuario o
agente autorizado como un nuevo comando de entrada al dominio propietario y
conservará la causación. Intelligence no enviará comandos a otros dominios.

### 4.13 Publication Service

Servicio Go desplegable externamente:

- páginas públicas;
- enlaces compartidos;
- instantáneas y vistas vivas;
- comentarios y documentos de retorno;
- revocación, caducidad y límites de acceso.

Solo almacenará proyecciones sanitizadas. No controlará equipos ni consultará la
base privada.

### 4.14 ReefOps Authorizer

Servicio de plataforma pequeño, escrito en Go, compatible con la interfaz de
autorización externa del APIM.

Responsabilidades:

- traducir ruta, método y metadata a `acción + recurso`;
- validar que los identificadores de recurso tengan el formato esperado;
- consultar el PDP seleccionado;
- devolver allow/deny y un `decision_id`;
- emitir auditoría y métricas;
- aplicar timeouts y denegación por defecto;
- construir el `ActorContext` interno firmado;
- ocultar al APIM los detalles específicos de la API y el modelo de OpenFGA.

No contendrá reglas de acuariofilia ni consultará bases de datos de dominio.
Recibirá su modelo de relaciones a través del PDP, alimentado por eventos.

Este adaptador evita acoplar las rutas del APIM a un proveedor concreto de
AuthZ.

## 5. Privado, compartido y público

### 5.1 Privado

- Reside en la instancia principal.
- Es la fuente de verdad.
- Incluye originales, notas internas, controles y credenciales.
- Nunca se consulta desde el servicio público.

### 5.2 Público

La instancia crea una proyección inmutable y versionada, la revisa y la publica.
El servicio externo puede servirla incluso si ReefOps local está apagado.

Adecuado para:

- perfil público del acuario;
- galería;
- habitantes seleccionados;
- bitácora pública;
- parámetros deliberadamente publicados.

### 5.3 Compartido

Se soportarán dos modalidades:

**Instantánea**

- Información cerrada en una fecha.
- Continúa disponible aunque el núcleo esté desconectado.
- Ideal para informes y consultas puntuales.

**Vista viva**

- La instancia local envía actualizaciones de la proyección.
- Puede recibir comentarios o documentos.
- Si el núcleo está desconectado, muestra el último estado sincronizado y su
  fecha.
- Nunca permite control directo de equipos desde el servicio público.

### 5.4 Flujo de publicación

```text
Datos privados
    │
    ▼
Selector de campos y recursos
    │
    ▼
Redacción y eliminación de metadatos
    │
    ▼
Vista previa y confirmación
    │
    ▼
Proyección firmada y versionada
    │
    ▼
Servicio de publicación
```

La revocación elimina el acceso en el servicio de publicación. Se conservará
localmente una auditoría de lo que estuvo disponible.

## 6. Perfiles de despliegue

### 6.1 Desarrollo local

Kubernetes integrado de Docker Desktop sobre Apple Silicon, usando el contexto
`docker-desktop`:

- charts Helm de ReefOps;
- dependencias gestionadas mediante charts fijados;
- Angular y servicios de dominio;
- PostgreSQL;
- SeaweedFS Community;
- NATS JetStream;
- adaptador local de Home Assistant;
- Envoy Gateway;
- Linkerd;
- ZITADEL;
- OpenFGA;
- servicio de IA local opcional;
- capturador o simulador de dispositivos.

Compose podrá existir para pruebas unitarias de un servicio, pero no será la
definición canónica de la plataforma. Los manifiestos Kubernetes y Helm serán la
fuente de verdad del despliegue.

Las imágenes se construirán para `linux/arm64`. CI construirá también
`linux/amd64` para asegurar la migración futura.

### 6.2 Instalación local sencilla

- Un único host.
- Kubernetes de Docker Desktop inicialmente.
- Charts Helm y valores versionados.
- Envoy Gateway y HTTPS local.
- Descubrimiento en red opcional.
- Actualización mediante Helm con copia previa.
- Copia de seguridad exportable a otro dispositivo.
- Persistencia respaldada fuera del ciclo de vida del clúster.

Un reset o recreación de Kubernetes de Docker Desktop puede eliminar recursos
del clúster. Las copias de PostgreSQL, SeaweedFS, ZITADEL y OpenFGA deberán salir
periódicamente de la VM de Docker Desktop.

### 6.3 Servidor local o profesional

- Componentes separados en uno o varios nodos.
- Base de datos y objetos en almacenamiento resistente.
- Integración con identidad corporativa opcional.
- Varias instalaciones y usuarios.
- Monitorización y copias externas.

### 6.4 Nube

- Las mismas imágenes de contenedor.
- Los mismos charts y APIs de Kubernetes.
- Base de datos PostgreSQL administrada o autogestionada.
- Almacenamiento S3 compatible.
- Trabajadores escalables.
- Proxy o balanceador externo.
- Secretos proporcionados por la plataforma.

La migración podrá sustituir componentes gestionados —PostgreSQL, S3, identidad
o observabilidad— mediante valores, sin cambiar los dominios.

### 6.5 Híbrido

- Núcleo y automatizaciones en la instalación.
- Publicación, acceso compartido y notificaciones externas en nube.
- IA ejecutada exclusivamente en la instalación local.
- Sincronización selectiva, nunca réplica implícita de todos los datos.

### 6.6 Namespaces

Distribución inicial:

```text
reefops-system       Envoy Gateway, operadores y plataforma
reefops-secret-delivery ESO sintético y trust root público
reefops-domains      dominios Go
reefops-intelligence visión, agentes y procesamiento
reefops-data         PostgreSQL, SeaweedFS y NATS
reefops-identity     ZITADEL, OpenFGA y ReefOps Authorizer
reefops-observability métricas, trazas y logs
linkerd              control plane de Linkerd
```

Los namespaces no son fronteras de dominio suficientes por sí mismos. Se
complementarán con ServiceAccounts, NetworkPolicies, políticas de Linkerd,
credenciales y permisos de base separados.

### 6.7 Topología de entrada

```text
Internet / LAN
      │
      ▼
APIM
      │ valida TLS, token, límites y ruta
      │
      ├──▶ ReefOps Authorizer ──▶ PDP
      │
      ▼
Adaptador del dominio / Query API
      │ valida ActorContext firmado
      ▼
Caso de uso
      │
      └──▶ outbox ──▶ NATS ──▶ otros dominios
```

El IdP se expondrá mediante una ruta propia con los requisitos de protocolo que
necesite. Las APIs de administración del APIM, PDP, NATS, PostgreSQL y Linkerd,
y la de MQTT si se incorpora, permanecerán internas.

### 6.8 Imágenes y recursos mínimos

Servicios Go:

- compilación estática cuando sea viable;
- imagen `scratch` o distroless;
- ejecución non-root;
- filesystem de solo lectura;
- sin shell ni package manager;
- límites y requests medidos;
- una imagen reutilizable con distintos comandos cuando no perjudique el
  aislamiento.

Servicios Python:

- runtime mínimo separado de imágenes de entrenamiento;
- modelos montados o descargados como artefactos versionados;
- imágenes diferentes para CPU y aceleración;
- dependencias fijadas;
- workers escalables a cero cuando no se utilicen.

Componentes de terceros se instalarán con características opcionales
desactivadas. Se medirá CPU, memoria, almacenamiento y tiempo de arranque en el
Mac mini antes de aceptar una dependencia permanente.

## 7. Configuración y portabilidad

La aplicación abstraerá:

- URL de PostgreSQL;
- proveedor de almacenamiento de objetos;
- cola y sistema de trabajos;
- adaptador de Home Assistant y, opcionalmente, broker MQTT;
- proveedores de correo y notificaciones;
- runtime local de IA por capacidad;
- servicio de publicación;
- autenticación externa.

Los adaptadores se seleccionarán mediante configuración. La lógica de dominio
no dependerá de SDK específicos de nube.

## 8. Identidad y autorización

- ZITADEL self-hosted compatible con OIDC/OAuth.
- OpenFGA para autorización fina.
- APIM como único PEP para tráfico norte-sur.
- ReefOps Authorizer como adaptador entre rutas y decisiones del PDP.
- Roles y relaciones por instalación, sistema, recurso y acción.
- Tokens y sesiones revocables.
- Enlaces compartidos representados como sujetos o concesiones limitadas.
- Federación de identidad configurable.
- MFA o passkeys exigibles para operaciones sensibles.
- Service accounts diferentes para servicios, cámaras, integraciones y agentes.
- Mínimo privilegio y denegación por defecto.
- Modelos de autorización versionados.
- Auditoría de decisiones y cambios de relaciones.

La autorización norte-sur se comprobará antes de llegar al dominio. El dominio
validará el contexto firmado y sus propias invariantes.

### 8.1 Diseño de APIs autorizables en el gateway

Para que el APIM pueda decidir sin conocer el dominio:

- el recurso protegido aparecerá preferentemente en la ruta;
- cada operación tendrá una acción declarada en metadata;
- las rutas evitarán comandos que mezclen recursos con permisos diferentes;
- los identificadores del body que cambien el objeto protegido requerirán una
  autorización adicional o un endpoint rediseñado;
- las cargas masivas se dividirán o llevarán una lista de decisiones
  verificables;
- WebSockets y streams autorizarán conexión y, cuando proceda, suscripciones.

Ejemplos:

```text
GET  /systems/{system_id}                 system:view
POST /systems/{system_id}/measurements    measurement:create
POST /organisms/{organism_id}/move        organism:move
GET  /shares/{share_id}                   share:view
```

### 8.2 Listas, búsquedas y proyecciones

Autorizar una URL no basta para filtrar una colección. Para operaciones como
“listar todos los acuarios visibles”:

- el dominio de proyecciones mantendrá índices de acceso alimentados por
  eventos de autorización;
- o consultará una API de lookup del PDP;
- la Query API aplicará el conjunto permitido antes de devolver resultados;
- el APIM seguirá autorizando el derecho general a ejecutar la consulta.

Nunca se devolverá una colección completa para filtrarla en Angular.

### 8.3 Tráfico que no pasa por el APIM

Eventos NATS, trabajos, cron, integraciones de dispositivo y process managers
no atraviesan el APIM. Utilizarán:

- identidad de workload mediante Linkerd/ServiceAccount;
- ACL del bus;
- subjects y streams permitidos;
- actor y delegación propagados en el sobre del mensaje;
- validación de capacidad para acciones sensibles;
- auditoría de consumidor y efecto.

El `ActorContext` de una petición no se reutilizará indefinidamente. Los
comandos asíncronos incluirán una concesión o decisión con alcance y expiración
adecuados al flujo.

### 8.4 Defensa frente a bypass

- NetworkPolicy y Linkerd impedirán acceso norte-sur directo a los dominios.
- Solo el APIM podrá alcanzar los adaptadores HTTP públicos.
- Los health checks utilizarán puertos o rutas separadas sin casos de uso.
- El APIM eliminará cabeceras de identidad entrantes.
- Los adaptadores verificarán firma, audiencia y expiración del contexto.
- Se rotarán claves sin interrumpir contextos todavía válidos.
- Una decisión de autorización tendrá identificador y quedará correlacionada
  con comando, evento y auditoría.

## 9. Seguridad

- HTTPS incluso en despliegues locales cuando sea viable.
- Secretos fuera de imágenes y repositorio.
- Cifrado de credenciales y tokens almacenados.
- Separación de originales privados y derivados publicados.
- Protección CSRF, XSS, SSRF, inyección y subida de archivos.
- Validación y análisis de archivos antes de procesarlos.
- Límites de tamaño, frecuencia y consumo.
- Auditoría de accesos y cambios sensibles.
- Dependencias e imágenes escaneadas.
- Actualizaciones firmadas o verificables.
- Sin telemetría obligatoria.

Las cámaras, controladores y el servidor MCP no serán accesibles desde Internet
por defecto.

La cadena de suministro incluirá:

- registro OCI privado o GHCR con imágenes referenciadas por digest;
- builds multi-arquitectura reproducibles mediante BuildKit;
- SBOM con Syft o equivalente;
- análisis de imágenes y configuración con Trivy;
- firma y verificación con Cosign antes del despliegue;
- actualización controlada de dependencias mediante Dependabot;
- políticas de admisión para impedir imágenes no aprobadas en perfiles
  profesionales.

Los accesos a secretos, cambios de claves, excepciones de red y modificaciones
de políticas de seguridad serán auditados.

## 10. Copias de seguridad y recuperación

Una copia coherente incluirá:

- PostgreSQL;
- objetos y documentos;
- configuración no secreta;
- claves necesarias para descifrar;
- manifiesto de versiones;
- comprobación de integridad.

Requisitos:

- copia manual y programada;
- destino local y remoto opcional;
- cifrado antes de salir de la instalación;
- política de retención;
- restauración en otra máquina;
- prueba periódica de restauración;
- copia previa a actualizaciones y migraciones;
- exportación completa en formatos documentados.

Los backups no se considerarán válidos hasta haber verificado que pueden
restaurarse.

En Docker Desktop, los volúmenes del clúster residen dentro de la VM y no
constituyen una copia. El destino primario inicial será una carpeta SMB
dedicada en un NAS local, montada fuera de la VM y protegida mediante
SMB Encryption, acceso nominal, invitado denegado, ABSE y ABE. Los artefactos
se cifrarán además con `age` antes de escribirse. La clave `age` y el material
de recuperación no residirán en el mismo NAS. Se añadirá un segundo medio para
completar 3-2-1 y la estrategia se endurecerá cuando el sistema controle vida
animal crítica.

El destino externo es un puerto operativo configurable, no una dependencia de
la topología ReefOps. Development usa actualmente un QNAP, pero otra instalación
podrá exportar el mismo paquete cifrado a otro NAS, disco, servicio cloud o
medio offline sin sustituir SeaweedFS ni modificar el `ObjectStore` de CNPG.

Cada backup y restauración conservará `operation_id`, alcance, manifiesto,
versiones, hashes, destino, cifrado, operador, tiempos, resultado y prueba de
lectura. Una restauración de ensayo generará evidencia sin confundirse con una
restauración productiva.

## 11. Actualizaciones y migraciones

- Versiones semánticas de la aplicación.
- Migraciones de base de datos incluidas y auditables.
- Comprobación de compatibilidad antes de actualizar.
- Copia previa automática.
- Actualización gradual de trabajadores.
- Posibilidad de volver a la imagen anterior cuando no haya una migración
  irreversible.
- Canal estable y canal de pruebas opcional.
- Soporte para instalaciones temporalmente desconectadas.

## 12. Observabilidad

- OpenTelemetry Collector como punto de recepción y exportación.
- Prometheus y Alertmanager para métricas y alertas.
- Grafana para paneles.
- Loki para logs y Tempo para trazas distribuidas.
- Logs estructurados con identificador de solicitud, correlación y trabajo.
- Métricas de API, base de datos, colas, sensores e inferencias.
- Estado de cámaras, integraciones y publicación.
- Historial de trabajos fallidos.
- Panel local de salud.
- Exportación opcional hacia herramientas externas.
- Redacción de secretos y datos sensibles.

La observabilidad local deberá funcionar sin enviar información a terceros.
La telemetría técnica tendrá muestreo y retención propios, pero nunca será el
único registro de una decisión funcional. Logs, métricas y trazas enlazarán con
`correlation_id` cuando exista.

Los identificadores de correlación pertenecen a logs, trazas y evidencias
estructuradas, no a labels Prometheus de cardinalidad no acotada.

El despliegue será incremental. La puerta inicial de plataforma instalará
Prometheus Operator, Prometheus, Alertmanager, Grafana, `kube-state-metrics` y
el exportador del nodo antes de los servicios stateful. Loki se incorporará
después de fijar almacenamiento, retención y redacción de logs. OpenTelemetry
Collector y Tempo precederán al primer servicio que emita trazas, pero no se
desplegarán sin consumidores. El alcance y la aceptación de la primera puerta
se definen en [observabilidad mínima](observabilidad.md).

La dependencia es siempre unidireccional: monitores, reglas y dashboards
dependen del componente observado. OpenBao, Envoy Gateway, SeaweedFS,
CloudNativePG y sus operandos nunca dependerán de que Prometheus, Grafana o su
configuración estén disponibles para reconciliarse o prestar servicio.

## 13. Disponibilidad y degradación

| Dependencia no disponible | Comportamiento esperado |
|---|---|
| Internet | Operación local normal; se encolan publicaciones y avisos externos |
| Servicio público | La parte privada funciona; se reintenta sincronización |
| Runtime local de IA | Registro y análisis manual disponibles; se reanudan trabajos al recuperarse |
| Cámara | Alerta técnica; no se infiere ausencia de organismos |
| Sensor | Se marca dato obsoleto; no se sustituye por un valor inventado |
| Home Assistant | API y tareas siguen funcionando; se encolan solo las órdenes seguras permitidas |
| Almacenamiento de objetos | Se bloquean nuevas cargas sin perder metadatos |
| Base de datos | La aplicación entra en modo seguro y no confirma escrituras |
| ZITADEL | Tokens emitidos pueden validarse hasta su expiración; no se crean sesiones nuevas |
| OpenFGA | Operaciones protegidas fallan cerradas; endpoints de salud siguen disponibles |
| Envoy Gateway | APIs no accesibles desde el exterior; los procesos internos continúan |
| Linkerd control plane | El tráfico existente continúa mientras las identidades sigan vigentes |
| NATS | Cada dominio conserva su outbox y reanuda publicación al recuperarse |
| Mac mini | Interrupción total; la recuperación depende de backups externos |

## 14. Modelo de datos y contratos

- Identificadores globalmente únicos para permitir migraciones y sincronización.
- Fechas en UTC y zona horaria de presentación por instalación.
- Unidades almacenadas con dimensión y unidad original.
- Historial temporal para posiciones, máscaras, fichas y configuraciones.
- Borrado lógico cuando sea necesaria trazabilidad.
- Eventos de integración versionados.
- API versionada.
- Esquemas de importación y exportación documentados.
- Archivos identificados por hash para integridad y deduplicación.

## 15. Límites de módulos

Dominios lógicos iniciales —módulos, no necesariamente procesos—:

1. Identidad y organizaciones.
2. Instalaciones y sistemas.
3. Topología y espacio.
4. Habitantes y bienestar.
5. Bioseguridad y salud.
6. Parámetros y mediciones.
7. Operación, tareas y bitácora.
8. Productos, nutrición y dosificación.
9. Equipos e integraciones.
10. Medios y visión.
11. Alertas y emergencias.
12. Inventario y costes.
13. Compartición y publicación.
14. Informes y portabilidad.
15. Experimentos.
16. Conocimiento e inteligencia.
17. Notificaciones.

Cada módulo será propietario de sus invariantes y publicará operaciones y
eventos explícitos.

Propiedad de capacidades transversales:

| Capacidad | Propietario | Regla de frontera |
|---|---|---|
| Catálogo de especies y conocimiento | Conocimiento e inteligencia | Publica versiones; habitantes conserva la referencia y snapshot utilizado |
| Experimentos | Experimentos | Consume hechos y construye su propia serie; no consulta tablas ajenas |
| Notificaciones y escalado | Notificaciones | Consume alertas y preferencias; publica estados de entrega |
| Escrituras offline y conflictos | Adaptador de sincronización | Reentrega comandos originales del cliente al punto de entrada; no modifica agregados |
| IA conversacional y RAG | Conocimiento e inteligencia | Lee únicamente proyecciones propias autorizadas |
| Análisis de imágenes | Medios y visión | Publica observaciones; no diagnostica ni modifica habitantes |
| MCP | Adaptador de entrada de plataforma | Consulta proyecciones autorizadas y envía comandos al propietario |
| Voz y asistentes | Adaptador de entrada/salida | No constituyen fuente de verdad ni ejecutan inferencia remota |

Experimentos, Intelligence, Vision y Notifications pueden comenzar como
workers o módulos del mismo despliegue. Su separación lógica y de datos será
obligatoria aunque compartan binario.

### 15.1 Topología de eventos, no de dependencias

No existirá un grafo de dependencias entre dominios. Existirá un catálogo de
eventos y suscripciones:

```text
Identity ── user-created ───────────────┐
                                        ├──▶ Sharing
Systems ── aquatic-system-created ──────┤
                                        ├──▶ Reporting
Measurements ── measurement-recorded ───┤
                                        ├──▶ Alerts
Inhabitants ── organism-added ──────────┤
                                        ├──▶ Welfare
Vision ── anomaly-detected ─────────────┘
```

El productor conoce el significado del hecho que publica, pero no quién lo
consume. El consumidor conoce el contrato del evento, pero no el modelo interno
del productor.

El catálogo deberá mostrar:

- propietario del evento;
- versión;
- semántica;
- esquema;
- política de compatibilidad;
- datos personales o sensibles;
- retención;
- productores;
- consumidores conocidos a efectos operativos, sin convertirlos en
  dependencias del productor.

### 15.2 Anatomía interna de un módulo Go

Cada dominio expresará primero su vocabulario. Ejemplo:

```text
core/
  systems/
    aquarium.go
    vessel.go
    hydraulic_connection.go
    volume.go
    commands/
      create_system.go
      connect_vessels.go
    queries/
      get_system_overview.go
    events/
      system_created.go
      vessel_connected.go
    ports/
      repository.go
      event_publisher.go
    adapters/
      postgres/
      http/

  inhabitants/
    organism.go
    population.go
    placement.go
    commands/
    queries/
    events/
    ports/
    adapters/

  measurements/
    measurement.go
    parameter.go
    quality.go
    commands/
    queries/
    events/
    ports/
    adapters/
```

Los nombres superiores son conceptos de acuariofilia. `http`, `postgres`,
`grpc` o `nats` aparecen en el borde de cada módulo y no determinan su
estructura.

No será obligatorio replicar carpetas vacías. La estructura crecerá según la
complejidad real del dominio.

### 15.3 Capas dentro del dominio

Cada módulo podrá contener:

1. **Modelo:** entidades, value objects, políticas e invariantes.
2. **Aplicación:** comandos, consultas y coordinación de casos de uso.
3. **Puertos:** contratos que el módulo necesita u ofrece.
4. **Adaptadores:** HTTP, PostgreSQL, NATS, MQTT u otros detalles.

Reglas:

- el modelo no importa frameworks, transporte ni persistencia;
- aplicación depende del modelo y de puertos;
- adaptadores dependen de aplicación y puertos;
- un adaptador nunca se convierte en API interna para otro dominio;
- no se introducirá una carpeta global de utilidades de negocio;
- código compartido contendrá únicamente primitivas técnicas estables.

### 15.4 API pública de un módulo

Cada dominio expondrá hacia usuarios o adaptadores de entrada:

- comandos;
- consultas;
- errores de negocio documentados.

Y hacia el resto de la plataforma únicamente:

- eventos de integración versionados.

Sus entidades, servicios, repositorios, comandos y consultas no serán
importables desde otro dominio.

Cuando un dominio necesite información originada en otro, consumirá eventos y
mantendrá una proyección local con los atributos estrictamente necesarios. Esa
proyección:

- pertenece al consumidor;
- puede tener un modelo diferente del productor;
- es eventualmente consistente;
- conserva la posición o versión del evento aplicado;
- puede reconstruirse mediante replay o una exportación controlada;
- no permite escribir en el dominio productor.

### 15.4A Entradas y comandos externos

La restricción de eventos se aplica a la comunicación **entre dominios**. Un
usuario, dispositivo o API podrá enviar un comando al dominio propietario de la
operación:

```text
HTTP: RegisterUser
        │
        ▼
Identity
        │
        └──▶ identity.user-created.v1
```

Otro dominio no llamará `RegisterUser`. Si una reacción necesita solicitar una
acción adicional, se modelará mediante:

- un evento que el propietario pueda interpretar como desencadenante;
- un process manager externo a los dominios;
- o un comando asíncrono dirigido, si se adopta expresamente ese patrón.

La opción predeterminada será evento. Los comandos asíncronos no se confundirán
con hechos: podrán rechazarse y tendrán destinatario explícito.

### 15.5 Propiedad de datos

Aunque varios dominios utilicen el mismo clúster PostgreSQL:

- cada dominio tendrá una base de datos distinta;
- cada tabla tendrá un módulo propietario;
- solo el propietario realizará escrituras;
- no existirán claves foráneas entre dominios;
- no existirán joins entre dominios;
- una migración no leerá ni modificará tablas ajenas;
- informes y pantallas combinadas usarán proyecciones alimentadas por eventos;
- los adaptadores no accederán a repositorios internos de otros dominios;
- el outbox pertenecerá a la transacción que originó el evento.

Una separación física futura podrá mover el esquema de un módulo a otra base
sin cambiar el contrato funcional.

Cada dominio tendrá:

- base de datos propia;
- owner, usuario de migración y usuario runtime propios;
- permisos únicamente sobre sus objetos;
- migraciones propias;
- repositorios propios;
- outbox propio;
- inbox o registro de eventos procesados;
- backup y restauración coordinables.

La base compartida será una optimización de despliegue, no una integración.

### 15.6 Comunicación

| Situación | Mecanismo |
|---|---|
| Caso de uso dentro de un dominio | Llamada interna |
| Comando de usuario hacia su dominio | REST/OpenAPI |
| Evento entre dominios | Outbox + NATS JetStream |
| Vista que combina dominios | Proyección/read model alimentado por eventos |
| Flujo de varios pasos | Process manager/saga basada en eventos |
| Telemetría de dispositivos | API/eventos de Home Assistant hacia el adaptador; MQTT opcional |

MCP, voz e Intelligence no son atajos para comunicación entre dominios. Actúan
como clientes o adaptadores de entrada: leen proyecciones creadas para su caso
de uso y envían comandos al dominio propietario a través de la misma frontera
autorizada que el resto de clientes.
| Imágenes o vídeo | Referencia a objeto S3 compatible |
| Herramientas para agentes | MCP |

No habrá RPC, REST ni interfaces Go entre dominios. gRPC o ConnectRPC podrán
utilizarse únicamente dentro de una unidad de dominio o como transporte técnico
para servicios de inferencia sin autoridad de negocio.

### 15.6A Composición de lecturas

Angular no encadenará consultas a numerosos dominios para construir cada
pantalla. Los paneles utilizarán read models específicos:

- resumen de instalación;
- estado de un sistema;
- ficha enriquecida de un habitante;
- bandeja de alertas;
- vista compartida;
- informe profesional.

Un servicio de proyecciones consumirá eventos y construirá estas vistas. La
respuesta indicará, cuando sea relevante, la hora o versión hasta la que está
actualizada.

Las proyecciones son descartables y reconstruibles. No son fuente de verdad ni
aceptan comandos de negocio.

### 15.6B Process managers y sagas

Los procesos que abarcan varios dominios no residirán dentro de uno de ellos.
Un process manager:

- escucha eventos;
- conserva el estado mínimo del flujo;
- emite eventos o comandos asíncronos;
- gestiona tiempos de espera y compensaciones;
- es idempotente;
- no modifica directamente datos de los dominios.

Ejemplo de incorporación:

```text
Inhabitants: organism-admission-started
                    │
                    ▼
Admission process manager
    ├── espera quarantine-completed
    ├── espera compatibility-reviewed
    ├── espera placement-approved
    └── emite organism-admission-ready
```

Cada dominio decide de forma autónoma cómo reaccionar y conserva sus propias
invariantes.

### 15.7 Contratos versionados

- OpenAPI para interfaces HTTP.
- Protobuf para servicios técnicos y eventos binarios cuando resulte
  conveniente.
- JSON Schema para eventos JSON.
- Versionado explícito de eventos y herramientas MCP.
- Compatibilidad hacia atrás durante ventanas publicadas.
- Consumer-driven contract tests entre desplegables.
- Ejemplos y fixtures versionados.
- Identificadores globales, correlation ID y causation ID.
- Idempotency key en comandos que puedan reintentarse.

Los contratos generarán clientes cuando sea posible; no se duplicarán modelos a
mano entre Angular, Go y Python.

### 15.7A Semántica y entrega de eventos

Los eventos de integración describirán hechos consumados en pasado:

- `identity.user-created.v1`;
- `systems.aquatic-system-created.v1`;
- `inhabitants.organism-added.v1`;
- `measurements.measurement-recorded.v1`;
- `vision.anomaly-detected.v1`.

No se usarán nombres imperativos como `create-user` para un hecho.

La infraestructura asumirá entrega **al menos una vez**:

- cada consumidor implementará inbox o deduplicación;
- los handlers serán idempotentes;
- no se asumirá un orden global;
- cuando importe el orden se utilizará una clave de agregado y una secuencia;
- los fallos se reintentarán con backoff;
- los mensajes agotados pasarán a una dead-letter stream;
- el replay será una operación soportada;
- los efectos externos usarán idempotency keys;
- la publicación se realizará mediante transactional outbox.

Un evento contendrá suficiente información para que los consumidores previstos
actualicen sus proyecciones sin consultar al productor. Esto no significa
publicar el modelo completo: se aplicarán minimización de datos, privacidad y
contratos específicos.

### 15.7B Consistencia eventual

La plataforma aceptará expresamente que:

- una proyección puede ir por detrás de su fuente;
- una pantalla combinada puede mostrar versiones temporalmente diferentes;
- un consumidor puede estar desconectado y recuperar eventos después;
- una operación entre dominios puede permanecer pendiente;
- una compensación puede sustituir a un rollback distribuido.

Las interfaces mostrarán estados como pendiente, sincronizando, confirmado,
fallido o parcialmente disponible cuando ocultar la consistencia eventual pueda
inducir a error.

Las invariantes que exijan consistencia inmediata deberán residir dentro de un
único dominio. Si una regla no puede tolerar consistencia eventual, se
reconsiderarán los límites antes de introducir coordinación síncrona.

### 15.8 Estructura del repositorio

Se utilizará inicialmente un monorepo:

```text
reefops/
  apps/
    web/                       # Angular
    core/                      # Go: monolito modular
    device-gateway/            # Go
    publication/               # Go
    mcp-server/                # Go o comando del Core

  intelligence/
    agent/                     # Python
    vision/                    # Python: inferencia y entrenamiento
    scientific/                # Python: cálculo especializado

  media/
    capture/                   # Go + FFmpeg/GStreamer

  contracts/
    openapi/
    protobuf/
    events/
    mcp/

  database/
    platform/                  # bootstrap y roles

  deploy/
    compose/
    local/
    cloud/

  docs/
    architecture/
    decisions/
    domain/
```

Dentro de `apps/core`, las carpetas superiores serán dominios, no capas
técnicas:

```text
apps/core/
  systems/
  spatial/
  inhabitants/
  welfare/
  biosecurity/
  measurements/
  husbandry/
  nutrition/
  equipment/
  alerts/
  sharing/
  reporting/
  platform/                    # arranque y composición
```

`platform` arranca procesos, registra adaptadores y conecta cada dominio con el
bus. No permite que un dominio invoque a otro ni contiene reglas de negocio.

### 15.9 Pruebas de arquitectura

CI comprobará:

- ciclos entre módulos;
- imports prohibidos;
- acceso a paquetes `internal` ajenos;
- propiedad de migraciones;
- claves foráneas o joins entre esquemas de dominio;
- llamadas HTTP, RPC o interfaces Go entre dominios;
- compatibilidad de contratos;
- eventos sin versión;
- endpoints no declarados en OpenAPI;
- dependencias técnicas dentro del modelo;
- escrituras cruzadas entre esquemas.

Las reglas se expresarán como código o scripts ejecutables, no solo como
convenciones escritas.

### 15.10 Criterios para extraer un microservicio

Un módulo podrá separarse cuando aparezca una o varias condiciones:

- necesita escalar de forma independiente;
- requiere GPU o hardware distinto;
- tiene dependencias incompatibles;
- debe aislar fallos o riesgos de seguridad;
- tiene un ciclo de despliegue propio;
- lo mantiene un equipo autónomo;
- necesita residencia de datos diferente;
- genera contención medible en el proceso o base compartidos;
- debe ejecutarse en una ubicación física distinta.

Antes de extraerlo deberá tener:

- contrato estable;
- propiedad clara de datos;
- eventos versionados;
- observabilidad;
- idempotencia;
- política de reintentos;
- estrategia de migración;
- pruebas de contrato;
- plan de operación y recuperación.

### 15.11 Servicios inicialmente separados

| Unidad | Lenguaje | Motivo |
|---|---|---|
| Angular Web | TypeScript | Experiencia web y PWA |
| ReefOps Core | Go | Dominio, API, rendimiento y despliegue local |
| Device Gateway | Go | Protocolos, concurrencia y operación en el borde |
| Media Capture | Go | Streaming y supervisión de procesos multimedia |
| Vision | Python | Ecosistema de visión y ML |
| Intelligence | Python | Agentes, RAG y modelos |
| Publication | Go | Superficie pública pequeña y eficiente |
| MCP Server | Go | Frontera segura y contratos del Core |
| ReefOps Authorizer | Go | Adaptación Envoy → OpenFGA y contexto firmado |

Core y MCP podrán comenzar en un mismo binario con comandos diferentes. Vision,
Intelligence y Publication sí tendrán límites de proceso desde el principio.

### 15.12 Tecnologías de datos

| Necesidad | Elección inicial | Evolución posible |
|---|---|---|
| Transacciones | PostgreSQL | PostgreSQL administrado o clúster |
| Geometría | PostGIS | Servicio de cálculo especializado |
| Telemetría | PostgreSQL particionado | TimescaleDB |
| Embeddings | pgvector | Motor vectorial dedicado |
| Objetos | SeaweedFS mediante S3 | Ceph RGW o S3 compatible |
| Buffer de borde | SQLite | Replicación especializada |
| Eventos | NATS JetStream | Clúster NATS |
| Dispositivos | Home Assistant local | MQTT o gateway especializado si se justifica |
| Caché | Ninguna obligatoria | Valkey |
| Analítica masiva | No incluida | ClickHouse |
| Búsqueda | PostgreSQL | OpenSearch si se demuestra necesario |

La introducción de TimescaleDB, ClickHouse, OpenSearch o un vector store
independiente requerirá métricas, carga de operación y un plan de backup.

### 15.13 Identificadores de trazabilidad

Toda interacción utilizará identificadores globales:

| Campo | Significado |
|---|---|
| `environment_id` | Entorno de confianza que origina y puede procesar la operación |
| `trace_id` | Recorrido técnico distribuido de una interacción |
| `span_id` | Operación técnica concreta dentro de una traza |
| `correlation_id` | Flujo funcional completo |
| `causation_id` | Mensaje o acción que causó el actual |
| `event_id` | Evento de integración único |
| `command_id` | Comando único e idempotente |
| `request_id` | Solicitud de entrada |
| `operation_id` | Operación administrativa, importación, backup o restauración |
| `job_id` | Ejecución persistente o planificada |
| `inference_id` | Ejecución concreta de IA o visión |
| `device_command_id` | Orden dirigida a un dispositivo |
| `authz_decision_id` | Decisión del PEP/PDP aplicada |
| `deployment_id` | Cambio desplegado en la plataforma |
| `actor_id` | Usuario, agente, integración o dispositivo responsable |
| `tenant_id` | Propietario u organización |
| `installation_id` | Instalación afectada |
| `system_id` | Sistema acuático afectado, cuando proceda |
| `aggregate_id` | Entidad o agregado modificado |
| `schema_version` | Versión del contrato |

`environment_id` será obligatorio en requests, comandos, eventos, outbox,
inbox, DLQ, jobs, inferencias, órdenes de dispositivo, notificaciones,
publicaciones, decisiones de autorización, auditoría, backups, restores y
despliegues. No se inferirá únicamente desde un UUID o desde datos del payload.
La clave funcional de búsqueda será `(environment_id, correlation_id)`.

`correlation_id` se propagará a través de eventos, process managers,
trabajadores, inferencias, notificaciones y publicación. `causation_id`
permitirá reconstruir el árbol causal, no solo una lista cronológica.

Los identificadores especializados no sustituyen a `correlation_id`: se
enlazarán a él y a su causa inmediata. No se generará un nuevo
`correlation_id` al cruzar un evento, cola, cron, APIM o herramienta.

### 15.14 Sobre de evento

Los eventos utilizarán un sobre común:

```json
{
  "environment_id": "development",
  "event_id": "019...",
  "event_type": "measurements.measurement-recorded",
  "schema_version": 1,
  "occurred_at": "2026-07-24T15:00:00Z",
  "recorded_at": "2026-07-24T15:00:01Z",
  "producer": "measurements",
  "tenant_id": "019...",
  "installation_id": "019...",
  "system_id": "019...",
  "aggregate_id": "019...",
  "aggregate_version": 42,
  "actor": {
    "type": "user",
    "id": "019..."
  },
  "request_id": "019...",
  "command_id": "019...",
  "authz_decision_id": "019...",
  "correlation_id": "019...",
  "causation_id": "019...",
  "trace_id": "..."
}
```

El payload de dominio estará separado del sobre. Los campos opcionales no se
rellenarán con valores inventados.

El consumidor comprobará `environment_id` antes de escribir inbox o ejecutar
un efecto. Un mismatch se rechazará o pondrá en cuarentena; nunca producirá
dosificación, orden física, notificación, publicación o mutación cross-env.

Se distinguirá:

- `occurred_at`: cuándo ocurrió realmente el hecho;
- `recorded_at`: cuándo lo recibió o registró la plataforma;
- `effective_at`: desde cuándo resulta aplicable, si difiere;
- tiempo del dispositivo y calidad de sincronización, para telemetría.

### 15.15 Auditoría funcional

Cada dominio conservará un audit trail append-only para operaciones sensibles:

- creación, modificación, corrección y archivado;
- cambios de permisos y visibilidad;
- mediciones y recalificaciones de calidad;
- dosis, tratamientos y alimentación;
- movimientos de organismos;
- reglas, automatizaciones y órdenes a equipos;
- creación y revocación de enlaces;
- exportaciones y accesos profesionales;
- recomendaciones y confirmaciones de IA;
- cambios de configuración y seguridad;
- altas, bajas, sesiones sensibles y cambios relevantes en ZITADEL;
- escrituras de relaciones y versiones de modelo en OpenFGA;
- decisiones denegadas y elevaciones de privilegio;
- despliegues, migraciones, rotaciones de claves y cambios de red;
- backups, restauraciones, importaciones y exportaciones;
- entregas, visitas y descargas de contenido compartido o público;
- envío, recepción y confirmación de notificaciones;
- replay, re-drive, descarte y edición administrativa de mensajes.

Cada entrada contendrá:

- actor real e identidad delegada;
- origen: web, voz, MCP, API, dispositivo, importación o proceso;
- intención o comando;
- recurso y versión anterior;
- campos modificados o diff seguro;
- resultado;
- motivo proporcionado;
- instante;
- correlación y causación;
- política o permiso que autorizó la acción;
- IP o dispositivo cuando resulte apropiado y legal.

Las correcciones añadirán una nueva versión. No sobrescribirán ni eliminarán
silenciosamente la anterior.

### 15.16 Historial temporal y procedencia

Los datos relevantes conservarán:

- validez temporal;
- instante de registro;
- autor o fuente;
- método;
- instrumento o dispositivo;
- precisión y confianza;
- versión;
- estado de revisión;
- relación con el original.

Esto se aplicará a:

- mediciones;
- posiciones y máscaras;
- fichas de habitantes;
- topologías;
- procedimientos;
- configuraciones;
- imágenes derivadas;
- resultados de laboratorio;
- conocimiento utilizado por agentes.

Un dato derivado enlazará sus entradas. Por ejemplo:

```text
fotografía original
   └──▶ corrección de color
          └──▶ segmentación
                 └──▶ máscara confirmada
                        └──▶ métrica de crecimiento
                               └──▶ alerta
```

Cada transformación conservará herramienta, versión, parámetros y responsable
de su confirmación.

### 15.17 Trazabilidad de IA y agentes

Cada ejecución registrará:

- agente, modelo y versión;
- proveedor o runtime local;
- prompt de sistema versionado;
- instrucciones aplicables;
- contexto y recursos autorizados consultados;
- versiones o fragmentos de las fuentes;
- herramientas invocadas y argumentos redactados;
- resultados de herramientas;
- respuesta estructurada;
- costes, tokens y duración cuando estén disponibles;
- confianza o incertidumbre declarada;
- revisión, corrección y decisión humana;
- acción posterior causada por la respuesta.

Los razonamientos internos no se utilizarán como mecanismo de auditoría. La
auditoría se basará en entradas, fuentes, decisiones estructuradas, herramientas
y resultados observables.

Una recomendación aceptada no se convertirá directamente en un hecho manual:
quedará registrada la cadena `recomendación → confirmación → comando → evento`.

### 15.18 Trazabilidad de dispositivos y automatizaciones

Para cada orden:

- regla, usuario o proceso que la originó;
- configuración y versión de la regla;
- lecturas utilizadas;
- dispositivo y canal;
- orden solicitada;
- instante de envío;
- acuse de recibo;
- resultado observado;
- reintentos y errores;
- anulación manual;
- evento o alerta causante.

Se distinguirá entre:

- orden emitida;
- orden recibida;
- orden ejecutada según el dispositivo;
- efecto físicamente confirmado.

No se asumirá que enviar una orden demuestra su ejecución.

### 15.19 Almacén y consulta de auditoría

Inicialmente:

- cada dominio escribe su auditoría funcional dentro de su almacenamiento;
- publica eventos de auditoría sanitizados;
- un dominio de Audit construye una proyección global de solo lectura;
- OpenTelemetry conserva trazas técnicas;
- los logs estructurados ayudan al diagnóstico, pero no son la fuente de
  auditoría.

La interfaz permitirá buscar por:

- actor;
- fecha;
- sistema;
- organismo;
- recurso;
- comando;
- evento;
- correlación;
- tipo de cambio;
- origen;
- resultado.

Desde cualquier alerta, entrada o cambio se podrá navegar hacia antecedentes y
efectos posteriores.

### 15.20 Integridad y resistencia a manipulación

- Registros de auditoría append-only a nivel de aplicación.
- Permisos de base que impidan su modificación por usuarios ordinarios.
- Hash de objetos y documentos.
- Firmas o HMAC para proyecciones publicadas y mensajes sensibles.
- Encadenamiento de hashes por lotes de auditoría cuando el riesgo lo
  justifique.
- Copias externas cifradas.
- Detección de huecos en secuencias de agregado y consumidores.
- Alertas ante cambios administrativos o de retención.
- Relojes y fuentes temporales monitorizados.

No se afirmará inmutabilidad criptográfica si el administrador de la
infraestructura conserva capacidad para alterar todos los componentes. Se
documentará el nivel real de garantía de cada despliegue.

### 15.21 Retención, privacidad y eliminación

Trazabilidad no justificará conservar indefinidamente todos los datos
personales. Cada categoría tendrá:

- finalidad;
- ámbito;
- periodo de retención;
- base de autorización;
- política de exportación;
- política de eliminación o anonimización.

Cuando exista una solicitud válida de eliminación:

- se anonimizarán o seudonimizarán actores cuando sea compatible con la
  obligación de auditoría;
- se eliminarán objetos y contenido personal no necesario;
- se conservará un tombstone mínimo si hace falta mantener integridad;
- se propagará la eliminación mediante un evento;
- las proyecciones y caches aplicarán la eliminación;
- se registrará que el proceso se completó sin conservar el contenido borrado.

Los eventos evitarán incluir datos personales innecesarios, ya que su
replicación dificulta la eliminación.

### 15.22 Reproducción de una cadena completa

La plataforma deberá poder responder:

1. ¿Quién o qué inició la acción?
2. ¿Con qué permisos?
3. ¿Qué datos y versiones utilizó?
4. ¿Qué decisión tomó el dominio?
5. ¿Qué cambió?
6. ¿Qué eventos se publicaron?
7. ¿Qué consumidores los procesaron?
8. ¿Qué proyecciones, alertas o acciones resultaron?
9. ¿Hubo reintentos, fallos o compensaciones?
10. ¿Qué confirmó, corrigió o descartó una persona?

Esta reconstrucción deberá ser posible por `correlation_id` sin depender de
buscar manualmente en logs de varios contenedores.

### 15.23 Trazabilidad de entrega, reintentos y reproducción

Cada consumidor mantendrá inbox o registro equivalente con `environment_id`,
`event_id`, consumer, versión, primer y último intento, resultado y efecto
producido. La unicidad e idempotencia se evaluarán dentro del entorno. Un retry
conservará la identidad del mensaje y añadirá el número de intento; no simulará
un hecho nuevo.

El paso a dead-letter, el re-drive y el replay serán operaciones auditadas. Se
registrarán filtro, rango, motivo, operador, versión del consumidor y
resultados. Las proyecciones podrán reconstruirse sin volver a ejecutar efectos
físicos, notificaciones o integraciones no idempotentes.

Para publicaciones y avisos se distinguirá:

```text
solicitado → proyectado → enviado → aceptado por destino → entregado/visitado
```

No se afirmará entrega cuando solo exista aceptación del proveedor. Los enlaces
compartidos registrarán creación, alcance, acceso, expiración y revocación sin
guardar más datos personales del visitante de los necesarios.

### 15.24 Trazabilidad de identidad, autorización y delegación

Toda acción protegida enlazará:

- sujeto autenticado y actor efectivo;
- delegante, agente o service account, cuando exista;
- sesión o credencial mediante un identificador no secreto;
- `authz_decision_id`;
- modelo y versión de OpenFGA;
- revisión del estado o tuplas empleadas cuando el PDP lo permita;
- acción, recurso, resultado y política del Gateway;
- expiración del `ActorContext`;
- entorno del issuer, audience, Gateway, store OpenFGA, recurso y actor
  efectivo, que deberán coincidir.
Las escrituras de relaciones se originarán en comandos auditados o consumidores
identificables. Se podrá recorrer:

```text
hecho de dominio
  → evento
  → actualización de relación
  → decisión de autorización
  → comando permitido o denegado
```

Los tokens, secretos y atributos sensibles no se copiarán al audit trail.

### 15.25 Trazabilidad de plataforma y cadena de suministro

El despliegue de una versión será relacionable con:

- `environment_id` y revisión GitOps del entorno;
- commit y ejecución CI;
- SBOM, análisis y firma;
- digest de cada imagen;
- chart y valores no secretos;
- migraciones ejecutadas;
- cambios de contratos y modelos de autorización;
- operador o automatización;
- resultados de rollout, rollback y verificaciones.

Una alerta o regresión podrá asociarse al despliegue vigente. Los cambios
directos de emergencia serán detectados como drift y deberán reconciliarse o
documentarse.

### 15.26 Criterios de aceptación de trazabilidad

La trazabilidad se probará con escenarios de extremo a extremo. Como mínimo:

1. una dosis manual autorizada;
2. una orden automática y su confirmación física;
3. una detección visual que genera alerta y revisión humana;
4. una invitación profesional creada, usada y revocada;
5. una publicación y su retirada;
6. un evento fallido, enviado a DLQ y reproducido;
7. una modificación de permisos;
8. un despliegue con migración y rollback;
9. un backup y una restauración de ensayo;
10. una eliminación propagada a proyecciones y objetos.

Para cada escenario, una consulta por `(environment_id, correlation_id)`
recuperará la cadena ordenada, evidenciará huecos y redactará datos según
permisos. Las pruebas de contrato verificarán que ningún productor omite campos
obligatorios del sobre. También intentarán consumir, autorizar, descifrar,
restaurar o reproducir datos entre development y production; todos esos casos
deberán fallar cerrados y dejar evidencia en el entorno receptor.

### 15.27 Contratos espaciales y de simulación

Topología y espacio será propietario de escenas, sistemas de coordenadas,
geometrías, posiciones, máscaras y sus versiones temporales. Medios y visión
solo será propietario de originales, derivados y observaciones visuales.

Una observación de segmentación incluirá referencia al medio, calibración,
modelo, máscara propuesta, confianza y sistema de coordenadas. Al aceptarla,
Topología y espacio creará una nueva versión con procedencia; nunca se
sobrescribirá una máscara manual o confirmada.

Los cálculos de luz, flujo, sombra, proximidad y crecimiento serán jobs
versionados. Sus contratos incluirán:

- snapshot de escena y topología de entrada;
- motor, algoritmo y versión;
- parámetros, malla, resolución y condiciones de contorno;
- resultado almacenado como objeto y resumen consultable;
- incertidumbre y limitaciones;
- `job_id`, correlación y hashes de entradas y salidas.

Los resultados son estimaciones y no hechos físicos hasta ser contrastados con
mediciones u observaciones.

### 15.28 Escritura offline y sincronización

La PWA podrá crear offline únicamente operaciones declaradas sincronizables,
inicialmente mediciones, entradas de bitácora y finalización de tareas. El
cliente generará `command_id`, `occurred_at`, versión base del agregado y una
marca de tiempo monotónica local cuando esté disponible.

Al recuperar conexión, el adaptador de sincronización:

1. autentica de nuevo al actor;
2. conserva el comando y correlación originales;
3. entrega el comando idempotente al punto de entrada del propietario;
4. registra `received_at` y calidad del reloj por separado;
5. devuelve aceptado, duplicado, rechazado o conflicto;
6. exige resolución explícita cuando no exista una política determinista.

No habrá last-write-wins silencioso. Dosis, tratamientos, permisos,
automatizaciones, control de equipos y cambios estructurales no serán
ejecutables offline en la primera versión.

### 15.29 Ciclo de vida de modelos locales

El catálogo local de modelos conservará:

- nombre, función, arquitectura, versión y licencia;
- procedencia, URL de descarga y digest de pesos;
- firma y resultado del análisis de artefactos;
- requisitos de RAM, CPU, GPU y runtime;
- compatibilidad y evaluación por versión de ReefOps;
- fecha de instalación, activación, retirada y operador;
- datasets y métricas de evaluación permitidos;
- versión de tokenizer, embedding y parámetros relevantes.

Un modelo nuevo pasará por descarga, verificación, evaluación local, activación
controlada y rollback. Cambiar un modelo de embeddings generará una
reindexación trazable; no se mezclarán vectores incompatibles. Los modelos y
configuraciones necesarios para interpretar una inferencia histórica se
conservarán o se archivará evidencia suficiente para reproducirla.

### 15.30 Canales externos de voz y asistentes

Un canal externo podrá transportar una orden o reproducir una respuesta, pero
no recibirá el contexto completo ni ejecutará la inferencia de ReefOps. El
adaptador aplicará minimización, consentimiento, scopes, confirmación reforzada
y redacción antes de enviar cualquier contenido.

La transcripción y la síntesis de voz serán locales. Las plataformas que exijan
reconocimiento o generación remotos no serán compatibles con este perfil; solo
podrán transportar texto o eventos ya minimizados y autorizados.

## 16. Fases arquitectónicas

### Fase A — Desarrollo local

- Monorepo modular de producto, catálogo de plataforma separado y repositorio
  GitOps privado.
- Kubernetes de Docker Desktop con provisionador `kind`.
- Helm como definición canónica.
- Core modular en Go con Screaming Architecture.
- PostgreSQL.
- PostGIS y pgvector.
- SeaweedFS Community.
- Angular, API Go y trabajadores Go.
- Contratos OpenAPI y eventos versionados.
- NATS JetStream.
- Adaptador local de Home Assistant cuando exista el primer caso IoT.
- Envoy Gateway sobre Gateway API.
- Linkerd con mTLS.
- ZITADEL, OpenFGA y ReefOps Authorizer.
- Servicios Python mínimos para IA/visión cuando sean necesarios.
- Inferencia mediante API desacoplada y runtime Metal en macOS cuando aporte
  ventaja.
- OpenBao para secretos runtime y SOPS + `age` solo para bootstrap cifrado.
- Flux CD y Taskfile para reconciliación y operación repetible.
- OpenTelemetry Collector, Prometheus, Alertmanager, Grafana, Loki y Tempo.
- Registro OCI, SBOM, escaneo y firma de imágenes.
- Backups cifrados fuera de la VM de Docker Desktop.
- Contratos de evento en CI, inbox/outbox y visor de correlaciones.
- Imágenes `arm64` y `amd64`.
- Pruebas automatizadas y datos de demostración.

### Fase B — Instancia local instalable

- Imágenes versionadas multi-arquitectura.
- HTTPS y configuración guiada.
- Backup fuera de la VM, restauración y actualización.
- Home Assistant e integraciones locales; MQTT opcional.
- Modo sin conexión de la PWA.
- Políticas de red, Linkerd y autorización en modo deny-by-default.

### Fase C — Publicación híbrida

- Servicio externo de publicación.
- Proyecciones públicas y compartidas.
- Conexión saliente segura.
- Comentarios y documentos de retorno.
- Revocación y auditoría.

### Fase D — Despliegues profesionales y cloud

- Configuración distribuida.
- Identidad federada.
- Servicios administrados opcionales.
- Escalado de trabajadores.
- Reutilización de charts en un Kubernetes multi-nodo.
- Sustitución opcional de PostgreSQL, S3, ZITADEL u observabilidad por servicios
  administrados.

## 17. Decisiones pendientes

- Librerías concretas de Angular y estrategia de estado.
- Mecanismo concreto de trabajos persistentes.
- Tecnología del servicio externo de publicación.
- Protocolo de conexión saliente y recepción de comentarios.
- Límites de datos de una vista viva.
- Selección y dimensionamiento de modelos locales de IA y visión para el M4.
- Soporte inicial de cámaras, sensores y controladores.
- ConnectRPC frente a gRPC puro para servicios técnicos de inferencia y medios.
- Atlas frente a Goose para migraciones.
- Alcance inicial de TimescaleDB frente a particionado nativo.
- Destino físico externo para backups del Mac mini.
- Dominio DNS local y autoridad certificadora para desarrollo.
- Rotación del trust anchor Linkerd sin interrupción antes de producción.
- Evolución del modelo inicial OpenFGA después del primer corte vertical.
- Estrategia de actualización para instalaciones completamente desconectadas.
- Revisión de RPO/RTO antes de controlar vida animal; desarrollo parte de RPO
  24 horas y RTO 4 horas.
- Proveedor de DNS local y mecanismo de exposición LAN.
- Segundo medio externo independiente del QNAP para completar 3-2-1.
- Canales y proveedores de email, push y avisos críticos.
- SAI, watchdog y canal de alerta fuera de banda según criticidad.
- Nombre de la organización y plan de GitHub.
- Credencial de bootstrap y deploy key de solo lectura para el GitOps privado;
  catálogo público de plataforma sin credencial.
- Checks obligatorios concretos tras la primera ejecución de Actions.
- Necesidad de hooks Codex después de observar reglas repetidamente incumplidas.

## 18. Plataforma operativa mínima

El primer entorno completo incluirá:

| Área | Elección inicial |
|---|---|
| Orquestación | Kubernetes integrado de Docker Desktop (`docker-desktop`) |
| Empaquetado | Helm |
| GitOps | Flux CD |
| Bootstrap y operación | Taskfile + scripts idempotentes |
| Entrada/APIM | Envoy Gateway |
| AuthN | ZITADEL |
| AuthZ | OpenFGA + ReefOps Authorizer |
| Malla | Linkerd |
| Eventos | NATS JetStream |
| Dispositivos | Home Assistant; MQTT opcional |
| Secretos | OpenBao con Raft y TLS local; SOPS + `age` solo para bootstrap cifrado |
| Observabilidad | OTel Collector, Prometheus, Alertmanager, Grafana, Loki, Tempo |
| Objetos | SeaweedFS mediante S3 |
| Imágenes y charts | GHCR, BuildKit, Trivy, Syft y Cosign |
| Dependencias | Dependabot |
| IA local | API local de inferencia; Metal fuera de la VM cuando proceda; sin proveedores cloud |
| Backups | Cifrados y almacenados fuera de la VM |

No forman parte del camino crítico inicial Helmfile, ClickHouse, OpenSearch,
una malla multiclúster ni alta disponibilidad simulada en un único Mac. Se
añadirán cuando aparezca una necesidad operativa que los justifique.

## 19. Repositorio de infraestructura y orden de reconciliación

Durante la transición, la infraestructura reside junto al código. Su destino
canónico será `reefops-platform`, mientras que `reefops-gitops` contendrá
únicamente composición concreta por clúster:

```text
infra/
├── bootstrap/              # prerrequisitos y entrada imperativa mínima
├── clusters/
│   └── local/              # reconciliación raíz de la instalación
├── infrastructure/         # servicios comunes del clúster
├── platform/               # identidad, autorización y buses
└── applications/           # releases de ReefOps
```

Flux aplicará capas explícitamente dependientes:

1. fuentes y configuración de Flux;
2. namespaces, quotas y políticas base;
3. `cert-manager`, OpenBao y entrega de secretos;
4. observabilidad mínima interna;
5. Gateway API y controlador inerte de Envoy Gateway;
6. almacenamiento y servicios de datos;
7. NATS;
8. Linkerd;
9. ZITADEL, OpenFGA y ReefOps Authorizer;
10. aplicaciones ReefOps internas, todavía sin entrada norte-sur;
11. plano de datos y rutas protegidas de Envoy Gateway.

Cada capa tendrá health checks y no desbloqueará la siguiente hasta estar
preparada. Las versiones de charts e imágenes estarán fijadas; no se usarán
ramas, tags o dependencias flotantes para despliegues reproducibles.

La primera dependencia real ya separa dos reconciliaciones: infraestructura
crea namespaces y almacenamiento; plataforma depende de ella y despliega
OpenBao. Ambas usan `wait` y la segunda no se desbloquea hasta que la primera
está preparada. La inicialización sellada de OpenBao es una transición
operativa explícita y no se confunde con health de aplicaciones consumidoras.

### 19.1 Entornos y réplicas

`clusters/<cluster>/environments/<entorno>` contendrá únicamente selección de
componentes, valores y patches del entorno. Los componentes reutilizables no se
copiarán entre entornos ni clústeres. Crear una réplica consistirá en:

1. disponer de un Kubernetes compatible;
2. crear identidades OpenBao y SOPS de bootstrap propias;
3. ejecutar el bootstrap;
4. apuntar Flux al path del nuevo clúster;
5. restaurar datos solo cuando corresponda.

Las claves, datos y credenciales nunca se clonarán implícitamente con la
infraestructura.

### 19.2 Reversión y recuperación

| Estado | Fuente de recuperación |
|---|---|
| Manifiestos y configuración | Revert de Git + reconciliación Flux |
| Imágenes | Digest conservado en el registro OCI |
| Secretos runtime | Backup OpenBao + recovery material offline |
| Material de bootstrap cifrado | SOPS + clave `age` externa |
| PostgreSQL, SeaweedFS y volúmenes | Backup verificado + restauración |
| Migraciones | Estrategia expand/contract y procedimiento por versión |
| Cambios físicos | Compensación explícita; nunca replay automático |

Un procedimiento de disaster recovery deberá demostrar la creación de un
clúster vacío, el bootstrap, la reconciliación completa y la restauración de
datos sin pasos manuales no documentados.

## 20. Ayudas de desarrollo y agentes

Las convenciones para `AGENTS.md`, skills, agentes especializados, MCP y
controles ejecutables se definen en
[Desarrollo asistido](desarrollo-asistido.md). Estas ayudas no forman parte del
runtime de ReefOps ni pueden relajar las decisiones de esta arquitectura.

## 21. Toolchains de desarrollo

Las herramientas del host, sus criterios de aislamiento y su proceso de
actualización se definen en
[Entorno de desarrollo](entorno-desarrollo.md). `Brewfile` será la declaración
reproducible del host macOS; las dependencias de aplicación permanecerán
fijadas dentro de sus workspaces.
