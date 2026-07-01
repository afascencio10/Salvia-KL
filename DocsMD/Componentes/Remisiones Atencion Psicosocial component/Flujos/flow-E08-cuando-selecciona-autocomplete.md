━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona o limpia el autocomplete de profesional asignada
   Tipo: User Interaction
   Código: E-08
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  (A) clic en SuggestionItem del dropdown
  (B) clic en "✕" del SelectedTag


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — caso (A) seleccionar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Guardar selección

  autocompleteSelected['profesional_asignada'] = { value: opt.value, label: opt.label }
  autocompleteText['profesional_asignada']     = ""
  autocompleteSuggestions                      = []

PASO 2 — Activar filtro y recargar

  activeFilters['agent_id'] = opt.value   // general_user_i_code
  currentPage               = 1
  selectedRemisiones        = []
  loading                   = true

PASO 3 — fetchRemisiones() con todos los filtros activos

  GET /api/v1/psychosocial-support/list
    &filter_agent_id={opt.value}
    (+ demás filter_*)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — caso (B) limpiar
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Limpiar estado autocomplete

  autocompleteSelected['profesional_asignada'] = null
  autocompleteText['profesional_asignada']     = ""

PASO 2 — Quitar filtro y recargar

  delete activeFilters['agent_id']
  currentPage        = 1
  selectedRemisiones = []

PASO 3 — fetchRemisiones()

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — cláusula WHERE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND ps.agent_id = {filter_agent_id}
```

Filtra por `psychosocial_support.agent_id` (profesional asignado directamente).
Remisiones asignadas solo vía dupla (sin `agent_id`) no aparecen con este filtro.
