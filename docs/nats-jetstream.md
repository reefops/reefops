# NATS JetStream

## Propósito y frontera

NATS JetStream será el transporte durable de eventos de integración entre
dominios. Es plataforma y no contiene reglas de acuariofilia. Un dominio solo
publica contratos versionados propios y no conoce consumidores, streams ni
implementaciones internas de otros dominios.

La entrega funcional combina outbox transaccional, publicación idempotente,
JetStream, inbox durable y handlers idempotentes. Un ACK confirma persistencia
en el inbox, no la repetición de un efecto físico. Dosis, órdenes a equipos y
notificaciones conservan una barrera durable adicional frente a replay.

## Contrato inicial

- subjects `reefops.events.<domain>.<event>.v<major>`;
- `Nats-Msg-Id` igual al identificador estable del evento;
- sobre con `event_id`, `event_type`, `event_version`, `occurred_at`,
  `producer`, `correlation_id`, `causation_id`, actor/delegación y payload;
- un stream `REEFOPS_EVENTS` que captura `reefops.events.>`;
- retención por límites, almacenamiento en fichero y réplica única en
  development, sin presentar el único Mac como HA;
- consumidores pull durables, ACK explícito y backoff acotado;
- cuentas y ACL separadas para sistema, publicación y consumo; ninguna
  aplicación recibe permisos sobre `>`.

## Seguridad y operación

NATS vive en `reefops-system`, usa `ClusterIP`, Pod Security `restricted`, un
PVC `reefops-hostpath-retain` y NetworkPolicy deny-by-default. Solo workloads
con identidad y allowlist explícitas acceden al puerto cliente. Monitorización
y health permanecen internas. No se crea LeafNode, MQTT, WebSocket, NodePort,
LoadBalancer ni ruta Gateway.

Las credenciales runtime se custodiarán en OpenBao cuando aparezca el primer
productor funcional. El gate de plataforma no versiona usuarios ni
contraseñas funcionales.

## Aceptación

La puerta exige revisión GitOps exacta, chart e imagen fijados, servidor listo,
JetStream habilitado, PVC retenido y ausencia de exposición. La prueba crea un
stream y consumidor sintéticos, publica un evento con correlación y causación,
comprueba deduplicación por `Nats-Msg-Id`, redelivery por falta de ACK, ACK
explícito, persistencia tras reinicio y cleanup. Prometheus debe descubrir el
servidor y alertar por indisponibilidad, almacenamiento próximo al límite y
consumidores con retraso.

El gate demuestra transporte durable en el host activo. No demuestra DR tras
perder Docker Desktop; backup y restore del estado JetStream se incorporarán
antes de que el bus transporte efectos sobre vida animal.
