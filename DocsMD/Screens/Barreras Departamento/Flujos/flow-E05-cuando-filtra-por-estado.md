# flow-E05 — Cuando filtra por estado

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por estado
   Tipo: User Interaction
   Funciones: onFilterSelect() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Sin debounce (select)
PASO 2 — currentPage = 0
PASO 3 — Labels UI ≠ valor enviado:
         Abierta → OPEN
         En gestión → En Gestion
         Articulada → Articulada
         Gestionada → MANAGED
         Todos → no envía status
PASO 4 — GET API y actualizar lista
```
