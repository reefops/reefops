# Runbook de ESO con OpenBao

## Objetivo

Desplegar y aceptar External Secrets Operator —ESO— como adaptador declarativo
de mínimo privilegio sobre la autoridad OpenBao de development. Este
procedimiento no crea otra autoridad ni entrega secretos de dominio.

## Precondiciones

- contexto Kubernetes actual `docker-desktop`;
- OpenBao activo, inicializado, no sellado y con auditoría;
- backup y ensayo de recuperación de OpenBao completados;
- repositorios `reefops-platform` y `reefops-gitops` limpios;
- commit de plataforma validado y publicado;
- token raíz original disponible solo durante las dos ceremonias que lo
  requieren.

No se generan claves nuevas. El token se introduce mediante lectura silenciosa,
se exporta únicamente para el proceso hijo y se elimina al terminar.

## Secuencia limpia

### 1. Validar y configurar OpenBao

Desde `reefops-platform` en el commit que contiene la integración:

```sh
task validate
read -s "BAO_TOKEN?Token raíz de OpenBao: "
echo
export BAO_TOKEN
task openbao-configure
unset BAO_TOKEN
```

El wrapper fija contexto, namespace, port-forward, CA y SNI; compara el
`cluster_id` del endpoint con el del pod activo y no acepta un endpoint
aportado por el operador.

### 2. Promover mediante GitOps

Un PR en `reefops-gitops` actualiza
`clusters/local/workloads/platform-source.yaml` al SHA completo de plataforma.
Flux obtiene ese commit mediante pull y aplica, en orden:

1. configuración de trust y red;
2. controlador ESO namespaced;
3. `SecretStore`, identidad y secreto sintético.

GitHub Actions no se conecta al clúster y no existe despliegue push.

### 3. Verificar la integración

Cuando todas las Kustomizations estén preparadas:

```sh
read -s "BAO_TOKEN?Token raíz de OpenBao: "
echo
export BAO_TOKEN
task openbao-verify-eso
unset BAO_TOKEN
```

La prueba:

1. valida contexto, entorno, CA, SNI y `cluster_id`;
2. comprueba el token raíz sin imprimirlo;
3. espera salud de ESO, `SecretStore` y `ExternalSecret`;
4. captura política y valor sintético;
5. rota el valor y comprueba el refresco;
6. aplica temporalmente una política `deny`;
7. confirma fallo cerrado y que no llega un valor bloqueado;
8. restaura y vuelve a comprobar política y valor originales;
9. exige deltas positivos de login y lectura en el audit device.

La salida de aceptación es:

```text
ESO y OpenBao verificados con refresco, revocación, auditoría y recuperación.
```

## Evidencia

La evidencia local no sensible se añade a:

```text
~/.local/state/reefops/eso-openbao/operations.jsonl
```

Una operación satisfactoria incluye `result=success`,
`restoration=success`, deltas de auditoría positivos, revisiones local y Flux
idénticas, UIDs y generaciones, actor, `cluster_id`, correlación y causación.
El fichero local es evidencia operativa mutable, no almacenamiento WORM.

## Fallo seguro y errores resueltos

El trap restaura política y valor si la prueba ya había comenzado a mutar el
estado. Un fallo de restauración cambia el resultado a error y debe bloquear
consumidores nuevos hasta revisión manual. Cada error muestra fase y ubicación
de evidencia sin exponer credenciales.

Quedan protegidas mediante validaciones ejecutables estas regresiones reales:

- **entorno rechazado incorrectamente**: los wrappers usan el JSONPath exacto
  de `reefops.io/environment`; una discrepancia real sigue fallando cerrada;
- **script Bash no ejecutable**: `task validate` ejecuta `bash -n` sobre todos
  los scripts y tests antes de ShellCheck;
- **contrato JSON de OpenBao**: `bao policy read -format=json` se consume desde
  el campo `policy`, tanto al capturar como al comprobar restauración;
- **error opaco**: la evidencia y stderr identifican la fase fallida.

Si un intento falla antes de `state-capture`, no existe estado que restaurar.
Si falla después, debe aparecer `restoration=success`; cualquier otro valor
requiere detenerse y revisar la autoridad antes de repetir.

## Criterios de cierre

- todas las Kustomizations Flux están `Ready=True`;
- `SecretStore/openbao` y `ExternalSecret/openbao-smoke-test` están preparados;
- no existen objetos RBAC cluster-wide de ESO;
- chart, imagen y commit están fijados por digest/SHA;
- la última evidencia tiene resultado y restauración satisfactorios;
- `task validate` pasa en producto, plataforma y GitOps.
