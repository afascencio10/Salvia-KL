## 2026-11-02 — REP-12: la fecha del reporte muestra el día de Bogotá

La columna «Fecha Reporte» del Listado de Reportes salía en «—» porque el listado no traía `creationDateCO`. Ahora `GET /api/v1/victim-contacts/list` devuelve el día de creación en hora de Bogotá (`AAAA-MM-DD`, sobre el `timestamp` ya guardado, sin `AT TIME ZONE`) y la celda lo muestra sin `new Date()`: el día no cambia por la zona horaria del navegador. Probado en 5 zonas, en los bordes de día y de año y con un reporte real del asistente público. Cierra H11 del plan de QA.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/report_list_item.go` | `ReportListItem` agrega `CreationDateCO` (`creationDateCO`); `creationDate` se conserva |
| `src/internal/repository/victim_contacts_list_repository.go` | Selecciona `TO_CHAR(victim_contact_creation_date, 'YYYY-MM-DD') AS creation_date_co` y lo mapea a la fila |
| `scripts/e2e/wp2-lib.js · wp2-rep12-fecha-reporte.js` | Nuevos — prueba backend + BD + pantalla en 5 zonas horarias (19 comprobaciones) |

## 2026-11-02 — Estado «Borrador iniciado» y hitos de asignación (WP-1 T7 · T6)

> Zona compartida con WP-2 (Listado Reportes): los cambios son aditivos y se anotan para integración.

- **Estado derivado «Borrador iniciado»** (T7): un reporte abierto con un caso de Registro V3 vinculado en `bo`/`in` se lee «Borrador iniciado»; con el caso en `te`/`cd` se lee «Cerrado — válido». Sin columna nueva; regla única en `related-tables.md`. Aplica a Listado Reportes y Mis Reportes (columna Estado y filtro `estado_reporte`).
- **Reasignación masiva** (T6): la asignación y su hito «Reporte asignado»/«Reporte reasignado» van en una transacción.

| Archivo | Cambio |
|---|---|
| `src/internal/repository/victim_contacts_list_repository.go` | `reportStatusFrom` con `linkedStatus`, columna `linked_status` y condiciones del filtro `estado_reporte` (`borrador_iniciado`; `cerrado_valido` incluye el caso completo) |
| `src/frontend/js/components/casos-component.js` · `src/frontend/css/casos-component.css` | etiqueta y badge «Borrador iniciado» |
| `src/frontend/html/salvia/reportes/list_reportes.html` · `get_my_reportes.html` | opción de filtro «Borrador iniciado» y `case-timeline.js` |
| `src/internal/repository/victim_contacts_reassign_repository.go` · `src/salvia/service/victim_contacts_reassign_service.go` | `AssignWithTimeline` y descripción del hito |
