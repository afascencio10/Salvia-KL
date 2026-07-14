---
okf_version: "1.0"
type: Project
title: "Índice del Data Dictionary"
description: "Inventario de todos los modelos de datos, esquemas de base de datos y catálogos que forman el dominio del sistema."
owner: "@platform-team"
status: active
tags: [database, data, index]
last_updated: "2026-07-02"
---

# Data Dictionary — Índice

Este directorio contiene la definición de los modelos de datos de SALVIA. El motor es **PostgreSQL**; los modelos del dominio nuevo se definen como structs GORM (`src/internal/`, `src/salvia/`) y se migran vía `AutoMigrate` + SQL en `src/migrations/`. El código legacy accede vía pgx (`src/security/dao`).

## Modelos Registrados

| Modelo | Tipo | Engine | Tags | Estado |
| :--- | :--- | :--- | :--- | :---: |
| [error-catalog](/.knowledge/2-data-dictionary/error-catalog.md) | `Error_Catalog` | N/A | `errors, api` | ✅ Active |
| [contact-attempts](/.knowledge/2-data-dictionary/contact-attempts-model.md) | `Data_Model` | PostgreSQL | `psicosocial, 3x3, contact-attempts, pii` | 🔜 Draft |
| [team-contact](/.knowledge/2-data-dictionary/team-contact-model.md) | `Data_Model` | PostgreSQL | `psicosocial, session, brownfield` | 🔜 Draft |
| [closure-form](/.knowledge/2-data-dictionary/closure-form-model.md) | `Data_Model` | PostgreSQL | `psicosocial, 3x3, closure, dynamic-form` | 🔜 Draft |
<!-- Añade aquí cada modelo: | [nombre](/.knowledge/2-data-dictionary/<nombre>-model.md) | `Data_Model` | PostgreSQL | `tags` | ✅ Active | -->

> [!NOTE]
> **Migración asociada al flujo 3x3 psicosocial:** `ALTER TABLE salvia.psychosocial_support ADD COLUMN next_contact_attempt_at TIMESTAMPTZ;` (fecha/hora del próximo intento de contacto, acción 3a). Ver [contact-attempts-model](/.knowledge/2-data-dictionary/contact-attempts-model.md).

> [!NOTE]
> Los modelos de datos concretos (`PsychosocialSupport`, casos, remisiones, seguimientos, barreras, usuarios...) todavía **no están documentados**: se documentan de forma incremental al tocar cada módulo (adopción brownfield, `AGENTS.md` §8), no todos por adelantado.

## Convenciones

- Cada tabla/colección tiene su propio archivo `.md` con esquema completo.
- Los nombres de tablas usan `snake_case` y plural (`users`, `payments`, `audit_logs`).
- Los campos PII están marcados con el tag `pii` en el frontmatter.
- Las migraciones se documentan en archivos separados tipo `Migration`.
