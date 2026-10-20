# Listado Reportes (roles `op` y `sv`) — Índice

Pantalla **nueva** — implementa la funcionalidad "Listado reportes" descrita en FigJam (canal `l25jrzyf`, sección "Gestion reportes"): ve todos los reportes, con filtros y reasignación por el supervisor. Reemplaza en este documento la planeación anterior (que proponía adaptar `get_victim_cases_sv.html`) — se descartó a favor de reutilizar `casos-component`, ya probado en producción. Ver veredicto completo en el chat de la sesión 2026-10-10.

> **Estado:** ✅ Implementado y probado — 2026-10-13. 155 pruebas de Playwright en verde más un E2E repetible de un reporte nuevo (ver [`QA/gestion-reportes-qa-plan.md`](QA/gestion-reportes-qa-plan.md)); resumen del cambio en [`changelogOct2026.md`](changelogOct2026.md).

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`listado-reportes-interface.md`](listado-reportes-interface.md) | Árbol de interfaz y extensión de `casos-component`/`reasignar-casos-modal` |
| [`listado-reportes-events.md`](listado-reportes-events.md) | Inventario de eventos |
| [`related-tables.md`](related-tables.md) | Columnas nuevas en `victim_contact`, tabla nueva `victim_contact_attempt`, GAPS |
| [`QA/gestion-reportes-qa-plan.md`](QA/gestion-reportes-qa-plan.md) | Plan de QA compartido por las tres piezas + resultados, bugs encontrados y decisiones |
| [`QA/E01-carga-pantalla-qa-cases.md`](QA/E01-carga-pantalla-qa-cases.md) | Casos de E-01 (acceso, columnas, filas, enlace legado, menú) |
| [`QA/E02-E03-E07-E08-filtros-busqueda-orden-qa-cases.md`](QA/E02-E03-E07-E08-filtros-busqueda-orden-qa-cases.md) | Casos de E-02, E-03, E-07 y E-08 |
| [`QA/E04-E05-E06-ver-y-reasignar-qa-cases.md`](QA/E04-E05-E06-ver-y-reasignar-qa-cases.md) | Casos de E-04, E-05 y E-06 |
| [`QA/gestion-reportes-e2e-qa-cases.md`](QA/gestion-reportes-e2e-qa-cases.md) | E2E de un reporte nuevo: formulario público → Listado → Mis Reportes → modal, repetible |
| [`changelogOct2026.md`](changelogOct2026.md) | Registro del cambio de octubre 2026 |

## Resumen de eventos

| # | Evento | Flujo |
|---|---|---|
| E-01 | Cuando carga la pantalla | [flow-E01](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E-02 | Cuando filtra por estado del reporte | [flow-E02](./Flujos/flow-E02-cuando-filtra-estado-reporte.md) |
| E-03 | Cuando filtra por asignación | [flow-E03](./Flujos/flow-E03-cuando-filtra-asignacion.md) |
| E-04 | Cuando presiona "Ver Reporte" en una fila | [flow-E04](./Flujos/flow-E04-cuando-presiona-ver-reporte.md) |
| E-05 | Cuando presiona "Reasignar Reportes" | [flow-E05](./Flujos/flow-E05-cuando-presiona-reasignar-reportes.md) |
| E-06 | Cuando confirma la reasignación | [flow-E06](./Flujos/flow-E06-cuando-confirma-reasignacion.md) |
| E-07 | Cuando busca por texto | [flow-E07](./Flujos/flow-E07-cuando-busca-por-texto.md) |
| E-08 | Cuando ordena por fecha de reporte | [flow-E08](./Flujos/flow-E08-cuando-ordena-por-fecha.md) |

> Paginación **no se re-documenta aquí**: es comportamiento genérico ya cubierto por `casos-component` (E-06 del componente) y se hereda sin cambios. Búsqueda (E-07) y ordenamiento (E-08) sí se documentan aquí — aunque el *mecanismo* de UI es genérico, su traducción a SQL/opciones es específica de `victim_contact` (corrección de planeación, 2026-10-11 — ver auditoría contra la imagen de funcionalidades de FigJam).

## Relación con otras piezas

| Pieza | Relación |
|---|---|
| [`Mis Reportes`](../Mis%20Reportes/mis-reportes-index.md) | Misma composición, pre-filtrada al agente actual, sin reasignación |
| [`Gestión Reporte Modal`](../../Componentes/Gestion%20Reporte%20Modal/gestion-reporte-modal-index.md) | Popup que abre el botón "Ver Reporte" — toda la lógica de validación previa, esquema 3x3 y cierre vive ahí, no en esta pantalla |
| `/salvia/primer_contacto` | Pantalla legada equivalente — coexiste (GAP-4, resuelto) |

## Enlace desde `/salvia/primer_contacto` (GAP-4, confirmado 2026-10-11 — implementado)

El usuario confirmó coexistencia **con** un enlace de salida: se agrega un link visible en la pantalla legada `/salvia/primer_contacto` (templates `get_victim_contacts.html` / `get_victim_contacts_sv.html`, ver `VictimContactFacade.VictimContactGET`) que navega a `/salvia/reportes`. Es un cambio mínimo, categoría B (pantalla existente):

- Un `<a href="/salvia/reportes">` o botón equivalente en la cabecera de esa pantalla — mismo patrón visual que sus botones de header actuales.
- Sin lógica nueva, sin llamada a backend — navegación simple, no requiere flujo propio (`dev-docs-system.md`: "eventos con lógica trivial... no lo necesitan").
- No se crea una carpeta `DocsMD/Screens/` nueva para `/salvia/primer_contacto` solo por esto — esa pantalla no tiene documentación MD propia hoy y agregarla completa está fuera del alcance de este requerimiento.

**Implementado:** `<a href="/salvia/reportes" class="btn btn-secondary">Ver nuevo Listado de Reportes</a>` en `get_victim_contacts.html` (op) y `get_victim_contacts_sv.html` (sv).

**Bug previo corregido (2026-10-13):** al probar el enlace se vio que `/salvia/primer_contacto` quedaba **en blanco para todos los roles**: `VictimContactGET` declaraba `victimCaseRes` sin asignarlo y los templates hacen `JSON.parse({{.victimCases}})` — `JSON.parse("")` lanza y Vue no monta. Se arregló con `victimCaseRes = "[]"` (decisión del usuario; esa pantalla ya no lista casos).

## Entrada en el menú lateral

Se agregó el ítem **Listado Reportes** (`/salvia/reportes`) en `sidebar.html` para `op` y `sv` (decisión del usuario, 2026-10-13). El ítem "Reportes" que ya existía apunta a `/salvia/casos`, no a esta pantalla.

## Checklist de completitud

- [x] ¿El ciclo de vida inicial está cubierto? → E-01
- [x] ¿Todo filtro nuevo específico de reportes está cubierto? → E-02, E-03
- [x] ¿La acción de ver detalle está cubierta? → E-04
- [x] ¿La reasignación masiva está cubierta? → E-05, E-06
- [x] ¿Los cambios de BD están identificados? → `related-tables.md`
- [x] ¿GAPs resueltos? → GAP-2, GAP-3 y GAP-4 confirmados por el usuario (2026-10-11). GAP-1 (orígenes del reporte) cubierto por la implementación: el detalle lee `form1`, `form2` y `form2_v2` según existan (2 951 de 3 503 reportes son solo-`form1` y 552 tienen `form2`); confirmar con negocio si existe un 4.º origen. GAP-6/GAP-7 siguen diferidos junto con M-10 (`related-tables.md`).
