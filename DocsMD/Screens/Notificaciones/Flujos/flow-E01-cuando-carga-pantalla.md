# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · loadOficios()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:   rol del usuario en sesión    → window.NotifConfig.currentRole (inyectado por Go)
  currentUserId: ID del usuario en sesión     → window.NotifConfig.currentUserId
}

PASO 1 — Mostrar el contenedor principal de la app
  document.getElementById('app').style.display = 'block'

PASO 2 — Verificar que el rol tiene acceso a la pantalla

SI currentRole NO está en ['op', 'an']:
  → Asignar loadError = 'Rol no autorizado para acceder a esta pantalla.'
  → La pantalla muestra el bloque de error (v-if loadError)
  → TERMINAR ejecución

SI currentRole es válido:
  → CONTINÚA FLUJO GENERAL

PASO 3 — Activar estado de carga
  isLoading = true
  loadError = null

PASO 4 — Construir URL del API según rol

SI currentRole === 'op':
  → url = '/api/v1/entity-letters?limit=100&page=0&agentId={currentUserId}'

SI currentRole === 'an':
  → url = '/api/v1/entity-letters?limit=100&page=0&notificationUserId={currentUserId}'

PASO 5 — Consultar oficios del usuario

GET {url}

→ resultado: array de entity_letter o error

PASO 6 — Manejar respuesta del API

SI status === 200:
  → Extraer items: response si es array, o response.items si es objeto
  → Ordenar items por createdAt DESC
  → Mapear cada item con mapApiToOficio():
      - Construir nombre completo de la víctima (victimName + victimLastName)
      - Extraer sectores de la barrera
      - Determinar letterPriority (alta | normal)
      - Calcular canManage según el estado del oficio y el rol:
          op  → canManage si estado es 'por_proyectar' o 'en_correccion'
          an  → canManage si estado es 'para_revisar', 'aprobacion_juridica',
                'para_radicar' o 'radicado'
  → Asignar resultado a this.oficios
  → isLoading = false
  → Vue renderiza la tabla con los oficios cargados

SI status === 401:
  → Redirigir a /static/landing.html
  → TERMINAR ejecución

SI otro status:
  → loadError = 'No se pudo cargar la lista de oficios. Intenta de nuevo.'
  → isLoading = false
  → Vue muestra el bloque de error

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                        | Paso afectado |
|------------------------------------------------------------|---------------|
| ¿Qué campos exactos devuelve el endpoint GET?              | PASO 5        |
| ¿El endpoint soporta paginación server-side? limit=100     | PASO 4        |
| puede quedarse corto si hay más de 100 oficios por usuario |               |
```
