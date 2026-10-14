# QA — M-11 Cuando cierra el modal

Casos de prueba de [`flow-M11`](../Flujos/flow-M11-cuando-cierra-modal.md): las tres formas de cerrar (✕, backdrop, botón "Cerrar" del paso `cerrado`), el bloqueo durante un guardado, la limpieza de estado y el aviso `updated` al padre, que recarga la tabla.

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md).

**Método.** M-11 no escribe en backend: todo se prueba por UI, **dentro de `Listado Reportes`** (el padre que escucha `updated`). Para ver el efecto del evento se observa la request de recarga `GET /api/v1/victim-contacts/list` que dispara el padre. Los casos M11-C7 y M11-C8 verifican que la tabla refleja lo gestionado en el modal. El flujo de `Mis Reportes` equivalente se cubre en `MR-E02-C2` y `MR-E02-C3`.

---

## Casos de prueba

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M11-C1 | Cerrar con "✕" | Modal abierto en cualquier paso | El diálogo desaparece del DOM y el padre recarga la tabla (se observa un `GET /api/v1/victim-contacts/list`) | UI + request interceptada |
| M11-C2 | Cerrar con el backdrop | Modal abierto | Clic en el fondo oscuro cierra; clic dentro del diálogo (cuerpo o header) **no** cierra | UI |
| M11-C3 | Botón "Cerrar" del paso `cerrado` | Reporte cerrado (p. ej. por `ni` desde la UI) | Cierra el modal y dispara la recarga de la tabla (1 `GET` del listado) | UI + request interceptada |
| M11-C4 | `updated` se emite siempre | Abrir un reporte y cerrarlo sin tocar nada | Igualmente ocurre 1 recarga del listado (no se distingue si hubo cambios) | UI + request interceptada |
| M11-C5 | No se puede cerrar con un guardado en curso | `page.route` retrasa ~2 s la respuesta de `/intento` o `/cerrar`; hacer clic en "✕" y en el backdrop mientras tanto | El modal sigue abierto durante el guardado; al llegar la respuesta se puede cerrar con normalidad | UI + request retrasada |
| M11-C6 | Estado limpio al abrir otro reporte | Reporte A en `postContacto` con la caja de confirmación abierta y texto escrito; cerrar; abrir el reporte B (sin intentos) | B abre en su propio paso inicial (`validacionPrevia`), sin caja de confirmación, textarea vacío y sin intentos de A | UI |
| M11-C7 | La tabla refleja un cierre | Reporte abierto en la tabla ("Asignado" o "Sin asignar"); cerrarlo por "No información mínima" y cerrar el modal | La fila de ese reporte muestra el estado "Cerrado - inválido" sin recargar la página a mano | UI |
| M11-C8 | La tabla refleja un intento | Reporte sin intentos ("Sin intentos" en "Esquema 3x3"); registrar 1 fallido y cerrar el modal | La columna "Esquema 3x3" de esa fila pasa a "1/9 intentos · 0/3 días" | UI |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TU-M11-1` — cerrar por ✕, backdrop y "Cerrar"; `updated` siempre | M11-C1, C2, C3, C4 | `gestion-reporte-modal/cierre-modal.spec.ts` |
| `TU-M11-2` — guardado en curso bloquea el cierre | M11-C5 | `gestion-reporte-modal/cierre-modal.spec.ts` |
| `TU-M11-3` — estado limpio entre reportes | M11-C6 | `gestion-reporte-modal/cierre-modal.spec.ts` |
| `TU-M11-4` — la tabla refleja cierres e intentos al cerrar el modal | M11-C7, C8 | `gestion-reporte-modal/cierre-modal.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| M11-C1 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C2 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C3 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C4 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C5 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C6 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C7 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |
| M11-C8 | ✅ | `gestion-reporte-modal/cierre-modal.spec.ts` | — |

**Resumen:** 8 de 8 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
