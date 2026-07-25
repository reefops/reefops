# Publicación del repositorio de producto

## Objetivo

Reemplazar el historial privado transitorio de `reefops` por un único commit
raíz publicable. El repositorio resultante contiene producto, requisitos,
arquitectura, contratos y controles de ingeniería, pero no estado deseado de
clúster, material SOPS ni copias de plataforma.

## Contenido excluido

- `.sops.yaml`;
- todo `infra/`;
- configuración concreta de clústeres o instalaciones;
- claves, secretos, backups y datos de acuarios;
- skills y agentes dedicados exclusivamente a operar GitOps.

Los componentes reutilizables permanecen en `reefops-platform`; composición,
versiones desplegadas y topología permanecen en `reefops-gitops`.

## Procedimiento

1. Actualizar esta decisión y retirar del árbol todo contenido excluido.
2. Ejecutar validación, Gitleaks, Trivy y comprobación de licencia.
3. Crear fuera del workspace un bundle de todos los refs, cifrado con la
   identidad `age` local. No conservar un bundle en claro.
4. Crear un commit huérfano con el árbol validado.
5. Confirmar que `main` tiene un solo commit raíz y nunca contiene los paths
   excluidos. Una rama automática creada durante la ventana solo es admisible
   si desciende de ese root y supera los mismos escaneos.
6. Reemplazar `main` remoto mediante force-with-lease mientras el repositorio
   sigue privado y sin protección de rama.
7. Comprobar por API un SHA antiguo. Si GitHub todavía sirve el objeto
   huérfano, verificar que no existan issues, PRs, releases o artefactos y
   borrar/recrear únicamente el repositorio de producto con el mismo nombre.
8. Eliminar refs y reflogs locales antiguos después de verificar el bundle.
9. Cambiar a público, habilitar seguridad gratuita y proteger `main`.

## Recuperación

El force-push elimina las referencias remotas al historial anterior, pero no
es un mecanismo de restauración de datos. El bundle cifrado permite reconstruir
el Git privado con la identidad `age` local. Restaurarlo exige crear un
repositorio privado separado; nunca se importará en el remoto público.

GitHub puede retener objetos inaccesibles después de un force-push. Si un SHA
antiguo sigue siendo consultable, la publicación no se considera terminada:
se recrea el repositorio para obtener un nuevo identificador y almacenamiento
sin esos objetos. Como el historial previo no contenía secretos detectados, no
es necesario un procedimiento de revocación adicional.
