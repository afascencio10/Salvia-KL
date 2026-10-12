# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Función: mounted() · loadCities() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentUserId: CommonSession.UserICode  (el API mine usa la sesión)
  assignedDepartmentId: CommonSession.AssignedDepartmentID  (solo para el select de ciudades)
}

PASO 1 — Mostrar #app

PASO 2 — Cargar ciudades del departamento
  SI assignedDepartmentId vacío → cities = []
  SI NO → GET /api/v1/locations/cities?department_id=:assignedDepartmentId
          → cities = [{ value, label }]

PASO 3 — GET /api/v1/barriers-v2/mine
  Auth: cookie de sesión
  Backend:
    SI no autenticado → 401
    SI rol ≠ en → 403
    SI AssignedDepartmentID vacío → 403
    WHERE deleted_at IS NULL
      AND user_territory_link_id = :UserICode
      AND status IN ('OPEN', 'En Gestion')
    LEFT JOIN victim_case, entity_branch, security.city
    LEFT JOIN LATERAL case_task ToDo articulacion_territorial del usuario
    ORDER BY created_at DESC

PASO 4 — Mapear filas (nombre, doc, teléfono, ciudad, taskId vacío si no hay tarea)

SI 401 → landing
SI !ok → mensaje de error
SI 200 → barreras = mapped items
```
