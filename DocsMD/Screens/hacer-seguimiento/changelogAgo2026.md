## 2026-08-20 — Remisión psicosocial con status abierto y remitente

Cuando `processFollowUpSubmission` (E-04) crea una remisión a Atención Psicosocial al derivar a `atencion_psico`, el registro nace con `status = abierto` (ya no `ACTIVE`) y persiste `submitted_by` (i_code del agente que llenó el formulario) y `submitted_by_team` (`general_user_team` resuelto vía `agentLightRepo`). Si el lookup del equipo falla, la remisión igual se crea y `submitted_by_team` queda vacío (WARN, no aborta). El default de `PsychosocialSupportService.Create` se alineó al mismo status.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/form_service.go` | Create de `psychosocial_support` con `abierto`, `submitted_by` y `submitted_by_team` |
| `src/salvia/service/psychosocial_support_service.go` | Default de Create vacío pasa a `abierto` |
| `src/salvia/service/psychosocial_support_service_test.go` | Assert del default actualizado a `abierto` |
| `DocsMD/Screens/hacer-seguimiento/Flujos/flow-E04-cuando-se-procesa-submission.md` | PASO 4 documenta remitente, equipo y fallback |
| `DocsMD/Screens/hacer-seguimiento/related-tables.md` | Columnas `status`, `submitted_by`, `submitted_by_team` |
| `DocsMD/Screens/hacer-seguimiento/hacer-seguimiento-index.md` | E-04 menciona status y remitente al crear |
| `DocsMD/Screens/hacer-seguimiento/QA/E04-tarea-validacion-psicosocial-qa-cases.md` | C1 espera `abierto` y campos de remitente |
| `DocsMD/Componentes/Remisiones Atencion Psicosocial component/remisiones-psicosocial-component-events.md` | GAP de `submitted_by_team` cerrado para este origen |

## 2026-08-13 — Tarea de validación de remisión psicosocial

Cuando `processFollowUpSubmission` (E-04) crea una remisión a Atención Psicosocial (`salvia.psychosocial_support`) al derivar a `atencion_psico`, ahora también crea una `salvia.case_task` sin asignar (`category: "Psicosocial"`, `type: "validar_remision"`, `description: "Validar que remision a psicosocial es valida"`) con `psychosocial_support_id` apuntando a la remisión recién creada, para que el equipo psicosocial confirme si procede. Probado en backend (3 casos: cumple criterios / no cumple / exclusión por medidas_emergencia) — ver `QA/E04-tarea-validacion-psicosocial-qa-cases.md`.

Pendiente: `case-task-modal.js` y `case-task-history.js` todavía no reconocen el `type` nuevo — la tarea no se puede completar ni ver su detalle desde la UI todavía.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/salvia/service/form_service.go` | `processFollowUpSubmission` crea `case_task` de validación al derivar a `atencion_psico` |
| `DocsMD/Screens/hacer-seguimiento/Flujos/flow-E04-cuando-se-procesa-submission.md` | Documentado el paso nuevo + gap de UI en tabla GAPS |
| `DocsMD/Screens/hacer-seguimiento/hacer-seguimiento-index.md` | Descripción de E-04 actualizada |
| `DocsMD/Screens/hacer-seguimiento/related-tables.md` | Agregada tabla `salvia.case_task` |
| `DocsMD/Screens/hacer-seguimiento/QA/E04-tarea-validacion-psicosocial-qa-cases.md` | Nuevo — casos de prueba backend + cobertura |
