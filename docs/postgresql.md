# PostgreSQL y CloudNativePG

## Propósito y frontera

PostgreSQL será la fuente transaccional de ReefOps. CloudNativePG (CNPG) es un
operador de plataforma: no es un dominio, no contiene reglas de acuariofilia y
ningún módulo dependerá de sus CRD. Los adaptadores de persistencia hablarán
PostgreSQL mediante `pgx`; la sustitución futura por un servicio administrado no
alterará el modelo funcional.

Compartir un clúster físico no autoriza integración por base de datos. Cada
dominio tendrá:

- base de datos, owner, rol de migración y rol runtime propios;
- migraciones, tablas, outbox, inbox y proyecciones bajo su propiedad;
- ausencia de grants, claves foráneas, vistas, funciones y joins hacia otro
  dominio;
- credenciales y política OpenBao independientes;
- backup físicamente coordinado y capacidad de exportación lógica por dominio.

Una proyección que combine hechos pertenecerá al consumidor y se alimentará por
eventos. No leerá tablas del productor. Esta decisión hace separable cada base
en otro clúster sin cambiar contratos funcionales.

## Versiones y cadena de suministro

Development utilizará la serie soportada CloudNativePG 1.30 con PostgreSQL
18.4. Kubernetes `docker-desktop` 1.34.1 está dentro del rango soportado por
CNPG 1.30. El operador, su chart, PostgreSQL y los plugins se fijarán por versión
y digest, nunca por tags flotantes.

El operando será la imagen oficial `ghcr.io/cloudnative-pg/postgis` en variante
`standard-trixie`. Esta variante hereda pgvector de la imagen PostgreSQL
standard y añade PostGIS. Se prefiere a los volúmenes de extensiones porque
Kubernetes 1.34 no habilita `ImageVolume` por defecto; ReefOps no modificará
feature gates de Docker Desktop para esta etapa. En Kubernetes 1.35 se
reevaluará separar las extensiones en imágenes OCI de solo lectura.

El backup utilizará el plugin CNPG-I Barman Cloud, no la integración Barman
in-tree deprecada. La versión inicial será 0.13 o la revisión estable compatible
que se fije al implementar. El plugin residirá junto al operador y se instalará
después de cert-manager.

Fuentes primarias:

- [versiones soportadas de CloudNativePG](https://cloudnative-pg.io/docs/1.30/supported_releases/);
- [instalación de CloudNativePG](https://cloudnative-pg.io/docs/1.30/installation_upgrade/);
- [imágenes PostgreSQL oficiales](https://github.com/cloudnative-pg/postgres-containers);
- [imágenes PostGIS oficiales](https://github.com/cloudnative-pg/postgis-containers);
- [plugin Barman Cloud](https://cloudnative-pg.io/plugin-barman-cloud/docs/intro/).

## Topología de development

Development tendrá un único `Cluster` CNPG en `reefops-data` y una única
instancia PostgreSQL. No se desplegarán réplicas en el mismo Mac para aparentar
alta disponibilidad. El operador y Barman Cloud vivirán en
`reefops-database-system`, separados de los datos.

PGDATA tendrá inicialmente 20 GiB RWO sobre
`reefops-hostpath-retain`, con `reclaimPolicy: Retain` y protección frente a
podas Helm/Flux. No habrá un volumen WAL separado: en un solo disco no añade
independencia de fallo y complica la recuperación. La StorageClass de Docker
Desktop no permite expansión; aumentar capacidad exigirá backup, recreación y
restore verificados.

PostgreSQL tendrá TLS, Pod Security `restricted`, UID no root, recursos
acotados y servicios exclusivamente `ClusterIP`. No se crearán Ingress,
`NodePort`, `LoadBalancer`, `externalIPs`, `hostNetwork` ni `hostPort`.
PgBouncer se diferirá hasta que exista un consumidor y métricas reales de
conexiones.

## Identidad, secretos y acceso

OpenBao seguirá siendo la autoridad de credenciales runtime. ESO materializará
Secrets namespaced mediante identidades Kubernetes de mínimo privilegio. Git,
Flux, Helm, logs, evidencias y dashboards no contendrán contraseñas ni cadenas
de conexión.

CNPG conservará sus certificados operativos según su mecanismo nativo. Cada
dominio futuro recibirá una identidad distinta y una NetworkPolicy explícita
por namespace y ServiceAccount. Pertenecer a ReefOps no concede acceso a
PostgreSQL. El acceso administrativo se realizará mediante Kubernetes y
`kubectl cnpg`, quedará auditado y no requerirá publicar el puerto 5432.

El primer despliegue no creará una base funcional compartida. La aceptación
creará y eliminará una base y roles sintéticos. Las bases de dominio aparecerán
solo junto con el módulo propietario y sus migraciones.

## Extensiones y migraciones

La aceptación demostrará que PostgreSQL puede crear y usar:

- PostGIS, con geometría, SRID e índice GiST;
- pgvector, con almacenamiento y consulta de distancia;
- particionado nativo para una tabla temporal sintética.

Habilitar una extensión en una base de dominio seguirá siendo una migración de
ese dominio. Que el binario esté disponible no concede `CREATE EXTENSION` al
rol runtime. El rol de migración y el runtime serán diferentes.

Atlas frente a Goose sigue siendo una decisión del primer servicio, no del
operador. La plataforma garantiza el motor y la recuperación; el dominio
garantiza la compatibilidad de sus migraciones.

## Backup, WAL y recuperación

Barman Cloud archivará WAL continuamente y realizará un backup físico diario en
un bucket dedicado de SeaweedFS. Usará una identidad S3 propia custodiada en
OpenBao. La retención inicial será de siete días; las cuatro copias semanales
externas se derivarán mediante el proceso de custodia.

SeaweedFS está dentro de la misma VM de Kubernetes y por sí solo no es un
destino de DR. Tras cada backup aceptado, un proceso local exportará inventario
y objetos Barman, verificará checksums, cifrará con `age` y escribirá el paquete
en un destino externo allowlisted fuera de la VM. La ruta concreta pertenece a
la configuración privada de plataforma. La identidad privada `age` no se
almacenará en el NAS.

La puerta inicial exige:

1. backup físico `Completed` y WAL archivado;
2. restauración en un segundo `Cluster` CNPG aislado;
3. comparación de datos, extensiones y marcador de consistencia;
4. eliminación del clúster de ensayo sin borrar el backup fuente;
5. exportación cifrada al QNAP y verificación de manifiesto;
6. evidencia encadenada con actor, autorización, revisiones, digests, LSN,
   timeline, tiempos, PVC, correlación, causación y resultado.

El RPO provisional sigue siendo 24 horas y el RTO cuatro horas. El archivado
WAL puede ofrecer un RPO menor, pero no se prometerá hasta medirlo. La
recuperación desde SeaweedFS no demuestra supervivencia a la pérdida conjunta
de Docker Desktop. El riesgo seguirá abierto hasta rehidratar el backup externo
en un S3 vacío y restaurar desde él.

## Consistencia con SeaweedFS

PostgreSQL y S3 no comparten una transacción distribuida. Toda operación de
medios usará estados explícitos y compensables:

1. reservar identidad lógica y clave de objeto;
2. escribir el objeto de forma idempotente con checksum;
3. confirmar metadata en la transacción del dominio y emitir el outbox;
4. detectar y reconciliar objetos huérfanos o filas pendientes.

Los backups coordinados compartirán `backup_set_id`, `environment_id`, corte
temporal, revisión GitOps y checksums. Para restaurar se elegirá primero el
PITR PostgreSQL y después el inventario de objetos compatible. Nunca se
declarará consistente una pareja solo porque sus ficheros tengan fechas
parecidas.

## Observabilidad y aceptación

Prometheus descubrirá las métricas nativas de CNPG y del plugin. Grafana
mostrará estado, conexiones, transacciones, bloqueos, checkpoints, WAL,
capacidad y antigüedad/resultado del último backup. Las alertas cubrirán
instancia ausente, PVC próximo a agotarse, backup vencido o fallido, archivado
WAL detenido y errores de reconciliación.

La aceptación verificará:

- revisiones GitOps exactas y digests efectivos;
- operador, plugin, Cluster, certificados y Secrets preparados;
- una instancia, PVC retenido y ausencia de exposición;
- SQL transaccional, constraints, rollback y concurrencia básica;
- PostGIS, pgvector y particionado;
- reinicio conservando PVC y datos;
- backup, restore aislado, cleanup y evidencia;
- aprobación reforzada antes de restaurar sobre un destino no vacío.

Logs y métricas son diagnóstico técnico. La auditoría de cambios funcionales
seguirá perteneciendo a cada dominio.
