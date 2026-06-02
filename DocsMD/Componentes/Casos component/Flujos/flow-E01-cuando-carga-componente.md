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

      'victim_info'       → case.names + " " + case.lastNames
                              Doc: case.docNumber  ·  Tel: case.victimPhone
                              // docNumber ← victim_case.victim_case_victim_doc_number
                              // victimPhone ← victim_case_form2_victim_phone, fallback victim_contact_form1_phone
      'operator'          → case.ownerNames + " " + case.ownerLastNames  (o "—" si vacío)
      'registration_date' → case.creationDate  formateada DD/MM/YYYY
      'risk_level'        → RiskBadge  usando case.riskStatus
                              // victim_case_form2.victim_case_form2_risk_level (1-4)
                              // 1=bajo, 2=moderado, 3=alto, 4=extremo
      'team'              → case.caseTeam  (victim_case_team; o "—" si vacío)
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
    COALESCE(
        NULLIF(vf2.victim_case_form2_victim_phone::text, ''),
        (SELECT vcf1.victim_contact_form1_phone::text
         FROM salvia.victim_contact_form1 vcf1
         WHERE vcf1.victim_contact_form1_victim_contact = vc.victim_case_victim_contact
         LIMIT 1),
        ''
    )                                  AS victim_phone,
    vc.victim_case_creation_date,
    vc.victim_case_status,
    vc.agent_id,                                        // nuevo campo en victim_case
    gup.general_user_profile_names     AS owner_names,
    gup.general_user_profile_last_names AS owner_last_names,
    gu.general_user_team               AS owner_team,
    COALESCE(vc.victim_case_team, '')  AS case_team,
    CASE vf2.victim_case_form2_risk_level
        WHEN 1 THEN 'bajo'
        WHEN 2 THEN 'moderado'
        WHEN 3 THEN 'alto'
        WHEN 4 THEN 'extremo'
    END                              AS risk_status,
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

  -- Perfil del agente asignado (nombres y equipo)
  LEFT JOIN security.general_user gu
         ON gu.general_user_i_code = vc.agent_id
  LEFT JOIN security.general_user_profile gup
         ON gup.general_user_profile_id = gu.general_user_general_user_profile

  -- Nivel de riesgo del tamizaje (formulario 2 del caso)
  LEFT JOIN salvia.victim_case_form2 vf2
         ON vf2.victim_case_form2_victim_case = vc.victim_case_id


PASO 6 — Aplicar filtro según filter_key

  SEGÚN filter_key:

    CASO 'casos_nuevos':
      → WHERE (vc.victim_case_creation_date AT TIME ZONE '{tz}')::date
            = (NOW() AT TIME ZONE '{tz}')::date
      // Casos creados hoy — Ver flow-E09

    CASO 'riesgo':
      → WHERE vf2.victim_case_form2_risk_level = {nivel}
      // filter_value: 'bajo'→1, 'moderado'→2, 'alto'→3, 'extremo'→4 — Ver flow-E10

    CASO 'equipo':
      → WHERE vc.victim_case_team = filter_value
      // Equipo del caso, no del agente — Ver flow-E11

    CASO 'persona_asignada':
      → WHERE vc.agent_id = filter_value         // filter_value = icode del agente
      // No requiere JOIN con rel_case_owner_victim_case; el agente está en victim_case.agent_id

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
        victimPhone:       victim_phone,       // form2 al crear caso; fallback contact form1
        creationDate:      victim_case_creation_date,
        status:            victim_case_status,
        ownerNames:        owner_names,        // del agente en victim_case.agent_id → general_user_profile
        ownerLastNames:    owner_last_names,
        ownerTeam:         owner_team,         // equipo del agente (general_user)
        caseTeam:          case_team,          // victim_case.victim_case_team
        riskStatus:        risk_status,        // victim_case_form2_risk_level: 1=bajo, 2=moderado, 3=alto, 4=extremo
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

| Variable / decisión                                                                                         | Paso afectado |
|-------------------------------------------------------------------------------------------------------------|---------------|
| Ruta exacta del endpoint backend (puede diferir de /api/v1/cases/list)                                     | PASO 2, 5     |
| Nombre exacto de la columna en BD: ¿agent_id? Confirmar con migración o DDL de salvia.victim_case          | PASO 5, 6     |
| Tabla y campo exacto del perfil del agente: ¿security.general_user_profile.general_user_i_code?            | PASO 5        |
| ¿El campo de teléfono de la víctima existe en victim_case o en victim_case_form1?                          | PASO 5        |
| Cuando hay múltiples follow_up_v2 no cerrados por caso, ¿cuál se usa para riskStatus?                      | PASO 5        |
| Valores exactos del enum riskStatus en follow_up_v2 (¿'alto','medio','bajo'?)                              | PASO 6, 9     |
| Tamaño de página por defecto: ¿20 registros es correcto?                                                   | PASO 1, 8     |
| ¿El campo agent_id en victim_case puede ser NULL? Si es NULL, el caso no tiene agente asignado             | PASO 5, 6     |

> 📌 DECISIÓN DE DISEÑO: El filtro 'persona_asignada' ya no usa salvia.rel_case_owner_victim_case.
> El JOIN con esa tabla se eliminó de la query base. El agente asignado se resuelve directamente
> desde victim_case.agent_id → security.general_user_profile.
