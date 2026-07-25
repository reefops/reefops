# Gestión de secretos

## 1. Autoridad y alcance

OpenBao es la autoridad local para secretos runtime y credenciales de
automatización. Pertenece a plataforma, no a ningún dominio de negocio.

Se distinguen cuatro clases:

| Clase | Autoridad | Entrega |
|---|---|---|
| Runtime de workloads | OpenBao local | Identidad Kubernetes; volumen en memoria o Secret Kubernetes cuando el consumidor lo exija |
| Bootstrap de OpenBao/Flux | Custodia offline | SOPS + `age` solo cuando deba declararse cifrado |
| CI de GitHub | OpenBao local | Réplica revocable en GitHub Secrets |
| Datos, medios y backups | Almacén local correspondiente | Nunca GitHub Secrets |

Los servicios de dominio no importarán un SDK de OpenBao ni conocerán dónde se
custodia un secreto. La plataforma lo montará como archivo o lo entregará al
adaptador técnico correspondiente.

External Secrets Operator —ESO— es el adaptador declarativo para componentes
que exigen un Secret Kubernetes. No es una segunda autoridad: cada
`ExternalSecret` selecciona claves concretas, tiene identidad y política
propias y conserva el alcance namespaced.

La NetworkPolicy solo permite clientes desde namespaces autorizados
explícitamente con `reefops.io/openbao-access=true`; pertenecer a ReefOps no
concede acceso lateral al gestor.

La deploy key de lectura de `reefops-gitops` es una credencial de bootstrap:
se genera localmente, se registra como read-only y se instala en `flux-system`.
Puede regenerarse y revocarse sin restaurar datos de OpenBao.

## 2. OpenBao local

OpenBao se desplegará mediante GitOps en el repositorio `reefops-platform`.
Usará almacenamiento persistente y audit device separado. En el Mac Mini
inicial se ejecutará como instancia standalone; esto no se presentará como alta
disponibilidad.

El almacenamiento será Raft integrado incluso con un único nodo. Raft aporta
snapshots verificables y un camino de evolución a varios nodos, pero un único
Mac Mini sigue siendo un único dominio de fallo. No se simulará alta
disponibilidad con varias réplicas sobre el mismo host.

OpenBao no expondrá HTTP en claro. cert-manager, instalado como componente de
plataforma independiente, mantendrá una CA local y un certificado interno para
`openbao.reefops-secrets.svc`. El trust root se distribuirá únicamente a
adaptadores y operadores autorizados. La clave de la CA vive como Secret
Kubernetes gestionado por cert-manager; su backup y rotación forman parte de la
recuperación de plataforma y nunca se copia a GitHub.

El audit device de fichero se declarará en el HCL del servidor y escribirá en
el PVC dedicado. OpenBao 2.6 no permite crearlo mediante API: esta separación
evita que una credencial administrativa pueda redirigir auditoría a un fichero
o endpoint arbitrario. El script de configuración comprobará que `file/`
existe, pero no intentará crearlo, modificarlo ni activar opciones inseguras.
Un cambio de audit device requiere Git, reconciliación y unseal o `SIGHUP`.
Las variables que el chart deriva de TLS no se redefinen en valores locales;
el render validado rechaza nombres de entorno duplicados antes de promocionar
OpenBao.
El nodo local no usa `service_registration "kubernetes"` ni permisos para
mutar pods: Kubernetes descubre el servicio estable y OpenBao mantiene el
mínimo RBAC posible. Tampoco conserva opciones HCL retiradas por la versión
fijada.

cert-manager renovará el certificado leaf, pero OpenBao solo relee sus ficheros
TLS mediante `SIGHUP`. La recarga será una operación local auditada que compara
el serial servido antes y después. La CA tendrá backup `age` independiente del
snapshot Raft. El certificado CA tiene una vigencia larga y se renueva con la
misma clave; la clave CA no rotará automáticamente. Su sustitución exige una
transición explícita de doble confianza.

Dependencias y orden:

1. namespace, almacenamiento y políticas de red;
2. cert-manager y sus CRD;
3. CA local y certificado TLS de OpenBao;
4. OpenBao;
5. inicialización y unseal mediante procedimiento local;
6. auth de Kubernetes y políticas de mínimo privilegio;
7. integración de entrega a workloads;
8. aplicaciones que consumen secretos.

Señales de salud:

- pod preparado;
- almacenamiento montado;
- certificado TLS vigente y verificado;
- instancia inicializada y no sellada;
- audit device operativo;
- login Kubernetes de prueba con permisos mínimos;
- lectura de una ruta sintética sin exponer su valor.

La puerta operativa completa, incluidos backup, evidencia y restauración, se
mantiene en [Estado del proyecto](estado-proyecto.md). Que HelmRelease y
StatefulSet estén reconciliados no basta: una instancia no inicializada,
sellada, sin auditoría o sin snapshot recuperable bloquea todos los
consumidores.

La prueba de bootstrap utiliza dos ServiceAccounts sin token montado
permanentemente: `openbao-smoke-test`, limitada a metadatos de una ruta
sintética, y `openbao-backup`, limitada a snapshots Raft. Cada operación
solicita un token Kubernetes de corta duración y obtiene un token OpenBao
efímero. El token inicial solo configura la autoridad y vuelve después a
custodia offline.

La instalación Helm no espera readiness en el primer arranque porque una
instancia sellada no puede estar preparada antes de la inicialización manual.
Esta excepción solo cubre la instalación: ninguna aplicación consumidora se
desbloquea hasta comprobar `initialized`, `sealed=false`, audit y autenticación.

Fallo seguro:

- una instancia sellada impide entregar nuevos secretos;
- los consumidores no reciben valores vacíos ni defaults inseguros;
- no se reinicializa automáticamente;
- no se imprime unseal material en logs;
- no se habilita `unsafe_allow_api_audit_creation`;
- una pérdida de OpenBao no se recupera mediante Git revert.

Recuperación:

- backup cifrado del almacenamiento fuera de la VM;
- backup de CA y metadatos de certificados conforme al runbook de plataforma;
- custodia offline de recovery/unseal material;
- restauración ensayada y documentada;
- rotación posterior de credenciales potencialmente expuestas;
- evidencia de actor, autorización, backup, versión y resultado.

El procedimiento ejecutable y sus comprobaciones se definen en
[Runbook de recuperación de OpenBao](runbooks/openbao-recuperacion.md).

Los StatefulSets conservarán sus PVC al eliminarse o escalarse. OpenBao usará
usuario no root, seccomp `RuntimeDefault`, privilege escalation deshabilitada,
capabilities mínimas y requests/limits explícitos. Un cambio de esos controles
deberá renderizarse y probarse antes de promoverse.

## 3. Entrega con ESO

La instancia inicial de ESO se ejecuta en `reefops-secret-delivery` y aplica:

- chart OCI e imagen por digest;
- controlador namespaced sin ClusterRoles, webhook ni cert-controller;
- `SecretStore` namespaced con TLS validado mediante la CA pública de OpenBao;
- ServiceAccount exclusiva para autenticación Kubernetes, sin token
  permanente;
- TokenRequest restringido por `resourceNames`;
- política OpenBao limitada a una clave sintética;
- NetworkPolicies default-deny y aperturas mínimas hacia DNS, API Kubernetes y
  OpenBao.

El orden de implantación es deliberado: publicar la plataforma, ejecutar
`openbao-configure` sobre la autoridad activa, promover el commit completo en
GitOps y ejecutar `openbao-verify-eso`. La última prueba rota un valor
sintético, deniega temporalmente la lectura, confirma fallo cerrado y restaura
política y valor. Nunca imprime el token ni el valor.

El [runbook de ESO y OpenBao](runbooks/eso-openbao.md) contiene comandos,
criterios de aceptación, evidencia y diagnóstico seguro.

## 4. GitHub Secrets

Antes de crear un GitHub Secret se intentará, por este orden:

1. `GITHUB_TOKEN` con permisos mínimos;
2. OIDC y credencial efímera;
3. GitHub App;
4. deploy key restringida solo para GitOps; plataforma pública sin credencial;
5. réplica de un secreto de OpenBao.

La sincronización se ejecutará en el host local y siempre iniciará la conexión
hacia GitHub. Un runner hospedado no tendrá ruta a OpenBao ni al clúster.

El contrato de sincronización exige:

- identificador lógico, ruta y versión de OpenBao;
- repositorio y, opcionalmente, environment destino fijados en la allowlist;
- una única versión activa; una versión histórica no puede republicarse;
- allowlist explícita de nombres de secretos CI;
- lectura interactiva o identidad local de corta duración;
- envío por entrada estándar a `gh secret set`;
- ninguna variable, argumento, fichero temporal o salida con el valor;
- registro de actor, versión, destino, instante y resultado sin valor ni hash
  susceptible de facilitar ataques;
- revocación en GitHub cuando se retire el secreto de OpenBao.

La misma operación admite `delete` para retirar la réplica GitHub sin leer el
valor. Cambiar el destino requiere un cambio revisado de la allowlist.
Antes de sincronizar, se verifica que la versión allowlisted sea la versión
actual de OpenBao y que no esté eliminada o destruida. El workflow consumidor
debe leer también la variable no secreta `<NOMBRE>_OPENBAO_VERSION`; la CI
rechaza cualquier job consumidor que no tenga una condición fail-closed exacta
contra la versión activa fijada en la allowlist.

GitHub Secrets no alojará secretos runtime, credenciales de dispositivos,
claves privadas SOPS, material de unseal, contexto del acuario ni backups.

## 5. Trazabilidad

La evidencia no sensible de una sincronización tendrá:

- `operation_id`;
- actor y método de autenticación local;
- decisión de autorización;
- `secret_id` y versión OpenBao;
- owner, repositorio, aplicación de GitHub y environment;
- fecha de inicio y final;
- resultado y error redactado;
- `correlation_id` y `causation_id`.

La auditoría funcional y el audit device de OpenBao son fuentes distintas. Los
logs técnicos ayudan al diagnóstico, pero no reemplazan ninguna de ellas.

## 6. Rotación y replay

Una rotación crea una versión nueva, actualiza consumidores y después revoca la
anterior. Reintentar la sincronización de la misma versión es idempotente.
Reproducir evidencia no vuelve a leer ni publicar el valor.

La rotación de un secreto CI no despliega aplicaciones. Si afecta a un
artefacto o release, el workflow y el digest resultante conservarán su propia
trazabilidad.
