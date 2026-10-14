# flow-E03 — Cuando filtra por entidad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por entidad
   Tipo: User Interaction
   Funciones: onFilterText() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Debounce 350 ms
PASO 2 — currentPage = 0
PASO 3 — GET …?entidad=trim (ILIKE entity_branch_name) page=0
PASO 4 — Actualizar items / total
```
