# flow-E04 — Cuando cambia de página

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de página
   Tipo: User Interaction
   Funciones: goToPage() · loadOficios()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  page: número de página destino (0-based)
}

PASO 1 — Validar rango
  SI page < 0 O page >= numPages → TERMINAR

PASO 2 — Actualizar estado
  currentPage = page
  window.scrollTo(0, 0)

PASO 3 — Consultar la página en BD (mantiene tab y filtros activos)

GET /api/v1/entity-letters?page={page}&limit=5&{filtros}

PASO 4 — Reemplazar oficios con response.items de la nueva página
```
