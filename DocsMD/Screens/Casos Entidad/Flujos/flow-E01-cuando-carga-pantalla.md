# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · reloadScreen() · loadEntityCases()
   Estado: planeado (ajustar código actual)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:       rol en sesión              → __CasosEntidadConfig / CommonSession
  entityBranchId:    sede del usuario et        → general_user.entity_branch_id (sesión)
}

PASO 1 — Mostrar app + document.title

PASO 2 — Verificar rol

SI currentRole !== 'et':
  → loadError = 'Rol no autorizado…'
  → TERMINAR

SI currentRole === 'et':
  → CONTINÚA

PASO 3 — Resolver sede del usuario

SI !entityBranchId (sesión):
  → loadError = 'Tu usuario no tiene una sede (entity_branch) asignada.'
  → TERMINAR

SI entityBranchId presente:
  → CONTINÚA

PASO 4 — Resolver organización padre para el header
  Consultar entity_branch + entity:
    entity_branch WHERE entity_branch_id = :entityBranchId
    JOIN entity ON entity_id
  → entityName  = entity.entity_name
  → sectorName  = label(entity.entity_sector)   // Justicia, Salud, …
  → (opcional) inyectar en config del facade o devolver en meta del listado

PASO 5 — Inicializar filtros
  filters.document = ''
  currentPage = 0
  pageSize = 5
  // SIN filters.city · SIN filters.entityId · SIN catálogo de entities

PASO 6 — Cargar listado
  → SUB-FLUJO: Cargar listado paginado

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Cargar listado paginado     │
│  (E01, E02, E04)                        │
└─────────────────────────────────────────┘

  S1. isLoading = true

  S2. GET /api/v1/entity-cases
      Auth: sesión + get_casos_entidad
      query: {
        document:  filters.document,   // opcional
        page:      currentPage,
        pageSize:  5
      }
      // La sede NO viene del cliente: el backend usa session.EntityBranchId

      Backend:
        WHERE ec.deleted_at IS NULL
          AND ec.entity_branch_id = :sessionEntityBranchId
          AND [document ILIKE '%…%' si aplica]
        ORDER BY ec.updated_at DESC
        LIMIT 5 OFFSET page*5

  S3. items = response.items
      totalItems = response.total
      entityName / sectorName desde response.meta (si vienen)
      isLoading = false

  → FIN SUB-FLUJO

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| Migración + poblar general_user.entity_branch_id                 | PASO 3        |
| Exponer EntityBranchId en CommonSession al login                 | PASO 3        |
| Ajustar API para filtrar por sede de sesión (no entityId query)  | PASO 6        |
```
