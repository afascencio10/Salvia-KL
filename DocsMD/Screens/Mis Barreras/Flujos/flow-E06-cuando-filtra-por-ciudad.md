# flow-E06 — Cuando filtra por ciudad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por ciudad
   Tipo: User Interaction
   Función: watcher filters (deep) · filters.cityId
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: filters.cityId  (value del select; vacío = Todas)

PASO 1 — currentPage = 0

PASO 2 — filteredBarreras en cliente:
  SI cityId vacío → pasa
  SI NO → String(b.cityId) === String(cityId)

Las opciones del select vienen de GET /api/v1/locations/cities?department_id=AssignedDepartmentID (E01).
No hay fetch adicional al cambiar el select.
```
