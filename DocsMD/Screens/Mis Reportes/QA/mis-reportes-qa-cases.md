# QA — Mis Reportes · E-01 Carga y E-02 Ver reporte

Casos de prueba de [`flow-E01`](../Flujos/flow-E01-cuando-carga-pantalla.md) y [`flow-E02`](../Flujos/flow-E02-cuando-presiona-ver-reporte.md) de `GET /salvia/mis-reportes` (solo rol `op`). La pantalla es la composición de `Listado Reportes` **acotada a los reportes del usuario**, sin reasignación: aquí solo se prueba lo que **cambia** respecto al listado (scope fijo, columnas y filtros que sobran, acceso) y que ese scope **sobrevive a cualquier filtro**. Filtros, búsqueda y orden en sí ya están cubiertos en [`E02-E03-E07-E08`](../../Listado%20Reportes/QA/E02-E03-E07-E08-filtros-busqueda-orden-qa-cases.md) y no se repiten; lo que ocurre dentro del modal, en los archivos de `gestion-reporte-modal`.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](../../Listado%20Reportes/QA/gestion-reportes-qa-plan.md). Los alias `R1…R10` son los del pool de la sección B del plan; el conjunto de "Mis Reportes de `OP1`" es R2, R4, R5, R6 y R10.

**Ruta.** La implementada es `/salvia/mis-reportes` (el diseño decía `/salvia/reportes/mis-reportes`); el usuario confirmó dejarla así el 2026-10-13. Los MD de diseño de Mis Reportes (`interface`, `events`, `flow-E01`) aún citan la ruta anterior.

**Método.** Backend primero (cómo compone el listado el scope `persona_asignada` con los demás parámetros), después UI con la sesión `OP1`, aislando el pool con el marcador `QAGR-<runId>` — `TB-LR-E07-1` corre antes que los `TU-MR-*`. Los casos que modifican datos (`MR-E01-C8`, `MR-E02-*`) usan reportes propios del test asignados a `OP1`, no los del pool.

**Riesgo principal (GAP-5, Q1).** `flow-E02` de Listado Reportes (PASO 1) hacía `activeFilter = {key:'estado_reporte', value}`, que **reemplazaría** el filtro base `persona_asignada` y mostraría reportes de otros agentes. Los casos `MR-E01-C5` y `MR-E01-C6` existen para detectar exactamente eso.

---

## Casos de prueba

### Backend — `GET /api/v1/victim-contacts/list` con scope de agente

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| MR-E01-C7 | El scope se compone con estado, búsqueda y paginación | Pool; `filter_key=persona_asignada&filter_value=<OP1>`, con sesión `OP1` y con sesión `SV` | (a) + `search=<marcador>` → R2, R4, R5, R6, R10 (`total` = 5); (b) + `filter_estado_reporte=cerrado_invalido` → R4, R5, R10 (no R8: sin agente); (c) + `filter_estado_reporte=sin_asignar` → vacío; (d) + `page_size=2` → 3 páginas (2+2+1) con solo reportes de `OP1`. Nunca aparece R3 (de `OP2`) ni un reporte sin agente | API |

### UI — acceso y estructura (E-01)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| MR-E01-C1 | `op` abre la pantalla | Sesión `OP1` en `/salvia/mis-reportes` | 200; encabezado "Mis Reportes"; columnas "Víctima / ID Reporte", "Tipo", "Esquema 3x3", "Estado", "Fecha Reporte" — **sin** "Asignado a"; sin casillas de selección ni botón de reasignación; filtros: solo "Estado del reporte" y la búsqueda — **sin** "Asignación"; el orden ofrece solo "Fecha de Reporte DESC/ASC" | UI |
| MR-E01-C2 | Acceso restringido a `op` | Sesión `SV`, `RO`, `PS`, `AD` y sin sesión | `sv` también es rechazado (la pantalla es solo `op`): distinto de 200 (401/403, Q8) en `page.request.get`; sin sesión → login | UI (request) |
| MR-E01-C3 | Solo mis reportes | Pool aislado, sesión `OP1` | Se ven R2, R4, R5, R6 y R10 (de cualquier estado, también cerrados); no se ven R3 (de `OP2`) ni R1, R8, R9 (sin agente) | UI |
| MR-E01-C4 | El scope sale de la sesión | Misma pantalla; abrirla también con un query string manipulado (`?filter_value=<icode de OP2>`, `?agent=<icode de OP2>`) | La request de listado lleva `filter_key=persona_asignada` y `filter_value` = icode de `OP1` en ambos casos; las filas son las mismas | UI + request interceptada |
| MR-E01-C5 | El scope sobrevive al filtro de estado | Pool aislado; elegir en "Estado del reporte": "Asignado", "Cerrado - inválido", "Cerrado - válido", "Sin asignar" y volver a "Todos" | "Asignado" → solo R2 (no R3); "Cerrado - inválido" → R4, R5, R10 (no R8); "Cerrado - válido" → R6; "Sin asignar" → vacío; "Todos" vuelve a los 5 propios, **no** a todos los reportes del sistema | UI + request interceptada |
| MR-E01-C6 | El scope sobrevive a la búsqueda y al orden | Pool aislado | Buscar el nombre de R3 (de `OP2`) → sin resultados; buscar el de R2 → solo R2; cambiar el orden a ASC y DESC reordena únicamente esas 5 filas; en ningún caso aparece un reporte ajeno | UI |
| MR-E01-C9 | Entrada "Mis Reportes" en el menú lateral — **caso nuevo, decisión del usuario 2026-10-13** | Sesión `OP1` y `SV` en `/salvia/lista-casos` | `op` tiene `<a href="/salvia/mis-reportes">Mis Reportes</a>` y llega a la pantalla; `sv` no la tiene | UI |
| MR-E01-C8 | La reasignación mueve el reporte de lista | Reporte propio `RX` asignado a `OP1`; `SV` lo reasigna a `OP2` por API; `OP1` recarga | `RX` desaparece de la lista de `OP1`. Y el caso inverso: un reporte sin asignar que `SV` reasigna a `OP1` aparece al recargar | UI + API |

### UI — ver y gestionar desde Mis Reportes (E-02)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| MR-E02-C1 | "Ver Reporte" abre el modal con acciones | Reporte propio sin intentos; clic en "Ver Reporte" | Abre el modal en `validacionPrevia` con "No información mínima" y "Procede recontacto" habilitados (op gestiona sin modo solo lectura) | UI |
| MR-E02-C2 | Un intento se refleja en la fila | Reporte propio sin intentos ("Sin intentos"): "Procede recontacto" → "Intento Fallido" → cerrar el modal | La fila pasa a "1/9 intentos · 0/3 días" sin recargar la página a mano (el modal avisa `updated` y la pantalla recarga) | UI |
| MR-E02-C3 | Cerrar un reporte no lo saca de la lista | Reporte propio abierto; cerrarlo por "No información mínima" y cerrar el modal | La fila sigue en la lista con el estado "Cerrado - inválido" (el scope es por asignación, no por estado) | UI |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-MR-E01-1` — el listado compone el scope con estado, búsqueda y paginación | MR-E01-C7 | `gestion-reportes-api/listado.spec.ts` |
| `TU-MR-E01-1` — acceso y estructura de la pantalla | MR-E01-C1, C2 | `mis-reportes/mis-reportes.spec.ts` |
| `TU-MR-E01-2` — solo los reportes del usuario, scope tomado de la sesión | MR-E01-C3, C4 | `mis-reportes/mis-reportes.spec.ts` |
| `TU-MR-E01-3` — el scope sobrevive a filtros, búsqueda y orden | MR-E01-C5, C6 | `mis-reportes/mis-reportes.spec.ts` |
| `TU-MR-E01-4` — reasignar saca o mete reportes de la lista | MR-E01-C8 | `mis-reportes/mis-reportes.spec.ts` |
| `TU-MR-E01-5` — entrada en el menú lateral | MR-E01-C9 | `mis-reportes/mis-reportes.spec.ts` |
| `TU-MR-E02-1` — ver y gestionar desde Mis Reportes; la fila se actualiza | MR-E02-C1, C2, C3 | `mis-reportes/mis-reportes.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| MR-E01-C7 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| MR-E01-C1 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E01-C2 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E01-C3 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E01-C4 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E01-C5 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E01-C6 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E01-C9 | ✅ | `mis-reportes/mis-reportes.spec.ts` | — |
| MR-E01-C8 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E02-C1 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E02-C2 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |
| MR-E02-C3 | ✅ | `mis-reportes/mis-reportes.spec.ts` | BUG-13 |

**Resumen:** 12 de 12 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](../../Listado%20Reportes/QA/gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
