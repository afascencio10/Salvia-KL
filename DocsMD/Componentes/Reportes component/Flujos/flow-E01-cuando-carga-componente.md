━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  columns:         columnas         → prop :columns (Array<{ key, label }>) del padre
  hiddenColumns:   columnas ocultas → prop :hiddenColumns (Array<string>) del padre (opcional)
  buttons:         botones de fila  → prop :buttons (Array<{ id, label }>) del padre (opcional, default [])
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reportes-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar estado interno del componente

  searchName     = ""
  searchPhone    = ""
  sortBy         = "registration_date"
  sortOrder      = "desc"
  currentPage    = 1
  pageSize       = 20
  totalReports   = 0
  reports        = []
  loading        = true
  loadError      = null

  visibleColumns = columns.filter(col => !hiddenColumns.includes(col.key))
  // SI buttons.length === 0 → no renderizar columna "Acciones"


PASO 2 — Consultar backend con paginación inicial

  GET /api/v1/reports/list
    &sort={sortBy}
    &order={sortOrder}
    &page={currentPage}
    &page_size={pageSize}

  // Sin search_name ni search_phone en carga inicial

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al cargar los reportes'
    → Mostrar ErrorState en lugar de la tabla
    → TERMINAR ejecución

  SI respuesta ok:
    → reports      = data.reports
    → totalReports = data.total
    → loading      = false
    → CONTINÚA PASO 3


PASO 3 — Calcular propiedades derivadas

  totalPages      = Math.ceil(totalReports / pageSize)
  filteredReports = reports


PASO 4 — Renderizar tabla y controles de paginación

  PARA CADA report EN filteredReports:
    Renderizar fila con visibleColumns (ver reportes-component-interface.md)

    Renderizar columna de acciones (SI buttons.length > 0):
      PARA CADA btn EN buttons:
        → Renderizar botón btn.label
        → Al presionar → emitActionClicked(btn.id, report)  (→ E-05)

  Renderizar PaginationBar (→ E-04)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/reports/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  search_name:   filtro por nombre        → query param (opcional)
  search_phone:  filtro por teléfono      → query param (opcional)
  sort:          'registration_date'
  order:         'asc' | 'desc'
  page:          default 1
  page_size:     default 20
}


PASO 5 — Construir query base

  SELECT
    vc.victim_contact_id,
    vc.victim_contact_i_code,
    vc.victim_contact_creation_date,
    vc.victim_contact_names,
    vc.victim_contact_last_names,
    vc.victim_contact_status,
    vcf2.victim_contact_form2_victim_col_phone,
    vcf2.victim_contact_form2_best_contact_time,
    vcf2.victim_contact_form2_facts_description,
    rt.victim_case_form2_enums_code          AS report_type_code,
    rt.victim_case_form2_enums_name          AS report_type_label,
    wrc.victim_case_form2_enums_name         AS will_receive_call_label,
    (
      SELECT COALESCE(
        json_agg(
          json_build_object(
            'code', e.victim_case_form2_enums_code,
            'name', e.victim_case_form2_enums_name
          )
          ORDER BY e.victim_case_form2_enums_name ASC
        ),
        '[]'::json
      )
      FROM   salvia.rel_victim_case_form2_enums_victim_contact_form2 rel
      JOIN   salvia.victim_case_form2_enums e
             ON e.victim_case_form2_enums_id = rel.rel_victim_case_form2_enums_victim_contact_form2_victim_case_form2_enums
      WHERE  rel.rel_victim_case_form2_enums_victim_contact_form2_victim_contact_form2 = vcf2.victim_contact_form2_id
    ) AS adjustments_gbv_json

  FROM salvia.victim_contact vc

  INNER JOIN salvia.victim_contact_form2 vcf2
          ON vcf2.victim_contact_form2_victim_contact = vc.victim_contact_id

  -- Excluir reportes que ya tienen caso
  LEFT JOIN salvia.victim_case vcase
         ON vcase.victim_case_victim_contact = vc.victim_contact_id

  LEFT JOIN salvia.victim_case_form2_enums rt
         ON rt.victim_case_form2_enums_id = vcf2.victim_contact_form2_report_type

  LEFT JOIN salvia.victim_case_form2_enums wrc
         ON wrc.victim_case_form2_enums_id = vcf2.victim_contact_form2_will_receive_call

  WHERE vcase.victim_case_victim_contact IS NULL    -- sin caso asociado
    AND vc.victim_contact_status = 'v'              -- solo reportes válidos


PASO 6 — Aplicar filtros de búsqueda (search_name, search_phone)

  SI search_name != "":
    AND (
      vc.victim_contact_names ILIKE '%' || search_name || '%'
      OR vc.victim_contact_last_names ILIKE '%' || search_name || '%'
      OR (vc.victim_contact_names || ' ' || vc.victim_contact_last_names) ILIKE '%' || search_name || '%'
    )

  SI search_phone != "":
    AND vcf2.victim_contact_form2_victim_col_phone::text ILIKE '%' || search_phone || '%'

  // Si ambos vienen con valor → cláusulas AND (intersección)


PASO 7 — Aplicar ordenamiento

  ORDER BY vc.victim_contact_creation_date {order}


PASO 8 — Calcular total y aplicar paginación

  total  = COUNT(*) con filtros de PASO 5 y 6
  offset = (page - 1) * page_size
  → LIMIT page_size OFFSET offset


PASO 9 — Construir y retornar response

  200 {
    reports: [ { id, icode, creationDate, names, lastNames, adjustmentsGBV,
                reportTypeCode, reportTypeLabel, willReceiveCallLabel,
                victimPhone, bestContactTime, factsDescription }, ... ],
    total, page, pageSize
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  REFERENCIA — Código existente equivalente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

La condición "sin caso + status válido" ya existe en:

  salvia_daos.GetVictimContactsWithoutVictimCase(page, "v", ...)
  → LEFT JOIN victim_case ... WHERE victim_case IS NULL AND status = 'v'

Usado en facades con filtro `fcv` (get_victim_cases.html, get_victim_contacts).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                                         | Paso afectado |
|-------------------------------------------------------------------------------------------------------------|---------------|
| Ruta exacta del endpoint backend (propuesta: `/api/v1/reports/list`)                                       | PASO 2, 5     |
| ¿Filtrar solo reportes con `authorizationAnswer = 'y'`?                                                     | PASO 5        |
| Nombres exactos de columnas en `victim_case_form2_enums` (code/name vs icode)                              | PASO 5        |
| Tamaño de página por defecto: ¿20 registros es correcto?                                                   | PASO 1, 8     |

> 📌 **DECISIÓN APLICADA:** Solo reportes sin `victim_case` y con `victim_contact_status = 'v'`.
> Los invalidados (`'i'`) quedan fuera de alcance de este componente.
