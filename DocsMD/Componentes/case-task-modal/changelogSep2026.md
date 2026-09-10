# Changelog — case-task-modal — septiembre 2026

## 2026-09-10 — Tipo `gestion_mecanismo` (Casos Mecanismo)

El rol `mc` registra cada gestión desde el mismo modal: `openMecanismo({ mechanismId, caseId, scope })` abre el formulario municipal, departamental o nacional. El guardado usa `POST /api/v1/case-tasks/gestion-mecanismo` y **siempre crea** una `case_task` Done (`result` = resumen del formulario). No completa una ToDo previa.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | Formulario, `openMecanismo`, POST dedicado |
| `src/salvia/service/case_task_service.go` | `SaveGestionMecanismo` |
| `DocsMD/Componentes/case-task-modal/*` | Interface, usage, flujos E02/E09 |

## 2026-09-02 — Articulación: seguimiento sin default

«¿El compromiso requiere seguimiento?» abre sin Sí ni No seleccionados. La fecha de seguimiento solo aparece si eligen Sí.

## 2026-09-02 — Tipo `articulacion_territorial` (Mis Barreras)

El Enlace Territorial gestiona barreras desde el mismo modal: formulario «Gestionar barrera» (fecha de contacto, compromiso, seguimiento condicional, sugerencias opcionales). `openArticulacion` admite tarea existente o contexto de barrera sin `taskId`. El guardado usa `POST /api/v1/case-tasks/articulacion-territorial`.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | Formulario, `openArticulacion`, POST dedicado |
| `src/salvia/service/case_task_service.go` | `SaveArticulacionTerritorial` |
| `DocsMD/Componentes/case-task-modal/*` | Interface, usage, flujos E02/E09 |
