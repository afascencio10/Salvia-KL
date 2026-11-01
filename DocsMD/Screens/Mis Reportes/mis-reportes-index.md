# Mis Reportes (rol `op`) — Índice

Pantalla **nueva** — "Mis reportes" de la funcionalidad "Gestión reportes" en FigJam: *"solo op, ve/gestiona reportes asignados a él"*. Es `Listado Reportes` pre-filtrada al agente actual, sin reasignación — ver diferencias exactas en [`mis-reportes-interface.md`](mis-reportes-interface.md).

> **Estado:** ✅ Implementado y probado — 2026-10-13. 155 pruebas de Playwright en verde (ver [`../Listado Reportes/QA/gestion-reportes-qa-plan.md`](../Listado%20Reportes/QA/gestion-reportes-qa-plan.md)); resumen del cambio en [`changelogOct2026.md`](changelogOct2026.md). **Re-verificado contra `develop` el 2026-10-31:** el código de las pantallas y del backend de reportes no cambió desde entonces; la única brecha es la columna "Fecha Reporte", que hoy sale en "—" (H11 del plan de QA).

## Archivos de esta pantalla

| Archivo | Descripción |
|---|---|
| [`mis-reportes-interface.md`](mis-reportes-interface.md) | Árbol de interfaz y diferencias vs. `Listado Reportes` |
| [`mis-reportes-events.md`](mis-reportes-events.md) | Inventario de eventos |
| [`related-tables.md`](related-tables.md) | Puntero — mismas tablas que `Listado Reportes` |
| [`QA/mis-reportes-qa-cases.md`](QA/mis-reportes-qa-cases.md) | Casos de prueba, cobertura y resultados de esta pantalla |
| [`changelogOct2026.md`](changelogOct2026.md) | Registro del cambio de octubre 2026 |

## Resumen de eventos

| # | Evento | Flujo |
|---|---|---|
| E-01 | Cuando carga la pantalla | [flow-E01](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E-02 | Cuando presiona "Ver Reporte" en una fila | [flow-E02](./Flujos/flow-E02-cuando-presiona-ver-reporte.md) |

> Filtro de estado, paginación y el modal de gestión se heredan sin cambios de `casos-component`. Búsqueda y ordenamiento usan la misma traducción específica de `victim_contact` ya documentada en `Listado Reportes` → [flow-E07](../Listado%20Reportes/Flujos/flow-E07-cuando-busca-por-texto.md) / [flow-E08](../Listado%20Reportes/Flujos/flow-E08-cuando-ordena-por-fecha.md) — no se duplica aquí.

## Ruta, permiso y menú

| Aspecto | Valor |
|---|---|
| Ruta | `GET /salvia/mis-reportes` (el diseño original decía `/salvia/reportes/mis-reportes`; el usuario confirmó dejar la implementada, 2026-10-13) |
| Handler | `MyReportsGET` en `VictimContactFacade.go` (handler propio, no una rama del de Listado Reportes) |
| Permiso | `get_my_reports`: solo `op`. `sv` y los demás roles reciben **401** |
| Menú lateral | Ítem **Mis Reportes** en `sidebar.html`, solo para `op` (decisión del usuario, 2026-10-13) |

## Checklist de completitud

- [x] ¿El ciclo de vida inicial está cubierto? → E-01
- [x] ¿La acción de ver detalle está cubierta? → E-02
- [x] ¿Cambios de BD? → Ninguno adicional a los de `Listado Reportes` (ver `related-tables.md`)
