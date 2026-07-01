━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Código: E-01
   Prerequisito: M-01 (schema dupla + psychosocial_support)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  mostrarCards:   mostrar cards resumen   → prop :mostrarCards (default false)
  defaultFilter:  filtro inicial fijo     → prop :defaultFilter (opcional)
                  shape: { agent_id?: string, dupla_id?: string }
  reasignacion:   modo reasignación       → prop :reasignacion (default false)
  pageSize:       registros/página          → prop :pageSize (default 20)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — remisiones-psicosocial-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar estado interno

  scopeFilter           = { ...defaultFilter }   // copia inmutable del prop
  activeFilters         = { ...defaultFilter }   // incluye agent_id o dupla_id si el padre los envió
  searchNumeroIdentidad = ""
  searchTelefono        = ""
  autocompleteText      = { profesional_asignada: "" }
  autocompleteSelected  = { profesional_asignada: null }
  autocompleteSuggestions = { profesional_asignada: [] }
  duplaOptions          = []
  equipoRemitenteOptions = []
  stats                 = { total: 0, pendienteAsignacion: 0, enProceso: 0, cerradas: 0 }
  sortBy                = "created_at"
  sortOrder             = "desc"
  currentPage           = 1
  pageSize              = prop pageSize || 20
  totalRemisiones       = 0
  remisiones            = []
  loading               = true
  loadError             = null
  selectedRemisiones    = []

  MAX_SESSIONS          = 6

  // Pre-selección UI según defaultFilter
  SI defaultFilter.agent_id:
    → activeFilters['agent_id'] = defaultFilter.agent_id
    → (opcional) precargar autocompleteSelected con nombre del agente

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

  Ejemplo pantalla "Mis remisiones" por psicólogo:
    &filter_agent_id={defaultFilter.agent_id}

  Ejemplo pantalla por dupla:
    &filter_dupla_id={defaultFilter.dupla_id}

  → remisiones, totalRemisiones


PASO 2B — Catálogo de duplas

  GET /api/v1/duplas
  → duplaOptions


PASO 2C — Equipos remitentes

  GET /api/v1/psychosocial-support/equipos-remitentes
  → equipoRemitenteOptions


PASO 2D — Stats para cards (solo si mostrarCards)

  GET /api/v1/psychosocial-support/stats
    (+ mismos filter_* que activeFilters en PASO 2A)

  200 {
    total: 8,
    pendienteAsignacion: 2,
    enProceso: 3,
    cerradas: 3
  }
  → stats


PASO 3 — Calcular derivados y finalizar carga

  totalPages = Math.ceil(totalRemisiones / pageSize)
  loading    = false


PASO 4 — Renderizar UI

  SI mostrarCards === true:
    → SummaryCardsRow con stats (encima de FilterBar)

  FilterBar → filtros E-03 … E-13

  Tabla → bloques por fila:

    Bloque REMISIÓN — submitted_by (solo lectura):
      → submittedByName = submitter_names + submitter_last_names
      → submitterTeam   = general_user_team del JOIN en submitted_by
      → SI submitted_by NULL → "—" en ambos campos

    (resto igual: badge/t tags "Por consultar", sesiones, dupla, etc.)

  PaginationBar → E-06

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/psychosocial-support/stats
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Mismos JOINs y filtros WHERE que el listado (PASO 5), sin paginación:

```sql
SELECT
  COUNT(*)::int AS total,
  COUNT(*) FILTER (
    WHERE ps.dupla_id IS NULL AND ps.agent_id IS NULL
  )::int AS pendiente_asignacion,
  COUNT(*) FILTER (
    WHERE ps.status = 'en_gestion'
  )::int AS en_proceso,
  COUNT(*) FILTER (
    WHERE ps.status = 'cerrado'
  )::int AS cerradas
FROM ... -- misma query base E-01
WHERE ps.deleted_at IS NULL
  AND ... -- mismos filter_* activos
```

Response:

```json
{
  "total": 8,
  "pendienteAsignacion": 2,
  "enProceso": 3,
  "cerradas": 3
}
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Filtros de alcance (defaultFilter)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Param | Condición |
|---|---|
| `filter_agent_id` | `ps.agent_id = {value}` |
| `filter_dupla_id` | `ps.dupla_id = {value}` |

Combinables con el resto de filtros UI (AND).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/psychosocial-support/list
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

(Query base sin cambios — ver versión anterior de este flujo)

`submitted_by` en SELECT/JOIN:
- Leer `ps.submitted_by` tal cual está en BD
- JOIN `submitter_gu` / `submitter_gup` para nombre y team
- Sin lógica de escritura en este endpoint


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  NOTAS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Tema | Decisión |
|---|---|
| Cards | Solo si `mostrarCards === true`; se recargan en cada `fetchRemisiones()` |
| defaultFilter | Persiste al limpiar filtros (E-05) |
| submitted_by | Solo lectura + JOIN; fuera de alcance cómo se guarda al crear |
| fetchStats | Mismos filtros que listado; no incluye paginación |
