# flow-E03 — Cuando aplica filtros

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando aplica filtros
   Tipo: User Interaction
   Funciones: watch filters · loadOficios()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filters.estado:    string  — estado del oficio (select)
  filters.identidad: string  — número de documento de la víctima
  filters.entidad:   string  — sector de la barrera (ILIKE en barrier_v2.sector)
  filters.radicado:  string  — número radicado del oficio
}

PASO 1 — Resetear página
  currentPage = 0

PASO 2 — Debounce de 350 ms (evita consultas excesivas al escribir)

PASO 3 — Construir URL con filtros y consultar BD

GET /api/v1/entity-letters?{params}

Mapeo filtro UI → query param API:
  filters.estado    → state
  filters.identidad → identidad
  filters.entidad   → entidad
  filters.radicado  → numeroRadicado

Los filtros se combinan con AND en la consulta SQL.
Se mantienen al cambiar de página.

PASO 4 — Actualizar tabla y paginador con response.total
```
