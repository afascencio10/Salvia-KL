---
okf_version: "1.0"
type: Test_Plan
title: "Plan de Pruebas: Flujo 3x3 de Atención Psicosocial"
description: "Estrategia de testing y matriz de casos (API + E2E) para el flujo 3x3 psicosocial: intentos de contacto, consentimiento, agendamiento de sesión, umbral diario, tope de 50 y cierre por imposibilidad."
owner: "@qa-squad"
status: active
tags: [qa, testing, psicosocial, 3x3]
dependencies:
  - Feature: 3-features/AtencionPsicosocial3x3/index.md
last_updated: "2026-07-09"
---

# Estrategia de Testing (Flujo 3x3 de Atención Psicosocial)

Herramienta oficial: **Playwright** (API y E2E). Cada caso se marca `PASS`/`FAIL` con fecha y comando al ejecutarse (AGENTS.md §7).

## Evidencia de ejecución (2026-07-09)

### 1. Suite HTTP aislada (`httptest` + stub, sin DB)
- **Comando**: `go test ./tests/psychosocial3x3/ -v` → `ok bitsflow/tests/psychosocial3x3` — **9/9 PASS**.
- Cubre: registro de rutas sin conflicto de Gin; GET historial (200 + `caseName`); GET 404; POST intento (201) y sin `was_answered` (400); PATCH consentimiento (200 + `requiresClosureForm`); POST sesión inmediata (201) y sin fecha (400); PUT próximo intento (200).
- **Build**: `go build . ./salvia/... ./internal/... ./tests/...` → OK. `node --check` de ambos componentes JS → OK.

### 2. E2E real contra Supabase (server en :8099, curl)
Proceso real `019f24bb-…-000000000002` ("Juliana Orozco"), estado inicial `abierto`. Datos de prueba **limpiados** al terminar (4 filas de `contact_attempts` borradas, sesión de `team_contact` borrada, `status` restaurado a `abierto`).

| Paso | Resultado | Estado |
| :--- | :--- | :---: |
| GET historial inicial | `attempts:[]`, `status: abierto`, `caseName: "Juliana Orozco"` | ✅ PASS |
| POST fallido ×1,×2 | `dailyFailedCount` 1→2, `status` sigue `abierto` | ✅ PASS |
| POST fallido ×3 | `dailyFailedCount=3`, **`dailyThresholdReached=true`** | ✅ PASS |
| POST exitoso | 201, **`processStatus=en_gestion`** | ✅ PASS |
| PATCH consentimiento `true` | `canScheduleSession=true`, `requiresClosureForm=false` | ✅ PASS |
| POST sesión inmediata | 201, fila en `team_contact` (`isPsicoSession=true`) + `redirectUrl` | ✅ PASS |
| GET historial final | 4 intentos con `sequenceNumber`, notas y badges correctos | ✅ PASS |
| PUT próximo intento | Tras aplicar la migración (columna creada), el `UPDATE` del endpoint corre OK como rol de la app (`UPDATE 1`, valor persistido) + handler 200 en httptest | ✅ PASS |

> **Migración aplicada:** la columna `next_contact_attempt_at` la agregó el rol `postgres` (SQL editor de Supabase / conexión privilegiada), ya que `salvia.psychosocial_support` es propiedad de `postgres` y `AutoMigrate` (rol `salvia_gorm`) no puede ALTERarla. Verificado en la DB real: la columna existe y el `UPDATE SET next_contact_attempt_at ...` que ejecuta `SetNextAttempt` devuelve `UPDATE 1` con el rol de la app. La única razón por la que no se hizo el `curl` HTTP en vivo del 5º endpoint fue que el arranque completo de la app quedó sin conexiones por saturación del pooler de Supabase (tope de 15 clientes en session mode) al correr en paralelo con otro server; no es un problema del código. La tabla `contact_attempts` sí se crea vía AutoMigrate (la app es su owner).

> **Datos de prueba limpiados:** al terminar el E2E, el proceso usado quedó restaurado (0 `contact_attempts`, `status='abierto'`, `next_contact_attempt_at=NULL`, sesión de `team_contact` borrada).

### 3. Cierre psicosocial (2026-07-10)

- **Seed aplicado** a Supabase (`seed_cierre_psicosocial.sql`, rol `salvia_gorm`): form `Cierre de proceso psicosocial` con **1 sección, 8 preguntas, 12 opciones, 5 condiciones de visibilidad** (verificado por conteo en la DB).
- **httptest** (`go test ./tests/psychosocial3x3/`) — incluye init-closure-form (201 + `preselectedMotivo`), reason inválido (400), close (200 + `status=cerrado`), submission faltante (400). **Todos PASS**.
- **Round-trip contra DB real** (SQL exacto de `InitClosureForm`+`CloseProcess`): crear `form_submission`+`answer(Motivo)`, leer Motivo, mapear estado:

  | Motivo | Estado | Resultado |
  | :--- | :--- | :---: |
  | `imposibilidad_contacto_3x3` | `cerrado` | ✅ PASS |
  | `no_consentimiento` | `en_devolucion` | ✅ PASS |

  Datos de prueba limpiados (submissions/answers borradas, `status` restaurado a `abierto`).
- No se hizo `curl` HTTP en vivo de los endpoints de cierre por la misma saturación del pooler de Supabase (arranque completo de la app sin conexiones); no es problema de código.

## Pruebas de API (Backend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Playwright Script (`.spec.ts`) |
| :--- | :--- | :--- | :--- |
| `TC-3X3-API-01` | POST contact-attempts con `was_answered=true` | `201`; `processStatus=en_gestion`; fila en `contact_attempts` | `tests/api/psychosocial_3x3.spec.ts` |
| `TC-3X3-API-02` | POST contact-attempts fallido solo con `was_answered=false` (sin nota) | `201`; fila creada; proceso sigue `abierto` | `tests/api/` |
| `TC-3X3-API-03` | POST contact-attempts con `attempt_at` explícito editado | `201`; `attemptAt` coincide con el enviado | `tests/api/` |
| `TC-3X3-API-03b` | POST con `note` en exitoso y en fallido | `201`; `note` persistido en ambos | `tests/api/` |
| `TC-3X3-API-04` | 3er intento fallido del mismo día | `counters.dailyThresholdReached=true` | `tests/api/` |
| `TC-3X3-API-05` | Intento fallido nuevo día (día 2) reinicia conteo diario | `dailyFailedCount=1`, `dailyThresholdReached=false` | `tests/api/` |
| `TC-3X3-API-06` | Intento nº 51 | `409 RESOURCE_CONFLICT` (tope 50) | `tests/api/` |
| `TC-3X3-API-07` | 9 intentos en 3 días distintos | `counters.closureEligible=true` | `tests/api/` |
| `TC-3X3-API-08` | 9 intentos en 2 días distintos | `closureEligible=false` | `tests/api/` |
| `TC-3X3-API-09` | PATCH consent `false` sobre intento exitoso | `200`; `requiresClosureForm=true` | `tests/api/` |
| `TC-3X3-API-10` | PATCH consent sobre intento fallido | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-3X3-API-11` | POST sessions `immediate=true` | `201`; fila `team_contact` (`is_psico_session=true`); `redirectUrl` presente | `tests/api/` |
| `TC-3X3-API-12` | POST sessions `immediate=false` con fecha/hora | `201`; `scheduledDate`/`scheduledTime` guardados; sin `redirectUrl` | `tests/api/` |
| `TC-3X3-API-13` | POST sessions sin consentimiento aceptado | `409 RESOURCE_CONFLICT` | `tests/api/` |
| `TC-3X3-API-14` | PUT next-attempt con fecha válida | `200`; `psychosocial_support.next_contact_attempt_at` actualizado | `tests/api/` |
| `TC-3X3-API-15` | Cualquier endpoint con rol distinto de `ps`/`ts` | `403 AUTH_INSUFFICIENT_ROLE` | `tests/api/` |
| `TC-3X3-API-16` | Endpoint con `psicosocialId` inexistente | `404 RESOURCE_NOT_FOUND` | `tests/api/` |
| `TC-3X3-API-17` | POST init-closure-form `reason=imposibilidad_contacto_3x3` | `201`; `submission_id`; `lockedMotivo=imposibilidad_contacto_3x3` | `tests/api/` |
| `TC-3X3-API-18` | POST init-closure-form `reason` inválido | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-3X3-API-19` | POST close con submission de motivo `imposibilidad_contacto_3x3` | `200`; `status=cerrado` | `tests/api/` |
| `TC-3X3-API-20` | POST close con submission de motivo `no_consentimiento` | `200`; `status=en_devolucion` | `tests/api/` |
| `TC-3X3-API-21` | POST close con submission sin Motivo respondido | `400 VALIDATION_FAILED` | `tests/api/` |

## Pruebas E2E / UI (Frontend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Playwright Script (`.spec.ts`) |
| :--- | :--- | :--- | :--- |
| `TC-3X3-UI-01` | Abrir modal, "Sí contestó" | Se abre modal de Consentimiento Informado | `tests/e2e/psychosocial_3x3.spec.ts` |
| `TC-3X3-UI-02` | Consentimiento "Sí acepta" → "De inmediato" | Redirige al formulario de atención (o placeholder 🔴) | `tests/e2e/` |
| `TC-3X3-UI-03` | Consentimiento "Sí acepta" → "Agendar" | Aparecen inputs de calendario; guarda sesión | `tests/e2e/` |
| `TC-3X3-UI-04` | Consentimiento "No acepta" | Abre formulario de cierre (o placeholder 🔴) | `tests/e2e/` |
| `TC-3X3-UI-05` | "No contestó" registra y cierra directo (sin segundo modal) | Se cierra el modal con overlay de éxito; no aparece paso de destinatario | `tests/e2e/` |
| `TC-3X3-UI-05b` | Escribir nota en modal ¿Contestó? y elegir "Sí contestó" | `note` viaja en el POST del intento exitoso | `tests/e2e/` |
| `TC-3X3-UI-05c` | Banner "N intento(s) sin respuesta" con intentos previos | Banner visible con el conteo correcto | `tests/e2e/` |
| `TC-3X3-UI-06` | 3er fallido del día | Aparece modal de acciones sin alerta intermedia | `tests/e2e/` |
| `TC-3X3-UI-07` | Acción 3a | Abre calendario de próximo intento; guarda `next_contact_attempt_at` | `tests/e2e/` |
| `TC-3X3-UI-08` | Acción 3b con totalCount=50 | Botón "Añadir intento" deshabilitado | `tests/e2e/` |
| `TC-3X3-UI-09` | Acción 3c sin cumplir 9-en-3-días | Botón de cierre por imposibilidad oculto | `tests/e2e/` |
| `TC-3X3-UI-10` | Fecha/hora del contacto editada | Se envía `attempt_at` editado (screenshot only-on-failure) | `tests/e2e/` |
| `TC-3X3-UI-11` | "No acepta" consentimiento | Abre `dinamic-form` de cierre con Motivo=`No consentimiento` preseleccionado (editable) | `tests/e2e/` |
| `TC-3X3-UI-12` | "Cerrar por imposibilidad" (3c) | Abre `dinamic-form` con Motivo=`Imposibilidad del contacto (3x3)` preseleccionado (editable) | `tests/e2e/` |
| `TC-3X3-UI-13` | ¿Llamada efectiva? = Sí | Se muestran Contenido, Plan de orientación y Temas trabajados | `tests/e2e/` |
| `TC-3X3-UI-14` | Hay nuevos hechos = verdadero | Descripción y Fecha se muestran y quedan obligatorias | `tests/e2e/` |
| `TC-3X3-UI-15` | Completar formulario de cierre | Proceso pasa a cerrado/en_devolucion; card refresca | `tests/e2e/` |

## Pruebas Unitarias (Go)

- `dailyFailedCount` cuenta por `attempt_at::date` y no por `created_at`.
- `maxAttemptsReached` en el límite exacto (49 → false, 50 → true).
- `closureEligible` en fronteras (8/3d, 9/2d, 9/3d).
- Transición `abierto → en_gestion` solo en contacto exitoso; fallido no cambia estado.

## Pruebas de Integración

- Registro de intento persiste en `contact_attempts` y emite evento de timeline del caso.
- Agendar sesión persiste en `team_contact` con `is_psico_session=true` y `form_submission_id` nulo.

## Casos Fuera de Alcance (🔴 dependen de specs futuros)

- Contenido del **formulario de atención** (redirect de sesión inmediata).
- Contenido del **formulario de cierre** psicosocial (no acepta / imposibilidad) y la transición efectiva a `cerrado`.
