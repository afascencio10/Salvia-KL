━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se abre con una tarea
   Tipo: User Interaction
   Función: open(taskId)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: el componente padre vía this.$refs.taskModal.open(taskId)

INPUT: {
  taskId:  UUID de la case_task a gestionar  → pasado por el padre
}


PASO 1 — Mostrar el modal y resetear estado previo

  visible       = true
  tarea         = null
  saveError     = null
  errorTarea    = null
  cargandoTarea = true
  guardando     = false


PASO 2 — Cargar la tarea del API

  GET /api/v1/tasks/{taskId}
  → tarea: { id, type, status, case_id, barrier_id, entity_letter_id, ... }

  SI error o status !== 200:
    → errorTarea    = mensaje de error
    → cargandoTarea = false
    → TERMINAR ejecución  // modal permanece visible mostrando el error

  SI ok:
    → tarea         = resultado
    → cargandoTarea = false
    → CONTINÚA FLUJO GENERAL


┌─────────────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Inicializar formulario según tarea.type             │
└─────────────────────────────────────────────────────────────────┘

  SEGÚN tarea.type:

    CASO 'gestion_llamada':
      form = {
        departamentoId: null,
        ciudadId:       null,
        municipioId:    null,
        entidadId:      null,
        entidadNombre:  "",
        funcionario:    "",
        descripcion:    "",
        generaOficio:   false,
        asunto:         "",
        rutaKofax:      "",
      }
      municipios = []
      entidades  = []

    CASO 'proyectar_oficio':
      form = {
        departamentoId: null,
        ciudadId:       null,
        municipioId:    null,
        entidadId:      null,
        entidadNombre:  "",
        funcionario:    "",
        asunto:         "",
        rutaKofax:      "",
      }
      municipios = []
      entidades  = []

    CASO 'comite_caso':
      form = {
        decisiones:                  [],
        observacionesOficio:         "",
        observacionesRecomendaciones: "",
        nivelMecanismo:              null,
        observacionesMecanismo:      "",
      }

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


PASO 3 — Vue re-evalúa reactivamente

  → El título del modal muestra el label correspondiente a tarea.type
  → El Body renderiza el formulario del tipo correcto
  → BtnConfirmar queda deshabilitado (formulario vacío → formularioValido = false)

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                             | Paso afectado |
|-----------------------------------------------------------------|---------------|
| Shape exacto del objeto tarea devuelto por GET /api/v1/tasks/:id| PASO 2        |
