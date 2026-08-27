## 2026-08-14 — Nuevo tipo de tarea `justificar_remision`

Cuando el psicólogo/a devuelve una remisión psicosocial (`ValidateRemision`, `isValid=false`), la `case_task` `justificar_remision` que se crea (ver changelog de `hacer-seguimiento` del 2026-08-13 y de `Remision Psicosocial Detalle`) ahora se puede completar desde este modal: el agente ve el motivo de la devolución (`tarea.description`, de solo lectura) y escribe su respuesta en un campo requerido (máx. 500 caracteres).

Al completarla (`CompleteWithFormData` → `sideEffectsJustificarRemision`, `case_task_service.go`):
- La `psychosocial_support` vuelve a `status = 'abierto'`.
- Se crea una nueva `case_task` `validar_remision` (`ToDo`) — asignada al `psychologist_id` de la `dupla` de la remisión, o al `professional_id` individual si no hay dupla, o sin asignar si la remisión no tiene ninguno de los dos.
- Se registra `case_timeline_event` "Remisión Justificada".

Probado con Playwright — ver `DocsMD/Screens/Remision Psicosocial Detalle/QA/validar-remision-justificar-qa-cases.md`.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | Nueva rama `justificar_remision`: título, ícono, validación, `_initForm`, `_buildFormData`, template (motivo + textarea respuesta) |
| `src/salvia/service/case_task_service.go` | Nuevo case en `CompleteWithFormData` + `sideEffectsJustificarRemision` (reabre remisión, crea `validar_remision`, timeline event) |
| `DocsMD/Componentes/case-task-modal/case-task-modal-interface.md` | Árbol, tabla de efectos de lado, validaciones y límite de caracteres actualizados |
| `DocsMD/Componentes/case-task-modal/case-task-modal-usage.md` | Tipos soportados + shape de `formData` para `justificar_remision` |

---

## 2026-08-14 — Límite de caracteres en campos de texto

Se agregó `maxlength` a los campos que no lo tenían: `entidadNombre` (150), `descripcion` (500), `rutaKofax` (255) y las 3 observaciones de `comite_caso` (500 c/u). `funcionario` (50) y `asunto` (80) ya lo tenían. Probado con Playwright — ver `QA/limite-caracteres-qa-cases.md`.

De paso: `detalle-caso.page.ts` (qa-salvia) tenía `abrirTaskModal()` roto — dependía de un panel dev que quedó comentado en `get_case_detail_sv.html` al reemplazarse por el tab real "Tareas" (`case-tasks.js`). Se actualizó para usar el flujo real.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-task-modal.js` | `maxlength` agregado a 6 campos (2 duplicados en `gestion_llamada`/`proyectar_oficio`) |
| `DocsMD/Componentes/case-task-modal/case-task-modal-interface.md` | Nueva tabla "Límite de caracteres por campo" |
| `DocsMD/Componentes/case-task-modal/QA/limite-caracteres-qa-cases.md` | Nuevo — casos de prueba + cobertura |
| `qa-salvia/pages/detalle-caso.page.ts` | `abrirTaskModal()` reescrito para usar el tab "Tareas" real |
| `qa-salvia/helpers/db.ts` | Nuevos helpers `assignTask()` y `setTaskDescription()` |
| `qa-salvia/tests/case-task-modal/limite-caracteres.spec.ts` | Nuevo — 2 tests |
