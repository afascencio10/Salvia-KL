## 2026-07-23 — Registrar gestión propia (Enlace Territorial) + sincronización de la doc con la implementación real

El MD de esta pantalla describía un estado mock (datos hardcoded, sin llamadas a API) que ya no correspondía al código — la pantalla real ya carga desde `GET /api/v1/barriers-v2/:id/detail` y tiene 3 tabs (Información General, Tareas, Timeline). Se actualizó `barrera-detalle-interface.md` y `barrera-detalle-index.md` para reflejar esa realidad, y se agregó la funcionalidad nueva: el rol Enlace Territorial (`en`) puede registrar, desde el tab "Tareas", una gestión propia como tarea ya completada — sin pasar por una tarea pendiente previa — restringido a barreras de su mismo departamento asignado. Ver plan completo en `DocsMD/Otros/temp/req-registrar-gestion-propia-enlace.md`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/case_timeline_event.go` | Constantes `TimelineTypeGestionPropia`, `TimelineEventGestionPropia`, `TimelineIconGestionPropia` |
| `src/salvia/service/case_task_service.go` | `CreateGestionPropia()` — crea `case_task` ya `Done` + evento de timeline + transición `OPEN → En Gestion` |
| `src/salvia/controller/case_task_controller.go` | Endpoint `POST /api/v1/case-tasks/gestion-propia` con validación de rol (`en`) y departamento |
| `src/security/dao/GeneralUserDAO.go` | Campo `GeneralUserAssignedDepartment` (DTO, PgDB, FieldDefinitions, `GetGeneralUser`, create/update) |
| `src/common/utils/CommonSession.go` | Campo `AssignedDepartmentID` en la sesión |
| `src/security/facades/GeneralUserFacade.go` | Propaga `AssignedDepartmentID` al login |
| `src/salvia/facades/BarreraDetalleFacade.go` | Inyecta `currentUserAssignedDepartment` al template |
| `src/frontend/html/salvia/barriers/barrera_detalle.html` | Computed `puedeRegistrarGestionPropia`, prop nueva a `<case-tasks>` |
| `src/frontend/js/components/case-tasks.js` | Botón + modal "Registrar gestión propia" |
| `src/frontend/js/components/case-task-history.js` | Rama de detalle para `tarea.type === 'gestion_propia'` |
| `src/frontend/html/security/general_user/set_general_user.html` | Select "Departamento asignado" (solo rol Enlace Territorial) |
| Migración SQL | `security.general_user.general_user_assigned_department` (varchar 20) |
