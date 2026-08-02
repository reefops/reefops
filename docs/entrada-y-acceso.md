# Entrada norte-sur y acceso

## 1. Decisión

ReefOps usará Kubernetes Gateway API como contrato de entrada y Envoy Gateway
como controlador y plano de datos. No se desplegará un controlador `Ingress`
clásico en paralelo.

La incorporación se divide deliberadamente en dos puertas:

1. **fundación inerte**: CRD de Gateway API y Envoy Gateway, controlador,
   métricas y reconciliación GitOps, sin `Gateway`, listeners, rutas ni servicios
   de plano de datos;
2. **entrada protegida**: `Gateway`, TLS, rutas y políticas de seguridad, después
   de desplegar ZITADEL, OpenFGA, ReefOps Authorizer y Linkerd y de decidir
   nombres, resolución local y modelo de certificados.

Instalar el controlador no autoriza por sí solo a publicar una interfaz. Toda
exposición debe aparecer en Git y superar la segunda puerta.

## 2. Reparto de responsabilidades

| Elemento | Responsabilidad |
|---|---|
| Gateway API | Contrato Kubernetes para listeners y rutas |
| Envoy Gateway | Reconciliar el contrato y gestionar el plano de datos Envoy |
| `SecurityPolicy` | Aplicar AuthN y delegar AuthZ externa |
| ZITADEL | Autenticar personas y emitir tokens OIDC |
| ReefOps Authorizer | Consultar OpenFGA y emitir el `ActorContext` interno firmado |
| OpenFGA | Decidir relaciones y permisos sobre recursos ReefOps |
| cert-manager y OpenBao PKI | Emitir y custodiar material TLS según el perfil |
| Servicios de dominio | Verificar el `ActorContext`, aplicar invariantes y auditar |

Envoy Gateway es plataforma, no un dominio de negocio. No conoce modelos
internos de acuarios, habitantes o instalaciones. El contrato ruta-acción lo
mantiene ReefOps Authorizer. Los dominios no duplican políticas relacionales ni
consultan OpenFGA.

## 3. Perímetros de publicación

La futura entrada distinguirá tres superficies sin mezclar sus políticas:

- **privada**: propietario y operadores autenticados;
- **compartida**: acceso revocable y limitado para veterinario, tienda o
  soporte, con identidad o enlace delegado según el caso;
- **pública**: proyecciones expresamente publicadas, nunca acceso directo al
  modelo privado.

Grafana, Prometheus, Alertmanager, OpenBao y las superficies administrativas
permanecen privadas. Grafana no se publicará hasta verificar OIDC, autorización,
TLS y atribución. Prometheus, Alertmanager y OpenBao no tendrán rutas públicas.

La navegación humana privada usará OIDC interactivo con redirect, callback
HTTPS y cookie segura. Las APIs para Angular, agentes y workloads usarán bearer
tokens JWT validados para emisor y audiencia concretos. Una proyección pública
no simulará un usuario autenticado.

Los enlaces compartidos usarán credenciales opacas de alta entropía,
revocables y almacenadas mediante hash. Se redactarán de logs y `Referer`; las
respuestas aplicarán `no-store` y `noindex`. La revocación inmediata y la
resistencia a enumeración formarán parte de la aceptación.

La política base será denegar. Una ruta protegida deberá declarar:

- hostname, listener, backend y clasificación;
- flujo de autenticación;
- acción y tipo de recurso para autorización;
- límites de tamaño, tiempo y tasa;
- propietario operativo y evidencia de aceptación.

## 4. Ownership y aislamiento

`GatewayClass`, controlador, `Gateway` y `SecurityPolicy` base son propiedad
exclusiva de plataforma. Las rutas HTTP/gRPC viven en el namespace del backend
y se adjuntan solo a listeners cuyo `allowedRoutes` seleccione explícitamente
ese namespace. Un `ReferenceGrant` será excepcional, vivirá en el namespace
propietario del objeto referenciado y requerirá su aprobación.

La entrada protegida inicial admite solo `HTTPRoute` y `GRPCRoute`. TCP/UDP
necesitarían otra puerta: el flujo OIDC y la autorización externa HTTP/gRPC no
protegen esas superficies de forma equivalente.

Solo plataforma puede crear o modificar `SecurityPolicy`. Una política de ruta
no podrá eliminar el baseline del `Gateway`: RBAC/admission impondrán la
estrategia de attachment y merge. Conflicto, política ausente o condición
distinta de `Accepted=True` implican fallo cerrado.

NetworkPolicy y Linkerd impedirán llegar directamente a los adaptadores
públicos. El gateway será su único origen permitido. Los backends rechazarán un
`ActorContext` falsificado, expirado, de otra audiencia o de otro entorno.

## 5. Cabeceras y trazabilidad

Envoy generará siempre `request_id`. ReefOps Authorizer validará y normalizará
un `correlation_id` externo; si falta o es inválido generará otro y registrará
su procedencia. No se usará como identidad ni como clave única.

Envoy eliminará antes de enrutar el `Authorization` externo y toda cabecera de
actor, delegación, decisión, `ActorContext`, `Forwarded`, `X-Forwarded-*`,
request/correlation/causation y trazado. Regenerará mediante configuración
central únicamente las necesarias y propagará desde Authorizer una allowlist
cerrada con formatos y tamaños limitados. El tratamiento de `traceparent` y
`baggage` será explícito y probado; no lo decide cada ruta.

Para cada petición protegida se conservarán, con redacción y retención:

- `request_id`, `correlation_id`, `causation_id`, procedencia e intento;
- listener, ruta, método y plantilla de recurso, sin secretos ni tokens;
- sujeto autenticado, delegación y audiencia;
- `authz_decision_id`, acción, recurso, decisión y motivo normalizado;
- revisión del gateway y versiones de políticas;
- backend, resultado y latencias de gateway, AuthN, AuthZ y servicio.

Access logs y trazas permiten diagnóstico, pero no sustituyen la auditoría
funcional. ReefOps Authorizer será propietario de un audit append-only para
decisiones permitidas y denegadas, incluso si no llegan al dominio, con
integridad, redacción y retención declaradas. `authz_decision_id` enlazará esas
evidencias.

Cada retry conserva correlación, causación y el `request_id` lógico, recibe un
`attempt_id` o span propio y registra número de intento y resultado. Envoy no
reintentará automáticamente operaciones no idempotentes. Una mutación
reintentable exigirá
`Idempotency-Key` end-to-end y registro durable del resultado; nunca se repetirá
una dosis, orden física o notificación por un retry de transporte.

## 6. Fundación inerte

La primera puerta instalará mediante Flux:

- chart OCI de Envoy Gateway fijado por versión y digest;
- imagen del controlador fijada por digest;
- CRD de Gateway API y extensiones de Envoy incluidas en el chart fijado;
- namespace dedicado, recursos, seguridad de pod y servicio `ClusterIP`;
- métricas internas consumibles por Prometheus;
- pruebas de ausencia de `Gateway`, `HTTPRoute`, `GRPCRoute`, `TCPRoute`,
  `UDPRoute` e `Ingress`;
- pruebas de ausencia de `LoadBalancer` y `NodePort`.

No se instala todavía el plano de datos Envoy: este solo nace al declarar un
`Gateway`. Su imagen se fijará por digest en la segunda puerta.

Actualizar CRD es una migración de API. Requiere revisar compatibilidad,
conversión, APIs almacenadas y rollback; retirar el HelmRelease no elimina
automáticamente CRD ni objetos existentes.

## 7. Criterios de aceptación de la fundación

La puerta se considera operativa solo cuando:

1. la documentación se integra antes que la infraestructura;
2. Flux aplica exactamente el commit completo de `reefops-platform` validado;
3. las CRD están `Established` y el controlador preparado;
4. chart e imagen coinciden con los digests aprobados;
5. no existen rutas, gateways, plano de datos ni servicios expuestos;
6. el controlador usa `ClusterIP`, usuario no root, seccomp, recursos limitados
   y capacidades eliminadas;
7. Prometheus descubre sus métricas sin hacerlas públicas;
8. la evidencia append-only registra `environment_id`, contexto e identidad del
   clúster, actor, operación, correlación, causación, revisión local y aplicada
   en cada reconciliación Flux, digests de chart, imagen y manifiesto, UIDs,
   tiempos, fase, error, restauración y resultado sin secretos;
9. `task validate` pasa en producto, plataforma y composición GitOps.

La evidencia se encadena con SHA-256, se conserva al menos un año y entra en el
backup cifrado operativo. SBOM, firma y procedencia se exigirán cuando chart o
imagen se repliquen en GHCR; en esta fase se verifica el digest del publicador y
queda diferida esa réplica.

## 8. Fallo cerrado de la segunda puerta

La entrada deniega sin fallback permisivo ni caché obsoleta cuando:

- en una ruta cuyo flujo exige JWT, falta el token o está expirado;
- firma, emisor, audiencia o scope son incorrectos;
- JWKS o ZITADEL no están disponibles y no hay una clave todavía válida dentro
  de la ventana aprobada;
- Authorizer u OpenFGA fallan, exceden timeout o responden mal;
- falta una decisión inequívoca de permiso;
- el `ActorContext` es inválido, está expirado o pertenece a otro entorno;
- `SecurityPolicy` no está aceptada, el `Gateway` no está `Programmed=True` o
  la ruta no tiene `Accepted=True` y `ResolvedRefs=True`.

La aceptación correlacionará cada fallo, comprobará el status esperado, la
ausencia de llamada al caso de uso y el audit de denegación. También probará
negativamente el acceso directo y las cabeceras falsificadas.

## 9. Decisiones de la segunda puerta development

La primera entrada protegida queda limitada al host operador:

1. usa `reefops.localhost` y `identity.reefops.localhost`, resueltos por el
   cliente local; no anuncia nombres en la LAN;
2. el `Gateway` conserva Service `ClusterIP` y se alcanza inicialmente mediante
   `port-forward`; publicar por LAN exige una puerta posterior;
3. cert-manager emite TLS desde una CA development separada cuya confianza se
   instala solo en clientes autorizados;
4. el listener HTTPS vive en `reefops-gateway-system` y admite rutas únicamente
   de namespaces etiquetados `reefops.io/gateway-access=protected`;
5. ZITADEL usa Authorization Code con PKCE para Angular y audiencia distinta
   para APIs; no se habilita password grant;
6. ReefOps Authorizer implementa Envoy `ext_authz` v3 mediante gRPC, recibe una
   allowlist cerrada de metadata de ruta y devuelve únicamente decisión, sujeto y un
   `ActorContext` firmado de vida corta;
7. development aplica límites locales en Envoy sin un almacén distribuido; no
   despliega el rate-limit latente hasta necesitar cuotas compartidas;
8. access logs se redactan y escriben a stdout; el audit funcional append-only
   pertenece al Authorizer y se conservará en PostgreSQL;
9. Envoy Proxy queda fijado por digest y elimina `Authorization`, cabeceras de
   actor y forwarded entrantes antes de regenerar su allowlist;
10. `correlation_id` debe ser UUID, se genera si falta, y no hay retries
    automáticos para métodos mutadores.

La primera prueba de la segunda puerta usa exclusivamente
`acceptance.synthetic.resource.view` y no publica un caso de uso de dominio. La
metadata confiable aporta sólo el identificador de contrato y los claims ya
validados; acción, relación, plantilla y extracción del recurso se resuelven en
la allowlist compilada del Authorizer. La respuesta ALLOW no se emite hasta que
el Check OpenFGA, la firma Ed25519 y el append de auditoría han terminado.

ZITADEL, OpenFGA y sus consolas no reciben rutas administrativas genéricas. La
ruta de identidad expone únicamente lo requerido por OIDC y login. OpenFGA no
tiene ruta Gateway.

Grafana sigue usando un `port-forward` efímero autenticado por Kubernetes; no
forma parte de esta puerta.
