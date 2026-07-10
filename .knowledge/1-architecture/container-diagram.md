---
okf_version: "1.0"
type: ADR
title: "Diagrama de Contenedores (C4 Level 2)"
description: "Vista de los contenedores técnicos de SALVIA: binario Go/Gin con Vue embebido y la base de datos PostgreSQL."
owner: "@architecture-guild"
status: active
tags: [architecture, c4, diagram, containers]
last_updated: "2026-07-02"
---

# Diagrama de Contenedores (C4 — Level 2)

## Descripción

SALVIA se despliega como un **único binario** (`salvia.exe`) que embebe la SPA de Vue y expone la API REST. Este diagrama muestra sus componentes internos y su relación con la base de datos.

## Diagrama

```mermaid
C4Container
    title SALVIA — Diagrama de Contenedores

    Person(user, "Profesional / Administrador", "Accede vía navegador web")

    System_Boundary(salvia, "SALVIA (salvia.exe — servicio Windows 'SalviaApp')") {
        Container(spa, "SPA Frontend", "Vue 3 (vue3-sfc-loader), Bootstrap 4, DataTables", "UI de casos, remisiones y seguimientos. Servida como assets embebidos (go:embed)")
        Container(api, "API REST", "Go, Gin — /api/v1", "Controllers, servicios y repositorios. Middleware de sesión y RBAC")
        Container(session, "Gestor de Sesión", "gin-contrib/sessions + gorilla/securecookie", "Sesiones server-side por cookie firmada (no JWT)")
    }

    ContainerDb(database, "PostgreSQL", "pgx (legacy) + GORM (nuevo dominio)", "Persistencia de casos, usuarios, remisiones, seguimientos, barreras")

    Rel(user, spa, "Usa", "HTTPS/TLS")
    Rel(spa, api, "Llama endpoints", "HTTPS/JSON")
    Rel(api, session, "Valida sesión", "In-process")
    Rel(api, database, "Lee/Escribe", "TCP/5432 SSL")
```

## Contenedores

| Contenedor | Tecnología | Responsabilidad | Puerto |
| :--- | :--- | :--- | :---: |
| SPA Frontend | Vue 3 + `vue3-sfc-loader` + Bootstrap 4 / DataTables | UI, componentes, consumo de API. Sin paso de build (SFC cargadas en runtime) | Servida por el binario |
| API REST | Go 1.24 + Gin | Routing `/api/v1`, controllers, servicios, repositorios | TLS (certs) o `PORT` env |
| Gestor de Sesión | gin-contrib/sessions + gorilla/securecookie | Autenticación por sesión server-side (cookie firmada) | In-process |
| PostgreSQL | pgx v4/v5 + GORM (`driver/postgres`) | Persistencia relacional; migraciones vía GORM AutoMigrate + SQL en `src/migrations/` | 5432 |

## Comunicación entre Contenedores

| Origen → Destino | Protocolo | Autenticación |
| :--- | :--- | :--- |
| Navegador → SPA/API | HTTPS/TLS (`certs/salvia.crt`) | Cookie de sesión (HttpOnly, firmada) |
| API → PostgreSQL | TCP 5432 con `sslmode=require` | Connection string desde `config/db_config.json` |

> [!NOTE]
> El acceso a datos convive en dos capas: consultas directas con **pgx** (código legacy bajo `security/dao`, usuario `salvia_legacy`) y **GORM** para el dominio nuevo (`internal/repository`, `salvia/`, usuario `salvia_gorm`). Ver [ADR-001](/.knowledge/1-architecture/adrs/adr-001-stack-go-gin-vue-postgres.md).

> [!WARNING]
> El frontend Vue va **embebido en el binario** (`go:embed frontend/*`): cualquier cambio de UI requiere recompilar y redesplegar `salvia.exe`. No hay servidor de assets independiente.
