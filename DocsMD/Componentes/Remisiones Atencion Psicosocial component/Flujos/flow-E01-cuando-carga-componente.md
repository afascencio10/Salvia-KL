━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Código: E-01
   Prerequisito: M-01 + M-02 (schema dupla, psychosocial_support, team_contact)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  mostrarCards:   mostrar cards resumen   → prop :mostrarCards (default false)
  defaultFilter:  filtro inicial fijo     → prop :defaultFilter (opcional)
                  shape: { professional_id?: string, dupla_id?: string }
  reasignacion:   modo reasignación       → prop :reasignacion (default false)
  pageSize:       registros/página          → prop :pageSize (default 20)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — remisiones-psicosocial-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar estado interno

  scopeFilter           = { ...defaultFilter }
  activeFilters         = { ...defaultFilter }
  searchNumeroIdentidad = ""
  searchTelefono        = ""
  autocompleteText      = { profesional_asignada: "" }
  autocompleteSelected  = { profesional_asignada: null }
  autocompleteSuggestions = { profesional_asignada: [] }
  duplaOptions          = []
  equipoRemitenteOptions = []
  stats                 = { total: 0, abierto: 0, enGestion: 0, enDevolucion: 0, cerrado: 0 }
  sortBy                = "created_at"
  sortOrder             = "desc"
  currentPage           = 1
  pageSize              = prop pageSize || 20
  totalRemisiones       = 0
  remisiones            = []
  loading               = true
  loadError             = null
  selectedRemisiones    = []

  SESSION_DOTS          = 6
  SESSION_LABEL         = "Sesiones (4 - 6)"   // quemado en UI

  // Pre-selección UI según defaultFilter
  SI defaultFilter.professional_id:
    → activeFilters['professional_id'] = defaultFilter.professional_id
    → (opcional) precargar autocompleteSelected con nombre del profesional

  SI defaultFilter.dupla_id:
    → activeFilters['dupla_id'] = defaultFilter.dupla_id
    → pre-seleccionar dropdown dupla_asignada cuando duplaOptions cargue


PASO 2 — Cargar datos en paralelo

  requests = [
    fetchRemisiones(),           // PASO 2A — siempre
    fetchDuplas(),               // PASO 2B — siempre
    fetchEquiposRemitentes()     // PASO 2C — siempre
  ]

  SI mostrarCards === true:
    → requests.push(fetchStats())   // PASO 2D

  Promise.all(requests)

  SI alguna promesa falla:
    → loading = false; loadError = mensaje; TERMINAR


PASO 2A — Listado paginado

  GET /api/v1/psychosocial-support/list
    ?page={currentPage}
    &page_size={pageSize}
    &sort=created_at
    &order=desc
    (+ todos los activeFilters como filter_*)

  Ejemplo pantalla "Mis remisiones" por profesional:
    &filter_professional_id={defaultFilter.professional_id}

  Ejemplo pantalla por dupla:
    &filter_dupla_id={defaultFilter.dupla_id}

  Cada ítem incluye:
    → submittedByName     (JOIN submitted_by)
    → submittedByTeam     (columna ps.submitted_by_team)
    → professionalId      (columna ps.professional_id)
    → sessionCount        (subquery team_contact — ver PASO 4)
    → dupla, status, datos de caso, etc.

  → remisiones, totalRemisiones


PASO 2B — Catálogo de duplas

  GET /api/v1/duplas
  → duplaOptions


PASO 2C — Equipos remitentes

  GET /api/v1/psychosocial-support/equipos-remitentes
  → DISTINCT submitted_by_team (sin JOIN a general_user)
  → equipoRemitenteOptions


PASO 2D — Stats para cards (solo si mostrarCards)

  GET /api/v1/psychosocial-support/stats
    (+ mismos filter_* que activeFilters en PASO 2A)

  200 {
    total: 8,
    abierto: 2,
    enGestion: 3,
    enDevolucion: 0,
    cerrado: 3
  }
  → stats


PASO 3 — Calcular derivados y finalizar carga

  totalPages = Math.ceil(totalRemisiones / pageSize)
  loading    = false


PASO 4 — Renderizar UI

  SI mostrarCards === true:
    → SummaryCardsRow con 5 cards:
        • Total remisiones → stats.total
        • Abiertos       → stats.abierto
        • En gestión     → stats.enGestion
        • En devolución  → stats.enDevolucion
        • Cerrados       → stats.cerrado

  FilterBar → filtros E-03 … E-13

  Tabla → bloques por fila:

    Bloque CASO — sin cambios (víctima, riesgo, links)

    Bloque REMISIÓN — simplificado (sin badge ni tags):
      → submittedByName  = JOIN submitted_by → general_user_profile
      → submitterTeam    = ps.submitted_by_team (columna directa)
      → createdAt
      → "Ver remisión →"
      → NO renderizar ReferralTypeBadge ni ExtraTags

    Bloque ESTADO Y ASIGNACIÓN:
      → status badge (abierto | en_gestion | en_devolucion | cerrado)
      → label fijo: "Sesiones (4 - 6)"
      → 6 puntos (.rps-dot); pintar los primeros N donde:
           N = sessionCount de la remisión
           sessionCount = COUNT(team_contact)
             WHERE psicosocial_id = ps.id
               AND is_psico_session = true
               AND is_completed = true
               AND deleted_at IS NULL
      → dupla / professionalId asignación (sin cambio de lógica de visualización)

  PaginationBar → E-06

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/psychosocial-support/stats
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Mismos JOINs y filtros WHERE que el listado, sin paginación:

```sql
SELECT
  COUNT(*)::int                                            AS total,
  COUNT(*) FILTER (WHERE ps.status = 'abierto')::int       AS abierto,
  COUNT(*) FILTER (WHERE ps.status = 'en_gestion')::int    AS en_gestion,
  COUNT(*) FILTER (WHERE ps.status = 'en_devolucion')::int AS en_devolucion,
  COUNT(*) FILTER (WHERE ps.status = 'cerrado')::int       AS cerrado
FROM ... -- misma query base del listado
WHERE ps.deleted_at IS NULL
  AND ... -- mismos filter_* activos
```

Response:

```json
{
  "total": 8,
  "abierto": 2,
  "enGestion": 3,
  "enDevolucion": 0,
  "cerrado": 3
}
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Filtros de alcance (defaultFilter)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Param | Condición |
|---|---|
| `filter_professional_id` | `ps.professional_id = {value}` |
| `filter_dupla_id` | `ps.dupla_id = {value}` |

Combinables con el resto de filtros UI (AND).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/psychosocial-support/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SELECT adicional / subquery por fila:

```sql
(
  SELECT COUNT(*)::int
  FROM salvia.team_contact tc
  WHERE tc.psicosocial_id = ps.id
    AND tc.is_psico_session = true
    AND tc.is_completed = true
    AND tc.deleted_at IS NULL
) AS session_count
```

Campos de remitente:
- `ps.submitted_by` → JOIN nombre
- `ps.submitted_by_team` → columna directa (sin JOIN para team)

Asignación:
- `ps.professional_id` (reemplaza agent_id)
- `ps.dupla_id` → JOIN dupla para nombres psicóloga / trab. social


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  NOTAS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Tema | Decisión |
|---|---|
| Cards | 5 cards: Total + 4 status de PsychosocialSupportStatus; solo si `mostrarCards === true` |
| defaultFilter | Persiste al limpiar filtros (E-05); usa `professional_id` |
| submitted_by_team | Columna denormalizada; evita JOIN para equipo remitente |
| Sesiones UI | Label quemado "Sesiones (4 - 6)"; puntos = team_contact completadas |
| REMISIÓN | Sin badge tipo ni tags — eliminados post-reunión |
