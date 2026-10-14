# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · reloadScreen() · loadCities() · loadBarreras()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:  CommonSession.CurrentRole
  (departamento NO viaja al cliente; el API usa AssignedDepartmentID)
}

PASO 1 — document.title + mostrar #app

PASO 2 — Verificar rol
SI currentRole !== 'en':
  → loadError = 'Rol no autorizado…'
  → TERMINAR

PASO 3 — Inicializar filtros vacíos, currentPage = 0, pageSize = 10
  Cargar ciudades: GET /api/v1/locations/cities?department_id=assignedDepartmentId

PASO 4 — SUB-FLUJO: Cargar listado paginado

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Cargar listado paginado     │
│  (E01–E06)                              │
└─────────────────────────────────────────┘

  S1. isLoading = true; loadError = null

  S2. GET /api/v1/barriers-v2/department
      query: page, pageSize=10,
             docNumber?, entidad?, cityId?, status?
      Auth: sesión cookie
      Backend:
        SI no autenticado → 401
        SI rol ≠ en → 403
        SI AssignedDepartmentID vacío → 403
        WHERE deleted_at IS NULL
          AND department_id = :AssignedDepartmentID
          AND [filtros opcionales]
        ORDER BY created_at DESC
        LIMIT 10 OFFSET page*10

  S3. SI 401 → landing
      SI 403 sin departamento → mensaje de asignar departamento
      SI !ok → 'No se pudo cargar la lista de barreras. Intenta de nuevo.'
      SI 200 → items, total, isLoading = false
```
