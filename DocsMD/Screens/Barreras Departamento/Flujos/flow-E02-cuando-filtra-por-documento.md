# flow-E02 — Cuando filtra por documento

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por documento
   Tipo: User Interaction
   Funciones: onFilterText() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Debounce 350 ms
PASO 2 — currentPage = 0
PASO 3 — GET /api/v1/barriers-v2/department?docNumber=trim&page=0&pageSize=10
         (+ entidad, ciudad, status si están activos)
PASO 4 — Actualizar items / total
```
