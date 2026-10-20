# QA — M-08 / M-09 Cierre post-contacto (no competencia · decisión de la persona)

Casos de prueba del paso `postContacto` tras un contacto exitoso: cierre por no competencia ([`flow-M08`](../Flujos/flow-M08-cuando-confirma-cierre-no-competencia.md), `close_reason='co'`) y por decisión de la persona ([`flow-M09`](../Flujos/flow-M09-cuando-confirma-cierre-decision-persona.md), `close_reason='dp'`). M-09 es idéntico a M-08 salvo por el código, el texto del label y el texto del motivo: los casos de backend se **parametrizan sobre ambos códigos** en el mismo flujo de prueba.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md). "Apertura de caso" (M-10) está diferido: aquí solo se verifica que el botón **no** exista.

**Método.** Las precondiciones comunes de `/cerrar` (404, 409 ya cerrado, 401/403, motivo desconocido) están en `M02-C3…C6`. Para llegar a `postContacto` en UI se siembra un intento exitoso por BD (el tramo UI del exitoso ya se prueba en `M06-C3`). Cada test usa su propio reporte.

---

## Casos de prueba

### Backend — `POST /api/v1/victim-contacts/:icode/cerrar` con `co` / `dp`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M08-C1 | Cierre por no competencia | Reporte con 1 intento exitoso; body `{close_reason:'co', close_observation:'  Motivo QA  '}`; sesión `OP1` | 200 `{ok:true, closeReason:'co', category:'Inválido'}`. BD: `close_reason` = enum `co`, `close_observation` = `'Motivo QA'` (sin espacios laterales), `closed_date` ≈ ahora | API + BD |
| M08-C2 | Observación obligatoria | Body con `close_observation` ausente, `""` y `"   "` (para `co`) | 400 "la observación es obligatoria para este motivo" en los tres; la BD no cambia | API + BD |
| M08-C3 | Límite de longitud (Q6) | Observación de exactamente 500 caracteres y otra de 501 | 500 → 200 y persistida completa; 501 → 400 (no un error 500 por desbordar `varchar(500)`); en el caso de 501 la BD no cambia | API + BD |
| M08-C4 | Caracteres especiales | Observaciones: `O'Brien Núñez`, `"comillas dobles"`, `'simples'`, emoji, un salto de línea y `<script>alert(1)</script>` | 200 en todas; BD con el texto idéntico (solo trim externo); `GET` del detalle devuelve el mismo texto; la tabla `victim_contact` queda intacta (consultas parametrizadas) | API + BD |
| M09-C1 | Cierre por decisión de la persona | Igual que M08-C1 con `close_reason:'dp'` | 200 `closeReason:'dp'`, `category:'Inválido'`; BD con enum `dp` y la observación persistida | API + BD |
| M09-C2 | Observación obligatoria en `dp` | Igual que M08-C2 con `'dp'` | 400 en los tres casos; BD sin cambios | API + BD |
| M09-C3 | Límite y caracteres especiales en `dp` | Igual que M08-C3 y M08-C4 con `'dp'` | Mismos resultados (500 ok / 501 → 400 / texto idéntico) | API + BD |

### UI — paso `postContacto`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M08-C5 | Opciones disponibles | Modal en `postContacto` | Solo "Cierre por no competencia" y "Cierre por decisión de la persona". **No existe** "Apertura de caso" (`toHaveCount(0)`) | UI |
| M08-C6 | Abrir la confirmación de no competencia | Clic en "Cierre por no competencia" | Se ocultan las dos opciones y aparece la caja de confirmación con el label "Registre el motivo de no competencia:", un textarea vacío y los botones "Cancelar" y "Confirmar cierre" (deshabilitado) | UI |
| M08-C7 | Habilitación de "Confirmar cierre" | Caja de confirmación abierta | Vacío → deshabilitado; solo espacios → deshabilitado; con texto → habilitado; al borrar el texto → deshabilitado de nuevo | UI |
| M08-C8 | "Cancelar" | Caja abierta con texto escrito | Vuelven las dos opciones; no sale ninguna request de escritura; la BD no cambia | UI + requests + BD |
| M08-C9 | Confirmar el cierre por no competencia | Texto "Motivo QA" escrito | Un POST con body `{close_reason:'co', close_observation:'Motivo QA'}`; paso `cerrado`: "Categoría: Inválido", "Motivo: Inválido por Cierre por no competencia — Motivo QA", ícono gris y botón "Cerrar" | UI + request interceptada |
| M08-C10 | La vista local coincide con la cargada | Tras M08-C9, cerrar el modal y reabrir | El "Motivo:" es **idéntico** al que se mostró justo tras confirmar (la vista armada en el cliente = la que arma M-01 desde backend con `closeReasonLabel` + `closeObservation`) | UI |
| M08-C11 | Texto con HTML y emoji | Observación `<b>negrita</b> <script>window.__x=1</script> 😀` | Se muestra como texto plano en `cerrado` (sin negrita, sin ejecutar el script: `window.__x` indefinido), también al reabrir; el layout no se rompe | UI |
| M08-C12 | Cerrar el modal sin confirmar | Caja abierta con texto escrito; cerrar el modal con "✕" | No se persiste nada (BD sin `close_reason`); al reabrir, el modal está en `postContacto` con las dos opciones y sin caja ni texto | UI + BD |
| M08-C13 | Falla el guardado (Q3) | `page.route` responde 409 y, en otra corrida, 500 a `/cerrar` | No pasa a `cerrado`; el botón se reactiva; el modal se puede cerrar; la BD no cambia | UI + request simulada + BD |
| M08-C14 | Doble clic en "Confirmar cierre" | Caja con texto | Exactamente 1 POST a `/cerrar` | UI + conteo de requests |
| M09-C4 | Cierre por decisión de la persona de punta a punta | Clic en "Cierre por decisión de la persona" | Label "Registre la manifestación de la persona:"; al confirmar, body `{close_reason:'dp', close_observation:<texto>}` y "Motivo: Inválido por Cierre por decisión de la persona — <texto>"; reabrir muestra el mismo texto | UI + request interceptada |
| M09-C5 | Alternar entre las dos opciones | Elegir "no competencia", pulsar "Cancelar" y elegir "decisión de la persona" | El label corresponde a la opción vigente ("manifestación de la persona") y el cierre usa `dp`, no `co` | UI + request interceptada |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-M08-1` — `/cerrar` con `co` y `dp`: observación obligatoria, límite y caracteres especiales | M08-C1, C2, C3, C4, M09-C1, C2, C3 | `gestion-reportes-api/cerrar.spec.ts` |
| `TU-M08-1` — no competencia de punta a punta: opciones, habilitación, cierre y persistencia | M08-C5, C6, C7, C9, C10 | `gestion-reporte-modal/post-contacto.spec.ts` |
| `TU-M08-2` — cancelar, cerrar sin confirmar y texto con HTML | M08-C8, C11, C12 | `gestion-reporte-modal/post-contacto.spec.ts` |
| `TU-M08-3` — falla de guardado y doble clic | M08-C13, C14 | `gestion-reporte-modal/post-contacto.spec.ts` |
| `TU-M09-1` — decisión de la persona de punta a punta y alternancia de opciones | M09-C4, C5 | `gestion-reporte-modal/post-contacto.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| M08-C1 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M08-C2 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M08-C3 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | BUG-08 |
| M08-C4 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M09-C1 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M09-C2 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | — |
| M09-C3 | ✅ | `gestion-reportes-api/cerrar.spec.ts` | BUG-08 |
| M08-C5 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C6 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C7 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C8 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C9 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C10 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | BUG-09 |
| M08-C11 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C12 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C13 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M08-C14 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |
| M09-C4 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | BUG-09 |
| M09-C5 | ✅ | `gestion-reporte-modal/post-contacto.spec.ts` | — |

**Resumen:** 19 de 19 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
