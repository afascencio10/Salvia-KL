# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Funciones: mounted() · loadScreenData()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:   rol del usuario en sesión     → config inyectada por Go / sesión
  currentUserId: general_user_i_code sesión    → config / sesión
}


PASO 1 — Mostrar contenedor principal
  document.getElementById('app').style.display = 'block'   // o equivalente Vue

PASO 2 — Verificar acceso por rol

SI currentRole !== 'sv':
  → loadError = 'Rol no autorizado para administrar duplas.'
  → Mostrar bloque de error
  → TERMINAR ejecución

SI currentRole === 'sv':
  → CONTINÚA FLUJO GENERAL

PASO 3 — Activar estado de carga
  isLoading = true
  loadError = null

PASO 4 — Consultar datos iniciales en paralelo (o secuencial)

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Cargar profesionales activos │
└─────────────────────────────────────────┘

  GET /api/v1/duplas/profesionales
  // mismos criterios que RRM-03 profesionales (usuarios activos ps/ts)

  Consulta referencia:
  ```sql
  SELECT
    gu.general_user_i_code AS icode,
    TRIM(CONCAT(gup.general_user_profile_names, ' ',
                gup.general_user_profile_last_names)) AS full_name,
    r.role_code AS role
  FROM security.general_user gu
  JOIN security.general_user_profile gup
    ON gup.general_user_i_code = gu.general_user_i_code
  JOIN security.rel_role_general_user rrgu
    ON rrgu.general_user_i_code = gu.general_user_i_code
  JOIN security.role r
    ON r.role_i_code = rrgu.role_i_code
  WHERE gu.general_user_status = 'e'
    AND r.role_code IN ('ps', 'ts')
  ORDER BY r.role_code ASC, full_name ASC
  ```

  → Separar en:
      psychologists[]  = role === 'ps'
      socialWorkers[]  = role === 'ts'

  → FIN SUB-FLUJO

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Cargar duplas activas        │
└─────────────────────────────────────────┘

  GET /api/v1/duplas/activas
  // lista enriquecida (no confundir con GET /api/v1/duplas del catálogo id/name)

  Consulta referencia:
  ```sql
  SELECT
    d.id,
    d.name,
    d.psychologist_id,
    d.social_worker_id,
    TRIM(CONCAT(ps_gup.general_user_profile_names, ' ',
                ps_gup.general_user_profile_last_names)) AS psychologist_name,
    TRIM(CONCAT(ts_gup.general_user_profile_names, ' ',
                ts_gup.general_user_profile_last_names)) AS social_worker_name
  FROM salvia.dupla d
  JOIN security.general_user_profile ps_gup
    ON ps_gup.general_user_i_code = d.psychologist_id
  JOIN security.general_user_profile ts_gup
    ON ts_gup.general_user_i_code = d.social_worker_id
  WHERE d.deleted_at IS NULL
  ORDER BY d.name ASC
  ```

  → duplas = response.items (o array)

  → FIN SUB-FLUJO

PASO 5 — Manejar errores de carga

SI alguna petición falla (401):
  → Redirigir a login / landing
  → TERMINAR

SI alguna petición falla (otro status):
  → loadError = 'No se pudo cargar la información de duplas. Intenta de nuevo.'
  → isLoading = false
  → Mostrar bloque de error
  → TERMINAR

SI ambas ok:
  → CONTINÚA FLUJO GENERAL

PASO 6 — Enriquecer profesionales con estado de dupla

  Construir mapa:
    assignedByUserId = {}
    PARA CADA d EN duplas:
      assignedByUserId[d.psychologistId] = { duplaId: d.id, duplaName: d.name }
      assignedByUserId[d.socialWorkerId] = { duplaId: d.id, duplaName: d.name }

  PARA CADA profesional EN psychologists ∪ socialWorkers:
    SI assignedByUserId[profesional.icode] existe:
      → profesional.enDupla = true
      → profesional.duplaId = ...
      → profesional.duplaName = ...
    SI NO:
      → profesional.enDupla = false
      → profesional.duplaId = null
      → profesional.duplaName = null

PASO 7 — Finalizar render

  isLoading = false
  → Vue renderiza:
      • Card Psicólogas (badge En dupla / Disponible)
      • Card Trabajadoras Sociales
      • Card Duplas activas (N) con filas Editar / Eliminar

→ FIN EJECUCIÓN ✓
```

---

## Shape de respuesta sugerida

### `GET /api/v1/duplas/profesionales`

```json
{
  "psychologists": [
    { "icode": "user-...", "fullName": "Alejandra Mora" }
  ],
  "socialWorkers": [
    { "icode": "user-...", "fullName": "Valentina Ospina" }
  ]
}
```

### `GET /api/v1/duplas/activas`

```json
{
  "items": [
    {
      "id": "uuid",
      "name": "Dupla 1",
      "psychologistId": "user-...",
      "psychologistName": "Alejandra Mora",
      "socialWorkerId": "user-...",
      "socialWorkerName": "Valentina Ospina"
    }
  ]
}
```

> **Nota:** `GET /api/v1/duplas` (sin sufijo) sigue siendo el catálogo `{ duplas: [{id,name}] }` del listado de remisiones. Esta pantalla usa `/duplas/activas`.

---

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión | Paso afectado |
|---|---|
| ¿Mostrar profesionales inactivos que aún figuren en una dupla? (hoy: lista de profesionales solo activos; duplas muestran nombres vía JOIN de perfil) | PASO 4 / 6 |
