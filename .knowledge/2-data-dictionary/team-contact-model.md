---
okf_version: "1.0"
type: Data_Model
title: "Modelo: team_contact"
description: "Sesión/contacto del equipo de Atención Psicosocial. Tabla existente (brownfield). En el flujo 3x3 almacena la sesión terapéutica agendada (inmediata o con fecha flexible) tras aceptar el Consentimiento Informado."
owner: "@platform-team"
status: active
tags: [database, data, psicosocial, team-contact, session, brownfield]
engine: PostgreSQL
schema: salvia
table: team_contact
code_refs:
  - "src/internal/models/team_contact.go"
last_updated: "2026-07-09"
---

# Modelo de Datos: `salvia.team_contact`

Tabla **ya existente** que registra contactos/sesiones del equipo de Atención Psicosocial. Documentada aquí de forma incremental (brownfield §8) porque el flujo 3x3 la usa para persistir la **sesión terapéutica agendada** (opción "sesión inmediata" o "reagendar" tras aceptar el Consentimiento Informado).

## Esquema

| Campo | Tipo | Null | Default | Descripción |
| :--- | :--- | :---: | :--- | :--- |
| `id` | `varchar(36)` | No | `gen_random_uuid()` | PK. |
| `case_id` | `varchar(36)` | No | — | Caso de víctima. |
| `form_submission_id` | `varchar(36)` | Sí | `NULL` | Submission del formulario de atención. **Vacío hasta que se responde** el formulario de la sesión (🔴 form futuro). |
| `dupla_id` | `varchar(36)` | Sí | `NULL` | Dupla asignada. |
| `psicosocial_id` | `varchar(36)` | Sí | `NULL` | Proceso `psychosocial_support` asociado. |
| `professional_id` | `varchar(36)` | Sí | `NULL` | Profesional que agenda/atiende. |
| `team` | `varchar(50)` | Sí | `NULL` | Equipo de la profesional. |
| `scheduled_date` | `timestamptz` | Sí | `NULL` | Fecha agendada de la sesión (ahora, si es inmediata; flexible, si se reagenda). |
| `scheduled_time` | `varchar(8)` | Sí | `NULL` | Hora agendada (`HH:MM` / `HH:MM:SS`). |
| `is_completed` | `boolean` | No | `false` | La sesión ya se realizó. |
| `completed_at` | `timestamp` | Sí | `NULL` | Momento de completado. |
| `status` | `varchar(20)` | Sí | `NULL` | Estado de la sesión. |
| `summary` | `text` | Sí | `NULL` | Resumen de la sesión. |
| `is_psico_session` | `boolean` | No | `false` | `true` cuando la fila representa una sesión psicosocial (la creada por el flujo 3x3). |
| `created_at` | `timestamptz` | No | `now()` | Auditoría. |
| `updated_at` | `timestamptz` | No | `now()` | Auditoría. |
| `deleted_at` | `timestamptz` | Sí | `NULL` | Soft delete (GORM). |

## Uso en el flujo 3x3

Al aceptar el Consentimiento Informado y elegir sesión:
- **Inmediata:** se crea una fila con `is_psico_session = true`, `is_completed = false`, `scheduled_date = ahora`, `form_submission_id = NULL`; luego se redirige al formulario de atención (🔴 futuro).
- **Reagendar:** se crea una fila con `is_psico_session = true`, `is_completed = false`, `scheduled_date`/`scheduled_time` = fecha/hora flexible elegida; sin redirección.

## Struct GORM (existente)

Definido en `src/internal/models/team_contact.go` (ver `TableName() = "salvia.team_contact"`). Este spec no cambia el struct; solo formaliza su uso en el flujo 3x3.

## Deuda documental reconocida

El resto de usos de `team_contact` (reasignación, listados psicosocial en `src/internal/repository/psychosocial_*`) queda como deuda documental medible; se documentará al tocar esos módulos (brownfield §8).
