# flow-E02 — Cuando filtra por documento

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por documento
   Tipo: User Interaction
   Funciones: onFilterDocument() · loadEntityCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filters.document: texto del input “Documento de identidad”  → v-model del usuario
  filters.city:     filtro de ciudad vigente                  → estado Vue
  filters.entityId: entidad activa del selector               → estado Vue (E07)
}

PASO 1 — Debounce de escritura (~350 ms)

PASO 2 — currentPage = 0

PASO 3 — SI !filters.entityId → TERMINAR (loadEntityCases no llama API)

PASO 4 — GET /api/v1/entity-cases
  query: {
    entityId, document (trim, ILIKE parcial), city, page=0, pageSize=5
  }

PASO 5 — Vue actualiza items / totalItems

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                      | Paso afectado |
|----------------------------------------------------------|---------------|
| ¿Normalizar documento (quitar puntos/espacios/guiones)?  | PASO 4        |
```
