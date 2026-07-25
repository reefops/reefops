# Observabilidad mínima de development

## 1. Decisión

La primera capa posterior a la entrega de secretos será una observabilidad
interna y ligera para el único clúster `development`. Su objetivo es hacer
diagnosticables las siguientes incorporaciones stateful; no sustituye la
auditoría funcional ni adelanta la exposición pública de interfaces.

La primera fase usará `kube-prometheus-stack` para desplegar:

- Prometheus Operator y sus CRD;
- una instancia de Prometheus;
- una instancia de Alertmanager;
- Grafana con dashboards declarativos;
- `kube-state-metrics`;
- el exportador del nodo;
- reglas y monitores mantenidos como código.

Loki y un recolector de logs se incorporarán después de disponer del
almacenamiento de objetos y de una política de redacción y retención probada.
OpenTelemetry Collector y Tempo se desplegarán antes del primer servicio
aplicativo que emita trazas. No se instalarán componentes vacíos únicamente
para completar el diagrama final.

## 2. Límites de la primera fase

- Toda la telemetría permanece en el Mac mini y no usa Grafana Cloud ni otros
  proveedores externos.
- Prometheus, Alertmanager y Grafana son `ClusterIP`; el acceso operativo
  temporal se realiza mediante `kubectl port-forward`.
- No existe ingreso norte-sur hasta desplegar Envoy Gateway, ZITADEL, OpenFGA y
  ReefOps Authorizer.
- Grafana es de solo lectura para el acceso operativo inicial. Dashboards,
  fuentes de datos y alertas se cambian mediante GitOps, no desde la interfaz.
- La retención inicial de métricas es de siete días, con límite adicional por
  tamaño para proteger el disco local.
- Se usa una sola réplica por componente. Ejecutar réplicas adicionales en un
  único nodo no aporta alta disponibilidad.
- Los volúmenes son persistentes, pero la telemetría técnica no entra todavía
  en el RPO de datos funcionales. Perderla no autoriza a perder auditoría.
- Las versiones de chart e imágenes quedan fijadas por digest y se actualizan
  mediante PR.

## 3. Señales mínimas

La capa debe permitir consultar como mínimo:

- salud y recursos del nodo y de Kubernetes;
- estado de pods, deployments, statefulsets, daemonsets, jobs y PVC;
- reconciliaciones de Flux;
- disponibilidad y caducidad de certificados;
- estado sellado y disponibilidad de OpenBao sin exponer secretos;
- disponibilidad y sincronización de External Secrets Operator;
- consumo de CPU, memoria y almacenamiento de la propia observabilidad.

Las primeras alarmas ReefOps cubrirán:

- un workload esperado no disponible;
- reinicios repetidos;
- PVC próximo a agotarse;
- certificado próximo a caducar;
- reconciliación Flux fallida o estancada;
- OpenBao sellado o no disponible;
- `SecretStore` o `ExternalSecret` no preparado;
- Prometheus sin poder evaluar reglas o Alertmanager no disponible.

## 4. Criterios de aceptación

La fase se considera operativa únicamente cuando:

1. Flux reconcilia una revisión completa y fija de `reefops-platform`;
2. todos los pods esperados están preparados y sujetos a límites de recursos;
3. Prometheus descubre el nodo, Kubernetes y los componentes de plataforma
   acordados;
4. Grafana consulta Prometheus y presenta los dashboards aprovisionados;
5. una regla sintética pasa de inactiva a activa y vuelve a inactiva;
6. Alertmanager recibe esa alerta sin depender de un receptor externo;
7. reiniciar Prometheus, Alertmanager y Grafana conserva su configuración y los
   datos dentro de las garantías declaradas;
8. las interfaces no son accesibles desde fuera del clúster salvo mediante un
   túnel operativo autenticado por Kubernetes;
9. la prueba guarda revisión local y Flux, UIDs, tiempos, resultado y
   `correlation_id`, sin valores secretos;
10. `task validate` pasa en producto, plataforma y composición GitOps.

## 5. Evolución

Después de esta puerta se desplegarán SeaweedFS, PostgreSQL y NATS. La
observabilidad se ampliará en cada incorporación con sus métricas, dashboards,
alarmas y prueba de degradación. Loki, OpenTelemetry Collector y Tempo se
añadirán cuando existan consumidores reales y una decisión de almacenamiento,
retención y redacción aplicable.
