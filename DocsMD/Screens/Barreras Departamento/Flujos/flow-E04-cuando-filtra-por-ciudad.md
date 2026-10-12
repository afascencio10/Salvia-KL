# flow-E04 — Cuando filtra por ciudad

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por ciudad
   Tipo: User Interaction
   Funciones: onFilterSelect() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Catálogo (E01): GET /api/v1/locations/cities?department_id=:AssignedDepartmentID
  → todas las ciudades del departamento, ordenadas por nombre

PASO 1 — Sin debounce (select). value = city_id o vacío (Todas)
PASO 2 — currentPage = 0
PASO 3 — GET …?cityId=  (igualdad exacta b.city_id)
PASO 4 — Actualizar items / total
```
