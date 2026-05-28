━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona o limpia una opción del autocomplete
   Tipo: User Interaction
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
  autocompleteText['persona_asignada']        = ""          // limpiar el input de texto


PASO A2 — Actualizar el filtro activo y resetear paginación

  activeFilter  = { key: 'persona_asignada', value: option.value }
  searchText    = ""
  currentPage   = 1
  loading       = true
  cases         = []


PASO A3 — Consultar backend con el nuevo filtro

  GET /api/v1/cases/list
    ?filter_key=persona_asignada
    &filter_value={option.value}    // icode del agente seleccionado
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI respuesta no ok:
    → loading   = false
    → loadError = data.error || 'Error al filtrar por persona asignada'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases      = data.cases
    → totalCases = data.total
    → loading    = false
    → Renderizar tabla y PaginationBar actualizados

  → FIN SITUACIÓN A ✓


┌──────────────────────────────────────────────────────────────┐
│  SITUACIÓN B — Limpiar la selección (presionar "✕")          │
└──────────────────────────────────────────────────────────────┘

PASO B1 — Limpiar el estado del autocomplete

  autocompleteSelected['persona_asignada']    = null
  autocompleteSuggestions['persona_asignada'] = []
  autocompleteText['persona_asignada']        = ""


PASO B2 — Resetear filtro activo al defaultFilter y paginación

  activeFilter  = defaultFilter    // vuelve al filtro con el que se montó el componente
  searchText    = ""
  currentPage   = 1
  loading       = true
  cases         = []


PASO B3 — Consultar backend con el defaultFilter

  GET /api/v1/cases/list
    ?filter_key={defaultFilter.key}
    &filter_value={defaultFilter.value}
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  SI respuesta no ok:
    → loading   = false
    → loadError = data.error || 'Error al cargar los casos'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases      = data.cases
    → totalCases = data.total
    → loading    = false
    → Renderizar tabla y PaginationBar actualizados

  → FIN SITUACIÓN B ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| ¿Al limpiar el autocomplete se vuelve al defaultFilter o a "sin filtro"?         | PASO B2       |
| ¿El dropdown de sugerencias se cierra al hacer clic fuera del componente?        | Interfaz      |
| ¿Se puede combinar el autocomplete activo con la búsqueda por ID/teléfono (E-03)? | PASO A3      |
