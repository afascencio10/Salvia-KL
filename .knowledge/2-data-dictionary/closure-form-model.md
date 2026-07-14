---
okf_version: "1.0"
type: Data_Model
title: "Formulario dinámico: Cierre de proceso psicosocial"
description: "Definición (form / section / questions / options / visibility_conditions) del formulario dinámico de Cierre de proceso psicosocial, renderizado por el componente dinamic-form. Se dispara desde los caminos de cierre del flujo 3x3 (no consentimiento / imposibilidad de contacto)."
owner: "@platform-team"
status: active
tags: [database, data, psicosocial, 3x3, closure, dynamic-form]
engine: PostgreSQL
schema: salvia
table: form
code_refs:
  - "src/cmd/seed/seed_cierre_psicosocial.sql"
last_updated: "2026-07-10"
---

# Formulario dinámico: Cierre de proceso psicosocial

Formulario dinámico (motor `dinamic-form`) que registra el cierre de un proceso de Atención Psicosocial. Se seedéa como filas en `salvia.form → form_section → question → option → visibility_condition` (mismo patrón que [Cierre de caso](/.knowledge/../src/cmd/seed/cierre_caso_db_seeding.md), `da8423ab…`). El motor recupera la estructura por `form_id`.

## Identificadores registrados

### A. Formulario (`salvia.form`)
- **ID:** `fcc7dc8d-835b-4559-9283-2ea8b8e6092b`
- **Nombre:** `Cierre de proceso psicosocial`
- **Descripción:** `Formulario para registrar el cierre de un proceso de Atención Psicosocial`
- **Estado:** `active`

### B. Sección (`salvia.form_section`)
- **ID:** `0821f9a1-3bf2-4349-8cb8-5b0c337716fe`
- **Formulario:** `fcc7dc8d-835b-4559-9283-2ea8b8e6092b`
- **Nombre:** `Cierre`
- **Orden:** `1`

### C. Preguntas (`salvia.question`)

| # | Descripción | `question_type` | UUID | Obligatoria (BD) | Visibilidad |
| :---: | :--- | :--- | :--- | :---: | :--- |
| 1 | ¿La llamada fue efectiva? | `single` | `c253f68c-866e-4093-ba80-a41bbb96b0ae` | Sí | siempre |
| 2 | Motivo de cierre | `single` | `1c6d97ac-dc83-4e73-b747-14cd93ad422d` | Sí | siempre |
| 3 | Contenido de la atención | `text` | `036f6945-d55a-496d-b343-b9c57eca954f` | Sí | si #1 = `si` |
| 4 | Plan de orientación | `multiple` | `96146a85-c194-470a-a2c8-c93413c88fc2` | No | si #1 = `si` |
| 5 | Temas trabajados durante la atención | `text` | `9628c09c-5e5d-4d2a-82ff-d944491044ac` | Sí | si #1 = `si` |
| 6 | Hay nuevos hechos de violencia | `boolean` | `555d2810-c0b1-4422-95cc-02479fc061f4` | Sí | siempre |
| 7 | Descripción de los hechos | `text` | `39c3f808-b100-42d6-83fd-78f7aec00f96` | No¹ | si #6 = `true` |
| 8 | Fecha | `date` | `5e0713eb-90ac-4183-bcdf-f72d1caeb489` | No¹ | si #6 = `true` |

¹ **Obligatorias condicionalmente:** las preguntas 7 y 8 se marcan `required = false` en BD, pero el frontend las **exige cuando son visibles** (es decir, cuando #6 = `true`), replicando el patrón de validación condicional del cierre de seguimientos.

### D. Opciones (`salvia.option`)

**¿La llamada fue efectiva?** (`c253f68c…`)
| Label | Value | UUID | Orden |
| :--- | :--- | :--- | :---: |
| Sí | `si` | `4dff9f5d-36cf-4db7-a767-52b0d3add5f2` | 1 |
| No | `no` | `8f1bdbbd-1dcc-4b4f-a66f-52897df0e4d5` | 2 |

**Motivo de cierre** (`1c6d97ac…`)
| Label | Value | UUID | Orden |
| :--- | :--- | :--- | :---: |
| Cumplimiento de objetivos | `cumplimiento_objetivos` | `53b8a4eb-d506-4215-870a-577e670e4907` | 1 |
| Cumplimiento esquema | `cumplimiento_esquema` | `2d17669f-6f61-4c02-9aea-e9267e567bf0` | 2 |
| No consentimiento | `no_consentimiento` | `9c25c5f0-13f2-499b-9c8f-5d39bdaf61ef` | 3 |
| Imposibilidad del contacto (3x3) | `imposibilidad_contacto_3x3` | `d2c6065a-9da8-443e-a1d6-2cd5eee47a87` | 4 |
| Desistimiento del proceso | `desistimiento_proceso` | `d054dd41-b927-4f85-b7bb-2949dc1576c9` | 5 |

**Plan de orientación** (`96146a85…`)
| Label | Value | UUID | Orden |
| :--- | :--- | :--- | :---: |
| Enrutamiento | `enrutamiento` | `d7dac006-bc8e-4d5b-83b3-e914073eb671` | 1 |
| Activación de ruta | `activacion_ruta` | `09fa69e1-9f65-4402-bf7c-4bad61aa450b` | 2 |
| Seguimiento | `seguimiento` | `9e1d08f3-68d3-4130-bbb1-c8440bd92b85` | 3 |
| Medidas de emergencia | `medidas_emergencia` | `2cb7d162-5b36-4516-a28b-4abf8a4adde9` | 4 |
| Plan de estabilización | `plan_estabilizacion` | `01d19bc2-9349-466a-aedf-bfd5fcc5fa4b` | 5 |

### E. Condiciones de visibilidad (`salvia.visibility_condition`)

`target_type = 'QUESTION'`, `operator = 'EQUALS'`.

| UUID | target_id (pregunta que se muestra) | trigger_question_id | trigger_value |
| :--- | :--- | :--- | :--- |
| `fae70d75-6d71-492f-bb36-54234f085bf7` | #3 Contenido (`036f6945…`) | #1 ¿Llamada? (`c253f68c…`) | `si` |
| `66aa281a-9a24-41d6-9c09-6cb73f44e485` | #4 Plan (`96146a85…`) | #1 ¿Llamada? (`c253f68c…`) | `si` |
| `61bd463f-c2a1-4da7-8bd0-029602e08e86` | #5 Temas (`9628c09c…`) | #1 ¿Llamada? (`c253f68c…`) | `si` |
| `b7e5bd86-97fc-4bc4-b618-4d4a3bdf71bc` | #7 Descripción (`39c3f808…`) | #6 Nuevos hechos (`555d2810…`) | `true` |
| `8ece8f82-75ec-4634-a44a-e51aa0bc0858` | #8 Fecha (`5e0713eb…`) | #6 Nuevos hechos (`555d2810…`) | `true` |

## Uso en el flujo 3x3

- El formulario se abre desde los dos caminos de cierre del [flujo 3x3](/.knowledge/3-features/AtencionPsicosocial3x3/index.md): **No acepta consentimiento** e **Imposibilidad de contacto** (3c). Ver [init-closure-form](/.knowledge/3-features/AtencionPsicosocial3x3/backend/endpoints/init-closure-form.md).
- Al abrir, el **Motivo de cierre** (#2) se **preselecciona (editable)** según el disparador (`no_consentimiento` o `imposibilidad_contacto_3x3`).
- Al completar el formulario, el proceso transiciona a `cerrado` o `en_devolucion` según el motivo — ver [complete-closure](/.knowledge/3-features/AtencionPsicosocial3x3/backend/endpoints/complete-closure.md).

## Seed

El seed vive en `src/cmd/seed/seed_cierre_psicosocial.sql` (patrón `seed_cierre_caso.sql`). Nota: `*.sql` está en `.gitignore`; el seed se aplica manualmente. Este documento es la fuente canónica de los UUIDs.
