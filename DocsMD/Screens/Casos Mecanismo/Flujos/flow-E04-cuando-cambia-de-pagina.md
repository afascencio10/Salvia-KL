# flow-E04 — Cuando cambia de página

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de página
   Tipo: User Interaction
   Funciones: changePage() · loadCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  page:             destino 0-based
  filters.document: vigente
  filters.city:     vigente
}

PASO 1 — Validar page ∈ [0, totalPages)

PASO 2 — currentPage = page; scroll top

PASO 3 — GET /api/v1/mecanismos/mine?document&city&page&pageSize=5
  (usuario desde sesión)

PASO 4 — Vue reemplaza items
```
