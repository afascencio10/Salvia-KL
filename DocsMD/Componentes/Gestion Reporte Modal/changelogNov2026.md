## 2026-11-02 — WP-1: historial del reporte, hitos de timeline y «Registro de caso» (T6 · T7, REP-06 · REP-03)

El modal de gestión muestra el historial del reporte, un encabezado con nombre, teléfono y radicado, la sección «Hechos del reporte» en todos los pasos y, tras un contacto exitoso de un reporte Víctima o Tercero, el botón «Registro de caso» que crea un borrador de Registro V3 (el reporte queda «Borrador iniciado»). Probado por API, SQL y Playwright; ver `QA/M13-…` y `QA/M10-M12-…`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/gestion-reporte-modal.js` · `…/reportes/gestion_reporte_modal.html` · `src/frontend/css/gestion-reporte-modal.css` | encabezado, «Hechos del reporte», historial (`historyVersion`), «Registro de caso», aviso Familiar/Hombre, paso `borradorIniciado` y `irAlBorrador()` |
| `src/frontend/html/salvia/reportes/list_reportes.html` · `get_my_reportes.html` | +1 `<script>` de `case-timeline.js` y la opción de filtro «Borrador iniciado» |
| `src/internal/repository/victim_contacts_management_repository.go` | `InsertAttempt` y `Close` en una transacción con su hito (con `FOR UPDATE`; `Close` recibe quién cierra); `GetDetail` extendido (radicado, tipo, borrador, cierre derivado); 409 si hay borrador |
| `src/salvia/service/victim_contacts_management_service.go` · `src/salvia/controller/victim_contacts_management_controller.go` | nueva firma de `Close`, `CanConvert` según el rol, mapeo del 409 por borrador |
| `src/internal/repository/victim_contacts_reassign_repository.go` · `src/salvia/service/victim_contacts_reassign_service.go` | `AssignWithTimeline` (asignación + hito «asignado»/«reasignado» en una transacción) |
| `src/internal/models/report_detail_ext.go` | Nuevo — `ReportDetailExt` (embebe `ReportDetail`; no se toca `report_list_item.go`) |
| `src/internal/repository/report_timeline.go` | Nuevo — `NewReportMilestoneEvent` y los hitos |
| `src/salvia/controllers/VictimContactControllerV2.go` · `report_milestones.go` | hito «Reporte recibido» (después del commit, best-effort) |
| `src/salvia/controller/victim_case_v3_from_report_controller.go` · `service/…_service.go` · `internal/repository/…_repository.go` | Nuevos — conversión (ver `Screens/Registro V3/changelogNov2026.md`) |
| `src/internal/repository/victim_contacts_list_repository.go` · `src/frontend/js/components/casos-component.js` · `src/frontend/css/casos-component.css` | estado derivado «Borrador iniciado» en el listado (ver `Screens/Listado Reportes/changelogNov2026.md`) |

M-10 pasa de diferido a implementado con otro diseño (no cierra el reporte; `ap` sigue bloqueado en `/cerrar`) y nace M-12.
