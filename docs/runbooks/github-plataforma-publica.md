# Publicación segura del catálogo de plataforma

## Objetivo

Publicar únicamente `reefops-platform`, mantener privados producto y GitOps y
retirar la credencial que deja de ser necesaria. El propietario operativo es
plataforma; GitHub conserva colaboración y Flux conserva reconciliación
pull-only.

## Precondiciones y orden

1. Auditar el árbol y todo el historial de plataforma con Gitleaks y Trivy.
2. Verificar licencia Apache-2.0 y ausencia de `.sops.yaml`, `clusters/`,
   secretos y configuración de instalaciones.
3. Confirmar que `reefops-gitops` fija exactamente el commit auditado.
4. Hacer público exclusivamente `reefops-platform`.
5. Promover en GitOps la fuente HTTPS pública sin `secretRef`.
6. Esperar a que `GitRepository/reefops-platform` esté `Ready` en ese commit.
7. Revocar la deploy key antigua, eliminar `platform-git-auth` y destruir la
   clave local obsoleta.
8. Aplicar protección de `main` y verificar la visibilidad de los tres
   repositorios.

La limpieza de credenciales no se ejecuta antes de verificar HTTPS para evitar
interrumpir la reconciliación. La operación es idempotente y no imprime claves.

## Fallo y recuperación

Si falla la publicación antes del paso 5, la fuente SSH existente continúa
funcionando. Si falla HTTPS, se conserva la deploy key y se revierte el commit
GitOps.

Una vuelta futura a repositorio privado exige crear una credencial nueva,
independiente y de solo lectura, promocionar primero la fuente autenticada y
verificarla antes de cambiar visibilidad. La clave retirada no se recupera ni
se reutiliza.

El cambio de visibilidad y las revocaciones quedan registrados por GitHub. La
promoción conserva commit y estado en Flux; los comandos operativos de limpieza
se ejecutan con actor autenticado y objetivos exactos.
