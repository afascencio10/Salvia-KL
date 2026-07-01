━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el filtro de número de identidad
   Tipo: User Interaction
   Código: E-03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en SearchFilter `numero_identidad` — debounce 400ms

INPUT: {
  queryText:  texto ingresado  → v-model searchNumeroIdentidad
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Debounce 400ms; queryText = trim(searchNumeroIdentidad)

PASO 2 — Resetear paginación y selección

  currentPage        = 1
  selectedRemisiones = []
  loading            = true

PASO 3 — Actualizar activeFilters

  SI queryText === "":
    → delete activeFilters['numero_identidad']
  SI NO:
    → activeFilters['numero_identidad'] = queryText

PASO 4 — fetchRemisiones() con todos los filtros activos

  GET /api/v1/psychosocial-support/list
    &filter_numero_identidad={queryText}   // omitir si vacío
    (+ demás filtros activos)

  → Renderizar tabla o EmptyState

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — cláusula WHERE adicional
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```sql
AND vc.victim_case_victim_doc_number ILIKE '%' || {filter_numero_identidad} || '%'
```

Mínimo recomendado: 3 caracteres antes de consultar (opcional, igual que autocomplete).
