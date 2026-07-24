# flow-E04 — Cuando cambia de página

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de página
   Tipo: User Interaction
   Funciones: changePage() · loadEntityCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  page:             página destino (0-based)
  totalPages:       computed
  filters.*:        entidad + documento + ciudad vigentes
}

PASO 1 — Validar page ∈ [0, totalPages) y entityId presente

PASO 2 — currentPage = page; scroll al top

PASO 3 — GET /api/v1/entity-cases con page = currentPage (mantiene filtros)

PASO 4 — Vue reemplaza items
```
