---
okf_version: "1.0"
type: Test_Plan
title: "Plan de Pruebas: Consulta de Caso para Agente Psicosocial"
description: "Estrategia de testing y matriz de casos de prueba para el acceso de solo lectura de ps/ts a la pantalla CaseDetail, sin regresión para sv/ro/op."
owner: "@qa-squad"
status: active
tags: [qa, testing, psicosocial, case-detail, rbac]
dependencies:
  - Spec: 3-features/ConsultaCasoPsicosocial/index.md
last_updated: "2026-07-03"
---

# Estrategia de Testing (Consulta de Caso Psicosocial)

## Pruebas de API (Backend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Playwright Script (`.spec.ts`) |
| :--- | :--- | :--- | :--- |
| `TC-CCP-API-01` | `GET /salvia/casos/:id/detalle` con sesión rol `ps` | `200 OK`, shell HTML se renderiza (hoy da `403`/redirección) | `tests/api/case-detail-psicosocial.spec.ts` |
| `TC-CCP-API-02` | `GET /salvia/casos/:id/detalle` con sesión rol `ts` | `200 OK` | `tests/api/case-detail-psicosocial.spec.ts` |
| `TC-CCP-API-03` (regresión) | `GET /salvia/casos/:id/detalle` con sesión rol `en` (Enlace territorial, sin este permiso) | Sigue denegado — el cambio en `Menu.go` no debe otorgar el permiso a roles no solicitados | `tests/api/case-detail-psicosocial.spec.ts` |
| `TC-CCP-API-04` (regresión) | `GET /salvia/casos/:id/detalle` con sesión rol `sv`, `ro`, `op` | Sigue en `200 OK` como hoy | `tests/api/case-detail-psicosocial.spec.ts` |

## Pruebas E2E / UI (Frontend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Playwright Script (`.spec.ts`) |
| :--- | :--- | :--- | :--- |
| `TC-CCP-UI-01` | `ps` entra por "Ver caso" desde Mis Remisiones Psicosocial | Navega a `/salvia/casos/:id/detalle`, carga sin error | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-02` | `ts` entra por "Ver caso" desde Mis Remisiones Psicosocial | Igual que `TC-CCP-UI-01` | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-03` | `ps`/`ts` en CaseDetail | Las 6 tabs (Info General, Derivaciones, Barreras, Seguimientos, Tareas, Timeline) están presentes y muestran datos | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-04` | `ps`/`ts` en CaseDetail | Botón "Reasignar" (caso) **ausente** en el DOM | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-05` | `ps`/`ts` en CaseDetail | Botón "+ Nuevo seguimiento" **ausente** en el DOM | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-06` | `ps`/`ts` en tab Seguimientos, con un seguimiento pendiente | Controles "▶ Iniciar", "✏️ Editar", "📋 Reasignar" **ausentes** para esa fila | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-07` | `ps`/`ts` en tab Tareas (y en Barreras al expandir una barrera con tareas) | Botón "Gestionar" **ausente**, incluso para una tarea asignada al `icode` del usuario `ps`/`ts` de prueba | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-08` | `ps`/`ts` en tab Derivaciones | Widgets "🧪 Dev — Probar case-task-modal/history" **ausentes** | `tests/e2e/case-detail-psicosocial.spec.ts` |
| `TC-CCP-UI-09` (regresión) | `sv` en CaseDetail | "Reasignar" (caso), "+ Nuevo seguimiento", "✏️ Editar"/"📋 Reasignar" seguimiento siguen visibles y funcionales | `tests/e2e/case-detail-regression.spec.ts` |
| `TC-CCP-UI-10` (regresión) | `ro` en CaseDetail, con seguimiento propio pendiente | "▶ Iniciar" visible y funcional; tarea propia asignada muestra "Gestionar" | `tests/e2e/case-detail-regression.spec.ts` |
| `TC-CCP-UI-11` (regresión) | `op` en CaseDetail, con seguimiento propio pendiente | "▶ Iniciar" y "📋 Reasignar" (dueño) visibles y funcionales | `tests/e2e/case-detail-regression.spec.ts` |

## Pruebas Unitarias (Jest/Vitest)

- `esConsultaSoloLectura` (computed en `get_case_detail_sv.html`) retorna `true` solo para `userRole === 'ps'` y `userRole === 'ts'`, `false` para cualquier otro rol.
- `puedeCompletar(task)` (`case-tasks.js`) retorna `false` para `userRole === 'ps'`/`'ts'` **independientemente** del valor de `task.assignedUserId`, incluso si coincide con `userId`.

## Pruebas de Integración (Supertest)

- N/A — este spec no modifica ningún endpoint del API JSON (ver deuda técnica reconocida en el [index.md](../index.md) del feature).

## Evidencia de Ejecución (obligatoria antes de marcar el feature como Done)

Según AGENTS.md §7, actualizar esta tabla con `PASS`/`FAIL`, fecha y comando exacto al ejecutar cada caso — no se acepta "los tests pasan" sin evidencia.

| ID Caso | Resultado | Fecha | Comando |
| :--- | :---: | :--- | :--- |
| `TC-CCP-API-01` … `TC-CCP-API-04` | 🔴 Bloqueado | 2026-07-03 | `PORT=8080 go run main.go` — el pool de conexiones de Supabase "develop" estaba en su límite (`max clients reached in session mode - pool_size: 15`). Se detuvo el servidor de inmediato para no sumar carga a un recurso compartido ya saturado; no se llegó a levantar sesión real de `ps`/`ts` vía navegador. |
| `TC-CCP-UI-01` … `TC-CCP-UI-11` | 🔴 Bloqueado | 2026-07-03 | Misma causa — sin servidor arriba no hay UI que ejercitar en navegador/Playwright. |
| Verificación equivalente de lógica (no sustituye Playwright, pero valida el código real) | ✅ PASS | 2026-07-03 | `go test ./salvia/config/... -run TestGetCaseDetailSvPermission_TEMP -v` (test temporal contra `PermissionsByRole["get_case_detail_sv"]` real: confirma `ps`/`ts` habilitados y `ad`/`en`/`an` siguen denegados; eliminado tras pasar) |
| Verificación equivalente — `esConsultaSoloLectura` y `puedeCompletar` | ✅ PASS | 2026-07-03 | `node -e '...'` reproduciendo exactamente la lógica añadida en `get_case_detail_sv.html`/`case-tasks.js`: `esConsultaSoloLectura` es `true` únicamente para `ps`/`ts` (11 roles probados); `puedeCompletar` retorna `false` para `ps`/`ts` incluso cuando `task.assignedUserId === userId` (caso adversarial que motivó el cambio) |
| Chequeo de sintaxis | ✅ PASS | 2026-07-03 | `node --check` sobre `case-tasks.js` y sobre el `<script>` principal de `get_case_detail_sv.html` (con placeholders `{{...}}` de Go template neutralizados) |
| `go build` del paquete tocado | ✅ PASS | 2026-07-03 | `go build ./salvia/...` desde `src/` |

> [!WARNING]
> **Pendiente real antes de dar el feature por terminado**: falta ejecutar `TC-CCP-API-01…04` y `TC-CCP-UI-01…11` contra un servidor realmente arriba (con sesión de un usuario `ps`/`ts` de prueba) — la verificación de lógica de arriba confirma que el código es correcto en aislamiento, pero no reemplaza probar el flujo real end-to-end (carga de página, render de tabs, ausencia de botones en el DOM). Reintentar cuando el pool de `develop` tenga cupo, o correr contra una base local/aislada si el proyecto habilita una.
