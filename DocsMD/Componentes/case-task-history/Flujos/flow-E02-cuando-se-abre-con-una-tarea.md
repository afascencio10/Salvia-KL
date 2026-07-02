━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se abre con una tarea
   Tipo: User Interaction
   Función: open(taskId)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: el componente padre (case-tasks.js, listado de tareas del
caso) vía this.$refs.taskHistory.open(taskId)

INPUT: {
  taskId:  UUID de la case_task ya completada  → pasado por el padre
}


PASO 1 — Mostrar el modal y resetear estado previo

  visible  = true
  tarea    = null
  error    = null
  cargando = true


PASO 2 — Cargar la tarea del API

  GET /api/v1/case-tasks/{taskId}
  → tarea: { id, type, description, status, assignedUserId, assignedUserName,
             completedAt, formData, ... }
  // assignedUserName ya viene resuelto por el backend (GetByID enriquecido)

  SI status === 401:
    → Redirigir a /static/landing.html
    → TERMINAR ejecución

  SI error o status !== 200:
    → error    = "No se pudo cargar la tarea."
    → cargando = false
    → TERMINAR ejecución  // modal permanece visible mostrando el error

  SI ok:
    → tarea    = resultado
    → cargando = false
    → CONTINÚA FLUJO GENERAL


PASO 3 — Vue re-evalúa reactivamente

  → El título del modal muestra el label correspondiente a tarea.type
  → El Body renderiza "Completado por", la fecha, la descripción y el
    detalle de formData según tarea.type (ver rama SEGÚN tarea.type en
    case-task-history-interface.md)

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ GAPS resueltos
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Lo que se creía un GAP                                          | Resolución |
|--------------------------------------------------------------------|-----------|
| tarea.assignedUserName no venía en GetByID                        | Se enriqueció `GetByID` en el backend (`taskToJSON` + `lookupUserName`, compartido con el listado). Confirmado en test E2E: el nombre llega resuelto cuando `assignedUserId` no está vacío. |

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| ¿Se valida tarea.status === 'Done' y se muestra algo distinto si no lo   | PASO 2        |
| está? El listado externo (case-tasks.js) debería filtrar solo Done,     |               |
| pero eso no está garantizado desde este componente.                    |               |
