# flow-E02 — Cuando filtra por documento

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por documento
   Tipo: User Interaction
   Función: watcher filters (deep)
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: filters.identidad  (texto)

PASO 1 — currentPage = 0

PASO 2 — filteredBarreras en cliente:
  SI identidad vacío → pasa
  SI NO → docNumber.includes(identidad)

No hay fetch adicional.
```
