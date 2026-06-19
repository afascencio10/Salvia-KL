━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario cambia el ordenamiento
   Tipo: User Interaction
   Código: E-03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en el `<select>` SortSelector de FilterBar
               Handler: `onSortChange(event)`

INPUT: {
  sortOption:   valor del `<select>` (criterio + dirección en una sola clave)
                'registration_date_desc' | 'registration_date_asc'
}

> **Orden inicial (E-01):** `registration_date_desc` — Fecha de registro DESC.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reportes-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Parsear la opción seleccionada

  SI sortOption no es uno de los 2 valores válidos:
    → TERMINAR ejecución (sin cambios)

  SI sortOption === sortBy + '_' + sortOrder (misma opción que la activa):
    → TERMINAR ejecución

  → sortBy    = 'registration_date'
  → sortOrder = parte después del último '_' (asc | desc)
  → CONTINÚA PASO 2


PASO 2 — Resetear paginación y recargar desde backend

  currentPage = 1
  reports     = []
  loading     = true

  GET /api/v1/reports/list
    ?search_name={searchName}
    &search_phone={searchPhone}
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI respuesta no ok:
    → loadError = data.error || 'Error al ordenar los reportes'
    → TERMINAR ejecución

  SI respuesta ok:
    → reports, totalReports, loading = false
    → CONTINÚA PASO 3


PASO 3 — Actualizar vista

  filteredReports = reports
  → Renderizar tabla + PaginationBar

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  UI — opciones del select
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| value                      | Texto visible              |
|----------------------------|----------------------------|
| registration_date_desc     | Fecha de registro DESC     |
| registration_date_asc      | Fecha de registro ASC      |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/reports/list  (PASO 7 de E-01)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ORDER BY vc.victim_contact_creation_date {order}

  → FIN EJECUCIÓN ✓
