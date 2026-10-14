## 2026-10-13 — Implementación y QA de Listado Reportes

Pantalla nueva `/salvia/reportes` (op y sv): lista todos los reportes con filtros de estado y asignación, búsqueda y orden por fecha. El supervisor reasigna reportes a un agente integral y «Ver Reporte» abre el modal de gestión. Reutiliza `casos-component` y `reasignar-casos-modal` con props nuevos. La QA con Playwright corrigió 17 bugs, entre ellos el detalle en 500, el 65 % de reportes legados que aparecían abiertos y `/salvia/primer_contacto` en blanco. Se agregó el menú lateral y el backend bloquea `ap`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/reportes/list_reportes.html` | Nuevo — pantalla: casos-component, reasignar-casos-modal y gestion-reporte-modal |
| `src/salvia/facades/VictimContactFacade.go` | Nuevo `ListReportsGET`; `victimCaseRes = "[]"` para que `/salvia/primer_contacto` deje de quedar en blanco |
| `src/salvia/facades/MainRouter.go` | Solo adiciones — ruta `/salvia/reportes` |
| `src/salvia/config/Menu.go` | Solo adiciones — permisos `get_list_reports`, `list_victim_contacts`, `manage_victim_contacts`, `reassign_victim_contacts` |
| `src/salvia/config/Enums.go` | Solo adiciones — templates `list_reportes` y `get_my_reportes` |
| `src/salvia/controller/victim_contacts_list_controller.go · service/…_list_service.go · internal/repository/…_list_repository.go` | Nuevo — `GET /api/v1/victim-contacts/list` |
| `src/salvia/controller/victim_contacts_reassign_controller.go · service/…_reassign_service.go · internal/repository/…_reassign_repository.go` | Nuevo — `POST /api/v1/victim-contacts/reasignar-bulk` (solo sv) |
| `src/internal/models/report_list_item.go` | Nuevo — `ReportListItem` (fila) y `ReportDetail` (detalle del modal) |
| `src/salvia/service/victim_contacts_labels.go` | Nuevo — etiquetas legibles de tipo de reporte, parentesco y ajustes razonables |
| `src/main.go` | Solo adiciones — cableado de repositorios, servicios y controllers; AutoMigrate de `victim_contact` |
| `src/migrations/add_victim_contact_management.sql` | Nuevo — columnas de `victim_contact`, tabla `victim_contact_attempt` y 5 motivos de cierre (no está en git) |
| `src/frontend/html/salvia/victim_contact/get_victim_contacts.html · get_victim_contacts_sv.html` | Enlace "Ver nuevo Listado de Reportes" |
| `src/frontend/html/sidebar/sidebar.html` | Entradas "Listado Reportes" (op, sv) y "Mis Reportes" (op) |
| `src/frontend/js/components/casos-component.js · case_component.html · casos-component.css` | Extensión aditiva — ver changelog de `Casos component` |
| `src/frontend/js/components/reasignar-casos-modal.js · reasignar_casos_modal.html` | Extensión aditiva (`entityKind`) — ver changelog de `Casos component` |
