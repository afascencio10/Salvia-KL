━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario cambia de página
   Tipo: User Interaction
   Código: E-04
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en PrevBtn ("← Anterior") o NextBtn ("Siguiente →")
               en la PaginationBar al pie de la tabla

INPUT: {
  newPage:   número de página solicitada   → currentPage ± 1
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reportes-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar que la página solicitada es alcanzable

  SI newPage < 1 || newPage > totalPages:
    → No hacer nada (el botón debería estar :disabled, pero se verifica igual)
    → TERMINAR ejecución

  SI newPage === currentPage:
    → No hacer nada
    → TERMINAR ejecución

  SI NO:
    → CONTINÚA PASO 2


PASO 2 — Actualizar estado

  currentPage = newPage
  loading     = true
  reports     = []       // limpiar filas mientras carga nueva página


PASO 3 — Consultar backend con la nueva página

  GET /api/v1/reports/list
    ?search_name={searchName}
    &search_phone={searchPhone}
    &sort={sortBy}
    &order={sortOrder}
    &page={currentPage}
    &page_size={pageSize}

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al cargar la página'
    → TERMINAR ejecución

  SI respuesta ok:
    → reports      = data.reports
    → totalReports = data.total
    → loading      = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista y hacer scroll al inicio de la tabla

  filteredReports = reports

  → Scroll al inicio del componente (.rc-wrapper)
  → Renderizar tabla con filteredReports actualizados
  → Actualizar PaginationBar con currentPage y totalPages

  → FIN EJECUCIÓN ✓
