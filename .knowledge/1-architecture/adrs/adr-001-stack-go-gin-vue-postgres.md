---
okf_version: "1.0"
type: ADR
title: "ADR-001: Stack Go/Gin + Vue embebido + PostgreSQL"
description: "Documenta el stack técnico real de SALVIA: backend Go con Gin y GORM, frontend Vue embebido en el binario, PostgreSQL, y autenticación por sesión server-side."
owner: "@architecture-guild"
status: accepted
tags: [architecture, adr, stack, go, vue, postgres]
last_updated: "2026-07-02"
---

# ADR-001: Stack Go/Gin + Vue embebido + PostgreSQL

> [!IMPORTANT]
> Este ADR **documenta una decisión ya materializada en el código** (ingeniería inversa, no una decisión nueva). Aprobado por el equipo el 2026-07-02 (`status: accepted`).

## Contexto

SALVIA es un sistema de gestión y atención psicosocial. Necesita servir una UI web y una API REST, persistir datos relacionales (casos, remisiones, seguimientos, barreras, usuarios) y desplegarse de forma sencilla en servidores Windows self-hosted. El stack observado en el repositorio es:

- **Backend**: Go 1.24, módulo `bitsflow`, framework **Gin** (`gin-gonic/gin`).
- **ORM / acceso a datos**: **GORM** (`gorm.io/driver/postgres`) para el dominio nuevo (`internal/repository`, `salvia/`) y **pgx** (`jackc/pgx`) para el código legacy (`security/dao`). Migraciones vía GORM `AutoMigrate` + SQL en `src/migrations/`.
- **Frontend**: **Vue 3** cargado en runtime con `vue3-sfc-loader` (sin paso de build, sin `package.json`), sobre Bootstrap 4 / AdminLTE y DataTables. Se **embebe en el binario** con `go:embed frontend/*`.
- **Autenticación**: **sesión server-side** con `gin-contrib/sessions` + `gorilla/securecookie` (cookie firmada). **No** se usa JWT.
- **Base de datos**: **PostgreSQL** (Supabase en develop; self-hosted en producción).
- **Empaque**: un **único binario** (`salvia.exe`) que sirve UI + API; corre como servicio Windows `SalviaApp`.

## Opciones Evaluadas

| Opción | Pros | Contras |
| :--- | :--- | :--- |
| Monolito Go/Gin con Vue embebido (elegida) | Un solo artefacto para desplegar; sin toolchain de Node en el server; despliegue trivial en Windows | UI acoplada al ciclo de build del backend; recompilar para cualquier cambio de front |
| SPA Vue con build (Vite) + API Go separada | Front desacoplado, HMR, ecosistema de build moderno | Requiere Node en CI/deploy; dos artefactos y CORS; más piezas móviles |
| JWT stateless en vez de sesión server-side | Escala horizontal sin estado compartido | Revocación compleja; no aporta valor para un despliegue single-node |

## Decisión

Hemos decidido **mantener el monolito Go/Gin con frontend Vue embebido vía `go:embed`, GORM+pgx sobre PostgreSQL, y autenticación por sesión server-side**, por ser lo ya implementado y lo mejor alineado con el despliegue single-node en Windows self-hosted.

## Justificación

El modelo de despliegue (un servicio Windows por entorno, sin orquestador) se beneficia de un artefacto único sin dependencias de runtime. La sesión server-side es suficiente para un backend de instancia única y evita la complejidad de revocación de JWT. La coexistencia pgx/GORM refleja una migración incremental del código legacy al nuevo dominio, no una decisión de green-field.

## Consecuencias (Trade-offs)

- **Positivo**: despliegue simple (copiar un `.exe`); sin toolchain de Node en el servidor; un solo proceso que operar y monitorear.
- **Negativo**: todo cambio de UI exige recompilar y redesplegar el binario; el escalado horizontal exigiría externalizar el store de sesiones; dos capas de acceso a datos (pgx + GORM) aumentan la carga cognitiva hasta completar la migración.

## Referencias

- [Diagrama de Contenedores (C4 L2)](/.knowledge/1-architecture/container-diagram.md)
- [Matriz de Entornos](/.knowledge/4-operations/environment-matrix.md)
- `src/main.go`, `src/go.mod`, `src/config/db_config.json`
