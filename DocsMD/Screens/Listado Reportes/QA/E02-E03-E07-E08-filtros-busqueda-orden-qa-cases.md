# QA — Listado Reportes · E-02 / E-03 / E-07 / E-08 Filtros, búsqueda y orden

Casos de prueba de los cuatro eventos que se traducen en parámetros de `GET /api/v1/victim-contacts/list`: filtro de estado ([`flow-E02`](../Flujos/flow-E02-cuando-filtra-estado-reporte.md)), filtro de asignación ([`flow-E03`](../Flujos/flow-E03-cuando-filtra-asignacion.md)), búsqueda por texto ([`flow-E07`](../Flujos/flow-E07-cuando-busca-por-texto.md)) y orden por fecha ([`flow-E08`](../Flujos/flow-E08-cuando-ordena-por-fecha.md)).

> **Estado: ejecutado el 2026-10-13** — resultados al final de este archivo. Datos, usuarios, helpers y decisiones y dudas `(Q#)`: [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md). Los alias `R1…R10` y los conjuntos esperados por filtro están en la sección B del plan. **Se asierta pertenencia / no pertenencia de esos alias**, no igualdad exacta (la BD tiene datos reales); todo se acota con el marcador `QAGR-<runId>` en `search`.

**Método.** Backend primero, por parámetro. Nombres de parámetro según `listado-reportes-interface.md` y `VictimContactsListFilters` (Q1): `filter_estado_reporte=<valor>`, `filter_key=asignacion&filter_value=<icode|__unassigned__>`, `search=<texto>`, `sort=registration_date&order=asc|desc`, `page`, `page_size`. El `flow-E02` se reescribió después de la ejecución con los nombres y códigos de la migración (`ni/nc/co/dp/ap`); estos casos los siguen desde el principio. **`TB-LR-E07-1` se ejecuta antes que cualquier flujo `TU-LR-*`/`TU-MR-*`.** Después, UI: los selectores y tiempos se confirman con `playwright-cli` antes de escribir los specs.

---

## Casos de prueba

### Backend — filtro de estado (E-02)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E02-C1 | `asignado` | Pool; `filter_estado_reporte=asignado` | Incluye R2 y R3; excluye R1, R9 (sin agente) y R4, R5, R6, R8, R10 (cerrados, aunque tengan agente) | API |
| LR-E02-C2 | `sin_asignar` | `filter_estado_reporte=sin_asignar` | Incluye R1 y R9; excluye R2, R3 (asignados) y todos los cerrados (R4, R5, R6, R8, R10) | API |
| LR-E02-C3 | `cerrado_invalido` | `filter_estado_reporte=cerrado_invalido` | Incluye R4 (`ni`), R5 (`nc`), R8 (`co`, sin agente) y R10 (`dp`) — los cuatro motivos inválidos; excluye R6 (`ap`) y los abiertos | API |
| LR-E02-C4 | `cerrado_valido` (Q2) | `filter_estado_reporte=cerrado_valido` | Incluye R6 (sembrado por BD con `ap`); excluye todos los demás | API |
| LR-E02-C5 | Valor vacío o ausente | Sin parámetro y con `filter_estado_reporte=` | Sin restricción adicional: los 10 reportes del pool aparecen en ambos casos | API |
| LR-E02-C6 | Combinación con otros filtros | `cerrado_invalido` + `filter_key=asignacion&filter_value=<OP1>`; luego `cerrado_invalido` + `search=<marcador>` | Primero R4, R5, R10 (no R8: sin agente); luego los cuatro cerrados inválidos. Los filtros se combinan con AND, ninguno anula al otro | API |

### Backend — filtro de asignación (E-03)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E03-C1 | Por agente | `filter_key=asignacion&filter_value=<OP1>` | Incluye R2, R4, R5, R6 y R10 (también los cerrados); excluye R3 (`OP2`) y los sin agente | API |
| LR-E03-C2 | "Sin asignar" | `filter_value=__unassigned__` | Incluye R1, R9 y R8 (cerrado pero sin agente: el filtro es por agente, no por estado); excluye R2, R3, R4, R5, R6, R10 | API |
| LR-E03-C3 | Agente inexistente | `filter_value=<UUID aleatorio>` | 200 con `reports: []` y `total: 0` (no 500) | API |

### Backend — búsqueda por texto (E-07)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E07-C1 | Por nombres y apellidos | `search=<marcador>`, luego el marcador en minúsculas y solo un fragmento | Los tres devuelven el pool completo: coincidencia parcial e insensible a mayúsculas, sobre nombres o apellidos | API |
| LR-E07-C2 | Por teléfono, en ambos orígenes | `search=<teléfono de R1>` (Form2) y `search=<teléfono de R9>` (Form1) | Cada uno devuelve su reporte: el teléfono se busca en `victim_contact_form2` **y** en `victim_contact_form1` | API |
| LR-E07-C3 | Por documento (solo Form1) | `search=<documento de R9>`; y un documento cualquiera con R1 (solo Form2) | R9 aparece; R1 no aparece por documento (limitación documentada: Form2 no tiene columna de documento) | API |
| LR-E07-C4 | Los campos se combinan con OR | `search=` un texto que solo coincide con el teléfono de R1 (ningún nombre lo contiene) | Devuelve R1 (a diferencia de `/salvia/primer_contacto`, que combina sus 4 campos con AND) | API |
| LR-E07-C5 | Sin coincidencias | `search=zzz-sin-resultados-<runId>` | 200, `reports: []`, `total: 0` | API |
| LR-E07-C6 | Caracteres especiales | `search` = `O'Brien`, `%`, `_`, `\`, `"`, `<>`, emoji | 200 en todos, sin errores de servidor (no se asierta la semántica de `%`/`_`) | API |
| LR-E07-C7 | Búsqueda + estado + paginación | `search=<marcador>` + `filter_estado_reporte=cerrado_invalido` + `page_size=2` | `total` = 4 (R4, R5, R8, R10); página 1 con 2 y página 2 con 2, sin repetidos | API |

### Backend — orden (E-08)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E08-C1 | Orden por fecha de reporte | Pool con fechas sembradas (R1 = hoy−1 … R10 = hoy−10); `sort=registration_date&order=desc` y luego `asc` | `desc`: R1 antes que R2 … antes que R10; `asc`: el orden inverso | API |
| LR-E08-C2 | Orden no soportado | `sort=next_follow_up` | Sin error de servidor (< 500): 200 con el orden por defecto o 400 (no está definido; la UI nunca lo pide) | API |

### UI — filtro de estado (E-02)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E02-C7 | Opciones del dropdown | `SV` en el listado | "Estado del reporte" ofrece "Todos" y 4 opciones con valores `asignado`, `sin_asignar`, `cerrado_valido`, `cerrado_invalido` (etiquetas "Asignado", "Sin asignar", "Cerrado - válido", "Cerrado - inválido"; se acepta guion corto o largo) | UI |
| LR-E02-C8 | Aplicar cada opción | Pool aislado por el marcador; marcar 2 casillas antes de filtrar | Por cada opción: sale la request con el parámetro del filtro, las filas son las del conjunto esperado, la página vuelve a 1 y las casillas marcadas se limpian (la reasignación no queda armada); "Todos" restaura el listado completo | UI + request interceptada |
| LR-E02-C9 | Estado + búsqueda | "Cerrado - inválido" seleccionado y luego el nombre de R4 en la búsqueda | Solo R4; cambiar el orden o la página no pierde ninguno de los dos filtros | UI |

### UI — filtro de asignación (E-03)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E03-C4 | Opción fija "Sin asignar" | Autocomplete "Asignación" vacío y luego con texto (p. ej. "xyz") | "Sin asignar" es la primera opción del desplegable en ambos casos | UI |
| LR-E03-C5 | Búsqueda de agentes (Q12) | Escribir 3 letras rápido | Una sola request (debounce de 400 ms) a `/api/v1/agents/search` con `q=<texto>` y `role=op`; las sugerencias son solo agentes `op` | UI + request interceptada + BD |
| LR-E03-C6 | Elegir "Sin asignar" | Clic en la opción fija | **No** sale ninguna request a `/api/v1/agents/search`; el listado se pide con `__unassigned__`, en página 1; las filas son las de LR-E03-C2 | UI + request interceptada |
| LR-E03-C7 | Elegir y limpiar un agente | Elegir `OP1` y luego limpiar el autocomplete | Con `OP1`: filas de LR-E03-C1; al limpiar, vuelve el listado completo | UI |

### UI — búsqueda (E-07) y orden (E-08)

| Caso | Descripción | Condición inicial | Resultado esperado | Método de verificación |
|---|---|---|---|---|
| LR-E07-C8 | Búsqueda con debounce | Caja "Buscar por nombres, apellidos, documento o teléfono"; escribir el marcador carácter a carácter | Una sola request tras la ráfaga (debounce de 400 ms), con `search` y `page=1`; salen las filas del pool; al borrar el texto, la request ya no lleva `search` | UI + request interceptada |
| LR-E07-C9 | Sin resultados | Buscar `zzz-sin-resultados-<runId>` | La tabla muestra su estado vacío (el texto existente de `casos-component`) sin banner de error | UI |
| LR-E08-C3 | Opciones de orden | Listado abierto | El `<select>` "Ordenar por" tiene exactamente "Fecha de Reporte DESC" y "Fecha de Reporte ASC"; no existe "Próximo seguimiento" | UI (`toHaveCount`) |
| LR-E08-C4 | Cambiar el orden | Pool aislado, orden DESC; cambiar a ASC y volver | A ASC: request con `order=asc`, página 1 y R10 aparece antes que R1; al volver a DESC, R1 va primero | UI + request interceptada |

---

## Cobertura

| Flujo de prueba | Casos cubiertos | Spec (a crear) |
|---|---|---|
| `TB-LR-E02-1` — filtro de estado: 4 valores, vacío y combinaciones | LR-E02-C1, C2, C3, C4, C5, C6 | `gestion-reportes-api/listado.spec.ts` |
| `TB-LR-E03-1` — filtro de asignación | LR-E03-C1, C2, C3 | `gestion-reportes-api/listado.spec.ts` |
| `TB-LR-E07-1` — búsqueda por texto (nombres, teléfono Form1/Form2, documento, OR, especiales, paginación) | LR-E07-C1, C2, C3, C4, C5, C6, C7 | `gestion-reportes-api/listado.spec.ts` |
| `TB-LR-E08-1` — orden por fecha | LR-E08-C1, C2 | `gestion-reportes-api/listado.spec.ts` |
| `TU-LR-E02-1` — dropdown de estado en la UI | LR-E02-C7, C8, C9 | `listado-reportes/filtros-busqueda-orden.spec.ts` |
| `TU-LR-E03-1` — autocomplete de asignación en la UI | LR-E03-C4, C5, C6, C7 | `listado-reportes/filtros-busqueda-orden.spec.ts` |
| `TU-LR-E07-1` — búsqueda en la UI | LR-E07-C8, C9 | `listado-reportes/filtros-busqueda-orden.spec.ts` |
| `TU-LR-E08-1` — selector de orden en la UI | LR-E08-C3, C4 | `listado-reportes/filtros-busqueda-orden.spec.ts` |

## Resultados de ejecución

**Fecha:** 2026-10-13 · **Build:** `SOG_SALVIA-develop` (árbol de trabajo, con los arreglos de esta ronda) en `http://localhost:9090` · **BD:** Supabase de pruebas · **Herramienta:** Playwright headless (`CI=1`), 1 worker, sin reintentos · **Repo de pruebas:** `qa-salvia/`

| Caso | Resultado | Prueba (spec en `qa-salvia/tests/`) | Hallazgo |
|---|---|---|---|
| LR-E02-C1 | ✅ | `gestion-reportes-api/listado.spec.ts` | BUG-07 |
| LR-E02-C2 | ✅ | `gestion-reportes-api/listado.spec.ts` | BUG-07 |
| LR-E02-C3 | ✅ | `gestion-reportes-api/listado.spec.ts` | BUG-07 |
| LR-E02-C4 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E02-C5 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E02-C6 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E03-C1 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E03-C2 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E03-C3 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C1 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C2 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C3 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C4 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C5 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C6 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E07-C7 | ✅ | `gestion-reportes-api/listado.spec.ts` | BUG-07 |
| LR-E08-C1 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E08-C2 | ✅ | `gestion-reportes-api/listado.spec.ts` | — |
| LR-E02-C7 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E02-C8 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E02-C9 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E03-C4 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E03-C5 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E03-C6 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E03-C7 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E07-C8 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E07-C9 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E08-C3 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |
| LR-E08-C4 | ✅ | `listado-reportes/filtros-busqueda-orden.spec.ts` | — |

**Resumen:** 29 de 29 casos ✅.

La columna *Hallazgo* indica el bug (BUG-xx) que ese caso detectó o que hoy vigila; causa raíz y corrección de cada uno están en [`gestion-reportes-qa-plan.md`](gestion-reportes-qa-plan.md#bugs-encontrados-y-corregidos).
