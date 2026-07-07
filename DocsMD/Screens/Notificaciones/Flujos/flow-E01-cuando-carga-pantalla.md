# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · loadOficios() · buildOficiosUrl()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:   rol del usuario en sesión    → window.NotifConfig.currentRole (inyectado por Go)
  currentUserId: ID del usuario en sesión     → window.NotifConfig.currentUserId
}

PASO 1 — Mostrar el contenedor principal de la app
  document.getElementById('app').style.display = 'block'

PASO 2 — Verificar que el rol tiene acceso a la pantalla

SI currentRole NO está en ['op', 'an', 'ro']:
  → Asignar loadError = 'Rol no autorizado para acceder a esta pantalla.'
  → La pantalla muestra el bloque de error (v-if loadError)
  → TERMINAR ejecución

SI currentRole es válido:
  → CONTINÚA FLUJO GENERAL

PASO 3 — Activar estado de carga
  isLoading = true
  loadError = null

PASO 4 — Construir URL del API con paginación y filtros (buildOficiosUrl)

Parámetros base:
  page  = currentPage (0-based)
  limit = itemsPerPage (5)

Según rol:
  op / ro → agentId={currentUserId}
  an      → notificationUserId={currentUserId}

Tab activo:
  currentTab === 'gestionar' → manageableOnly=true

Filtros opcionales (solo si tienen valor):
  state          ← filters.estado
  identidad      ← filters.identidad
  entidad        ← filters.entidad
  numeroRadicado ← filters.radicado

PASO 5 — Consultar oficios del usuario (paginado en BD)

GET /api/v1/entity-letters?{params}

→ resultado: { items, total, page, pageSize, pendingCount } o error

PASO 6 — Manejar respuesta del API

SI status === 200:
  → Extraer items de response.items
  → Mapear cada item con mapApiToOficio()
  → Asignar totalItems = response.total
  → Asignar pendingCount = response.pendingCount (badge del tab "Oficios por gestionar")
  → isLoading = false
  → Vue renderiza la tabla con la página actual

SI status === 401:
  → Redirigir a /static/landing.html
  → TERMINAR ejecución

SI otro status:
  → loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.'
  → isLoading = false
  → Vue muestra el bloque de error

PASO 7 — Carga en background de locaciones (loadLocations)
  GET /api/v1/locations/departments
  GET /api/v1/locations/cities
  (para modales de proyectar; no bloquea la tabla)
```
