# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · reloadScreen() · loadCases()
   Estado: implementado
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:     rol en sesión   → __CasosMecanismoConfig / CommonSession
  currentUserId:   icode sesión
}

PASO 1 — Mostrar app + document.title

PASO 2 — Verificar rol

SI currentRole !== 'mc':
  → loadError = 'Rol no autorizado para ver Casos Mecanismo.'
  → TERMINAR

SI currentRole === 'mc':
  → CONTINÚA

PASO 3 — Inicializar filtros
  filters.document = ''
  filters.city = ''
  currentPage = 0
  pageSize = 5

PASO 4 — Cargar listado
  → SUB-FLUJO: Cargar listado paginado

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Cargar listado paginado     │
│  (E01, E02, E03, E04)                   │
└─────────────────────────────────────────┘

  S1. isLoading = true

  S2. GET /api/v1/mecanismos/mine
      Auth: sesión + get_casos_mecanismo (rol mc)
      query: {
        document:  filters.document,   // opcional
        city:      filters.city,       // opcional
        page:      currentPage,
        pageSize:  5
      }
      // El usuario NO viene del cliente: el backend usa session.UserICode

      Backend:
        WHERE cam.deleted_at IS NULL
          AND cam.status = 'active'
          AND cam.assigned_user_id = :sessionUserICode
          AND [document ILIKE '%…%' si aplica]
          AND [city ILIKE '%…%' si aplica]
        ORDER BY cam.created_at DESC
        LIMIT 5 OFFSET page*5

  S3. items = response.items
      totalItems = response.total
      isLoading = false

  → FIN SUB-FLUJO

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno.
```
