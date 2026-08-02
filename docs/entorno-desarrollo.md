# Entorno de desarrollo

## 1. Política

Las herramientas del host se declararán en `Brewfile` y se instalarán con
`brew bundle`. El archivo expresa el conjunto requerido; las versiones exactas
de aplicaciones, imágenes y dependencias se fijarán dentro de cada workspace o
manifiesto.

No se instalarán globalmente Angular CLI, librerías Go ni paquetes Python de la
aplicación:

- Angular y TypeScript serán dependencias versionadas del workspace;
- Go utilizará módulos y herramientas fijadas por el repositorio;
- Python utilizará `uv`, `pyproject.toml` y entornos virtuales;
- los generadores de contratos se ejecutarán mediante tareas versionadas.

## 2. Toolchains

| Área | Herramienta del host | Uso |
|---|---|---|
| Backend | Go 1.26 | Core, gateways, workers y servicios de red |
| IA y visión | Python 3.12 + uv | Entornos reproducibles y compatibilidad ML |
| Frontend | Node.js 24 LTS + pnpm | Angular y herramientas TypeScript |
| Contratos y datos | Protobuf, Buf, sqlc, Goose | Eventos, APIs técnicas, acceso SQL tipado y migraciones PostgreSQL |
| Contenedores | Docker Desktop | Build y Kubernetes local |
| Kubernetes | kubectl 1.34, Helm, Flux, Kustomize integrado | Cliente alineado con el clúster local 1.34 |
| Secretos | OpenBao CLI, SOPS + age | Autoridad local y bootstrap cifrado |
| Calidad | golangci-lint, gofumpt, ShellCheck, pre-commit, ripgrep | Comprobaciones locales y validadores |
| Manifiestos | yq + kubeconform | Transformación y validación Kubernetes |
| Medios | FFmpeg | Captura y transformación local |
| GitHub | gh, actionlint | Repositorios, PR, Actions y validación de workflows |
| Cadena de suministro | cosign, syft, trivy, gitleaks, oras | Firma, SBOM, escaneo de árbol e historial y OCI |

Python 3.12 será la versión base inicial para maximizar compatibilidad con
librerías de visión y ML. La existencia de otro Python del sistema no cambia el
runtime del proyecto.

## 3. Instalación

```sh
brew bundle
task tools
task validate
```

`task tools` mostrará versiones y fallará si no puede resolver la toolchain
esperada.

## 4. Versionado y actualizaciones

- `Brewfile` no autoriza actualizaciones automáticas de dependencias del
  proyecto.
- Una actualización de Go, Python o Node requiere documentar compatibilidad,
  ejecutar tests y actualizar los archivos de versión del workspace.
- `kubectl` se mantendrá en la misma minor que el clúster local; actualizar
  Docker Desktop Kubernetes exige actualizar esta fijación y validar Flux.
- CI utilizará versiones explícitas, no el estado accidental del Mac.
- No se añadirá Rust, Java, CUDA ni otra toolchain hasta que un módulo
  documentado la necesite.
- Las herramientas retiradas se eliminarán del `Brewfile` después de comprobar
  que ningún workflow las utiliza.

Goose queda adoptado para migraciones PostgreSQL. La versión del host es una
ayuda de desarrollo; el módulo Go y la imagen migradora fijarán su dependencia
exacta. `sqlc` genera consultas tipadas, pero no aplica migraciones.

Los modelos OpenFGA mantienen un ciclo separado mediante FGA CLI, artefactos
OCI firmados y promoción de IDs en OpenBao.

## 5. Herramientas todavía no instaladas

- Angular CLI: se añadirá localmente al crear el workspace Angular.
- Dependencias Python de IA: se añadirán al crear el servicio y su
  `pyproject.toml`.
- Generador OpenAPI para Go: se fijará como tool del módulo Go; no se instalará
  un generador Java global.
