# flow-E04 — Cuando cambia de página

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de página
   Tipo: User Interaction
   Funciones: changePage() · loadEntityCases()
   Estado: planeado (sin cambio funcional mayor; sin city/entityId)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  page:             destino 0-based
  filters.document: vigente
  entityBranchId:   sesión (API)
}

PASO 1 — Validar page ∈ [0, totalPages)

PASO 2 — currentPage = page; scroll top

PASO 3 — GET /api/v1/entity-cases?document&page&pageSize=5
  (sede desde sesión)

PASO 4 — Vue reemplaza items
```
