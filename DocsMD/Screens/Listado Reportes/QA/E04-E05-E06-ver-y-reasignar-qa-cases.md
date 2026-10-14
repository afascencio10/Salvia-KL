# QA — Listado Reportes · E-04 / E-05 / E-06 Ver reporte y reasignar

Casos de prueba de "Ver Reporte" ([`flow-E04`](../Flujos/flow-E04-cuando-presiona-ver-reporte.md)), del botón de reasignación masiva ([`flow-E05`](../Flujos/flow-E05-cuando-presiona-reasignar-reportes.md)) y de la reasignación en sí, `POST /api/v1/victim-contacts/reasignar-bulk` ([`flow-E06`](../Flujos/flow-E06-cuando-confirma-reasignacion.md)). Lo que ocurre **dentro** del modal "Ver Reporte" está en los archivos del componente `gestion-reporte-modal`.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md).

**Método.** Backend primero (E-06), después UI. **Estos casos modifican asignaciones**, así que usan reportes propios del test, no el pool de filtros: `RA` (sin asignar, sin equipo), `RB` (asignado a `OP2`, con el equipo de `OP2`) y `RC` (asignado a `OP1`, con 2 fallidos sembrados). Se crean con el marcador `QAGR-<runId>` propio de cada spec y se borran al final. El agente destino de las reasignaciones es `OP1`.

---

## Casos de prueba

### Backend — `POST /api/v1/victim-contacts/reasignar-bulk` (E-06)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E06-C1 | Reasignación masiva | Sesión `SV`; body `{report_icodes:[RA, RB], agent_icode:<OP1>}` | 200 `{ok:true, updated:2, skipped:0}`. BD: `victim_contact_agent_id` = `OP1` en ambos y `victim_contact_team` = equipo de `OP1` — también en `RB`, cuyo equipo previo se pisa | API + BD |
| LR-E06-C2 | Un icode inexistente en la lista | Body con `[RA, <UUID aleatorio>]` | 200 con `updated:1`, `skipped:1`; `RA` sí queda asignado | API + BD |
| LR-E06-C3 | Payload incompleto | `report_icodes` vacío o ausente; `agent_icode` vacío o ausente | 400 en los cuatro; la BD no cambia | API + BD |
| LR-E06-C4 | Agente destino inválido | `agent_icode` inexistente; y un agente inactivo (`status ≠ 'e'`) si existe uno en BD | 400 en ambos; la BD no cambia | API + BD |
| LR-E06-C5 | Sin permiso (Q10, Q8) | Sesión `OP1`, `RO`, `PS`, `AD` y sin sesión, con un body válido | Distinto de 200 (401/403) en todos; la BD no cambia. La tabla de permisos dice que `op` no reasigna aunque `flow-E06` no valide el rol de forma explícita | API + BD |
| LR-E06-C6 | Idempotencia | Reasignar `RA` a `OP1` cuando ya está asignado a `OP1` | 200 `updated:1`, sin error | API + BD |
| LR-E06-C7 | Efecto en las vistas del listado | Tras C1: pedir el listado con `asignacion=<OP1>`, `asignacion=<OP2>` y los estados `asignado` / `sin_asignar` | `OP1` incluye `RA` y `RB`; `OP2` ya no incluye `RB`; `asignado` incluye `RA`; `sin_asignar` ya no lo incluye | API |
| LR-E06-C8 | Reasignar no toca la gestión | Reasignar `RC` (2 fallidos) a `OP2` | `intentos` intactos (2 filas), `close_reason` NULL, `diasCompletos` e `intentosHoy` sin cambios en el detalle | API + BD |

### UI — "Ver Reporte" (E-04)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E04-C1 | Abre el modal de esa fila | Pool aislado; clic en "Ver Reporte" de R2 | El diálogo abre con el radicado de R2 y sale `GET /api/v1/victim-contacts/{icode de R2}`; la columna de acciones ofrece un único botón por fila: "Ver Reporte" | UI + request interceptada |
| LR-E04-C2 | Funciona en cualquier estado | "Ver Reporte" de R1 (sin asignar) y de R4 (cerrado) | R1 abre en el paso que le corresponda (`validacionPrevia`); R4 abre en `cerrado`, solo lectura | UI |

### UI — botón de reasignación (E-05)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E05-C1 | Habilitación según la selección | Sesión `SV`, ninguna fila marcada | El botón de reasignación no abre el modal (oculto o deshabilitado); con ≥ 1 fila marcada queda disponible | UI |
| LR-E05-C2 | Etiqueta del botón (Q10) | `SV` con ≥ 1 fila marcada | Dice "Reasignar Reportes". Si dice "Reasignar Casos" (texto genérico del componente) se registra como hallazgo no bloqueante | UI |
| LR-E05-C3 | Selección sin la regla de equipo único | Marcar `RA` (sin equipo), `RB` (equipo de `OP2`) y `RC` (equipo de `OP1`); luego "seleccionar todo" en la página | Todas las filas quedan marcadas y no aparece ningún aviso de "equipos distintos". Depende de `REASSIGN_CROSS_TEAM_CONTINGENCY = true` (`casos-component.js:43`): si pasa a `false`, la regla de equipo único bloquearía reportes sin equipo | UI |
| LR-E05-C4 | Abre el modal con lo seleccionado | `RA` y `RB` marcados; clic en el botón | Se abre `reasignar-casos-modal` indicando que son 2 reportes (cantidad o lista visible; el texto exacto no se asierta) | UI |

### UI — confirmar la reasignación (E-06)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E06-C9 | Reasignación de punta a punta | `SV` marca `RA` y `RB`, abre el modal, elige a `OP1` y confirma | Un POST a `reasignar-bulk` con `{report_icodes:[RA, RB], agent_icode:<OP1>}`; el modal se cierra; la tabla se recarga y las casillas quedan limpias; "Asignado a" muestra el nombre de `OP1` en ambas filas; el estado de `RA` pasa de "Sin asignar" a "Asignado" | UI + request interceptada + BD |
| LR-E06-C10 | Cancelar | Modal de reasignación abierto con un agente elegido; cerrarlo sin confirmar | No sale ningún POST a `reasignar-bulk`; la BD no cambia | UI + requests + BD |
| LR-E06-C11 | Candidatos del modal (GAP-3) | Modal abierto con reportes de equipos distintos | Ofrece solo agentes `op` activos, sin filtrar por el equipo de los reportes: incluye a `OP1` y a `OP2` y no incluye usuarios `sv`/`ro` | UI + BD |
| LR-E06-C12 | Efecto visible en Mis Reportes | Tras C9, iniciar sesión como `OP1` y abrir "Mis Reportes" | `RA` y `RB` aparecen en la lista de `OP1` | UI |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-LR-E06-1` — reasignar: caso feliz, parcial, idempotente, efecto en el listado y sin tocar la gestión | LR-E06-C1, C2, C6, C7, C8 | `gestion-reportes-api/reasignar.spec.ts` |
| `TB-LR-E06-2` — reasignar: validaciones y permisos | LR-E06-C3, C4, C5 | `gestion-reportes-api/reasignar.spec.ts` |
| `TU-LR-E04-1` — "Ver Reporte" desde la fila | LR-E04-C1, C2 | `listado-reportes/ver-y-reasignar.spec.ts` |
| `TU-LR-E05-1` — selección y botón de reasignación | LR-E05-C1, C2, C3, C4 | `listado-reportes/ver-y-reasignar.spec.ts` |
| `TU-LR-E06-1` — reasignar de punta a punta, cancelar, candidatos y efecto en Mis Reportes | LR-E06-C9, C10, C11, C12 | `listado-reportes/ver-y-reasignar.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| LR-E06-C1 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C2 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C3 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C4 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C5 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C6 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C7 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E06-C8 | ✅ | `gestion-reportes-api/reasignar.spec.ts` | — |
| LR-E04-C1 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E04-C2 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E05-C1 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | BUG-14 |
| LR-E05-C2 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E05-C3 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | BUG-14 |
| LR-E05-C4 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E06-C9 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E06-C10 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E06-C11 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | — |
| LR-E06-C12 | ✅ | `listado-reportes/ver-y-reasignar.spec.ts` | BUG-13 |

**Resumen:** 18 de 18 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
