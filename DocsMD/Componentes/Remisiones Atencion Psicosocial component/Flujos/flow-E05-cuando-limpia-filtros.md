━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario presiona "Limpiar filtros"
   Tipo: User Interaction
   Código: E-05
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en ClearFiltersBtn "Limpiar filtros"


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Resetear filtros de UI (no el alcance fijo)

  activeFilters = { ...scopeFilter }   // restaurar solo defaultFilter del prop
  searchNumeroIdentidad = ""
  searchTelefono        = ""
  autocompleteText.profesional_asignada      = ""
  autocompleteSuggestions.profesional_asignada = []

  // Restaurar pre-selección según scopeFilter
  SI scopeFilter.agent_id:
    → activeFilters['agent_id'] = scopeFilter.agent_id
    → mantener autocompleteSelected si venía del defaultFilter

  SI scopeFilter.dupla_id:
    → activeFilters['dupla_id'] = scopeFilter.dupla_id

  SI NO scopeFilter.agent_id:
    → autocompleteSelected.profesional_asignada = null


PASO 2 — Resetear paginación y selección

  currentPage        = 1
  selectedRemisiones = []
  loading            = true


PASO 3 — Recargar listado y stats

  fetchRemisiones()   // con activeFilters (= scopeFilter)

  SI mostrarCards === true:
    → fetchStats()    // mismos filtros

PASO 4 — Re-renderizar FilterBar, cards (si aplica) y tabla

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  COMPORTAMIENTO POR TIPO DE PANTALLA
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Pantalla | defaultFilter | Tras "Limpiar filtros" |
|---|---|---|
| Listado general | `{}` | Lista completa sin filtros UI |
| Mis remisiones (psicólogo) | `{ agent_id }` | Solo remisiones de ese agent_id |
| Mis remisiones (dupla) | `{ dupla_id }` | Solo remisiones de esa dupla |
