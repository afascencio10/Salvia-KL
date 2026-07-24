# flow-E07 — Cuando cambia entidad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia entidad
   Tipo: User Interaction
   Funciones: onChangeEntity() · loadCitiesAndCases() · loadEntityCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  filters.entityId: valor del <select>   → v-model + @change
  entities:         catálogo E01
  filters.document: se CONSERVA
  filters.city:     se limpia (las opciones de ciudad son de la entidad nueva)
}

PASO 1 — filters.entityId ya actualizado por v-model; currentPage = 0
          applySelectedEntityMeta() → sectorName

PASO 2 — SI !entityId:
  → limpiar items/cities; empty “Selecciona una entidad…”
  → TERMINAR

PASO 3 — filters.city = '' (ciudades del select son de la entidad anterior)

PASO 4 — GET /api/v1/entities/{entityId}/cities → cities
         GET /api/v1/entity-cases?entityId&document&page=0&pageSize=5

PASO 5 — Vue renderiza listado

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| Relación usuario et ↔ entidad (hoy: catálogo completo + 1ª en E01) | E01 / E07   |
```
