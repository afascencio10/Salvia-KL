━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · cargarDatos() · cargarTodosAgentes() · cargarTareasCaso()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  caseICode:   icode del caso a mostrar    → variable JS inyectada por Go: CASE_ICODE = "{{.caseICode}}"
  userRole:    rol del usuario en sesión   → inyectado por Go: "{{.userRole}}"
  userICode:   icode del usuario en sesión → inyectado por Go: "{{.userICode}}"
  userTeam:    team del usuario en sesión  → inyectado por Go: "{{.userTeam}}"
  currentUser: nombre del usuario          → inyectado por Go: "{{.currentUser}}"
}

PASO 1 — Inicializar estado de carga
  cargando = true, error = null, errorTipo = null

PASO 2 — Disparar dos llamadas en paralelo: cargarDatos() y cargarTodosAgentes()

┌─────────────────────────────────────────────────────────┐
│  SUB-FLUJO A: cargarTodosAgentes()                      │
└─────────────────────────────────────────────────────────┘

  A1. GET /api/v1/equipo-operadores

      Respuesta: lista de agentes activos de todos los equipos
        [{ icode, fullName, team }, ...]

  A2. Guardar resultado en todosAgentes[]
      (se usa para resolver nombres de agentes en la UI: banner, seguimientos, historial)

  → FIN SUB-FLUJO A → CONTINÚA en paralelo con sub-flujo B

┌─────────────────────────────────────────────────────────┐
│  SUB-FLUJO B: cargarDatos()                             │
└─────────────────────────────────────────────────────────┘

  B1. GET /api/v1/casos/:caseICode/detalle

      El backend ejecuta en secuencia:
        1. Query victim_case WHERE victim_case_i_code = caseICode
        2. Resolve agentName: JOIN security.general_user + general_user_profile WHERE icode = victim_case.agent_id
        3. Resolve townName + deptName: town → city → department WHERE town_code = victim_case.town_code
        4. Count denunciasAnteriores: COUNT victim_case WHERE doc_number = victim_case.doc_number (excluye el actual)
        5. Query victim_case_form1 WHERE victim_case = victim_case_id
        6. Query victim_case_form2 WHERE victim_case = victim_case_id
             + enums relacionados (tipoViolencia, subtipoViolencia, ambitoViolencia, planAtencion, ajusteRazonable)
             + resumen (genero, nacionalidad, edadCalculada, tipoAgresorResumen)
             + territorioOcurrencia (city del municipio de los hechos)
        7. Query follow_up_v2 WHERE case_id = caseICode ORDER BY scheduled_date ASC
             + para cada seguimiento: query follow_up_attempts WHERE follow_up_id ORDER BY created_at ASC
        8. Query case_timeline_event WHERE case_id ORDER BY created_at DESC
        9. Query emergency_measure WHERE case_id
        10. Query psychosocial_support WHERE case_id
        11. Query economic_stabilization WHERE case_id
        12. Query barrier_v2 WHERE case_id
        13. Traducir enums via Locale["sp"] (tipoViolencia, subtipo, ámbito, género, nacionalidad, agresor, planAtencion, ajuste)
        14. Si victim_case tiene follow_up_id: query follow_up + follow_up_entries

  SI respuesta 404:
    → error = "Código HTTP 404 — No se encontró el caso con ID: {caseICode}"
    → errorTipo = 404
    → cargando = false
    → Pantalla muestra panel de error con icon 🔍 y opción volver al listado
    → TERMINAR ejecución

  SI otro error HTTP:
    → error = "Código HTTP {status} — Error del servidor..."
    → errorTipo = status
    → cargando = false
    → Pantalla muestra panel de error con icon ⚠️ y botón Reintentar
    → TERMINAR ejecución

  SI respuesta ok (200):
    → Poblar estado Vue:
        caso             = data.Case (+ caso.agentName = data.agentName)
        form1            = data.Form1
        form2            = data.Form2
        followUp         = data.FollowUp
        entries          = data.Entries
        followUpsV2      = data.FollowUpsV2
        townName         = data.townName
        deptName         = data.deptName
        denunciasAnteriores = data.denunciasAnteriores
        emergencyMeasures   = data.emergencyMeasures
        psychosocialSupports = data.psychosocialSupports
        economicStabilizations = data.economicStabilizations
        barriers         = data.barriers
        tipoViolencia    = data.tipoViolencia
        subtipoViolencia = data.subtipoViolencia
        ambitoViolencia  = data.ambitoViolencia
        nacionalidadResumen = data.nacionalidadResumen
        generoResumen    = data.generoResumen
        territorioOcurrencia = data.territorioOcurrencia
        edadCalculada    = data.edadCalculada
        tipoAgresorResumen = data.tipoAgresorResumen
        nombreIdentitario = data.nombreIdentitario || 'No registra'
        diversidadResumen = data.diversidadResumen
        planAtencion     = data.planAtencion
        ajusteRazonable  = data.ajusteRazonable
        hechosTimeline   = data.timelineEvents.filter(e => e.type === 'Hechos del caso')
        cargando = false
    → CONTINÚA FLUJO B

  B2. cargarTareasCaso()
      GET /api/v1/case-tasks?caseId={caseICode}
      → caseTasks[] (usado en banners de tareas pendientes, tab Tareas y por barrera en tab Barreras)

  B3. Ejecutar en paralelo: cargarOperadores() y cargarAgentesRO()
      (ambas usan el team del caso recién cargado — caso.victimCaseTeam)

  ┌─────────────────────────────────────────────────────────┐
  │  SUB-FLUJO B3a: cargarOperadores()                      │
  └─────────────────────────────────────────────────────────┘

    B3a-1. Determinar team del caso: caso.victimCaseTeam (puede ser vacío)

    B3a-2. GET /api/v1/operadores?team={teamCaso}
           (si teamCaso vacío: GET /api/v1/operadores sin filtro)

           Backend: usuarios activos con rol 'ro' filtrados por team
           Respuesta: [{ icode, fullName, team }, ...]

    B3a-3. Ordenar por fullName ASC → operadores[]
           (usado en modal Reasignar caso, dropdown de operadores)

    → FIN SUB-FLUJO B3a

  ┌─────────────────────────────────────────────────────────┐
  │  SUB-FLUJO B3b: cargarAgentesRO()                       │
  └─────────────────────────────────────────────────────────┘

    B3b-1. Determinar team del caso: caso.victimCaseTeam (puede ser vacío)

    B3b-2. GET /api/v1/agentes-ro?team={teamCaso}
           (si teamCaso vacío: GET /api/v1/agentes-ro sin filtro)

           Backend: usuarios activos con rol 'ro' filtrados por team
           Respuesta: [{ icode, fullName, team }, ...]

    B3b-3. Ordenar por fullName ASC → agentesRO[]
           (usado en modal Nuevo seguimiento y modal Reasignar seguimiento para rol sv)

    → FIN SUB-FLUJO B3b

  → FIN SUB-FLUJO B → CONTINÚA FLUJO GENERAL

PASO 3 — Vue renderiza la pantalla con los datos cargados
  - Computed nivelRiesgo: form2.riskLevel > 0 → usa form2; si no → form1.femicideRisk (1→Alto, 0→Bajo)
  - Computed proximoSeguimiento: primer followUpsV2 con status PENDIENTE y !is_completed
  - Computed ultimoAgente: caso.agentName → fallback victimCaseOwnerDescription (parseo texto legacy)
  - Computed puedeReasignar: userRole === 'sv'
  - Computed agentesParaReasignarSeg: sv → agentesRO; op → solo él mismo (un único elemento); ro → todosAgentes filtrado por team del caso
  - Tab inicial: tabActiva = 'info' (Sección Resumen abierta, Hechos y Completa cerradas)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| Diferencia entre operadores[] y agentesRO[]: ambos llaman al mismo endpoint con rol 'ro'. Parece redundante — podría unificarse. | B2a, B2b |
| hechosTimeline filtra type === 'Hechos del caso' — confirmar que ese es el valor exacto del campo `type` en case_timeline_event | B1 paso 8 |
