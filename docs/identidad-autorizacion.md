# Identidad y autorización

## 1. Decisión

ReefOps adopta esta combinación para la primera entrada protegida:

- ZITADEL como proveedor de identidad OIDC;
- OpenFGA como almacén y motor de relaciones autorizables;
- ReefOps Authorizer como adaptador `ext_authz`, emisor de `ActorContext` y
  propietario de la auditoría de decisiones;
- Envoy Gateway como punto de aplicación norte-sur;
- Linkerd para identidad y mTLS de workloads, sin sustituir AuthN ni AuthZ.

La decisión se toma antes de crear datos propios de ZITADEL u OpenFGA. Cambiar
después de ese punto exigiría migrar sujetos, clientes, relaciones y evidencias
de auditoría.

## 2. Alternativas descartadas

| Responsabilidad | Alternativa | Motivo para no elegirla ahora |
| --- | --- | --- |
| AuthN | Keycloak | Es una alternativa válida y más madura, pero añade un runtime Java y una superficie mayor; el modelo organizativo central y el runtime Go de ZITADEL encajan mejor. |
| AuthN | Authentik | Está orientado especialmente a SSO y proxy de aplicaciones; las organizaciones del producto no son su contrato central. |
| AuthN | Ory | Obliga a componer identidad, OAuth/OIDC y UI de login/consentimiento, aumentando servicios y código propio. |
| AuthZ | SpiceDB | Su control explícito de consistencia es potente, pero ZedTokens, operador y tuning añaden complejidad que el volumen inicial no justifica. |
| AuthZ | Cerbos u OPA | Evalúan bien políticas y atributos, pero no sustituyen de forma natural el grafo persistente de compartición y relaciones por recurso. |
| Entrada | Kong, Traefik o APISIX | No aportan una mejora que compense retirar Envoy Gateway, ya operativo y compatible con Gateway API, JWT y autorización externa. |

La decisión se revisará si aparecen LDAP/Kerberos y administración IAM
tradicional, consistencia causal de autorización a gran escala o requisitos de
APIM que Envoy Gateway no cubra.

## 3. Fuentes de verdad

| Dato o decisión | Propietario |
| --- | --- |
| Credenciales, MFA, sesiones y recuperación | ZITADEL |
| Sujeto OIDC y organización IAM | ZITADEL, proyectado por Identity |
| Organización funcional, miembros y vínculo con ID IAM externo | Identity de ReefOps |
| Relación autorizable sobre un recurso ReefOps | OpenFGA |
| Acción y recurso derivados de una ruta | ReefOps Authorizer mediante metadata allowlisted |
| Contraseña/código, contador, caducidad y revocación de un enlace | Sharing |
| Selección, redacción y snapshot de una vista compartida | Sharing |
| Snapshot público, indexación y retirada | Publication |
| Invariantes y aceptación de comandos | Dominio propietario |
| Auditoría de la decisión norte-sur | ReefOps Authorizer en PostgreSQL |

OpenFGA decide relaciones y permisos. No almacena secretos de enlace, no
incrementa contadores, no redacta contenido y no convierte una autorización en
un comando de dominio válido.

## 4. Flujo norte-sur

Angular utiliza Authorization Code con PKCE directamente contra ZITADEL. Envoy
Gateway no mantiene una sesión OIDC para la SPA:

1. Envoy valida firma, `issuer`, `audience`, expiración y scopes del access
   token.
2. Elimina `Authorization`, cabeceras de actor y forwarded no confiables antes
   de generar la allowlist interna.
3. Envía al Authorizer sujeto, organización activa, método, plantilla de ruta,
   acción declarada, tipo e identificador de recurso y `correlation_id`.
4. Authorizer valida formatos, consulta OpenFGA con store y modelo fijados y
   deniega ante timeout, error o respuesta incompleta.
5. Si permite, firma un `ActorContext` de vida corta que incluye al menos
   `actor_id`, `subject_id`, `organization_id`, delegación, acción, recurso,
   `decision_id`, `correlation_id`, `issued_at` y `expires_at`.
6. Envoy propaga únicamente el contexto firmado y las cabeceras allowlisted.
7. El dominio verifica el contexto y aplica además sus invariantes.

Los endpoints de descubrimiento, autorización, token y login que necesite
ZITADEL tendrán rutas explícitas. Su consola administrativa y OpenFGA no tendrán
una ruta genérica ni pública.

## 5. Organización activa y sujeto

ZITADEL administra el contenedor IAM usado para login y delegación. Identity de
ReefOps es propietario de la organización funcional y conserva el ID externo
de ZITADEL. Los eventos de provisión mantienen el mapeo; ningún dominio usa un
ID de ZITADEL como agregado de negocio.

Un usuario puede pertenecer a varias organizaciones funcionales. El token
identifica al sujeto, pero cada petición ReefOps selecciona exactamente una
organización activa:

- la organización debe estar presente en claims confiables o en un parámetro
  de sesión firmado;
- Authorizer comprueba la pertenencia antes de evaluar el recurso;
- el recurso debe pertenecer a esa misma organización, salvo un flujo de
  compartición explícito;
- no se agregan silenciosamente los permisos de todas las organizaciones del
  usuario;
- el `ActorContext` siempre contiene la organización efectiva.

Esto evita que una sesión multi-organización actúe accidentalmente sobre el
tenant incorrecto.

## 6. Modelo OpenFGA inicial

El primer modelo expresará tipos equivalentes a estos:

```text
user
organization
installation
system
resource
share
share_principal
service_account
```

Relaciones mínimas:

- `organization`: `owner`, `administrator`, `member`;
- `installation`: `organization`, `owner`, `administrator`, `caretaker`,
  `technician`, `observer`;
- `system`: `installation` y excepciones directas por rol;
- `resource`: `parent`, `inherits_view_from`, `viewer`, `contributor`;
- `share`: `target`, `grantee`, `bearer`, `viewer`, `contributor`;
- `service_account`: capacidades explícitas por recurso, nunca equivalencia
  automática con un usuario humano.

Los permisos iniciales serán aditivos y específicos, por ejemplo `can_view`,
`can_measure`, `can_operate`, `can_comment`, `can_upload` y `can_administer`.
No existirá un permiso genérico `manage` que mezcle lectura con acciones
físicas.

El modelo ejecutable inicial está en
[`authorization/reefops.fga`](../authorization/reefops.fga) y sus casos en
[`authorization/reefops.fga.yaml`](../authorization/reefops.fga.yaml). La CLI
oficial los ejecuta como parte de `task validate`. El fixture representa la
organización activa como una tupla limitada al caso porque el formato de test
de la CLI no acepta todavía contextual tuples; Authorizer deberá enviarla como
contextual tuple y nunca persistirla.

### 6.1 Herencia restringible

RF-072 permite heredar la visibilidad del sistema y restringirla en un recurso
concreto. Por ello, `resource.can_view` no será un `viewer OR can_view from
parent` incondicional.

La relación `inherits_view_from: [system]` será explícita. Crear un recurso con
herencia escribe esa relación; hacerlo privado la elimina. Así se puede cortar
la herencia sin introducir un `deny` global y sin hacer visible un recurso
privado por estar enlazado desde otro objeto.

### 6.2 Compartición temporal

Sharing crea un `share` con alcance cerrado. Para una persona autenticada usa
un `grantee`; para una URL anónima, Authorizer deriva un `share_principal` del
token opaco después de validar su hash.

- La condición temporal de OpenFGA puede limitar el tuple.
- Sharing comprueba además estado, contraseña/código y número de accesos dentro
  de una transacción propia.
- La vista o snapshot limita qué campos existen; OpenFGA no redacta objetos.
- Revocar elimina o invalida primero la capacidad en Sharing y su tuple; desde
  ese momento Authorizer falla cerrado.
- La caché de decisiones positivas no podrá sobrevivir a una revocación. En la
  primera versión no habrá caché positiva en Authorizer.

## 7. Versionado y escritura

- El store ID y `authorization_model_id` estarán fijados en configuración
  entregada por secretos/configuración GitOps.
- Cada modelo será inmutable, revisado y probado antes de activarse.
- Los dominios no escribirán OpenFGA directamente.
- Consumidores idempotentes proyectarán eventos de membresía, propiedad,
  compartición y revocación a tuples.
- Cada cambio conservará `correlation_id`, `causation_id`, actor, versión de
  modelo y resultado.
- La reconciliación tendrá una vía de reparación desde eventos/proyecciones,
  pero nunca desde consultas cruzadas a tablas internas de otro dominio.

## 8. Contextos síncronos y capacidades asíncronas

El `ActorContext` del Gateway solo autoriza la petición HTTP inmediata. No se
copiará como credencial reutilizable a NATS, jobs o workflows.

Una operación asíncrona sensible llevará una capacidad delegada distinta,
firmada, con audiencia de consumidor, acción y recurso cerrados, actor y
delegación, `correlation_id`, `causation_id`, emisión, expiración corta e ID de
decisión/modelo. El consumidor valida la capacidad y sus invariantes antes de
producir efectos. Reintentos y replay no amplían su alcance ni repiten una
acción física.

## 9. Prueba de arquitectura obligatoria

Antes de abrir el Gateway se probarán como mínimo estos casos:

1. un usuario pertenece a dos organizaciones y la organización activa impide
   un acceso cruzado accidental;
2. un miembro autorizado hereda acceso desde instalación a sistema y recurso;
3. un recurso elimina `inherits_view_from` y deja de ser visible pese al acceso
   sobre su padre;
4. un enlace válido permite solo el alcance seleccionado y caduca;
5. la revocación de un enlace deniega la siguiente petición;
6. OpenFGA o Authorizer no disponibles producen denegación;
7. un identificador malformado o metadata de ruta no allowlisted se deniega;
8. ninguna consola administrativa es alcanzable mediante Gateway.

La prueba congelará además el vocabulario inicial ruta–acción–tipo de recurso.
Cada ruta tendrá exactamente una extracción determinista; una acción genérica
no podrá sustituir permisos distintos como `view`, `comment`, `upload`,
`propose`, `operate`, `download` o `administer`.

Superar estas pruebas cierra el diseño. No autoriza todavía acciones físicas:
esas operaciones requerirán permisos separados, procedimiento, idempotencia y
las invariantes del dominio.
