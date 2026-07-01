━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario cambia de página
   Tipo: User Interaction
   Código: E-06
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en PrevBtn o NextBtn de PaginationBar

INPUT: {
  newPage:  currentPage ± 1
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Validar rango

  SI newPage < 1 || newPage > totalPages || newPage === currentPage:
    → TERMINAR

PASO 2 — Actualizar estado

  currentPage        = newPage
  selectedRemisiones = []
  loading            = true
  remisiones         = []

PASO 3 — fetchRemisiones() manteniendo filtros activos

  GET /api/v1/psychosocial-support/list
    ?page={newPage}
    &page_size={pageSize}
    (+ todos los filter_* de activeFilters)

  SI mostrarCards === true:
    → fetchStats()   // mismos filtros; stats no dependen de la página

PASO 4 — Scroll al inicio de .rps-wrapper y re-renderizar

  → FIN EJECUCIÓN ✓
