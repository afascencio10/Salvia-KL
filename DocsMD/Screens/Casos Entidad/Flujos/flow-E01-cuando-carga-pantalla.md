# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · reloadScreen() · loadCitiesAndCases() · loadEntityCases()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:  rol del usuario en sesión   → window.__CasosEntidadConfig.currentRole
  currentUser:  nombre del usuario          → window.__CasosEntidadConfig.currentUser
}

PASO 1 — Mostrar el contenedor principal
  document.getElementById('app').style.display = 'block'
  document.title = windowTitle

PASO 2 — Verificar acceso por rol

SI currentRole !== 'et':
  → loadError = 'Rol no autorizado para ver Casos Entidad.'
  → isLoading = false
  → Vue muestra el bloque de error
  → TERMINAR ejecución

SI currentRole === 'et':
  → CONTINÚA FLUJO GENERAL

PASO 3 — Inicializar estado de pantalla
  isLoading = true
  loadError = null
  filters.document = ''
  filters.city = ''
  filters.entityId = ''
  currentPage = 0
  items = []
  totalItems = 0
  sectorName = ''
  entities = []
  cities = []
  pageSize = 5

PASO 4 — Cargar catálogo completo de entidades

  GET /api/v1/entities
  Auth: sesión + permiso get_casos_entidad

  → resultado: [{ id, icode, name, sector, sectorName }, ...] o error

SI status === 401:
  → Redirigir a /static/landing.html
  → TERMINAR

SI status !== 200:
  → loadError = 'No se pudo cargar el catálogo de entidades.'
  → isLoading = false
  → TERMINAR

SI status === 200:
  → entities = response
  → CONTINÚA FLUJO GENERAL

PASO 5 — Preseleccionar primera entidad (temporal)

SI entities.length === 0:
  → filters.entityId = ''
  → isLoading = false
  → Empty: sin entidades en catálogo
  → TERMINAR

SI entities.length > 0:
  → filters.entityId = String(entities[0].id)
  → sectorName = entities[0].sectorName
  → CONTINÚA FLUJO GENERAL

  // Temporal: relación usuario et ↔ entidad aún no definida.
  // Cuando exista, reemplazar este paso por la entidad de la sesión.

PASO 6 — Cargar ciudades de la entidad + listado de casos
  → loadCitiesAndCases()

  S6.1 GET /api/v1/entities/{entityId}/cities
       → cities = [{ id, name }, ...]  (DISTINCT city desde entity_branch de esa entidad)

  S6.2 Ejecutar SUB-FLUJO: Cargar listado paginado

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Cargar listado paginado     │
│  (usado por E02, E03, E04, E07)         │
└─────────────────────────────────────────┘

  S1. Validar entidad activa
      SI !filters.entityId:
        → items = []; totalItems = 0; isLoading = false
        → FIN SUB-FLUJO

  S2. isLoading = true; loadError = null

  S3. GET /api/v1/entity-cases
      query: {
        entityId:  filters.entityId,     // required
        document:  filters.document,     // opcional — ILIKE parcial
        city:      filters.city,         // opcional — ciudad del CASO (exacta desde select)
        page:      currentPage,          // 0-based
        pageSize:  5
      }

      Backend:
        entity_case
          JOIN entity_branch → entity_id = :entityId
          JOIN victim_case (documento + ciudad del caso vía town→city)
          LEFT JOIN victim_case_form2 (riesgo)
        WHERE deleted_at IS NULL
          AND document ILIKE '%doc%' si aplica
          AND case_city.city_name ILIKE :city si aplica
        ORDER BY entity_case.updated_at DESC
        LIMIT 5 OFFSET page*5

  S4. Respuesta 200:
      → items = response.items
      → totalItems = response.total
      → sectorName = response.sectorName || meta de entidad seleccionada
      → isLoading = false
      → Vue renderiza filas + paginación

     401 → landing.html
     otro → loadError

  → FIN SUB-FLUJO

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| Relación usuario et ↔ entidad (hoy: primera del catálogo)        | PASO 5        |
| Ruta detalle caso / entidad para E05–E06                         | (otros flujos)|
```
