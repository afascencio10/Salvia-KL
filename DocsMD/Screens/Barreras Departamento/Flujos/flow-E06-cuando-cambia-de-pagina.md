# flow-E06 — Cuando cambia de página

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de página
   Tipo: User Interaction
   Funciones: changePage() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Validar page ∈ [0, totalPages)
PASO 2 — currentPage = page; scroll top
PASO 3 — GET /api/v1/barriers-v2/department con filtros vigentes + page + pageSize=10
PASO 4 — Reemplazar items
```
