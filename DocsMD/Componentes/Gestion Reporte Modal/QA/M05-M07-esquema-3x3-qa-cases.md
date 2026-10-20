# QA — M-05 / M-06 / M-07 Esquema de recontacto 3x3

Casos de prueba del paso `esquema3x3`: registrar un intento fallido ([`flow-M05`](../Flujos/flow-M05-cuando-registra-intento-fallido.md)), registrar un intento exitoso ([`flow-M06`](../Flujos/flow-M06-cuando-registra-intento-exitoso.md)) y cerrar por no contacto ([`flow-M07`](../Flujos/flow-M07-cuando-cierra-por-no-contacto.md)). Es el núcleo de la regla de negocio: **"día completo" = día (Bogotá) con ≥ 3 intentos fallidos**; el cierre por no contacto exige **3 días completos** (mínimo 9 fallidos) y nunca se dispara solo.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md).

**Método.** El esquema completo (3 días × 3) **solo es alcanzable sembrando intentos de días pasados por BD** (`sembrarIntentos`): la API siempre registra con la hora actual. Por eso los casos de backend siembran el estado previo y registran por API solo el último intento; los de UI hacen igual y ejercitan por la interfaz el tramo final. Las precondiciones comunes de `/cerrar` (404, 409 ya cerrado, 401/403, motivo inválido) están en `M02-C3…C6` y no se repiten para `nc`. Cada test usa un reporte propio.

---

## Casos de prueba

### Backend — `POST /api/v1/victim-contacts/:icode/intento`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M05-C1 | Primer intento fallido del día | Reporte abierto sin intentos; sesión `OP1` | 200 `{ok:true, intentosHoy:1, diasCompletos:0, quotaReached:false}`. BD: 1 fila con `result='fallido'`, `attempt_agent` = icode de `OP1`, `date` ≈ ahora (±1 min, según el convenio del Paso 0), `deleted_at` NULL, `i_code` único | API + BD |
| M05-C2 | Tercer fallido del día | 2 fallidos hoy; registrar el 3° | 200 con `intentosHoy:3`, `quotaReached:true`, `diasCompletos:1`. El reporte sigue abierto (`close_reason` NULL) | API + BD |
| M05-C3 | 4° y 5° fallido el mismo día | 3 fallidos hoy; registrar 2 más | 200 en ambos (sin límite duro); `intentosHoy` 4 y luego 5; `diasCompletos` se mantiene en 1 (un día no cuenta dos veces) | API + BD |
| M05-C4 | El esquema completo no cierra el reporte | 2 días completos sembrados (3+3) y 2 fallidos hoy; registrar el 3° | `diasCompletos:3`. BD: `close_reason` y `closed_date` siguen NULL — solo se habilita M-07 | API + BD |
| M05-C5 | Un día incompleto no suma | 2 días completos sembrados y 1 fallido hoy; registrar el 2° | `diasCompletos:2`, `intentosHoy:2`, `quotaReached:false` | API |
| M05-C6 | Reporte cerrado | Reporte cerrado por BD | 409; no se inserta ninguna fila | API + BD |
| M05-C7 | Inexistente, sin sesión y roles sin acceso | UUID aleatorio; sin cookie; sesión `RO`, `PS`, `AD` | 404; 401; distinto de 200 (401/403, Q8) — y en ningún caso se inserta una fila | API + BD |
| M05-C8 | `result` inválido (Q4) | Body `{}`, `{result:''}`, `{result:'otro'}` | 400 en los tres; no se inserta ninguna fila | API + BD |
| M05-C9 | Gestión de reportes ajenos (GAP-2) | Reporte asignado a `OP2`; registrar con `SV` y con `OP1` | 200 en ambos; `attempt_agent` es quien registró, no el agente asignado | API + BD |
| M06-C1 | Intento exitoso | Reporte abierto con 2 fallidos hoy; registrar `exitoso` | 200 `{ok:true}`. BD: 1 fila `result='exitoso'`. El detalle posterior mantiene `intentosHoy:2` y `diasCompletos:0` (el exitoso no cuenta en los contadores del 3x3) | API + BD |
| M06-C2 | Exitoso sobre reporte cerrado | Reporte cerrado por BD | 409; no se inserta ninguna fila | API + BD |
| M05-C21 | Reporte legado invalidado (`victim_contact_status='i'`) — **caso nuevo** | Reporte con status `i` y sin `close_reason`; intento `fallido` y `/cerrar` con `ni` | 409 en ambos (se lee como cerrado); 0 filas en `victim_contact_attempt`; no aparece un motivo de cierre nuevo sobre el legado | API + BD |

### Backend — `POST /api/v1/victim-contacts/:icode/cerrar` con `close_reason='nc'`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M07-C1 | Cierre con el esquema completo | 3 días × 3 fallidos sembrados (9 en total) | 200 `{ok:true, closeReason:'nc', category:'Inválido'}`. BD: `close_reason` = enum `nc`, `closed_date` ≈ ahora, `close_observation` NULL | API + BD |
| M07-C2 | Esquema incompleto → 409 | Un reporte por escenario: (a) 0 intentos; (b) 5+4 fallidos en 2 días (9 en total); (c) 3+3+2 (8); (d) 3+3+2 fallidos + 1 exitoso (9 intentos, 8 fallidos); (e) 3+3+3 con 1 borrado lógicamente (8 válidos) | 409 "el esquema 3x3 no está completo todavía" en los cinco; la BD no cambia. Es la defensa server-side del botón oculto | API + BD (tabla de escenarios) |
| M07-C3 | Intentos extra no estorban | 4+3+3 y 5+3+4 fallidos en 3 días | 200 en ambos (cuentan los días con ≥ 3, no el total exacto) | API + BD |
| M07-C4 | Frontera de día Bogotá en el cierre | 3 fallidos a las 23:50 Bogotá de D−3, 3 a las 00:10 Bogotá de D−2 (mismo día UTC que los anteriores) y 3 a mediodía de D−1 | 200: son 3 días Bogotá aunque solo 2 días UTC | API + BD |

### UI — paso `esquema3x3`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M05-C10 | Primer fallido por UI | Modal en `esquema3x3` sin intentos | Un POST a `/intento` con body `{result:'fallido'}`. Contador "Hoy: 1/3 · Días completos: 0/3"; el log muestra el grupo de hoy con badge "1/3" y una fila con hora `HH:MM` (Bogotá, ±2 min de la actual), "Fallido" y el nombre del agente; sin banners | UI + request interceptada |
| M05-C11 | Tercer fallido del día (cuota) | 2 fallidos hoy, sin días previos | Aparece el banner de cuota ("Puede continuar hoy o retomar mañana…"); badge del día "3/3" con la clase `cc-day--completo`; el banner de esquema completo **no** aparece (`diasCompletos` = 1); el reporte sigue abierto | UI |
| M05-C12 | Cuarto fallido del mismo día | 3 fallidos hoy | Los botones siguen activos; se agrega la fila; contador "Hoy: 4/3 · Días completos: 1/3"; badge "4/3" con la clase de completo | UI |
| M05-C13 | Esquema completo: banner | 2 días completos sembrados (3+3) y 2 fallidos hoy; registrar el 3° por UI | Aparece "Se completó el esquema 3x3 (3 intentos en 3 días distintos)" con el botón "Cerrar por no contacto"; el banner de cuota **no** aparece; "Intento Fallido/Exitoso" siguen disponibles; BD: reporte aún abierto | UI + BD |
| M05-C14 | Log agrupado por día | Fallidos sembrados: 3 en D−1, 2 en D−2, 1 en D−3 | 3 grupos en orden descendente (D−1 arriba); badges "3/3" (con clase de completo), "2/3" y "1/3" (sin ella); cada fila con hora, resultado y agente | UI + siembra por BD |
| M05-C15 | Zona horaria en el log | Un fallido sembrado a las 02:15 UTC (= 21:15 Bogotá del día anterior) | Se muestra en el grupo del día Bogotá anterior con hora "21:15" (no "02:15" ni el día UTC) | UI + siembra por BD |
| M05-C16 | Persistencia | Registrar 2 fallidos por UI, cerrar el modal y reabrir | Mismos contadores y mismas 2 filas del log (los datos vienen del backend, no del estado local) | UI |
| M05-C17 | Falla el guardado (Q3) | `page.route` responde 500 a `/intento`, con `fallido` y con `exitoso` | No se agrega fila ni cambian los contadores; con `exitoso` **no** pasa a `postContacto`; los botones se reactivan; la BD no cambia | UI + request simulada + BD |
| M05-C18 | Doble clic rápido en "Intento Fallido" (Q3) | Modal en `esquema3x3` | Una sola fila y un solo POST (el guardado en curso bloquea el segundo clic). Si salieran dos, es hallazgo: el flujo marca `saving=true` pero no dice que deshabilite los botones | UI + conteo de requests + BD |
| M05-C20 | El log usa siempre la zona de Bogotá aunque el navegador esté en otra zona — **caso nuevo** | 3 fallidos en el MISMO día Bogotá (D−2: 09:30, 09:40, 10:30) que en Tokio (UTC+9) caen a ambos lados de la medianoche; navegador con `timezoneId='Asia/Tokyo'` | Contador "Hoy: 0/3 · Días completos: 1/3"; UN solo grupo con badge "3/3" completo y fecha del día Bogotá; horas "9:30 a. m.", "9:40 a. m.", "10:30 a. m." | UI + siembra por BD |
| M05-C19 | `sv` registra en un reporte ajeno (GAP-2) | Reporte asignado a `OP2`; sesión `SV` | Los dos botones están disponibles; la fila registrada muestra el nombre de `SV` como agente | UI |
| M06-C3 | Intento exitoso por UI | Modal en `esquema3x3` | Un POST con body `{result:'exitoso'}`; tras una breve pausa (~400 ms) pasa a `postContacto`: "Contacto establecido", texto de ayuda y **2 botones** de decisión ("Cierre por no competencia", "Cierre por decisión de la persona"); sin recuadro de confirmación abierto | UI + request interceptada |
| M06-C4 | Persistencia del contacto | Tras M06-C3 sin decidir: cerrar el modal y reabrir | Abre en `postContacto` (no en `esquema3x3`). BD: 1 fila `exitoso` y `close_reason` NULL | UI + BD |
| M06-C5 | El exitoso interrumpe el esquema | (a) Esquema 3x3 completo (3 días × 3) y se registra un exitoso; (b) exitoso como primer intento, sin fallidos previos | En ambos pasa a `postContacto`; no se ve el banner de "Cerrar por no contacto" | UI + siembra por BD |
| M06-C6 | `postContacto` no muestra el esquema | Tras 2 fallidos y 1 exitoso | No hay contadores "Hoy …/3" ni log de intentos ni botones de intento en `postContacto` | UI |
| M07-C5 | Umbral del botón "Cerrar por no contacto" | 2 días completos sembrados y 1 fallido hoy (7 en total) | Con 8 fallidos (3+3+2) el botón no existe (`toHaveCount(0)`); el 9° (3° del día 3) lo hace aparecer | UI + siembra por BD |
| M07-C6 | Cierre por no contacto de punta a punta | 2 días completos sembrados (3+3); registrar 3 fallidos hoy por UI y pulsar "Cerrar por no contacto" | Un POST a `/cerrar` con body `{close_reason:'nc'}`; paso `cerrado`: "Categoría: Inválido", "Motivo: Inválido por Cierre por no contacto"; BD con `close_reason=nc` | UI + request interceptada + BD |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-M05-1` — intento fallido: contadores, cuota blanda y no-cierre | M05-C1, C2, C3, C4, C5 | `gestion-reportes-api/intentos.spec.ts` |
| `TB-M05-2` — intento: precondiciones, validación y permisos | M05-C6, C7, C8, C9, C21, M06-C2 | `gestion-reportes-api/intentos.spec.ts` |
| `TB-M06-1` — intento exitoso | M06-C1 | `gestion-reportes-api/intentos.spec.ts` |
| `TB-M07-1` — cerrar `nc` con esquema completo (incluye frontera Bogotá) | M07-C1, C3, C4 | `gestion-reportes-api/cerrar.spec.ts` |
| `TB-M07-2` — cerrar `nc` con esquema incompleto (409) | M07-C2 | `gestion-reportes-api/cerrar.spec.ts` |
| `TU-M05-1` — fallidos por UI: contadores y cuota | M05-C10, C11, C12 | `gestion-reporte-modal/esquema-3x3.spec.ts` |
| `TU-M05-2` — esquema completo: banner, umbral y cierre `nc` de punta a punta | M05-C13, M07-C5, M07-C6 | `gestion-reporte-modal/esquema-3x3.spec.ts` |
| `TU-M05-3` — log agrupado, zona horaria y persistencia | M05-C14, C15, C16, C20 | `gestion-reporte-modal/esquema-3x3.spec.ts` |
| `TU-M05-4` — falla de guardado y doble clic | M05-C17, C18 | `gestion-reporte-modal/esquema-3x3.spec.ts` |
| `TU-M05-5` — `sv` gestiona un reporte ajeno | M05-C19 | `gestion-reporte-modal/esquema-3x3.spec.ts` |
| `TU-M06-1` — exitoso → `postContacto`, reapertura e interrupción del esquema | M06-C3, C4, C5, C6 | `gestion-reporte-modal/esquema-3x3.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| M05-C1 | ✅ | `gestion-reportes-api/intentos.spec.ts` | BUG-06 |
| M05-C2 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C3 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C4 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C5 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C6 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C7 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C8 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C9 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M06-C1 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M06-C2 | ✅ | `gestion-reportes-api/intentos.spec.ts` | — |
| M05-C21 | ✅ | `gestion-reportes-api/intentos.spec.ts` | BUG-07 |
| M07-C1 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M07-C2 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M07-C3 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M07-C4 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | BUG-05 |
| M05-C10 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | BUG-10 |
| M05-C11 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M05-C12 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M05-C13 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M05-C14 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | BUG-15 |
| M05-C15 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | BUG-15 |
| M05-C16 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M05-C17 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M05-C18 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M05-C20 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | BUG-15 |
| M05-C19 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | BUG-10 |
| M06-C3 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M06-C4 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M06-C5 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M06-C6 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M07-C5 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |
| M07-C6 | ✅ | `gestion-reporte-modal/esquema-3x3.spec.ts` | — |

**Resumen:** 33 de 33 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
