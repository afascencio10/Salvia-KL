# flow-E02 — Cuando filtra por documento

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por documento
   Tipo: User Interaction
   Funciones: onFilterDocument() · loadEntityCases()
   Estado: planeado (ajustar: ya no depende de entityId del selector)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filters.document: texto del input
  entityBranchId:   sede de sesión (implícita en API)
}

PASO 1 — Debounce ~350 ms

PASO 2 — currentPage = 0

PASO 3 — GET /api/v1/entity-cases
  query: { document: trim, page: 0, pageSize: 5 }
  Backend acota siempre a session.EntityBranchId + ILIKE parcial en documento

PASO 4 — Vue actualiza items / totalItems

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                      | Paso afectado |
|----------------------------------------------------------|---------------|
| ¿Normalizar documento (puntos/espacios/guiones)?         | PASO 3        |
```
