━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  defaultFilter:   filtro inicial   → prop :defaultFilter del padre
                                      shape: { key: string, value?: string }
                                      keys válidos: 'casos_nuevos' | 'riesgo' | 'equipo' | 'persona_asignada'
  columns:         columnas         → prop :columns (Array<{ key, label }>) del padre
  hiddenColumns:   columnas ocultas → prop :hiddenColumns (Array<string>) del padre (opcional)
  buttons:         botones de fila  → prop :buttons (Array<{ id, label }>) del padre
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar estado interno del componente

  activeFilter   = defaultFilter              // { key, value? }
  searchText     = ""
  sortBy         = "registration_date"        // criterio de orden por defecto
  sortOrder      = "desc"
  cases          = []                         // lista vacía hasta recibir respuesta
  loading        = true
  loadError      = null

  visibleColumns = columns.filter(col => !hiddenColumns.includes(col.key))
  // SI hiddenColumns vacío → visibleColumns = columns completo


PASO 2 — Consultar backend con el filtro inicial

  GET /api/v1/cases/list?filter_key={activeFilter.key}&filter_value={activeFilter.value}&sort={sortBy}&order={sortOrder}

  // Cada objeto de caso devuelto por el backend incluye campos de dos fuentes:
  //   · salvia.victim_case    → id, i_code, names, lastNames, creationDate, status, docNumber
  //   · salvia.follow_up_v2  → agentNames, agentLastNames, team, riskStatus, nextFollowUpDate

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al cargar los casos'
    → Mostrar ErrorState en lugar de la tabla
    → TERMINAR ejecución

  SI respuesta ok:
    → cases   = data.cases   // Array de objetos de caso enriquecidos
    → loading = false
    → CONTINÚA PASO 3


PASO 3 — Calcular filteredCases inicial (propiedad computada)

  filteredCases = cases   // sin búsqueda activa, todos los casos son visibles
  // La propiedad computada filteredCases se recalcula automáticamente
  // cuando cambie searchText (→ ver E-03) o sortBy / sortOrder (→ ver E-04)


PASO 4 — Renderizar la tabla

  PARA CADA caso EN filteredCases:
    Renderizar fila con visibleColumns:

      'victim_info'       → case.names + " " + case.lastNames  /  case.i_code
      'operator'          → case.agentNames + " " + case.agentLastNames  (o "—" si vacío)
      'registration_date' → case.creationDate  formateada DD/MM/YYYY
      'risk_level'        → RiskBadge  usando case.riskStatus
      'team'              → case.team  (o "—" si vacío)
      'assigned_person'   → case.agentNames + " " + case.agentLastNames  (o "—" si vacío)
      'next_follow_up_date' →
          SI case.nextFollowUpDate existe: mostrar fecha formateada DD/MM/YYYY
          SI no existe: mostrar "—"

    Renderizar columna de acciones:
      PARA CADA btn EN buttons:
        → Renderizar botón btn.label
        → Al presionar → emitActionClicked(btn.id, case)  (→ Ver E-05)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filter_key:    key del filtro activo   → query param
  filter_value:  valor del filtro        → query param (opcional según filter_key)
  sort:          campo de ordenamiento   → query param: 'registration_date' | 'next_follow_up'
  order:         dirección               → query param: 'asc' | 'desc'
}


PASO 5 — Construir query base sobre salvia.victim_case

  Consultar salvia.victim_case  →  VictimCaseLight
  JOIN salvia.follow_up_v2 ON follow_up_v2.case_id = victim_case.victim_case_i_code
    (LEFT JOIN para no excluir casos sin seguimientos)


PASO 6 — Aplicar filtro según filter_key

  SEGÚN filter_key:

    CASO 'casos_nuevos':
      → WHERE NOT EXISTS (
            SELECT 1 FROM salvia.follow_up_v2
            WHERE case_id = victim_case_i_code
            AND   status  = 'REALIZADO'
        )
      → CONTINÚA PASO 7

    CASO 'riesgo':
      → WHERE follow_up_v2.risk_status = filter_value
      → CONTINÚA PASO 7

    CASO 'equipo':
      → WHERE follow_up_v2.team = filter_value
      → CONTINÚA PASO 7

    CASO 'persona_asignada':
      → WHERE follow_up_v2.agent_id = filter_value
      → CONTINÚA PASO 7

    DEFAULT (sin filtro o defaultFilter del padre):
      → Sin cláusula WHERE adicional
      → CONTINÚA PASO 7


PASO 7 — Resolver nextFollowUpDate por caso

  Para cada caso resultante, buscar el próximo seguimiento pendiente:

  SELECT scheduled_date
  FROM   salvia.follow_up_v2
  WHERE  case_id = victim_case_i_code
  AND    status  = 'PENDIENTE'
  ORDER  BY scheduled_date ASC
  LIMIT  1
  → nextFollowUpDate = scheduled_date del primer registro (o null si no hay)


PASO 8 — Aplicar ordenamiento

  SEGÚN sort:
    CASO 'registration_date':
      → ORDER BY victim_case.victim_case_creation_date {order}

    CASO 'next_follow_up':
      → ORDER BY nextFollowUpDate {order} NULLS LAST


PASO 9 — Construir y retornar response

  200 {
    cases: [
      {
        id:                 victim_case_id,
        i_code:             victim_case_i_code,
        names:              victim_case_victim_names,
        lastNames:          victim_case_victim_last_names,
        docNumber:          victim_case_victim_doc_number,
        creationDate:       victim_case_creation_date,
        status:             victim_case_status,
        agentNames:         follow_up_v2.agent_names,      // campo virtual del modelo
        agentLastNames:     follow_up_v2.agent_last_names,  // campo virtual del modelo
        team:               follow_up_v2.team,
        riskStatus:         follow_up_v2.risk_status,
        nextFollowUpDate:   (calculado en PASO 7 o null)
      },
      ...
    ]
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                  | Paso afectado |
|--------------------------------------------------------------------------------------|---------------|
| Ruta exacta del endpoint backend (puede diferir de /api/v1/cases/list)              | PASO 2, 5     |
| ¿El número de teléfono de la víctima está en victim_case o en victim_case_form1?    | PASO 5        |
| ¿La tabla usa paginación o carga la lista completa en una sola llamada?              | PASO 2, 9     |
| Cuándo hay múltiples FollowUpV2 por caso, ¿cuál se usa para agentNames y riskStatus?| PASO 5, 9     |
| Valores exactos del enum riskStatus en follow_up_v2 (¿'alto','medio','bajo'?)       | PASO 6, 9     |
| Valores exactos del campo team en follow_up_v2 (lista cerrada o libre?)             | PASO 6        |
