# Estado del proyecto

## 1. Propósito

Este documento separa decisiones, implementación y operación real. No sustituye
los requisitos, la arquitectura ni los runbooks. Se actualiza al completar una
etapa de plataforma y nunca presenta como operativo un componente meramente
documentado o reconciliado pero no preparado.

## 2. Estado verificado el 25 de julio de 2026

### Decidido y documentado

- límites de dominio mediante Screaming Architecture;
- integración entre dominios exclusivamente mediante eventos versionados;
- outbox, inbox, idempotencia, correlación, causación y auditoría funcional;
- Angular para la aplicación web, Go para Core y adaptadores de plataforma y
  Python para IA y visión local;
- Kubernetes, Flux, Helm/Kustomize y GitOps;
- Envoy Gateway, ZITADEL, OpenFGA, ReefOps Authorizer y Linkerd;
- PostgreSQL, PostGIS y pgvector;
- NATS JetStream; MQTT queda diferido y la frontera IoT preferente será Home
  Assistant;
- SeaweedFS Community accedido mediante el puerto S3;
- OpenBao como autoridad local de secretos runtime;
- External Secrets Operator como adaptador declarativo para consumidores que
  requieran Secret Kubernetes;
- repositorios `reefops`, `reefops-platform` y `reefops-gitops`;
- un único entorno development sobre el Kubernetes local, con arquitectura y
  composición preparadas para añadir otros entornos y promover por digest
  mediante GitOps cuando sean necesarios.

### Implementado y operativo

- organización y repositorios GitHub;
- validación, gobierno, análisis de secretos y actualización de dependencias del
  repositorio de producto;
- Kubernetes integrado de Docker Desktop en el contexto `docker-desktop`;
- Flux CD reconciliando el repositorio privado `reefops-gitops`;
- catálogo público `reefops-platform` consumido por Flux;
- namespaces base;
- cert-manager y PKI interna de OpenBao;
- StatefulSet de OpenBao con TLS, almacenamiento Raft, PVC separado de auditoría
  y políticas de red;
- OpenBao inicializado, no sellado, con audit device, autenticación Kubernetes,
  políticas operativas, backup cifrado y lectura sintética verificados;
- External Secrets Operator 2.8.0 desplegado en
  `reefops-secret-delivery`, con alcance namespaced, RBAC mínimo, TLS contra la
  CA pública de OpenBao y sin credenciales estáticas;
- integración ESO/OpenBao verificada de extremo a extremo con autenticación
  Kubernetes, refresco, revocación temporal, fallo cerrado, auditoría y
  restauración automática del estado sintético;
- restauración Raft ensayada en una instancia aislada, validada con la identidad
  original y retirada después mediante GitOps;
- ciclo de recuperación aislada automatizado hasta sellado cifrado de
  evidencias, cierre idempotente y verificación de limpieza por UID de PVC/PV.

El OpenBao operativo pertenece a development. No hay un entorno production
desplegado ni reservado en el clúster local. Cuando se cree, su autoridad
requerirá una instalación y ceremonia nuevas, sin copiar claves o datos de
development.

### Operativo con riesgo residual

OpenBao está operativo para el desarrollo local. El simulacro recuperó el mismo
`cluster_id`, eligió líder Raft, verificó el audit device y leyó los metadatos
de `ci/healthcheck`. El intento exitoso tuvo un RTO de 12 minutos y 22 segundos;
la campaña completa, incluido el diagnóstico del primer intento, duró 3 horas,
44 minutos y 36 segundos, dentro del RTO provisional de 4 horas. El snapshot
restaurado tenía 4 horas, 44 minutos y 5 segundos, dentro del RPO provisional
de 24 horas.

La custodia se completará en tres fases:

1. NAS: backups cifrados externos a la VM de Kubernetes, completado;
2. nube: copia *off-site* de un paquete cifrado localmente con una segunda
   credencial independiente, pendiente;
3. medio físico: copia cifrada en USB, desconectada y guardada fuera del
   emplazamiento, diferida hasta disponer del dispositivo.

La tercera fase no bloquea el desarrollo local ni la siguiente etapa de
plataforma. Hasta completarla se acepta explícitamente el riesgo residual: la
copia en nube aporta separación geográfica, pero no constituye una copia
offline y la estrategia 3-2-1 no está completa.

### No desplegado

- PostgreSQL;
- SeaweedFS;
- NATS JetStream;
- integración con Home Assistant;
- Envoy Gateway y Linkerd;
- ZITADEL, OpenFGA y ReefOps Authorizer;
- observabilidad;
- aplicaciones Angular, Go o Python;
- dominios funcionales.

## 3. Puertas de secretos

### 3.1 OpenBao

OpenBao se considera terminado para desarrollo local porque existe evidencia de
estos controles:

1. instancia inicializada una sola vez y no sellada;
2. material de unseal y token inicial fuera de Git, del workspace, del historial
   de shell, de logs y de la VM de Docker Desktop;
3. CA respaldada con cifrado `age` fuera de la VM;
4. audit device de fichero activo en su PVC dedicado;
5. auth Kubernetes habilitado y login sintético de mínimo privilegio probado;
6. políticas operativas y de backup aplicadas como código;
7. mount sintético `ci/` y lectura de metadatos verificada sin revelar valores;
8. snapshot Raft inicial cifrado fuera de la VM y digest comprobado;
9. restauración ensayada en una instancia aislada y retirada después por GitOps;
10. evidencia no sensible con actor, `operation_id`, `correlation_id`, versión,
    tiempos, decisión, resultado y hashes de artefactos cifrados;
11. Flux reconciliado y sin drift;
12. `task validate` superado en los repositorios afectados.

El material secreto inicial exige una ceremonia local y no se declara en Git.
La configuración posterior, scripts, políticas y comprobaciones sí serán
reproducibles y revisables.

La evidencia local append-only enlaza el cierre del simulacro con la operación
de verificación, conserva correlación y causación, revisiones, identidades de
los PVC, índices Raft, tiempos y hashes. El audit funcional final se conserva
cifrado en el QNAP. Los datos sensibles utilizados para abrir y autenticar la
instancia restaurada no forman parte de esa evidencia.

### 3.2 Entrega ESO/OpenBao

La entrega declarativa se considera operativa para development porque:

1. el chart OCI y la imagen de ESO están fijados por digest;
2. Flux aplica el mismo commit completo de plataforma validado localmente;
3. no existen `ClusterSecretStore`, ClusterRoles ni webhooks de ESO;
4. la identidad de lectura de OpenBao usa TokenRequest de Kubernetes y una
   política limitada a `ci/eso-smoke-test`;
5. `SecretStore` y `ExternalSecret` están preparados;
6. la rotación sintética llega al Secret destino sin imprimir el valor;
7. al sustituir temporalmente la política por `deny`, ESO falla cerrado y no
   reemplaza el último valor correcto;
8. política y valor originales se restauran y se comprueban;
9. el audit device registra nuevos logins y lecturas;
10. la evidencia enlaza actor, revisiones local y Flux, UIDs, generaciones,
    `cluster_id`, correlación, causación, resultado y restauración.

El procedimiento y los fallos conocidos resueltos se mantienen en el
[runbook de ESO y OpenBao](runbooks/eso-openbao.md).

## 4. Decisiones operativas iniciales

Para desarrollo local se fijan provisionalmente:

| Objetivo | Valor inicial |
|---|---|
| RPO | 24 horas |
| RTO | 4 horas |
| Backup programado | Diario y antes de cambios stateful |
| Retención | 7 diarios y 4 semanales |
| Cifrado | `age` antes de abandonar el proceso que genera la copia |
| Restauración de ensayo | Trimestral y después de cambios del procedimiento |

El destino primario es una carpeta SMB dedicada, cifrada en transporte y
montada en macOS fuera de la VM de Kubernetes. El acceso invitado está
denegado; únicamente una identidad operativa dedicada y el grupo administrativo
del NAS tienen lectura/escritura. ABSE y ABE están activados, y la carpeta no
usa Time Machine, sincronización, versiones anteriores ni papelera.

Los artefactos se cifrarán con `age` antes de escribirse en el montaje SMB. La
identidad privada `age`, el material Shamir y el token inicial de OpenBao no se
guardarán en esa carpeta ni en el NAS. El montaje remoto es el primer medio
externo. La segunda fase usará un proveedor cloud únicamente con cifrado previo
en el Mac; no se subirá la identidad privada `age` en claro ni la contraseña de
la envoltura al mismo proveedor. La tercera fase añadirá el medio USB offline
cuando esté disponible.

El asistente MCP del NAS se utiliza solo como adaptador operativo local, con
credencial limitada a estado y carpetas compartidas y acceso restringido al
host operador. La versión beta confirmó una lectura correcta, pero falló al
crear una carpeta sobre un volumen thin por interpretar incorrectamente su
espacio libre; la creación se realizó en la consola del NAS y se verificó
después por MCP. Este adaptador no es fuente de verdad ni sustituye GitOps.

Los objetivos se revisarán antes de que ReefOps ejecute acciones físicas o
proteja vida animal. En ese perfil, un RPO de 24 horas puede ser insuficiente.

## 5. Secuencia después de la entrega de secretos

La siguiente etapa stateful puede comenzar manteniendo registrada la tercera
fase de custodia como riesgo residual:

1. desplegar la
   [observabilidad mínima de development](observabilidad.md), comenzando por
   métricas, alertas y dashboards; logs y trazas se añadirán cuando tengan
   consumidores y almacenamiento definidos;
2. desplegar SeaweedFS y ejecutar el contrato S3 de ReefOps;
3. decidir y desplegar PostgreSQL mediante CloudNativePG con backup y
   restauración;
4. desplegar NATS JetStream;
5. completar entrada, identidad, autorización y malla;
6. implementar el primer corte vertical de instalaciones y sistemas acuáticos;
7. integrar Home Assistant cuando el corte vertical necesite dispositivos.

El primer corte funcional incluirá autorización, persistencia, auditoría,
outbox, evento versionado y trazas. No será un CRUD aislado de esas garantías.
