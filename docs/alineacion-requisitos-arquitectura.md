# Alineación entre requisitos y arquitectura

## 1. Propósito

Esta matriz demuestra que cada requisito funcional tiene un propietario
arquitectónico, un mecanismo técnico y una evidencia verificable. Complementa
los documentos de requisitos y arquitectura; no sustituye sus detalles.

Una fila cubierta significa:

1. el dominio propietario acepta comandos externos y conserva sus invariantes;
2. ningún otro dominio accede a sus tablas o APIs;
3. las reacciones entre dominios se producen mediante eventos versionados;
4. las lecturas combinadas utilizan proyecciones;
5. comando, evento, job y efecto mantienen correlación y causación;
6. existe al menos un criterio de aceptación automatizable.

## 2. Matriz de cobertura

| Requisitos | Propietario | Soporte técnico y datos | Evidencia mínima |
|---|---|---|---|
| RF-001, RF-002, RF-005A | Instalaciones y sistemas | PostgreSQL; agregados Installation y AquaticSystem; API Go | Creación, cambio de etapa y archivado reconstruibles por correlación |
| RF-003, RF-004, RF-005 | Instalaciones y sistemas | Recipientes, cámaras, conexiones y lotes de agua versionados | Eventos de cambio topológico y validación de volúmenes |
| RF-006, RF-007, RF-008, RF-008E | Topología y espacio | PostGIS; escenas y vistas 2D/3D temporales | Snapshot, sistema de coordenadas, actor y versión |
| RF-008A, RF-008B, RF-008C, RF-008D | Topología y espacio | Geometrías PostGIS, objetos derivados y series temporales | Linaje máscara manual/inferida/confirmada y métricas de evolución |
| RF-009, RF-009A, RF-009B | Topología y espacio | Jobs locales de luz, flujo y simulación; resultados en almacenamiento S3 | Hash de escena, motor, parámetros, incertidumbre y resultado |
| RF-010, RF-011, RF-014 | Parámetros y mediciones | PostgreSQL particionado; catálogo, observación y metrología | Método, unidad original, instrumento, calidad y tiempos |
| RF-012 | Alertas y emergencias | Consumidor de eventos, motor determinista y estado de alertas | Regla versionada, entradas, evaluación, apertura y cierre |
| RF-013 | Informes y portabilidad | Proyecciones analíticas alimentadas por eventos | Serie y correlación reproducibles desde datos versionados |
| RF-015 | Experimentos | Agregado Experiment y proyecciones propias de resultados | Hipótesis, baseline, cambios, interferencias y conclusión |
| RF-020 | Conocimiento e inteligencia | Catálogo local versionado y snapshots de fichas | Fuente, licencia, versión y fecha de actualización |
| RF-021, RF-022, RF-023, RF-025 | Habitantes y bienestar | Agregados Organism/Population; historial temporal | Alta, posición, observaciones y evaluaciones de bienestar |
| RF-024, RF-034 | Bioseguridad y salud | Cuarentenas, tratamientos y casos sanitarios | Responsable, producto, dosis, evolución y revisión profesional |
| RF-026, RF-027 | Habitantes y bienestar | Pasaporte, procedencia, movimientos y aclimatación | Cadena de custodia y cambios de ubicación |
| RF-030, RF-030A, RF-031, RF-032, RF-033 | Productos, nutrición y dosificación | Productos, lotes, planes, calculadoras y registros | Cálculo separado de ejecución; confirmación y consumo de lote |
| RF-040, RF-041 | Equipos e integraciones | Inventario técnico y planes de mantenimiento | Alta, configuración, intervención y estado |
| RF-042, RF-043 | Equipos e integraciones | Adaptador local de Home Assistant; Device Gateway/SQLite y MQTT opcionales | Orden solicitada/recibida/ejecutada/confirmada físicamente |
| RF-050, RF-051, RF-052 | Operación, tareas y bitácora | Entradas, procedimientos, tareas y cambios de agua | Autor, ocurrencia, registro, evidencias y finalización |
| RF-053, RF-054 | Alertas y emergencias | Incidentes, planes offline, escalado y timeline | Detección, reconocimiento, acciones y resolución |
| RF-060, RF-063, RF-063A, RF-063B, RF-063C | Conocimiento e inteligencia | RAG e inferencia local sobre proyecciones autorizadas | Modelo, fuentes, prompt versionado, propuesta y revisión |
| RF-061, RF-062 | Medios y visión | SeaweedFS mediante S3, captura guiada y workers Python locales | Original, transformaciones, calidad y resultado |
| RF-063D | Adaptador de voz | STT/TTS local; comandos por Gateway; respuestas minimizadas | Transcripción, dispositivo, actor, confirmación y comando |
| RF-063E, RF-063F, RF-063G | Adaptador MCP | MCP Go, APIM, ActorContext y proyecciones autorizadas | Agente, herramienta, argumentos redactados, decisión y efecto |
| RF-064 | Conocimiento e inteligencia | Guardrails, permisos y human-in-the-loop | Recomendación, confianza, confirmación, comando y evento |
| RF-065, RF-066 | Medios y visión | Cámaras locales, planificador y jobs persistentes | Cámara, captura, intervalo, ausencia de datos y ejecución |
| RF-067, RF-068, RF-069, RF-069A | Medios y visión | Detección, tracking y línea base mediante modelos locales | Modelo, evidencia, confianza, comparación temporal y observación |
| RF-069B | Alertas y emergencias | Consumo de observaciones visuales y reglas de persistencia | Observación, regla, deduplicación, alarma y notificación |
| RF-069C, RF-069D | Medios y visión | Revisión humana, datasets locales y retención | Confirmación/corrección, privacidad, versión y procedencia |
| RF-070 | Identidad y organizaciones | ZITADEL, proyección Identity y eventos de provisión | Sujeto, actor, alta/baja y sesión sensible |
| RF-071, RF-071A | Operación, tareas y bitácora | Asignaciones, procedimientos y evidencias de formación | Responsable, versión de procedimiento y aprobación |
| RF-072, RF-073, RF-074, RF-075, RF-076 | Compartición y publicación | OpenFGA, Authorizer y proyecciones compartidas | Alcance, decisión, acceso, aportación, expiración y revocación |
| RF-077, RF-078, RF-079 | Compartición y publicación | Publication Service y derivados públicos separados | Consentimiento, snapshot, publicación, visita y retirada |
| RF-080, RF-081 | Inventario y costes | Movimientos, lotes, compras y proyecciones de coste | Origen del movimiento, consumo, ajuste y cálculo |
| RF-090, RF-091 | Informes y portabilidad | Read models, workers de informe y exportación | Snapshot/vista viva, fuentes, generación, acceso y hash |
| RF-092 | Adaptador de sincronización | PWA, cola local y comandos idempotentes con versión base | occurred/received, duplicado, conflicto y resolución |
| RF-100 | Notificaciones | Consumidor de alertas, preferencias y adaptadores de canal | Solicitud, proveedor/canal, aceptación, entrega y escalado |

## 3. Reglas bidireccionales

La cobertura también se valida desde la arquitectura hacia el producto:

| Componente técnico | Requisito que lo justifica |
|---|---|
| Angular/PWA | Uso móvil, RF-050, RF-051 y RF-092 |
| PostgreSQL/PostGIS | RF-001–RF-015 y datos transaccionales/espaciales |
| SeaweedFS mediante el puerto S3 | RF-061, RF-065–RF-069D, RF-075 y RF-091 |
| NATS JetStream + outbox/inbox | Reacciones entre todos los dominios sin dependencias |
| Home Assistant + adaptador de dispositivos; MQTT opcional | RF-042 y RF-043 |
| Workers y planificador | RF-041, RF-051, RF-066, RF-091 y RF-100 |
| Python y runtimes locales | RF-060–RF-069D |
| Envoy Gateway | Entrada de APIs privadas, compartidas y públicas |
| ZITADEL | RF-070 y autenticación de RF-063E–RF-079 |
| OpenFGA + Authorizer | RF-070 y RF-072–RF-079 |
| Linkerd y NetworkPolicy | Protección frente a bypass y tráfico interno |
| Publication Service | RF-073–RF-079 |
| OpenTelemetry y Audit | Requisitos no funcionales de trazabilidad |
| OpenBao, SOPS de bootstrap, firma y SBOM | Seguridad, portabilidad y operación verificable |

## 4. Definition of Done arquitectónica

Un requisito no se considerará implementado hasta que:

- tenga agregado o proyección propietaria identificada;
- sus comandos y eventos estén versionados;
- exista autorización a nivel de acción y recurso;
- genere auditoría funcional y correlación distribuida;
- trate idempotencia, retry y degradación;
- aplique privacidad, retención y procedencia;
- tenga pruebas de contrato y aceptación;
- pueda exportarse o eliminarse conforme a su política;
- no introduzca imports, joins, foreign keys ni llamadas síncronas entre
  dominios.

CI extraerá los identificadores `RF-*` de requisitos y comprobará que todos
aparecen en esta matriz. También detectará referencias a requisitos
inexistentes.
