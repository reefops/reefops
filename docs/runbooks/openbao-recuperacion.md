# Runbook de recuperación de OpenBao

## Objetivo

Recuperar la autoridad local de secretos después de perder el PVC o el clúster.
Git, Flux y SOPS no sustituyen este procedimiento.

## Prerrequisitos

- snapshot Raft cifrado almacenado fuera de la VM;
- identidad `age` separada de los backups; para recuperación completa del
  emplazamiento, custodiada también offline;
- recovery/unseal material custodiado por separado;
- OpenBao de versión compatible desplegado, inicializado y no sellado;
- CA local y certificado de servicio recuperables o reemitibles;
- token local autorizado para snapshot;
- ventana de mantenimiento y bloqueo de rotaciones.

No se copiará ningún material de recuperación a GitHub, GitHub Secrets, logs o
argumentos de CI.

### Custodia progresiva

La custodia de recuperación se implanta en tres fases:

1. conservar los backups cifrados en el NAS externo a Kubernetes;
2. crear una copia *off-site* del material de recuperación, cifrada localmente
   con una credencial independiente antes de subirla al proveedor cloud;
3. copiar el paquete cifrado a un USB, desconectarlo y custodiarlo fuera del
   emplazamiento.

La fase 3 queda diferida hasta disponer del USB. Mientras tanto, no se afirmará
que existe una copia offline ni que se cumple completamente la regla 3-2-1. La
contraseña de la envoltura no se guardará en el NAS, el proveedor cloud, Git,
el workspace, el historial de shell ni junto al propio paquete.

## Inicialización de una instalación nueva

1. confirmar que el StatefulSet usa Raft, los PVC están enlazados y la política
   de red está activa;
2. confirmar que cert-manager, el Issuer y el Certificate están preparados;
3. acceder localmente al pod mediante un port-forward temporal usando la CA;
4. ejecutar `bao operator init` de forma interactiva y trasladar inmediatamente
   recovery/unseal material a custodia offline;
5. realizar el unseal sin copiar las claves a historial, ficheros del proyecto
   o variables persistentes;
6. confirmar que el audit device declarativo `file/` se activó durante el
   unseal;
7. autenticarse con el token inicial y ejecutar `task openbao-configure`;
8. crear identidades operativas de mínimo privilegio y revocar el uso habitual
   del token inicial;
9. verificar TLS, audit, mount `ci/`, auth Kubernetes y snapshot cifrado inicial.

La salida de `bao operator init` no será automatizada ni capturada por CI.
Si un cambio GitOps del HCL reinicia el pod después de inicializarlo, se
reutilizan dos claves Shamir existentes para unseal; nunca se vuelve a ejecutar
`bao operator init`.

## Backup y restauración de la CA

La CA de OpenBao se respalda separadamente del snapshot Raft. El backup:

```sh
   REEFOPS_OPENBAO_CA_BACKUP_DIR=<montaje-externo>/openbao/ca \
REEFOPS_OPENBAO_BACKUP_RECIPIENT=age1... \
task openbao-ca-backup
```

El Secret se reduce a tipo, nombre, namespace y datos; se cifra directamente
por pipe con `age` y nunca se escribe plaintext. Se conserva el SHA-256 del
artefacto cifrado.

Después de perder el clúster:

1. reconciliar únicamente infraestructura y cert-manager;
2. mantener suspendidas las reconciliaciones `reefops-openbao-pki` y
   `reefops-openbao`;
3. comprobar el digest del backup de CA;
4. restaurarlo con confirmación explícita mediante `task openbao-ca-restore`;
5. reconciliar PKI para reemitir el certificado leaf con la CA recuperada;
6. reconciliar OpenBao y continuar con el restore Raft;
7. verificar que los clientes siguen confiando en la misma CA.

Si no existe backup de CA, se debe tratar como una rotación de trust root:
generar una CA nueva, distribuir ambos roots durante la transición, reemitir el
leaf, verificar consumidores y retirar el root anterior. Nunca se sustituirá
silenciosamente.

## Renovación TLS

OpenBao recarga `tls_cert_file` y `tls_key_file` al recibir `SIGHUP`. Tras una
renovación de cert-manager, el operador ejecutará `task openbao-tls-reload`.
La tarea compara el serial servido con el serial deseado, envía la señal al
contenedor y falla si el endpoint TLS no presenta después el certificado nuevo.
La operación conserva actor, seriales, correlación y resultado, sin claves.

La clave CA no rota automáticamente. Su sustitución sigue el procedimiento de
doble confianza anterior; la renovación automática normal reemplaza el
certificado leaf y, con amplia antelación, puede renovar el certificado CA
conservando la misma clave.

## Backup

1. Autenticarse localmente en OpenBao con una identidad de backup.
2. Definir un directorio externo y un recipient público `age`.
3. Ejecutar:

   ```sh
   REEFOPS_OPENBAO_BACKUP_DIR=<montaje-externo>/openbao/raft \
   REEFOPS_OPENBAO_BACKUP_RECIPIENT=age1... \
   task openbao-backup
   ```

4. Verificar que se produjo un `.snap.age`, su SHA-256 cifrado y evidencia local.
   Después, comprobar de forma no destructiva el digest, descifrado efímero,
   estructura, checksums internos y metadatos Raft:

   ```sh
   REEFOPS_OPENBAO_VERIFY_FILE=<montaje-externo>/openbao/raft/openbao-<fecha>.snap.age \
   REEFOPS_OPENBAO_VERIFY_MANIFEST=<montaje-externo>/openbao/raft/openbao-<fecha>-<operación>.snap.age.manifest.json.age \
   REEFOPS_OPENBAO_VERIFY_MANIFEST_SHA256=<digest-manifiesto-verificado> \
   REEFOPS_OPENBAO_VERIFY_SHA256=<digest-verificado> \
   REEFOPS_OPENBAO_VERIFY_IDENTITY=/custodia/backup.agekey \
   task openbao-verify-backup
   ```

5. Copiar el fichero cifrado a un segundo medio y verificar el mismo digest.
6. Ejecutar trimestralmente una restauración de prueba en una instancia aislada.

El script crea un snapshot plaintext temporal con permisos `0600`, lo cifra de
inmediato y lo elimina al finalizar. El destino debe estar fuera de la VM de
Docker Desktop.

El destino concreto se conserva únicamente en la composición operativa
privada. Antes de cada backup se comprobará que:

- el punto es un filesystem SMB montado y no un directorio local vacío;
- el servidor y share coinciden con el destino esperado;
- el usuario operativo puede crear, leer y retirar un fichero sintético;
- SMB Encryption está habilitado en el NAS;
- el espacio y la cuota permiten completar la operación.

La identidad privada `age`, las claves Shamir y el token inicial no se copiarán
al share. El MCP de QNAP puede verificar metadatos y permisos, pero no
transporta backups ni sustituye estas comprobaciones locales.

## Restauración

La restauración reemplaza el estado de OpenBao y es destructiva. Hay dos modos:

- `in-place`: recupera sobre una instancia que conserva el mismo seal material
  y no usa `-force`;
- `disaster-recovery`: reconstruye tras perder PVC o clúster, exige confirmar
  que se custodia el seal material original y usa `-force` de forma explícita.

Durante la implantación del primer ensayo, el entrypoint ejecutable permanece
restringido a `isolated-recovery`. Los pasos de esta sección describen el
contrato futuro de recuperación real, pero no autorizan ni permiten todavía
aplicarlo sobre `reefops-secrets`.

1. detener consumidores y rotaciones;
2. conservar un snapshot de emergencia del estado recuperable actual;
3. comprobar versión, digest cifrado y procedencia del backup;
4. autenticar localmente con permiso de backup y restore;
5. ejecutar el entrypoint aprobado para la autoridad activa cuando exista.

Actualmente no hay tarea ejecutable para restaurar `reefops-secrets`. El motor
interno solo admite por allowlist el target aislado y no se publica en
Taskfile. Habilitar el activo exige entrypoints distintos para `in-place` y
`disaster-recovery`, un fence de mantenimiento verificable, backup preventivo
cifrado y verificado, aprobación de un solo uso ligada a la identidad completa
del target y tratamiento explícito de resultados inciertos. No se resolverá
relajando variables del motor del ensayo.

Después de cualquiera de los dos modos:

1. realizar unseal con el material correspondiente al snapshot;
2. autenticarse con una identidad restaurada de recuperación;
3. ejecutar:

   ```sh
   REEFOPS_CORRELATION_ID=<correlation-id-del-restore> \
   REEFOPS_CAUSATION_ID=<operation-id-del-restore> \
   REEFOPS_OPENBAO_VERIFY_PATH=ci/healthcheck \
   task openbao-verify-recovery
   ```

La verificación exige instancia inicializada y no sellada, audit device, auth
Kubernetes, mount `ci/` y metadata de una ruta sintética. Solo entonces se
considera recuperada la autoridad.

6. comprobar la evidencia y el snapshot previo;
7. rotar credenciales emitidas después del snapshot;
8. reactivar consumidores progresivamente;
9. registrar actor, autorización, snapshot, resultado y rotaciones.

`-force` queda encapsulado en el modo `disaster-recovery`, con doble
confirmación y evidencia. Un restore aplicado no demuestra que las credenciales
posteriores al snapshot sigan siendo válidas.

## Ensayo aislado

El ensayo trimestral se ejecuta en `reefops-openbao-recovery`, nunca en
`reefops-secrets`. La raíz `recovery/openbao` de `reefops-platform` permanece
opt-in y la composición privada la habilita mediante una Kustomization GitOps
temporal. Tiene namespace, Service, CA, certificado y PVC propios, sin Ingress
y con denegación de red por defecto.

Excepción Raft: el target reutiliza el `node_id` lógico conservado dentro del
snapshot (`reefops-local-0`). Un identificador distinto permite descifrar y
aplicar el estado, pero el nodo restaurado no se reconoce como miembro y no
puede elegir líder. El aislamiento no depende del `node_id`, sino del namespace,
Service, PKI, red y PVC independientes.

La configuración declarativa del audit device debe coincidir exactamente con
la persistida en el snapshot, incluida su descripción. Una descripción
específica del entorno de ensayo provoca que OpenBao rechace el post-unseal
como una modificación no permitida. El aislamiento del audit se obtiene con su
PVC, no cambiando el contrato lógico restaurado.

Antes de aplicar el snapshot se comprueba:

- contexto Kubernetes `docker-desktop`;
- endpoint y SNI exclusivos de recuperación;
- instancia activa en `reefops-secrets` preparada y sin cambios;
- target sin inicializar y PVC nuevos;
- ambos SHA-256 custodiados: snapshot y manifiesto;
- versión OpenBao productora idéntica a la restauradora;
- aprobación local acotada al cluster ID temporal, modo
  `disaster-recovery`, actor y expiración.

La inicialización del target produce material temporal que no se imprime, no se
versiona y se elimina tras el restore. Una vez aplicado el snapshot, el target
se sella y solo se abre con las claves Shamir originales. La comprobación final
registra como mínimo:

- commit GitOps y versión OpenBao;
- digest de snapshot y manifiesto;
- índice y término Raft;
- actor, aprobación, `correlation_id` y `causation_id`;
- hora inicial/final, RPO observado, RTO y resultado;
- identidad del clúster restaurado, audit device y mounts; el login Kubernetes
  de mínimo privilegio se vuelve a comprobar contra la autoridad activa;
- salud inalterada de la instancia activa.

Antes de inicializar se ejecuta `task openbao-recovery-preflight` desde
`reefops-platform`. La tarea comprueba contexto, salud del activo, certificado,
ServiceAccount sin token, ausencia de RBAC cluster-wide y target aún sin
inicializar. El restore aislado exige además:

```sh
export REEFOPS_OPENBAO_RESTORE_TARGET_SCOPE=isolated-recovery
export REEFOPS_OPENBAO_RESTORE_TARGET_NAMESPACE=reefops-openbao-recovery
export REEFOPS_OPENBAO_RESTORE_TARGET_SERVICE=openbao-recovery
export BAO_ADDR=https://127.0.0.1:18200
export BAO_TLS_SERVER_NAME=openbao-recovery.reefops-openbao-recovery.svc
```

El target se inicializa con `task openbao-recovery-init`. El material temporal
permanece en `tmpfs`; si una operación queda a medias, la misma tarea reanuda la
fase sin reinicializar.

La aprobación y la ejecución son actos separados. Con los digest ya
verificados, el operador ejecuta en un terminal:

```sh
REEFOPS_OPENBAO_RESTORE_SHA256=<digest-snapshot> \
task openbao-recovery-approve
```

La tarea muestra un challenge ligado al digest que debe escribirse
explícitamente. Solo después, y manteniendo los demás parámetros de backup y
restore del runbook, se ejecuta `task openbao-recovery-restore`. El ejecutor no
crea aprobaciones ni completa confirmaciones: exige el estado aprobado,
revisión GitOps y clúster originales, además de
`REEFOPS_ORIGINAL_SEAL_MATERIAL_CONFIRMED=true` y
`REEFOPS_CONFIRM_OPENBAO_RESTORE=force-restore-<digest>`.

Después del unseal original, la verificación y el cierre son explícitos:

1. ejecutar `task openbao-verify-recovery` indicando el state file del drill y
   usando como `causation_id` la operación de restore;
2. ejecutar `task openbao-recovery-evidence-seal` para cifrar por streaming el
   audit y el inventario de PVC/PV;
3. ejecutar `task openbao-recovery-close`; el cierre enlaza la verificación,
   valida los digest y calcula RPO/RTO sin usar la hora de cleanup;
4. retirar la Kustomization mediante PR GitOps;
5. ejecutar `task openbao-recovery-cleanup-verify` indicando el `drill_id`.

La limpieza no se considera completa hasta que desaparezcan HelmRelease,
StatefulSet, pods, PVC, los PV registrados por UID y sus VolumeAttachments.
Los PV con política distinta de `Delete` requieren una aprobación separada y
no se borran imperativamente. Nunca se eliminan recursos de `reefops-secrets`
como parte del ensayo.
