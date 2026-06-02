━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por casos nuevos
   Tipo: User Interaction
   Código: E-09  (antes parte de E-02)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en el Chip **"Casos nuevos"** (`filter.type === 'chip'`,
               `filter.key === 'casos_nuevos'`) dentro de FilterGroup.
               Handler: `toggleChipFilter('casos_nuevos')`

> **Definición de negocio:** "Casos nuevos" = casos **creados en el día de hoy**,
> según `salvia.victim_case.victim_case_creation_date`.
> Ya no usa la lógica legacy de "sin seguimiento REALIZADO".

INPUT: {
  (ninguno extra — el filtro es binario on/off)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Determinar si el chip se activa o desactiva

  SI activeChipKey === 'casos_nuevos':
    → Desactivar: activeChipKey = null
    → CONTINÚA PASO 2

  SI NO:
    → Activar: activeChipKey = 'casos_nuevos'
    → CONTINÚA PASO 2

  // activeChipKey es independiente de activeFilter y de agentId
  // El chip activo usa clase cc-chip--active en la UI


PASO 2 — Resetear estado de búsqueda y paginación

  searchText   = ""
  currentPage  = 1
  loading      = true
  loadError    = null
  cases        = []


PASO 3 — Consultar backend

  GET /api/v1/cases/list
    ?chip_filter=casos_nuevos          // cuando el chip está activo
    &filter_key=persona_asignada       // solo si prop agentId está definida (Mis casos)
    &filter_value={agentId}
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  // chip_filter es aditivo: se combina con agentId, search y ordenamiento

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al aplicar el filtro'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases      = data.cases
    → totalCases = data.total
    → loading    = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista

  filteredCases = cases

  SI filteredCases.length === 0:
    → EmptyState "No se encontraron casos"

  SI filteredCases.length > 0:
    → Renderizar tabla + PaginationBar

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (filter_key = casos_nuevos)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 5 — Agregar cláusula WHERE sobre la query base de E-01

  CASO chip_filter = 'casos_nuevos' (o filter_key legacy = 'casos_nuevos'):

    WHERE (vc.victim_case_creation_date AT TIME ZONE 'America/Bogota')::date
        = (NOW() AT TIME ZONE 'America/Bogota')::date

  // Solo casos registrados hoy en hora de Bogotá
  // Se combina con AND al resto de filtros (persona_asignada, search, etc.)


PASO 6 — Respuesta

  Idéntica a E-01 PASO 9 (cases[], total, page, pageSize).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas (E-09 implementado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Zona horaria "hoy" | `America/Bogota` |
| Combinación con `agentId` | `chip_filter=casos_nuevos` + `filter_key=persona_asignada` en paralelo |
| Estado activo del chip | `activeChipKey` + clase `cc-chip--active`; segundo clic desactiva |
