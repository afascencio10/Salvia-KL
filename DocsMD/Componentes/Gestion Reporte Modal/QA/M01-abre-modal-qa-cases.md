# QA — M-01 Cuando abre el modal

Casos de prueba derivados de [`flow-M01`](../Flujos/flow-M01-cuando-abre-modal.md): `GET /api/v1/victim-contacts/:icode` (detalle, intentos, contadores, `canManage`) y la elección del paso inicial del modal (`validacionPrevia` / `esquema3x3` / `postContacto` / `cerrado`).

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md).

**Método.** Backend primero (C1–C15, vía `page.request` con la sesión del rol indicado + verificación en BD), después UI (C16–C24) con Playwright. El modal se abre desde "Ver Reporte" en Listado Reportes, buscando el reporte de prueba por su marcador (ver plan, sección B).

---

## Casos de prueba

### Backend — `GET /api/v1/victim-contacts/:icode`

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M01-C1 | Detalle de un reporte tipo Víctima (V2) | Reporte nuevo `victima`, abierto, sin intentos; sesión `OP1` | 200. `icode` = el pedido; `reportType` = "Víctima"; `nombreMostrado` no vacío; `telefono`, `ajusteRazonable` y `descripcionHechos` = lo enviado al crear el reporte; `intentos` = `[]`; `intentosHoy` = 0; `diasCompletos` = 0; `closeReasonCode`/`closeObservation` vacíos; `canManage` = true | API + BD (contra `victim_contact`, `_form2`, `_form2_v2`) |
| M01-C2 | Tercero que reporta a nombre propio | Reporte `tercero-propio` | 200. `reporterNames`, `reporterLastNames`, `reporterEmail` y `parentesco` poblados con lo enviado | API + BD |
| M01-C3 | Tercero anónimo | Reporte `tercero-anonimo` | 200. `reporterNames`/`reporterLastNames` vacíos (el nombre no es obligatorio si es anónimo) | API |
| M01-C4 | Familiar y Hombre (sin nombres, decisión 7A) | Reportes `familiar` y `hombre` | 200 en ambos; `reportType` = "Familiar" / "Hombre"; `nombreMostrado` no vacío; `telefono` no es `"0"`; ningún campo contiene los textos `null` o `undefined` (Q11) | API |
| M01-C5 | Reporte inexistente | UUID aleatorio | 404; el cuerpo no trae datos de ningún reporte | API |
| M01-C6 | Sin sesión | Request sin cookie | 401 | API (contexto limpio) |
| M01-C7 | Rol sin acceso | Sesión `RO`, `PS`, `AD` | Respuesta distinta de 200 (401/403, Q8) y el cuerpo no incluye datos del reportante (datos personales) | API |
| M01-C8 | `canManage` no depende de la pertenencia (GAP-2) | Reporte asignado a `OP2`; luego uno sin asignar; sesión `SV` y sesión `OP1` | `canManage` = true en las cuatro combinaciones | API |
| M01-C9 | `intentosHoy` cuenta solo fallidos de hoy | Sembrar: 2 fallidos hoy, 1 exitoso hoy, 1 fallido hoy con `deleted_at`, 1 fallido de ayer | `intentosHoy` = 2 (excluye exitoso, borrado lógico y otros días) | API + siembra por BD |
| M01-C10 | `diasCompletos` = días con ≥ 3 fallidos | Un reporte por escenario: (a) 3 fallidos en 1 día; (b) 2+2+2 en 3 días; (c) 3+3+3; (d) 5 en un solo día; (e) 3+3+2; (f) 2 fallidos + 1 exitoso en el mismo día; (g) 3 fallidos de los que 1 está borrado lógicamente | (a) 1 · (b) 0 · (c) 3 · (d) 1 · (e) 2 · (f) 0 · (g) 0 | API + siembra por BD (tabla de escenarios) |
| M01-C11 | Frontera de día en Bogotá para `diasCompletos` | 3 fallidos a las 23:50 Bogotá del día D−2 y 3 a las 00:10 Bogotá del día D−1 (mismo día UTC, dos días Bogotá) | `diasCompletos` = 2 (agrupar por día UTC daría 1) | API + siembra por BD |
| M01-C12 | Frontera de "hoy" en Bogotá | 1 fallido a las 23:55 Bogotá de ayer y 1 a las 00:05 Bogotá de hoy | `intentosHoy` = 1, sin importar la hora UTC actual | API + siembra por BD |
| M01-C13 | Forma de la lista de intentos | Reporte con 3 intentos (2 fallidos, 1 exitoso) en fechas distintas | `intentos` ordenados por fecha ascendente; cada ítem trae `date` (ISO), `result` (`fallido`/`exitoso`) y `agentName` (nombre de quien lo registró, no su id) | API |
| M01-C14 | Reporte cerrado, por cada motivo | Reportes cerrados por BD con `ni`, `nc`, `co` (obs. "QA obs"), `dp` (obs. "QA obs") y `ap` | `closeReasonCode` = código; `closeReasonLabel` contiene "Inválido por No información Mínima" / "…Cierre por no contacto" / "…Cierre por no competencia" / "…Cierre por decisión de la persona" / "Válido - Apertura de Caso Salvia"; `closeCategory` = "Inválido" (×4) y "Válido" (`ap`); `closeObservation` = "QA obs" en `co`/`dp`, vacío en el resto | API + BD (tabla de motivos) |
| M01-C15 | Reporte legado invalidado (Q9) | Reporte con `victim_contact_status='i'` y descripción "QA legado", sin `close_reason` | 200. Se trata como cerrado: la UI muestra el paso `cerrado` de solo lectura (sin botones de acción) con "QA legado" como motivo | API + UI (el detalle se verifica en C21) |

### UI — apertura del modal y paso inicial

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| M01-C16 | Apertura: header inmediato y spinner | Reporte `victima`; respuesta del detalle retrasada ~1.5 s con `page.route` | Al hacer clic en "Ver Reporte": diálogo visible (`role="dialog"`, `aria-modal="true"`) con radicado, título, badge de tipo en mayúsculas y "✕" ya visibles; spinner mientras llega el detalle; al llegar desaparece y se ve el cuerpo | UI + request interceptada |
| M01-C17 | Paso `validacionPrevia` | Reporte abierto sin intentos | Se ven el bloque de información, el texto de ayuda y **exactamente 2 botones**: "No información mínima" y "Procede recontacto". No existen "Reporte de prueba" ni "Apertura de caso" | UI (`toHaveCount`) |
| M01-C18 | Bloque de información según los datos | (a) tercero propio; (b) tercero anónimo; (c) reporte sin teléfono | (a) "Quien reporta: {nombre} {apellido} ({parentesco})" y "Correo" visibles; (b) no aparecen ni "Quien reporta" ni "Correo"; (c) "Teléfono de contacto: No suministrado". Siempre: "Ajuste razonable solicitado" y la descripción completa de los hechos | UI |
| M01-C19 | Paso `esquema3x3` si hay ≥ 1 fallido sin exitoso ni cierre | Reporte con 2 fallidos hoy y 1 de ayer | Título "Esquema de Recontacto 3x3"; contador "Hoy: 2/3 · Días completos: 0/3" (los mismos números que devuelve la API); log con los 3 intentos | UI + API |
| M01-C20 | Paso `postContacto` si hay un exitoso sin cierre | Reporte con 3 fallidos y 1 exitoso; y variante con esquema 3x3 completo + 1 exitoso | Ambos casos abren en "Contacto establecido" (aunque `diasCompletos` ≥ 3); no se ve el contador ni el banner "Cerrar por no contacto" | UI + siembra por BD |
| M01-C21 | Paso `cerrado` | Reportes cerrados con `ni`, `co` (obs.), `ap` y el legado de C15 | Título "Reporte {icode} cerrado"; "Categoría: …" y "Motivo: …" coherentes con C14 (o con "QA legado"); botón "Cerrar"; **sin** botones de acción; ícono verde si "Válido", gris si "Inválido" | UI |
| M01-C22 | Precedencia: el cierre gana | Reporte cerrado que además tiene intentos fallidos y exitosos | Abre en `cerrado`, no en `postContacto` ni en `esquema3x3` | UI + siembra por BD |
| M01-C23 | Falla la carga del detalle (Q3) | `page.route` responde 404 y, en otra corrida, 500 | No queda spinner infinito; el modal sigue siendo cerrable; no se muestra ningún paso de gestión con datos vacíos | UI + request simulada |
| M01-C24 | No hay caché entre aperturas | Abrir un reporte, cerrar; sembrar 3 fallidos por BD; abrir el mismo reporte | Los contadores y el log reflejan el estado actual (el cálculo se hace siempre en backend al abrir) | UI + siembra por BD |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-M01-1` — detalle por tipo de reporte (Víctima, Tercero propio/anónimo, Familiar, Hombre) | M01-C1, C2, C3, C4 | `gestion-reportes-api/detalle.spec.ts` |
| `TB-M01-2` — errores y permisos del detalle | M01-C5, C6, C7, C8 | `gestion-reportes-api/detalle.spec.ts` |
| `TB-M01-3` — contadores, escenarios de `diasCompletos` y fronteras de día (siembra por BD) | M01-C9, C10, C11, C12, C13 | `gestion-reportes-api/detalle.spec.ts` |
| `TB-M01-4` — reporte cerrado por cada motivo y legado | M01-C14, C15 (parte API) | `gestion-reportes-api/detalle.spec.ts` |
| `TU-M01-1` — apertura: header, spinner, paso `validacionPrevia` y bloque de información | M01-C16, C17, C18 | `gestion-reporte-modal/apertura.spec.ts` |
| `TU-M01-2` — paso inicial según el estado sembrado (3x3, postContacto, cerrado, precedencia) | M01-C19, C20, C21, C22, C15 (parte UI) | `gestion-reporte-modal/apertura.spec.ts` |
| `TU-M01-3` — fallo de carga y reapertura sin caché | M01-C23, C24 | `gestion-reporte-modal/apertura.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| M01-C1 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01, BUG-02, BUG-03 |
| M01-C2 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01, BUG-02 |
| M01-C3 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C4 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01, BUG-02, BUG-03, BUG-04 |
| M01-C5 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C6 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C7 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C8 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C9 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C10 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01, BUG-05 |
| M01-C11 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01, BUG-05 |
| M01-C12 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C13 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C14 | ✅ | `gestion-reportes-api/detalle.spec.ts` | BUG-01 |
| M01-C15 | ✅ | `gestion-reporte-modal/apertura.spec.ts`, `gestion-reportes-api/detalle.spec.ts` | BUG-01, BUG-07 |
| M01-C16 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | BUG-11 |
| M01-C17 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |
| M01-C18 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |
| M01-C19 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |
| M01-C20 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |
| M01-C21 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | BUG-11 |
| M01-C22 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |
| M01-C23 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |
| M01-C24 | ✅ | `gestion-reporte-modal/apertura.spec.ts` | — |

**Resumen:** 24 de 24 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](../../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
