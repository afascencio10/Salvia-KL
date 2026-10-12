# flow-E06 — Cuando presiona "Gestionar"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Gestionar"
   Tipo: User Interaction
   Funciones: openGestionar() · case-task-modal.openMecanismo() · onTaskCompleted()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  item.id:     case_articulating_mechanism.id
  item.caseId: victim_case_i_code
  item.scope:  municipal | departamental | nacional
}

PASO 1 — Abrir modal
  this.$refs.taskModal.openMecanismo({
    mechanismId: item.id,
    caseId:      item.caseId,
    scope:       item.scope
  })
  → No hace GET de una tarea previa.
  → tarea en memoria: type = 'gestion_mecanismo'
  → Título según scope:
       municipal     → "Gestión del Comité Municipal"
       departamental → "Gestión del Comité Departamental"
       nacional      → "Gestión de la Instancia Nacional"
  → Formulario hardcodeado según scope (ver case-task-modal-interface.md)

PASO 2 — Usuario completa y pulsa "Guardar"
  POST /api/v1/case-tasks/gestion-mecanismo
  Body: { mechanismId, formData }

  Backend:
    SI no autenticado / rol ≠ mc          → 401 / 403
    SI mecanismo no existe                → 404
    SI assigned_user_id ≠ sesión          → 403
    SI status ≠ active                    → 409
    SI formData no coincide con scope     → 400
    SI OK:
      INSERT case_task
        type     = gestion_mecanismo
        category = Mecanismos
        status   = Done
        form_data, result (resumen), completed_at = now
        case_articulating_mechanism_id = mechanismId
      INSERT case_timeline_event
        category = Mecanismos
        type     = "Gestión del mecanismo"
        event_type = GESTION_MECANISMO
        description = result
        task_id + articulating_mechanism_id
        (si falla: WARN; la tarea queda Done)
      El case_articulating_mechanism permanece active.

PASO 3 — SI 200
  Modal emite @completed
  → toast "Gestión guardada."
  → loadCases()  (refresca lastComment / lastCommentAt)

PASO 4 — SI error
  Modal muestra saveError; el listado no cambia.
```
