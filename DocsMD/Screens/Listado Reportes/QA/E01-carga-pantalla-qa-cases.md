# QA — Listado Reportes · E-01 Cuando carga la pantalla

Casos de prueba de [`flow-E01`](../Flujos/flow-E01-cuando-carga-pantalla.md): acceso por rol a `GET /salvia/reportes`, estructura de la pantalla (columnas, botón de reasignación), contrato y contenido de `GET /api/v1/victim-contacts/list`, y el enlace nuevo desde `/salvia/primer_contacto`.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md). Los alias `R1…R10` son los del pool de la sección B del plan.

**Método.** Backend primero (contrato y paginación del listado, vía `page.request`), después UI. Los casos de UI aíslan el pool escribiendo el marcador `QAGR-<runId>` en la caja de búsqueda — por eso `TB-LR-E07-1` corre antes (ver plan). Los filtros, la búsqueda y el orden tienen su propio archivo; aquí solo la carga inicial.

---

## Casos de prueba

### Backend — `GET /api/v1/victim-contacts/list` (carga inicial)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E01-C4 | Contrato del listado | Pool cargado; sesión `SV` y luego `OP1`; sin parámetros | 200 `{reports, total, page, pageSize}` con `page` = 1 y `pageSize` = 20. Cada ítem trae las claves de `ReportListItem`: `id`, `i_code`, `names`, `lastNames`, `creationDate`, `reportType`, `totalFallidos`, `diasCompletos`, `reportStatus`, `assignedAgentICode`, `assignedAgentName`, `team`. `reportStatus` ∈ `sin_asignar`, `asignado`, `cerrado_valido`, `cerrado_invalido` | API |
| LR-E01-C5 | Listado sin autorización | Sin cookie; sesión `RO`, `PS`, `AD` | 401 sin sesión; con roles sin acceso distinto de 200 (401/403, Q8); en ningún caso se devuelven reportes ni datos personales | API |
| LR-E01-C10 | Paginación del endpoint | Pool (10 reportes, aislado con el marcador) con `page_size=2` | `total` = 10; página 1 y 2 devuelven reportes distintos, sin duplicados entre páginas; una página fuera de rango responde 200 con `reports: []` | API |

### UI — acceso y estructura

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E01-C1 | `sv` abre la pantalla | Sesión `SV` en `/salvia/reportes` | 200; encabezado "Reportes"; la tabla de `casos-component` muestra las columnas "Víctima / ID Reporte", "Tipo", "Esquema 3x3", "Estado", "Asignado a", "Fecha Reporte" (más la columna de selección y la de acciones); el botón de reasignación está presente | UI |
| LR-E01-C2 | `op` abre la pantalla | Sesión `OP1` en `/salvia/reportes` | 200; mismas columnas que C1; **sin** casillas de selección y **sin** botón de reasignación (`reasignacion=false` para `op`) | UI |
| LR-E01-C3 | Roles sin acceso a la pantalla | Sesión `RO`, `PS`, `AD`; y sin sesión | Distinto de 200 (401/403, Q8) en `page.request.get('/salvia/reportes')`; sin sesión → redirección/401 | UI (request) |
| LR-E01-C9 | Se ven todos los reportes | `OP1` en el listado, con el pool aislado | Aparecen también los reportes asignados a otros agentes (R3 de `OP2`) y los no asignados: el filtro base es "todos", no "los míos" | UI |

### UI — filas del pool

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E01-C6 | Celdas de cada fila | Pool aislado por el marcador | "Tipo" muestra Víctima/Tercero/Familiar/Hombre según el reporte; "Esquema 3x3": R1 → "Sin intentos", R2 → "5/9 intentos · 1/3 días", R5 → "9/9 intentos · 3/3 días"; "Asignado a": nombre del agente (R2) o "Sin asignar" (R1); "Fecha Reporte" con la fecha sembrada; "Víctima / ID Reporte" con nombre y radicado | UI |
| LR-E01-C7 | Estado por reporte | Pool aislado | R1 "Sin asignar", R2 "Asignado", R4 "Cerrado - inválido", R6 "Cerrado - válido" (Q2: R6 sembrado por BD); y el `reportStatus` de la API coincide (`sin_asignar`, `asignado`, `cerrado_invalido`, `cerrado_valido`) | UI + API |
| LR-E01-C8 | Reporte legado invalidado (Q9) | R7: `victim_contact_status='i'`, sin `close_reason` | Se muestra como cerrado/inválido; su botón "Ver Reporte" abre el modal en `cerrado` de solo lectura (ver `M01-C15`) | UI |
| LR-E01-C12 | Reporte sin nombres (Hombre) (Q11) | Reporte `hombre` recién creado, sin marcar (nombres vacíos), orden DESC por fecha | Aparece en la primera página con su radicado; ninguna celda vacía ni con "null"/"undefined"/"0" | UI |

### UI — enlace desde la pantalla legada

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E01-C11 | Enlace `/salvia/primer_contacto` → `/salvia/reportes` (GAP-4) | Sesión `SV` y luego `OP1` en `/salvia/primer_contacto` | El HTML servido trae `href="/salvia/reportes"`; la pantalla legada monta y el enlace se ve; al pulsarlo navega a `/salvia/reportes` y la pantalla carga. Solo se verifica el enlace, no la lógica legada. **Dos pruebas** (HTML y clic) | UI |
| LR-E01-C13 | Entrada "Listado Reportes" en el menú lateral — **caso nuevo, decisión del usuario 2026-10-13** | Sesión `SV`, `OP1` y `RO` en `/salvia/lista-casos` | `sv` y `op` tienen `<a href="/salvia/reportes">Listado Reportes</a>` en `#salviaSidebar` y al pulsarlo llegan a la pantalla; `ro` no la tiene (ni "Mis Reportes") | UI |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-LR-E01-1` — contrato del listado y permisos | LR-E01-C4, C5 | `gestion-reportes-api/listado.spec.ts` |
| `TB-LR-E01-2` — paginación del endpoint | LR-E01-C10 | `gestion-reportes-api/listado.spec.ts` |
| `TU-LR-E01-1` — acceso por rol y estructura de la pantalla | LR-E01-C1, C2, C3, C9 | `listado-reportes/carga.spec.ts` |
| `TU-LR-E01-2` — celdas y estados de las filas del pool | LR-E01-C6, C7, C8, C12 | `listado-reportes/carga.spec.ts` |
| `TU-LR-E01-3` — enlace desde `primer_contacto` (HTML y clic) | LR-E01-C11 | `listado-reportes/carga.spec.ts` |
| `TU-LR-E01-4` — entrada en el menú lateral | LR-E01-C13 | `listado-reportes/carga.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| LR-E01-C4 | ✅ | `gestion-reportes-api/listado.spec.ts` | BUG-02 |
| LR-E01-C5 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E01-C10 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E01-C1 | ✅ | `listado-reportes/carga.spec.ts` | — |
| LR-E01-C2 | ✅ | `listado-reportes/carga.spec.ts` | — |
| LR-E01-C3 | ✅ | `listado-reportes/carga.spec.ts` | — |
| LR-E01-C9 | ✅ | `listado-reportes/carga.spec.ts` | — |
| LR-E01-C6 | ✅ | `gestion-reportes-api/listado.spec.ts`, `listado-reportes/carga.spec.ts` | BUG-02, BUG-12 |
| LR-E01-C7 | ✅ | `gestion-reportes-api/listado.spec.ts`, `listado-reportes/carga.spec.ts` | BUG-07 |
| LR-E01-C8 | ✅ | `listado-reportes/carga.spec.ts` | BUG-07 |
| LR-E01-C12 | ✅ | `listado-reportes/carga.spec.ts` | BUG-02, BUG-12 |
| LR-E01-C11 | ✅ | `listado-reportes/carga.spec.ts` | BUG-16 |
| LR-E01-C13 | ✅ | `listado-reportes/carga.spec.ts` | — |

**Resumen:** 13 de 13 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
