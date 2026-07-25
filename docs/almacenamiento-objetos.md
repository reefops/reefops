# Almacenamiento de objetos

## Propósito y frontera

ReefOps almacenará fotografías, vídeo, documentos, modelos, resultados de
simulación y evidencias grandes mediante un puerto de aplicación
`ObjectStoragePort`. SeaweedFS es una implementación de plataforma y no un
dominio: ningún agregado, evento ni caso de uso importará sus tipos, consultará
su metadata interna o dependerá de su API administrativa.

El adaptador inicial hablará S3 y recibirá endpoint, región, credenciales y
opciones de direccionamiento por configuración. Este límite permite sustituir
SeaweedFS por Ceph RGW o por otro proveedor S3 sin cambiar el modelo funcional.
PostgreSQL conservará la identidad lógica, ownership, clasificación, hash,
tamaño, MIME, procedencia y estado de cada medio; una clave S3 nunca será por sí
sola una autorización ni la fuente funcional de verdad.

## Perfil local de development

El Mac mini ejecutará una única réplica de master, volume, filer y endpoint S3
en `reefops-data`. Master, catálogo del filer y datos tendrán PVC RWO
independientes sobre la StorageClass de development. El endpoint S3 y las
interfaces internas serán `ClusterIP`; no habrá Ingress, `NodePort`,
`LoadBalancer`, acceso anónimo ni consola administrativa pública.

No se crearán réplicas de datos sobre el mismo host para aparentar alta
disponibilidad. La colocación `000` y la réplica única son una limitación
explícita de development: la pérdida de Docker Desktop o del Mac puede destruir
todo el almacenamiento. La evolución multi-nodo deberá definir racks, centros
de datos, factor de réplica, quorum y anti-afinidad reales antes de declararse
HA.

La identidad S3 se generará una sola vez mediante una ceremonia local auditada,
se custodiará en OpenBao y se materializará con ESO en un Secret namespaced que
contenga el fichero de configuración exigido por SeaweedFS. Git, Helm, Flux,
logs y evidencias no contendrán access keys, secret keys ni su representación
codificada. La política OpenBao y el rol Kubernetes sólo permitirán a la
identidad ESO de `reefops-data` leer esa ruta concreta.

Los pods no necesitan hablar con la API Kubernetes: no montarán tokens de
ServiceAccount ni recibirán RBAC de descubrimiento. Pod Security, seccomp,
capabilities mínimas y NetworkPolicy restringirán el proceso. Los clientes S3
se autorizarán por namespace y ServiceAccount cuando existan; pertenecer a un
namespace ReefOps no concede acceso implícito.

## Contrato S3 inicial

Antes de habilitar un consumidor, una prueba ejecutará contra la versión y los
digests exactos:

- creación y borrado de un bucket sintético;
- `PUT`, `GET`, `HEAD`, `DELETE` y `ListObjectsV2`;
- carga multipart, cancelación y finalización;
- lectura por rango, metadata, tags y checksums;
- URL prefirmada de duración corta;
- persistencia tras reiniciar master, volume, filer y S3;
- exportación cifrada, restauración en un espacio aislado y comparación de
  contenido y metadata.

El contrato registrará resultado y versiones sin almacenar contenido privado
ni credenciales. `ETag` sólo se tratará como identificador opaco; el hash
funcional de ReefOps será un checksum explícito calculado sobre el contenido.
Versionado, lifecycle, notificaciones, políticas de bucket, retención y Object
Lock quedan prohibidos hasta añadir pruebas portables para cada capacidad.

No se crearán buckets funcionales antes de que exista su propietario. Medios
privados, derivados compartidos, publicaciones y evidencias tendrán espacios
lógicos separados. Publicar nunca cambiará el ACL del original: producirá un
objeto derivado y minimizado cuya referencia pertenecerá a la proyección
pública.

## Trazabilidad

Toda operación funcional de objeto mantendrá:

- `environment_id`, actor o servicio y delegación;
- decisión de autorización y política aplicada;
- `correlation_id` original y `causation_id` inmediato;
- identificador lógico, bucket lógico y clave redactada o seudonimizada;
- tamaño, MIME, checksum y procedencia;
- versión del adaptador, proveedor y resultado;
- motivo, retención y aprobación cuando haya publicación o eliminación.

La auditoría funcional vive en el dominio propietario y sus eventos. Métricas,
logs S3 y trazas técnicas ayudan a diagnosticar, pero no demuestran que una
persona estuviera autorizada a leer, compartir o borrar un medio.

Las reintentos serán idempotentes. Un fallo después de escribir pero antes de
confirmar se resolverá consultando el identificador y checksum esperados, no
creando claves nuevas sin control. La eliminación funcional tendrá estados
solicitada, autorizada, ejecutada y verificada; la expiración física no se
inferirá a partir de una respuesta HTTP aislada.

## Backup y recuperación

Los PVC y la replicación interna no son backup. Development realizará una
exportación lógica S3:

1. inventariar buckets permitidos y objetos con metadata y checksums;
2. descargar a un directorio temporal con permisos mínimos;
3. crear manifiesto y cadena de evidencia con revisión GitOps, digests,
   correlación, actor y tiempos;
4. cifrar con la clave `age` de backup antes de escribir en el QNAP;
5. eliminar el staging en claro;
6. restaurar periódicamente en un bucket aislado y comparar inventario,
   contenido y metadata;
7. destruir el espacio de ensayo tras conservar evidencia no sensible.

El RPO inicial de development será de 24 horas y el RTO objetivo de cuatro
horas. Son objetivos operativos, no garantías, hasta que un backup programado y
un ensayo completo los demuestren. El primer despliegue no se considera
aceptado sin una restauración sintética; tampoco habilita todavía backups de
datos reales.

Una restauración nunca sobrescribirá un bucket activo por defecto. Requerirá
destino aislado, digest del artefacto, revisión compatible y confirmación
explícita. Recuperar datos no recupera las filas PostgreSQL que los referencian;
cuando exista la base de datos, ambos backups compartirán un identificador de
consistencia y un runbook coordinado.

## Cadena de suministro y evolución

Se utilizará el chart oficial de SeaweedFS 4.39.0 y la imagen 4.39, fijados por
digest. Como el chart oficial se publica mediante repositorio Helm HTTP, la
automatización verificará el SHA-256 anunciado por upstream y reflejará el
paquete sin modificar en GHCR antes de que Flux lo consuma como OCI. El mirror
preservará versión, procedencia y digest upstream en la evidencia.

Una actualización exige repetir render, políticas, contrato S3, persistencia y
restore. El rollback declarativo puede restaurar manifiestos e imágenes, pero
no revierte formatos internos ni datos; una versión que cambie el formato de
disco necesita compatibilidad probada o migración export/import.

