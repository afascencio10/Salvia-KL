## 2026-11-02 — El timeline puede colgar de un reporte (WP-1 T5, REP-06 parte a)

Se reutilizan la tabla `salvia.case_timeline_event` y el componente `case-timeline` (no se crea otra tabla ni otro componente). Los timelines de casos no cambian.

| Archivo | Cambio |
|---|---|
| `src/internal/models/case_timeline_event.go` | `case_id` nulable (`default:null`), campo `victim_contact_id` (varchar(36), índice) y constantes de la categoría «Reportes» (tipos, íconos, colores, códigos) |
| `src/internal/repository/case_timeline_event_repository.go` | `GetByVictimContactID`; el SELECT se comparte con `GetByCaseID` (`COALESCE(case_id, '')`) |
| `src/salvia/controller/case_detail_controller.go` | `GET /api/v1/reportes/:icode/timeline-events` (op y sv) |
| `src/frontend/js/components/case-timeline.js` | prop `report-id`; `case-id` pasa a opcional; error si no llega ninguno |
| `src/salvia/service/form_service_v4_test.go` | el fake del repositorio suma el stub del método nuevo |

Los hitos que escriben en este timeline (recibido, asignado, intento, contacto logrado, cierre, borrador) son de WP-1 T6 — ver `Componentes/Gestion Reporte Modal/changelogNov2026.md`.
