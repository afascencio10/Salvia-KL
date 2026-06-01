━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  followUpId:   UUID del seguimiento     → plantilla Go ({{ .followUpId }})
  userICode:    i_code del usuario       → plantilla Go ({{ .userICode }})
  formId:       UUID fijo del formulario → hardcoded en data() ("2d0aeb46-...")
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — hacer_seguimiento.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar la app Vue

  Vue.createApp monta sobre #seguimiento-app
  → Se ejecuta mounted() → llama a loadFollowUp() y loadLocations() en paralelo (sin await entre sí)


PASO 2 — Llamar al backend para cargar el seguimiento

  GET /api/v1/follow-ups/{followUpId}/load?agent_id={userICode}&form_id={formId}

  SI respuesta no ok (status != 2xx):
    → this.loadError = data.error || 'Error al cargar el seguimiento'
    → TERMINAR ejecución   // DinamicForm no se monta; se muestra ErrorAlert

  SI respuesta ok:
    → data = JSON parseado { followUp, victimInfo, caseStatus, activeBarriers }
    → CONTINÚA PASO 3


PASO 3 — Poblar estado de la pantalla

  fu             = data.followUp
  submissionId   = fu.form_submission_id
  caseId         = fu.case_id
  caseICode      = fu.case_id
  victimInfo     = data.victimInfo
  caseRiskLevel  = data.victimInfo?.riskLevel || 0   // nivel de riesgo 1-4 del caso
  formState.currentBarriers = data.activeBarriers || []
  // activeBarriers: [{ id, barrierName }] — barreras activas (status != MANAGED) del caso
  // fijadas en la primera carga y reutilizadas en cargas posteriores


PASO 3b — Cargar ubicaciones geográficas en background (paralelo a loadFollowUp)

  Promise.all([
    GET /api/v1/locations/departments  → formState.statesColombia = [{ label, value }]
    GET /api/v1/locations/cities       → allCities = [{ label, value, departmentId }]
  ])
  → Sin loader, sin bloquear UI
  → Si alguna falla: se loggea advertencia, las opciones quedan vacías []
  → allCities es variable local (no en formState); se usa en E-05 para filtrar ciudades por barrera


PASO 4 — Calcular canEdit (prioridad en cascada)

  CONDICIÓN 1: caso cerrado (prioridad máxima)
    SI data.caseStatus === 'cd':
      → canEdit = false
      → TERMINAR cálculo   // no evaluar condiciones siguientes

  CONDICIÓN 2: seguimiento ya realizado
    SI fu.status === 'REALIZADO' && fu.completed_at existe:
      diffDays = (Date.now() - new Date(fu.completed_at)) / (1000 * 60 * 60 * 24)
      SI diffDays <= 5:
        → canEdit = true
      SI diffDays > 5:
        → canEdit = false
      → TERMINAR cálculo

  CONDICIÓN 3: cualquier otro estado
    → canEdit = true


PASO 5 — DinamicForm se monta con los props calculados

  DinamicForm recibe:
    :form-id       = formId        (hardcoded "2d0aeb46-...")
    :submission-id = submissionId  (del followUp)
    :can-edit      = canEdit       (calculado en PASO 4)
    :form-state    = formState     (objeto reactivo):
      psysocialRemisionState  → controla banner de remisión psicosocial
      shouldReassignCase      → controla visibilidad banner de reasignación
      reassingText            → texto del banner de reasignación
      canReassignHigh         → muestra pregunta "¿Confirmar reasignación a alto?"
      canReassignLow          → muestra pregunta "¿Confirmar reasignación a bajo?"
      statesColombia          → opciones de Q12 Departamento (Sección 4 repeater)
      newBarriers             → array por barrera: { cities, towns } para Q13 y Q14
      currentBarriers         → [{ id, barrierName }] para el repeater de Sección 3

  → El componente renderiza el formulario en modo lectura o edición según canEdit

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/follow-ups/:id/load
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


INPUT: {
  followUpId:  UUID del seguimiento   → path param :id
  agentId:     i_code del usuario     → query param agent_id
  formId:      UUID del formulario    → query param form_id
}


PASO 6 — Cargar el follow-up

  DB.follow_ups.FindByID({ followUpId })

  SI no existe:
    → 404 { error: "seguimiento no encontrado" }
    → TERMINAR

  SI fu.AgentID != agentId:
    → 403 { error: "seguimiento no asignado a este agente" }
    → TERMINAR

  SI fu.ScheduledDate > hoy:
    → 400 { error: "seguimiento no disponible aún" }
    → TERMINAR


PASO 7 — Resolver o crear el FormSubmission

  SI fu.form_submission_id existe:
    → submissionId = fu.form_submission_id

  SI fu.form_submission_id es vacío:
    → Crear nueva FormSubmission en DB ({ formId })
    → submissionId = nuevo ID
    → Actualizar fu.form_submission_id = submissionId en DB


PASO 8 — Cargar información de la víctima

  DB.victim_cases.FindByICode({ caseId: fu.case_id })
  → victimInfo = { Names, LastNames, Phone, GenderIdentity,
                   SexualOrientation, ContactPhone, Age, TownName,
                   riskLevel }   // COALESCE(victim_case_form2_risk_level, 0) — entero 1-4
  → Resolver locale: GenderIdentity y SexualOrientation traducidos al español


PASO 9 — Obtener status actual del caso

  DB.victim_cases.FindByICode({ iCode: fu.case_id })
  → caseStatus = vc.victim_case_status  (ej: "", "r", "ra", "c", "cd", "fc")

  SI falla la consulta:
    → caseStatus = ""   // no bloquea la respuesta; se loggea advertencia


PASO 10 — Cargar barreras activas (loadActiveBarriers)

  SI fu.active_barrier_ids existe y no está vacío:
    → ids = split(fu.active_barrier_ids, ",")
    → barriers = DB.barrier_v2.FindByIDs({ ids })
    // Usa los IDs fijados en la primera carga — no re-consulta barreras

  SI fu.active_barrier_ids está vacío (primera carga):
    → barriers = DB.barrier_v2.FindActiveByCaseID({ caseId: fu.case_id, status != 'MANAGED' })
    SI len(barriers) > 0:
      → Guardar IDs en fu.active_barrier_ids (comma-separated) en DB
    // Las barreras quedan fijas desde este momento para este seguimiento

  Por cada barrera:
    → barrierName = "{sector} — {description}"
       (si description vacía → usar specificBarriers como fallback)
    → activeBarriers.append({ id, barrierName })


PASO 11 — Responder al frontend

  200 {
    followUp:       { form_submission_id, case_id, status, completed_at, active_barrier_ids, ... },
    victimInfo:     { Names, LastNames, Phone, ..., riskLevel: 1|2|3|4 },
    caseStatus:     "cd" | "ra" | "r" | ...
    activeBarriers: [{ id, barrierName }, ...]
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                        | Paso afectado |
|--------------------------------------------------------------------------------------------|---------------|
| Si se debe mostrar un mensaje visual diferente cuando canEdit = false por cierre de caso   | PASO 5        |
| Si loadLocations termina después de que answers-updated se emite en carga inicial, las opciones de ubicación aparecen vacías en el primer render | PASO 3b       |
