---
okf_version: "1.0"
type: Project
title: "SALVIA - Índice OKF"
description: "SALVIA — Sistema de gestión y atención psicosocial (casos, remisiones, seguimientos, barreras). Índice principal de documentación OKF y mapa de navegación del repositorio."
owner: "@platform-team"
status: active
tags: [project, index, okf, navigation]
last_updated: "2026-07-02"
---

# Bienvenido al Repositorio

Esta documentación está estructurada bajo el **Open Knowledge Format (OKF)** y validada automáticamente por CI/CD y git hooks (`make hooks`).

## Mapa de Navegación

### 📋 0 — Producto
- [Glosario](/.knowledge/0-product/glossary.md) — Definiciones estándar de términos
- [PRD: Reporte consolidado de seguimientos](/.knowledge/0-product/prd-reporte-consolidado-seguimientos.md)

### 🏗️ 1 — Arquitectura
- [Diagrama de Contexto (C4 L1)](/.knowledge/1-architecture/system-context.md) — SALVIA como caja negra: actores y sistemas externos
- [Diagrama de Contenedores (C4 L2)](/.knowledge/1-architecture/container-diagram.md) — Monolito Go/Gin con Vue embebido y PostgreSQL
- **ADRs (Decisiones Arquitectónicas)**: usa la plantilla `.templates/adr/_adr.md.tmpl`
  - [ADR-001: Stack Go/Gin + Vue embebido + PostgreSQL](/.knowledge/1-architecture/adrs/adr-001-stack-go-gin-vue-postgres.md) — `accepted`

### 🗄️ 2 — Data Dictionary
- [Índice de Modelos](/.knowledge/2-data-dictionary/index.md) — Inventario de entidades de datos
- [Catálogo de Errores](/.knowledge/2-data-dictionary/error-catalog.md) — Códigos de error estándar

### ⚙️ 3 — Features
- **[Reporte consolidado de seguimientos](/.knowledge/3-features/ReporteConsolidadoSeguimientos/index.md)** — Exporta a Excel el historial de seguimientos de casos por rango de fechas (`draft`)
- **[Autenticación y Redirección de Sesión (Auth)](/.knowledge/3-features/Auth/index.md)** — Autenticación y navegación por roles
- **[Consulta de Caso para Agente Psicosocial](/.knowledge/3-features/ConsultaCasoPsicosocial/index.md)** — Acceso de solo lectura de `ps`/`ts` al detalle del caso vía "Ver caso" (`draft`)
- **[Flujo 3x3 de Atención Psicosocial](/.knowledge/3-features/AtencionPsicosocial3x3/index.md)** — Gestión de contacto 3x3 psicosocial: intentos (`contact_attempts`), consentimiento, agendamiento de sesión y cierre por imposibilidad (`draft`)

### 🚀 4 — Operations
- [Estrategia de Despliegue](/.knowledge/4-operations/deployment-strategy.md) — Blue/green, rollback
- [Matriz de Entornos](/.knowledge/4-operations/environment-matrix.md) — Dev/Staging/Prod
- [SLA/SLO](/.knowledge/4-operations/sla-slo.md) — Service Level Objectives
<!-- Añade aquí los runbooks: [Runbook: <Incidente>](/.knowledge/4-operations/runbook-<slug>.md) -->

### 📐 5 — Estándares
- [Definition of Done](/.knowledge/5-standards/definition-of-done.md) — Criterios de completitud
- [Coding Standards](/.knowledge/5-standards/coding-standards.md) — Convenciones de código
- [Testing Strategy](/.knowledge/5-standards/testing-strategy.md) — Pirámide de testing
- [Commit Convention](/.knowledge/5-standards/commit-convention.md) — Conventional Commits
- [Code Review Checklist](/.knowledge/5-standards/code-review-checklist.md) — Checklist para PRs
- [Security Checklist](/.knowledge/5-standards/security-checklist.md) — OWASP verificación

### 🔒 6 — Seguridad
- [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md) — 12 roles y matriz de permisos (espejo de `PermissionsByRole`)
- [Threat Model](/.knowledge/6-security/threat-model.md) — Modelo de amenazas (⚠️ adáptalo a tu proyecto)
- [Data Classification](/.knowledge/6-security/data-classification.md) — Clasificación PII/PHI
