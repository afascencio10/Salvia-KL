# Changelog — Mis Barreras — septiembre 2026

## 2026-09-02 — Columna Caso (Doc/Tel/ciudad) y filtro Ciudad

La columna Caso deja de mostrar el i_code: nombre, documento, teléfono y ciudad (como Barreras Departamento). Se agrega select de ciudades del departamento asignado. El Sí/No de seguimiento en el modal arranca sin selección.

**Archivos**

| Archivo | Cambio |
|---|---|
| `barrier_v2_mine_view.go` / `barrier_v2_repository.go` | `victimPhone`, `cityId`, `cityName` |
| `MisBarrerasFacade.go` | `assignedDepartmentId` |
| `barrieris.html` / `barriers.js` / `barriers.css` | Columna Caso + select ciudad |
| `case-task-modal.js` | `requiereSeguimiento` inicial `null` |

## 2026-09-02 — Rediseño: listado por barrera del enlace + modal Articulación territorial

La pantalla deja de listar `case_task` del usuario y pasa a listar `barrier_v2` activas con `user_territory_link_id` = sesión. Se quitan tabs y la columna Gestión Realizada. Gestionar abre `case-task-modal` (con o sin tarea previa). Guardar completa o crea+completa `articulacion_territorial` y, si hay seguimiento, crea ToDo `compromiso_con_enlace`.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/internal/models/barrier_v2.go` | Columna `user_territory_link_id` |
| `src/internal/models/barrier_v2_mine_view.go` | DTO del listado |
| `src/internal/repository/barrier_v2_repository.go` | SQL `FindMineWithRelations` |
| `src/salvia/controller/barrier_v2_gin_controller.go` | `GET /barriers-v2/mine` |
| `src/salvia/service/case_task_service.go` | `SaveArticulacionTerritorial` |
| `src/salvia/controller/case_task_controller.go` | `POST /articulacion-territorial` |
| `src/frontend/html/salvia/barriers/barrieris.html` | Tabla 4 columnas + case-task-modal |
| `src/frontend/js/components/barriers.js` | Carga mine; abre articulación |
| `src/frontend/js/components/case-task-modal.js` | Formulario Gestionar barrera |
| `DocsMD/Screens/Mis Barreras/*` | Documentación de pantalla |
