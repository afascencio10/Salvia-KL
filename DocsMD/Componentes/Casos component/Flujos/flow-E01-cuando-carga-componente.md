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
  currentPage    = 1
  pageSize       = 20                         // registros por página
  totalCases     = 0
  cases          = []
  loading        = true
  loadError      = null

  visibleColumns = columns.filter(col => !hiddenColumns.includes(col.key))
  // SI hiddenColumns vacío → visibleColumns = columns completo


PASO 2 — Consultar backend con el filtro inicial y paginación

  GET /api/v1/cases/list
    ?filter_key={activeFilter.key}
    &filter_value={activeFilter.value}
    &sort={sortBy}
    &order={sortOrder}
    &page={currentPage}
    &page_size={pageSize}

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al cargar los casos'
    → Mostrar ErrorState en lugar de la tabla
    → TERMINAR ejecución

  SI respuesta ok:
    → cases      = data.cases        // Array de objetos de caso (página actual)
    → totalCases = data.total        // total de casos que cumplen el filtro
    → loading    = false
    → CONTINÚA PASO 3


PASO 3 — Calcular propiedades derivadas

  totalPages     = Math.ceil(totalCases / pageSize)
  filteredCases  = cases              // la búsqueda local (E-03) opera sobre esta lista


PASO 4 — Renderizar tabla y controles de paginación

  PARA CADA caso EN filteredCases:
    Renderizar fila con visibleColumns:

      'victim_info'       → case.names + " " + case.lastNames  /  case.i_code
      'operator'          → case.ownerNames + " " + case.ownerLastNames  (o "—" si vacío)
      'registration_date' → case.creationDate  formateada DD/MM/YYYY
      'risk_level'        → RiskBadge  usando case.riskStatus
      'team'              → case.ownerTeam  (o "—" si vacío)
      'assigned_person'   → case.ownerNames + " " + case.ownerLastNames  (o "—" si vacío)
      'next_follow_up_date' →
          SI case.nextFollowUpDate existe: mostrar fecha formateada DD/MM/YYYY
          SI no existe: mostrar "—"

    Renderizar columna de acciones:
      PARA CADA btn EN buttons:
        → Renderizar botón btn.label
        → Al presionar → emitActionClicked(btn.id, case)  (→ Ver E-05)

  Renderizar PaginationBar:
    PrevBtn  :disabled si currentPage === 1        → changePage(currentPage - 1)
    PageInfo "Página {currentPage} de {totalPages}"
    NextBtn  :disabled si currentPage === totalPages → changePage(currentPage + 1)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filter_key:    key del filtro activo      → query param
  filter_value:  valor del filtro           → query param (opcional según filter_key)
  sort:          campo de ordenamiento      → query param: 'registration_date' | 'next_follow_up'
  order:         dirección                  → query param: 'asc' | 'desc'
  page:          número de página           → query param (default: 1)
  page_size:     registros por página       → query param (default: 20)
}


PASO 5 — Construir query base

  SELECT
    vc.victim_case_id,
    vc.victim_case_i_code,
    vc.victim_case_victim_names,
    vc.victim_case_victim_last_names,
    vc.victim_case_victim_doc_number,
    vc.victim_case_creation_date,
    vc.victim_case_status,
    u.general_user_profile_names     AS owner_names,
    u.general_user_profile_last_names AS owner_last_names,
    u.general_user_team              AS owner_team,
    fu.risk_status,
    (
      SELECT scheduled_date
      FROM   salvia.follow_up_v2
      WHERE  case_id = vc.victim_case_i_code
      AND    status  = 'PENDIENTE'
      AND    deleted_at IS NULL
      ORDER  BY scheduled_date ASC
      LIMIT  1
    ) AS next_follow_up_date

  FROM salvia.victim_case vc

  -- Operador/dueño activo del caso
  LEFT JOIN salvia.rel_case_owner_victim_case rel
         ON rel.victim_case_id = vc.victim_case_id
        AND rel.rel_case_owner_victim_case_status = 'a'

  -- Perfil del agente asignado (nombres y equipo)
  LEFT JOIN security.general_user_profile u
         ON u.general_user_i_code = rel.case_owner_id   // ⚠️ GAP: confirmar la tabla/campo exacto

  -- Último follow_up_v2 activo para el riesgo
  LEFT JOIN salvia.follow_up_v2 fu
         ON fu.case_id   = vc.victim_case_i_code
        AND fu.status   != 'CERRADO'
        AND fu.deleted_at IS NULL


PASO 6 — Aplicar filtro según filter_key

  SEGÚN filter_key:

    CASO 'casos_nuevos':
      → WHERE NOT EXISTS (
            SELECT 1 FROM salvia.follow_up_v2
            WHERE  case_id   = vc.victim_case_i_code
            AND    status    = 'REALIZADO'
            AND    deleted_at IS NULL
        )

    CASO 'riesgo':
      → WHERE fu.risk_status = filter_value

    CASO 'equipo':
      → WHERE u.general_user_team = filter_value

    CASO 'persona_asignada':
      → WHERE rel.case_owner_id = filter_value   // filter_value = icode del agente

    DEFAULT:
      → Sin cláusula WHERE adicional


PASO 7 — Aplicar ordenamiento

  SEGÚN sort:
    CASO 'registration_date':
      → ORDER BY vc.victim_case_creation_date {order}

    CASO 'next_follow_up':
      → ORDER BY next_follow_up_date {order} NULLS LAST


PASO 8 — Calcular total y aplicar paginación

  // Ejecutar la misma query sin LIMIT/OFFSET para obtener el total
  total = COUNT(*) de la query del PASO 5 con los filtros del PASO 6

  // Aplicar paginación a la query principal
  offset = (page - 1) * page_size
  → LIMIT page_size OFFSET offset


PASO 9 — Construir y retornar response

  200 {
    cases: [
      {
        id:                victim_case_id,
        i_code:            victim_case_i_code,
        names:             victim_case_victim_names,
        lastNames:         victim_case_victim_last_names,
        docNumber:         victim_case_victim_doc_number,
        creationDate:      victim_case_creation_date,
        status:            victim_case_status,
        ownerNames:        owner_names,        // del agente en rel_case_owner_victim_case activo
        ownerLastNames:    owner_last_names,
        ownerTeam:         owner_team,
        riskStatus:        risk_status,        // del follow_up_v2 más reciente no cerrado
        nextFollowUpDate:  next_follow_up_date  // null si no hay pendiente
      },
      ...
    ],
    total:    N,         // total de casos que cumplen el filtro (sin paginación)
    page:     P,         // página actual
    pageSize: PS         // registros por página
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                  | Paso afectado |
|--------------------------------------------------------------------------------------|---------------|
| Ruta exacta del endpoint backend (puede diferir de /api/v1/cases/list)              | PASO 2, 5     |
| Tabla y campo exacto del perfil del agente: ¿security.general_user_profile?         | PASO 5        |
| ¿El campo de teléfono de la víctima existe en victim_case o en victim_case_form1?   | PASO 5        |
| Cuando hay múltiples follow_up_v2 no cerrados por caso, ¿cuál se usa para riskStatus? | PASO 5     |
| Valores exactos del enum riskStatus en follow_up_v2 (¿'alto','medio','bajo'?)       | PASO 6, 9     |
| Tamaño de página por defecto: ¿20 registros es correcto?                            | PASO 1, 8     |
