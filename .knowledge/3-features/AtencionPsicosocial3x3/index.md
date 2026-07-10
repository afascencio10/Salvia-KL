---
okf_version: "1.0"
type: Project
title: "Feature: Flujo 3x3 de Atención Psicosocial"
description: "Flujo de gestión de contacto 3x3 para el módulo de Atención Psicosocial. La profesional (roles ps/ts) registra intentos de contacto sobre un proceso psychosocial_support; al contacto exitoso el proceso pasa a en_gestion y abre Consentimiento Informado (acepta → sesión inmediata o agendada en team_contact; no acepta → cierre); los intentos fallidos se acumulan por día y, al 3º fallido del día, se abre el modal de acciones (reprogramar próximo intento vía next_contact_attempt_at, añadir intento hasta 50, o cerrar por imposibilidad tras 9 intentos en 3 días distintos). No usa ni modifica follow_up_v2."
owner: "@platform-team"
status: active
tags: [feature, psicosocial, 3x3, contact-attempts, rbac, brownfield]
dependencies:
  - Spec: 6-security/rbac-matrix.md
  - Model: 2-data-dictionary/contact-attempts-model.md
  - Model: 2-data-dictionary/team-contact-model.md
last_updated: "2026-07-09"
---

# Feature: Flujo 3x3 de Atención Psicosocial

## Descripción

Análogo al esquema 3x3 de Seguimientos, pero para el módulo de **Atención Psicosocial** y con reglas propias. La profesional psicosocial (roles `ps` = Psicólogo, `ts` = Trabajador social) gestiona el contacto con la ciudadana desde un modal sobre un proceso `psychosocial_support` (estados `abierto → en_gestion → en_devolucion → cerrado`).

**Diferencia arquitectónica clave:** este flujo **NO usa ni modifica `follow_up_v2`** (el motor de Seguimientos). Opera de forma autónoma sobre tres tablas:
- `salvia.psychosocial_support` — el proceso (se le añade la columna `next_contact_attempt_at`).
- `salvia.contact_attempts` — **tabla nueva**; guarda **todos** los intentos (exitosos y fallidos). El módulo psicosocial **no** usa `follow_up_attempts`.
- `salvia.team_contact` — sesión terapéutica agendada (ya existe).

## Flujo funcional

### 1. Modal de contacto (¿Contestó?)
- Muestra "Contestó" / "No contestó", un selector **Fecha y hora del contacto** con default = ahora **editable**, y un campo de **nota** libre opcional (`note`) aplicable a ambos caminos (ej. "No contestó, celular apagado" / "Sí contestó, pide que lo llamen otro día"). Si hay intentos previos, muestra un banner "{n} intento(s) sin respuesta".
- **Contestó (was_answered=true):**
  1. Se registra el contacto exitoso en `contact_attempts`.
  2. El proceso `psychosocial_support` pasa a **`en_gestion`**.
  3. Se deja el registro en el timeline del caso.
  4. Se abre el modal de **Consentimiento Informado**.
     - **No acepta** → abre el **[Formulario de Cierre de proceso psicosocial](/.knowledge/2-data-dictionary/closure-form-model.md)** (`dinamic-form`) con Motivo preseleccionado (editable) en `no_consentimiento` → al completar, el proceso pasa a **`en_devolucion`**.
     - **Sí acepta** → opción **sesión de inmediato** o **reagendar**:
       - *Inmediato* → se crea la sesión en `team_contact` (`is_psico_session=true`, `scheduled_date=ahora`) y se **redirige al formulario de atención** (🔴 futuro).
       - *Reagendar* → input de calendario para fecha/hora flexible → crea la sesión en `team_contact` con esa fecha.
- **No contestó (intentos 1–2 del día):**
  1. Se registra el intento fallido con la **nota** libre opcional (`note`) capturada en el mismo modal, y se cierra (no hay segundo modal).
  2. Se deja el registro en el timeline.
  3. El caso **se mantiene `abierto`** (aún no en gestión).

### 2. Umbral de 3 intentos diarios (diferencia clave con Seguimientos)
El contador es **por día natural** (se reinicia cada día). Cuando hay **3 intentos fallidos en el mismo día**, el flujo **salta la alerta** de límite y abre directamente el **modal de acciones**:

- **3a. Agendamiento manual flexible (NO posponer automático):** a diferencia de Seguimientos (que ofrece +1/+2/+3 días automáticos), aquí la profesional elige **fecha y hora exactas** del próximo intento en un calendario/inputs (según sus turnos rotativos y disponibilidad de la ciudadana). Se persiste en **`psychosocial_support.next_contact_attempt_at`**.
- **3b. Añadir más intentos:** reabre el flujo de contacto para usar huecos de gestión del **mismo día**, hasta un **tope de 50 intentos** totales.
- **3c. Cierre por imposibilidad de contacto:** botón visible **solo** si se cumplen **≥ 9 intentos en ≥ 3 días distintos**. Abre el **[Formulario de Cierre de proceso psicosocial](/.knowledge/2-data-dictionary/closure-form-model.md)** con Motivo preseleccionado (editable) en `imposibilidad_contacto_3x3` → al completar, el proceso pasa a **`cerrado`**. Cierra **únicamente el proceso psicosocial** (no el caso de víctima).

## Documentos Relacionados

### Seguridad
- [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md) — permisos nuevos para `ps` y `ts` (registrar intento, consentimiento, agendar sesión, próximo intento, cierre).

### Data Dictionary
- [contact-attempts-model](/.knowledge/2-data-dictionary/contact-attempts-model.md) — entidad nueva
- [team-contact-model](/.knowledge/2-data-dictionary/team-contact-model.md) — sesión terapéutica (brownfield)
- [closure-form-model](/.knowledge/2-data-dictionary/closure-form-model.md) — formulario dinámico de cierre psicosocial (form/section/questions/options/visibility)

### Backend — Endpoints
| Endpoint | Spec | Estado |
| :--- | :--- | :---: |
| `GET /api/v1/psychosocial/{psicosocialId}/contact-attempts` | [list-contact-attempts.md](./backend/endpoints/list-contact-attempts.md) | 🔜 |
| `POST /api/v1/psychosocial/{psicosocialId}/contact-attempts` | [register-contact-attempt.md](./backend/endpoints/register-contact-attempt.md) | 🔜 |
| `PATCH /api/v1/contact-attempts/{id}/consent` | [set-consent.md](./backend/endpoints/set-consent.md) | 🔜 |
| `POST /api/v1/psychosocial/{psicosocialId}/sessions` | [schedule-session.md](./backend/endpoints/schedule-session.md) | 🔜 |
| `PUT /api/v1/psychosocial/{psicosocialId}/next-attempt` | [set-next-attempt.md](./backend/endpoints/set-next-attempt.md) | 🔜 |
| `POST /api/v1/psychosocial/{psicosocialId}/init-closure-form` | [init-closure-form.md](./backend/endpoints/init-closure-form.md) | 🔜 |
| `POST /api/v1/psychosocial/{psicosocialId}/close` | [complete-closure.md](./backend/endpoints/complete-closure.md) | 🔜 |

### Backend — Servicios
| Servicio | Spec | Estado |
| :--- | :--- | :---: |
| Psychosocial3x3Service | [psychosocial-3x3-service.md](./backend/services/psychosocial-3x3-service.md) | 🔜 |

### Frontend — Screens
| Pantalla | Spec | Estado |
| :--- | :--- | :---: |
| Remisión Temporal (`/salvia/remision-temporal/:id`, host del card) | [index.md](./frontend/screens/RemisionTemporal/index.md) | 🔜 |

### Frontend — Componentes
| Componente | Spec | Estado |
| :--- | :--- | :---: |
| psychosocial-contact-card (card que se monta en pantallas) | [index.md](./frontend/components/PsychosocialContactCard/index.md) | 🔜 |
| psychosocial-contact-modal (flujo de modales que hospeda el card) | [index.md](./frontend/components/PsychosocialContactModal/index.md) | 🔜 |

### Testing
| Plan | Spec | Estado |
| :--- | :--- | :---: |
| Test Plan | [atencion-psicosocial-3x3-test-plan.md](./testing/atencion-psicosocial-3x3-test-plan.md) | 🔜 |

## Historias de Usuario

| ID | Historia | Estado |
| :--- | :--- | :---: |
| US-3X3-01 | Como profesional psicosocial (`ps`/`ts`), quiero registrar si la ciudadana contestó o no —con fecha/hora editable— para dejar trazabilidad de cada intento de contacto. | 🔜 |
| US-3X3-02 | Como profesional, al lograr contacto exitoso quiero que el proceso pase a `en_gestion` y se abra el Consentimiento Informado, para continuar la atención. | 🔜 |
| US-3X3-03 | Como profesional, si la ciudadana acepta el consentimiento quiero elegir entre iniciar la sesión de inmediato o agendarla con fecha/hora flexible. | 🔜 |
| US-3X3-04 | Como profesional, si la ciudadana no acepta el consentimiento quiero cerrar el proceso psicosocial (vía formulario de cierre). | 🔜 |
| US-3X3-05 | Como profesional, tras 3 intentos fallidos en un mismo día quiero ver el modal de acciones sin la alerta intermedia. | 🔜 |
| US-3X3-06 | Como profesional, quiero fijar la fecha/hora exacta del próximo intento (3a) en lugar de un posponer automático. | 🔜 |
| US-3X3-07 | Como profesional, quiero poder añadir más intentos el mismo día hasta un tope de 50. | 🔜 |
| US-3X3-08 | Como profesional, tras 9 intentos en 3 días distintos quiero cerrar el proceso por imposibilidad de contacto. | 🔜 |
| US-3X3-09 | Como profesional (`ps`/`ts`), al dar "Ver remisión" en Mis Remisiones Psicosocial quiero llegar a una pantalla que cargue el card 3x3 del proceso (pantalla temporal `remision-temporal/:id`) mientras no exista la pantalla interna definitiva. | 🔜 |

## Reglas de Negocio (resumen normativo)

| ID | Regla |
| :--- | :--- |
| RN-01 | Todos los intentos (exitosos y fallidos) se registran en `contact_attempts`. Este módulo **no** escribe en `follow_up_attempts` ni en `follow_up_v2`. |
| RN-02 | Fecha/hora del contacto (`attempt_at`): default = ahora, editable por la profesional. |
| RN-03 | Contacto exitoso (`was_answered=true`) → `psychosocial_support.status = en_gestion`. Intentos fallidos mantienen `abierto`. |
| RN-04 | Umbral diario: `COUNT(*) WHERE psicosocial_id = X AND was_answered = false AND attempt_at::date = CURRENT_DATE` ≥ 3 → abre modal de acciones (salta alerta). Se reinicia cada día. |
| RN-05 | Tope duro: máximo **50** intentos totales por proceso. Al alcanzarlo, "Añadir intento" queda deshabilitado. |
| RN-06 | Cierre por imposibilidad (3c) visible solo si `COUNT(*) ≥ 9` **y** `COUNT(DISTINCT attempt_at::date) ≥ 3`. |
| RN-07 | 3a persiste la fecha/hora del próximo intento en `psychosocial_support.next_contact_attempt_at`. |
| RN-08 | La decisión de consentimiento se persiste en el `contact_attempts` exitoso (`consent_given`). |
| RN-09 | El cierre (no acepta / imposibilidad) afecta **solo** al proceso psicosocial; no cierra el caso de víctima. |
| RN-10 | Al registrar un intento (contestó o no), si había un `next_contact_attempt_at` cuya fecha es **igual o anterior** al día del intento, se limpia a `NULL` (el intento planeado ya se ejecutó, aunque el agente lo hiciera días después). Solo se conserva si lo programado es a **futuro**. |
| RN-11 | El historial se muestra **del más nuevo al más viejo**; la fila "PROGRAMADO" (próximo intento) va arriba. El `Intento #N` mantiene numeración cronológica (1 = más antiguo). |
| RN-12 | El cierre se hace vía el **Formulario de Cierre de proceso psicosocial** (`dinamic-form`, `form_id fcc7dc8d…`). El Motivo se preselecciona (editable) según el disparador; el cierre usa el Motivo final del formulario. Mapeo Motivo→estado: `imposibilidad_contacto_3x3`/`cumplimiento_*` → `cerrado`; `no_consentimiento`/`desistimiento_proceso` → `en_devolucion`. Solo se disparan hoy `no_consentimiento` e `imposibilidad_contacto_3x3`. |

## Gaps / Dependencias fuera de alcance (🔴 PENDIENTE DE DEFINIR)

> [!WARNING]
> Lo siguiente **no** se define en este spec; solo su punto de disparo/redirección.

- 🔴 **Formulario de atención de la sesión** — destino del redirect en "sesión inmediata" y de la sesión agendada (`team_contact.form_submission_id` se llena al responderlo).
- **Disparadores de cierre por cumplimiento/desistimiento** — el [Formulario de Cierre](/.knowledge/2-data-dictionary/closure-form-model.md) ya contempla esos motivos, pero por ahora **solo** se dispara desde los 2 caminos 3x3 (no consentimiento / imposibilidad). Un botón "Cerrar proceso" para los demás motivos queda para un incremento futuro.

## Nota de Brownfield (§8 AGENTS.md)

`psychosocial_support` y `team_contact` ya existen en código sin spec OKF. Este feature documenta (lazy specing) lo que necesita tocar: la nueva columna `next_contact_attempt_at` en `psychosocial_support` y el uso de `team_contact` para sesiones. El resto de esos modelos queda como deuda documental medible.
