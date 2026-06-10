━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario cambia de página
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en PrevBtn ("← Anterior") o NextBtn ("Siguiente →")
               en la PaginationBar al pie de la tabla

INPUT: {
  newPage:   número de página solicitada   → currentPage ± 1
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
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

  currentPage    = newPage
  loading        = true
  cases          = []       // limpiar filas mientras carga nueva página
  selectedCases  = []       // limpiar selección de reasignación (→ E-14); oculta botón "Reasignar Casos"


PASO 3 — Consultar backend con la nueva página

  GET /api/v1/cases/list
    ?filter_key={activeFilter.key}
    &filter_value={activeFilter.value}
    &sort={sortBy}
    &order={sortOrder}
    &page={currentPage}
    &page_size={pageSize}

  // El filtro activo y el ordenamiento se mantienen igual
  // Solo cambia el parámetro page

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al cargar la página'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases      = data.cases
    → totalCases = data.total   // debería ser el mismo valor, pero se actualiza por consistencia
    → loading    = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista y hacer scroll al inicio de la tabla

  filteredCases = cases   // searchText se mantiene; si hay texto activo, re-aplica el filtro local

  SI searchText !== "":
    term = searchText.toLowerCase()
    filteredCases = cases.filter(caso =>
      caso.i_code.toLowerCase().includes(term) ||
      caso.phone.toLowerCase().includes(term)
    )

  → Scroll al inicio del componente (.cc-wrapper)
  → Renderizar tabla con filteredCases actualizados
  → Actualizar PaginationBar con currentPage y totalPages

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| ¿El scroll al inicio es al tope del componente o al tope de la página?           | PASO 4        |
| ¿Se debe mostrar el número de página actual en la URL (query param)?             | PASO 2        |
| ¿La búsqueda local (searchText) se resetea al cambiar de página?                 | PASO 4        |
