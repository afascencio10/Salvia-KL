━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario cambia el ordenamiento
   Tipo: User Interaction
   Código: E-04
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en el `<select>` SortSelector de FilterBar
               Handler: `onSortChange(event)`

INPUT: {
  sortOption:   valor del `<select>` (criterio + dirección en una sola clave)
                'registration_date_desc' | 'registration_date_asc'
                'next_follow_up_desc'    | 'next_follow_up_asc'
}

> **Orden inicial (E-01):** `registration_date_desc` — Fecha de registro DESC.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Parsear la opción seleccionada

  SI sortOption no es uno de los 4 valores válidos:
    → TERMINAR ejecución (sin cambios)

  SI sortOption === sortBy + '_' + sortOrder (misma opción que la activa):
    → TERMINAR ejecución

  → sortBy    = parte antes del último '_'  (registration_date | next_follow_up)
  → sortOrder = parte después del último '_' (asc | desc)
  → CONTINÚA PASO 2


PASO 2 — Resetear paginación y recargar desde backend

  currentPage = 1
  cases       = []
  loading     = true

  GET /api/v1/cases/list
    ?sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}
    // + filtros activos (chip, dropdowns, agentId, search) sin modificar

  SI respuesta no ok:
    → loadError = data.error || 'Error al ordenar los casos'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases, totalCases, loading = false
    → CONTINÚA PASO 3


PASO 3 — Actualizar vista

  filteredCases = cases   // orden ya aplicado en SQL
  → Renderizar tabla + PaginationBar

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  UI — case_component.html (opciones del select)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| value                      | Texto visible              |
|----------------------------|----------------------------|
| registration_date_desc     | Fecha de registro DESC     |
| registration_date_asc      | Fecha de registro ASC      |
| next_follow_up_desc        | Próximo seguimiento DESC   |
| next_follow_up_asc         | Próximo seguimiento ASC    |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (PASO 7 de E-01)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 4 — Aplicar ORDER BY

  SEGÚN sort:
    CASO 'registration_date':
      → ORDER BY vc.victim_case_creation_date {order}

    CASO 'next_follow_up':
      → ORDER BY next_follow_up_date {order} NULLS LAST

  order = 'asc' | 'desc' (query param)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas (E-04)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Orden por defecto al cargar | `registration_date` + `desc` |
| Toggle al re-elegir mismo criterio | No — cada dirección es opción explícita |
| ¿Ordenamiento en cliente o servidor? | Servidor (paginación coherente) |
| NULL en próximo seguimiento | `NULLS LAST` en ambos sentidos |
