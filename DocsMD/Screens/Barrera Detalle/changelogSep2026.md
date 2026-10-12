# Changelog — Barrera Detalle — septiembre 2026

## 2026-09-02 — Gestión propia usa el modal de Articulación territorial

El botón "+ Registrar gestión propia" (Enlace Territorial, mismo departamento) deja el modal de tipo/descripción y abre `case-task-modal.openArticulacion` (el de Mis Barreras). El guardado pasa a `POST /api/v1/case-tasks/articulacion-territorial` sin `taskId`. Ese endpoint ahora autoriza crear si la barrera está asignada al usuario **o** si coincide el departamento del Enlace.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-tasks.js` | `abrirModalGestionPropia` → `openArticulacion`; se quita el modal propio |
| `src/frontend/html/salvia/barriers/barrera_detalle.html` | Props `entityBranchId`, `entityName`, `barrierDescription` |
| `src/salvia/service/case_task_service.go` | Auth de create: link **o** departamento |
| `src/salvia/controller/case_task_controller.go` | Pasa `AssignedDepartmentID`; 403 por departamento |
