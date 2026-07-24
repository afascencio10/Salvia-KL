# flow-E03 — Cuando filtra por ciudad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por ciudad
   Tipo: User Interaction
   Funciones: onFilterCity() · loadEntityCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filters.city:     nombre de ciudad del <select>   → opciones de GET /entities/:id/cities
  filters.document: filtro documento vigente
  filters.entityId: entidad activa
}

PASO 1 — currentPage = 0

PASO 2 — SI !filters.entityId → TERMINAR

PASO 3 — GET /api/v1/entity-cases
  query: { entityId, document, city (nombre exacto del select), page=0, pageSize=5 }
  Backend filtra por ciudad del CASO (victim_case → town → city).

PASO 4 — Vue actualiza items / totalItems
```
