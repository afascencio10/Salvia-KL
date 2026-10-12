# flow-E03 — Cuando filtra por ciudad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por ciudad
   Tipo: User Interaction
   Funciones: onFilterCity() · loadCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filters.city: texto del input
}

PASO 1 — Debounce ~350 ms

PASO 2 — currentPage = 0

PASO 3 — GET /api/v1/mecanismos/mine
  query: { document, city: trim, page: 0, pageSize: 5 }
  Backend: ILIKE parcial sobre el nombre de ciudad de la víctima
           (town → city), siempre acotado a assigned_user_id de sesión

PASO 4 — Vue actualiza items / totalItems
```
