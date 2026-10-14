# flow-E04 — Cuando cambia de página

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de página
   Tipo: User Interaction
   Función: goToPage(page)
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: page  (0-based)

SI page < 0 o page >= numPages → TERMINAR
currentPage = page
scroll al tope

paginatedBarreras = slice de filteredBarreras (5 por página)
```
