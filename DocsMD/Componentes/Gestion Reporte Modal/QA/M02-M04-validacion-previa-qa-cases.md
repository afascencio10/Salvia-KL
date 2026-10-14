# QA — M-02 "No información mínima" y M-04 "Procede recontacto"

Casos de prueba de los dos botones del paso `validacionPrevia`: [`flow-M02`](../Flujos/flow-M02-cuando-cierra-sin-informacion-minima.md) (escribe en backend) y [`flow-M04`](../Flujos/flow-M04-cuando-procede-recontacto.md) (transición local, sin backend).

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md).

**Método.** Backend primero: `POST /api/v1/victim-contacts/:icode/cerrar` con `close_reason='ni'`. Los casos M02-C3…C6 son **transversales al endpoint `/cerrar`** (existencia, ya cerrado, rol, motivo inválido): se prueban una sola vez aquí con `ni` y **no se repiten** en M07/M08/M09. Después UI con Playwright. Cada test usa su propio reporte `victima` recién creado (cerrarlo es irreversible).

---

## Casos de prueba

### Backend — `POST /api/v1/victim-contacts/:icode/cerrar` con `close_reason='ni'`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M02-C1 | Cierre por "No información mínima" | Reporte abierto sin intentos; sesión `OP1` | 200 `{ok:true, closeReason:'ni', category:'Inválido'}`. BD: `victim_contact_close_reason` apunta al enum con `code='ni'`, `closed_date` ≈ ahora (±1 min, según el convenio del Paso 0 del plan), `close_observation` NULL; 0 filas en `victim_contact_attempt` | API + BD |
| M02-C2 | `sv` cierra un reporte asignado a otro agente (GAP-2) | Reporte asignado a `OP2`; sesión `SV` | 200 con el mismo resultado de C1 | API + BD |
| M02-C3 | Reporte ya cerrado | Reporte cerrado (cualquier motivo, p. ej. `ni` por BD) | 409 "el reporte ya está cerrado"; `closed_date` y `close_reason` no cambian | API + BD |
| M02-C4 | Reporte inexistente | UUID aleatorio | 404 | API |
| M02-C5 | Sin sesión y roles sin acceso | Sin cookie; sesión `RO`, `PS`, `AD` | 401 sin sesión; con roles sin acceso distinto de 200 (401/403, Q8); en todos, la BD no cambia | API + BD |
| M02-C6 | Motivo ausente, vacío o desconocido (Q5) | Body `{}`, `{close_reason:''}`, `{close_reason:'zz'}` | 400 en los tres; la BD no cambia | API + BD |
| M02-C12 | Motivo diferido `ap` (apertura de caso) — **caso nuevo, decisión del usuario 2026-10-13** | Reporte abierto; body `{close_reason:'ap'}`; sesión `OP1`; y sin sesión | 400 con mensaje que menciona `'ap'`; el reporte NO queda "Cerrado - válido" (`close_reason` y `closed_date` siguen NULL); sin sesión sigue siendo 401 | API + BD |

### UI — paso `validacionPrevia`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M02-C7 | Cerrar por "No información mínima" | Modal abierto en `validacionPrevia` | Se emite un POST a `/cerrar` con body `{close_reason:'ni'}`. Pasa al paso `cerrado`: "Reporte {icode} cerrado", "Categoría: Inválido", "Motivo: Inválido por No información Mínima", ícono gris y botón "Cerrar"; ya no hay botones de acción | UI + request interceptada |
| M02-C8 | Doble clic rápido | Modal en `validacionPrevia` | Exactamente 1 POST a `/cerrar` y una sola transición a `cerrado` (el guardado en curso no se duplica) | UI + conteo de requests |
| M02-C9 | Persistencia | Tras C7: cerrar el modal y volver a abrir el mismo reporte | Abre directamente en `cerrado` con el **mismo texto** de "Motivo" y "Categoría" que mostró justo después de cerrar (la vista local coincide con la que arma M-01 desde el backend) | UI |
| M02-C10 | Falla el guardado (Q3) | `page.route` responde 409 y, en otra corrida, 500 a `/cerrar` | No muestra `cerrado`; el modal sigue en `validacionPrevia` con los botones reactivados y se puede cerrar; la BD no cambia | UI + request simulada + BD |
| M02-C11 | Ambos roles pueden gestionar (GAP-2) | Reporte sin asignar; abrir con `SV` y con `OP1` | Los dos ven los 2 botones de `validacionPrevia` habilitados; no hay modo solo lectura | UI |

### UI — M-04 "Procede recontacto"

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M04-C1 | Es una transición local | Modal en `validacionPrevia`; capturar todas las requests | Al pulsar "Procede recontacto" pasa a `esquema3x3` y **no** sale ninguna request de escritura (POST/PUT/PATCH/DELETE). BD: 0 intentos y sin cierre | UI + requests + BD |
| M04-C2 | Estado inicial del esquema 3x3 | Justo después de C1 | Título "Esquema de Recontacto 3x3"; contador "Hoy: 0/3 · Días completos: 0/3"; texto de ayuda de franjas horarias; botones "Intento Fallido" y "Intento Exitoso"; log con "No hay intentos registrados"; sin banner de cuota ni de esquema completo | UI |
| M04-C3 | Cerrar sin registrar intentos | Tras C1, cerrar el modal y reabrir el reporte | Vuelve a `validacionPrevia` (nada se persistió, ver GAP de `flow-M04`) | UI |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-M02-1` — cerrar con `ni` (OP1 y SV) | M02-C1, C2 | `gestion-reportes-api/cerrar.spec.ts` |
| `TB-M02-2` — precondiciones comunes de `/cerrar` (404, 409, 401/403, 400) | M02-C3, C4, C5, C6 | `gestion-reportes-api/cerrar.spec.ts` |
| `TB-M02-3` — motivo diferido `ap` | M02-C12 | `gestion-reportes-api/cerrar.spec.ts` |
| `TU-M02-1` — "No información mínima": cierre, persistencia y roles | M02-C7, C9, C11 | `gestion-reporte-modal/validacion-previa.spec.ts` |
| `TU-M02-2` — doble clic y falla de guardado | M02-C8, C10 | `gestion-reporte-modal/validacion-previa.spec.ts` |
| `TU-M04-1` — "Procede recontacto" sin backend y reapertura | M04-C1, C2, C3 | `gestion-reporte-modal/validacion-previa.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| M02-C1 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | BUG-06 |
| M02-C2 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M02-C3 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M02-C4 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M02-C5 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M02-C6 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M02-C12 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | BUG-17 |
| M02-C7 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M02-C8 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M02-C9 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M02-C10 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M02-C11 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M04-C1 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M04-C2 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |
| M04-C3 | ✅ | `gestion-reporte-modal/validacion-previa.spec.ts` | — |

**Resumen:** 15 de 15 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
