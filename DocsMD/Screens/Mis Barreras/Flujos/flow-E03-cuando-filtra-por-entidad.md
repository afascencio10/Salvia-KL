# flow-E03 — Cuando filtra por entidad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por entidad
   Tipo: User Interaction
   Función: watcher filters (deep)
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: filters.entidad  (texto)

PASO 1 — currentPage = 0

PASO 2 — filteredBarreras en cliente:
  SI entidad vacío → pasa
  SI NO → sector o entityBranchName contiene el texto (case-insensitive)

No hay fetch adicional.
```
