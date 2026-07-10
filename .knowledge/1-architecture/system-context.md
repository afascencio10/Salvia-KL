---
okf_version: "1.0"
type: ADR
title: "Diagrama de Contexto del Sistema (C4 Level 1)"
description: "Vista de alto nivel de SALVIA: actores que gestionan la atención psicosocial y los sistemas externos con los que se integra."
owner: "@architecture-guild"
status: active
tags: [architecture, c4, diagram, system]
last_updated: "2026-07-02"
---

# Diagrama de Contexto del Sistema (C4 — Level 1)

## Descripción

SALVIA es un sistema de **gestión y atención psicosocial**: registra casos, remisiones, seguimientos, barreras y duplas de atención. Este diagrama muestra quién interactúa con el sistema y con qué sistemas externos se integra.

## Diagrama

```mermaid
C4Context
    title SALVIA — Diagrama de Contexto

    Person(profesional, "Profesional psicosocial", "Registra y da seguimiento a casos, remisiones y barreras")
    Person(admin, "Administrador", "Gestiona usuarios, entidades, sedes y configuración del sistema")

    System(salvia, "SALVIA", "Aplicación web monolítica (Go/Gin + Vue embebido) para la gestión de atención psicosocial")

    System_Ext(postgres, "PostgreSQL", "Base de datos relacional. Supabase en develop; self-hosted (localhost) en producción")
    System_Ext(github, "GitHub Actions", "CI/CD: build del binario y despliegue automático a los runners self-hosted")

    Rel(profesional, salvia, "Gestiona casos y seguimientos", "HTTPS")
    Rel(admin, salvia, "Administra", "HTTPS")
    Rel(salvia, postgres, "Lee/Escribe (pgx + GORM)", "TCP/5432 SSL")
    Rel(github, salvia, "Compila y despliega el binario", "Self-hosted runner")
```

## Actores

| Actor | Descripción | Autenticación |
| :--- | :--- | :--- |
| Profesional psicosocial | Usuario operativo que gestiona casos, remisiones, seguimientos y barreras | Usuario/contraseña + sesión server-side |
| Administrador | Gestiona usuarios, entidades/sedes, roles y configuración | Usuario/contraseña + sesión server-side |

> [!NOTE]
> El detalle de roles y permisos por recurso es una decisión de negocio pendiente de documentar en la [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md).

## Sistemas Externos

| Sistema | Propósito | Protocolo | Notas |
| :--- | :--- | :--- | :--- |
| PostgreSQL | Persistencia principal | TCP 5432 (SSL) | Supabase (pooler) en `develop`; PostgreSQL local en `production` |
| GitHub Actions | Build y despliegue automático | Self-hosted runners (Windows) | Credenciales de BD inyectadas desde GitHub Secrets en el deploy |

## Fronteras del Sistema

> [!IMPORTANT]
> SALVIA es un **monolito auto-contenido**: un único binario Go (`salvia.exe`) sirve tanto la SPA de Vue (embebida vía `go:embed`) como la API REST bajo `/api/v1`. No hay microservicios ni gateway separado. El detalle técnico está en el [Diagrama de Contenedores](/.knowledge/1-architecture/container-diagram.md).
