━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona o limpia una opción del autocomplete
   Tipo: User Interaction
   Código: E-08
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado en dos situaciones:
  (A) El usuario hace clic en una SuggestionItem del SuggestionsDropdown
  (B) El usuario presiona "✕" (ClearBtn) en el SelectedTag

INPUT — Situación A (seleccionar):
{
  filterKey:   'persona_asignada'       → constante
  option:      { value: icode, label }  → sugerencia seleccionada
}

INPUT — Situación B (limpiar):
{
  filterKey:   'persona_asignada'       → constante
  option:      null                     → indica que se borra la selección
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


┌──────────────────────────────────────────────────────────────┐
│  SITUACIÓN A — Seleccionar una sugerencia                    │
└──────────────────────────────────────────────────────────────┘

PASO A1 — Fijar la selección en el estado del autocomplete

  autocompleteSelected['persona_asignada']    = option      // { value: icode, label }
  autocompleteSuggestions['persona_asignada'] = []          // cerrar el dropdown
  autocompleteOpen['persona_asignada']        = false
  autocompleteText['persona_asignada']        = ""          // limpiar el input de texto


PASO A2 — Actualizar el filtro activo y resetear paginación

  activeFilter  = { key: 'persona_asignada', value: option.value }
  searchText    = ""          // solo búsqueda E-03; chip/dropdown se mantienen
  currentPage   = 1
  loading       = true
  cases         = []


PASO A3 — Consultar backend con filtros combinados

  GET /api/v1/cases/list
    ?filter_key=persona_asignada
    &filter_value={option.value}              // icode del agente
    &chip_filter=casos_nuevos                  // si chip activo (E-09)
    &dropdown_filter_key=riesgo|equipo        // si dropdown activo (E-10/E-11)
    &dropdown_filter_value=...
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI respuesta no ok:
    → loading   = false
    → loadError = data.error || 'Error al filtrar por persona asignada'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases, totalCases, loading = false
    → Renderizar tabla y PaginationBar

  → FIN SITUACIÓN A ✓


┌──────────────────────────────────────────────────────────────┐
│  SITUACIÓN B — Limpiar la selección (presionar "✕")          │
└──────────────────────────────────────────────────────────────┘

PASO B1 — Limpiar el estado del autocomplete

  autocompleteSelected['persona_asignada']    = null
  autocompleteSuggestions['persona_asignada'] = []
  autocompleteOpen['persona_asignada']        = false
  autocompleteText['persona_asignada']        = ""


PASO B2 — Resetear filtro activo al defaultFilter y paginación

  activeFilter  = normalizeDefaultActiveFilter(defaultFilter)
  // Sin persona_asignada; chip/dropdown activos se conservan
  searchText    = ""
  currentPage   = 1
  loading       = true
  cases         = []


PASO B3 — Consultar backend sin persona_asignada

  GET /api/v1/cases/list
    ?chip_filter=...                          // filtros aditivos vigentes
    &dropdown_filter_key=...
    &search=...
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  → FIN SITUACIÓN B ✓


┌──────────────────────────────────────────────────────────────┐
│  Cierre del dropdown al clic fuera                           │
└──────────────────────────────────────────────────────────────┘

  Listener document.click en mounted (beforeUnmount lo remueve)
  SI clic fuera del contenedor [data-autocomplete-key] del filtro:
    → autocompleteOpen[key] = false
    → autocompleteSuggestions[key] = []


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas (E-08 implementado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Al limpiar autocomplete | `activeFilter` → `normalizeDefaultActiveFilter()` (sin persona_asignada) |
| Clic fuera del dropdown | Cierra sugerencias vía `handleDocumentClick` |
| Combinación con otros filtros | `fetchCases` envía chip_filter + dropdown_filter en paralelo con persona_asignada |
