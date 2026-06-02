━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por equipo
   Tipo: User Interaction
   Código: E-11  (antes parte de E-02)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio de valor en el DropdownFilter **"Por equipo"**
               (`filter.type === 'dropdown'`, `filter.key === 'equipo'`)
               Handler: `setDropdownFilter('equipo', value)`

> **No disponible en Mis casos:** si la prop `agentId` está definida, el dropdown
> no se muestra y el filtro no se envía al backend.

INPUT: {
  filterValue:   código de equipo seleccionado en el `<select>`
                 valor vacío "" → desactivar filtro equipo
                 opciones hardcodeadas en el padre (por ahora):
                   'Riesgo bajo' | 'Riesgo alto' | 'Hombres'
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si el filtro ya está activo

  SI agentId está definido:
    → No hacer nada (filtro equipo no permitido en Mis casos)
    → TERMINAR ejecución

  SI activeDropdown.key === 'equipo' Y activeDropdown.value === filterValue:
    → No hacer nada
    → TERMINAR ejecución

  SI filterValue === "":
    → activeDropdown = null (si key era 'equipo')
    → CONTINÚA PASO 2

  SI NO:
    → activeDropdown = { key: 'equipo', value: filterValue }
    → CONTINÚA PASO 2


PASO 2 — Resetear búsqueda y paginación

  searchText   = ""
  currentPage  = 1
  loading      = true
  loadError    = null
  cases        = []


PASO 3 — Consultar backend

  GET /api/v1/cases/list
    ?dropdown_filter_key=equipo
    &dropdown_filter_value={filterValue}   // Riesgo bajo | Riesgo alto | Hombres
    &chip_filter=casos_nuevos              // si chip activo (E-09)
    &dropdown_filter_key=riesgo            // si riesgo activo (E-10) — otro dropdown
    &search=...                            // si búsqueda activa (E-03)
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  // NO se envía filter_key=persona_asignada junto con equipo
  // En Mis casos (agentId) el dropdown equipo no aparece


PASO 4 — Actualizar vista

  filteredCases = cases
  Columna `team` → case.caseTeam (victim_case_team)
  → Renderizar tabla o EmptyState

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (dropdown_filter_key = equipo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

**Fuente de dato:** `salvia.victim_case.victim_case_team`  
**Modelo GORM:** `VictimCase.VictimCaseTeam` → columna `victim_case_team`

> Filtra por el equipo **del caso**, no por `general_user_team` del agente.

PASO 5 — Agregar cláusula WHERE

  CASO dropdown_filter_key = 'equipo' (filter_value presente)
  Y NO hay scope agentId (filter_key ≠ persona_asignada):

    WHERE vc.victim_case_team = {filter_value}

  // Coincidencia exacta: 'Riesgo bajo', 'Riesgo alto', 'Hombres'
  // Casos con victim_case_team NULL o vacío quedan excluidos


PASO 6 — Opciones del dropdown

  Hardcodeadas en `get_victim_cases.html`:

    { value: 'Riesgo bajo', label: 'Riesgo bajo' }
    { value: 'Riesgo alto', label: 'Riesgo alto' }
    { value: 'Hombres',     label: 'Hombres'     }


PASO 7 — Columna `team` en la tabla

  `caseTeam` ← `COALESCE(vc.victim_case_team, '')` en la respuesta JSON.
  UI: `${ caseObj.caseTeam || '—' }`


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas (E-11 implementado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Códigos de equipo | `Riesgo bajo`, `Riesgo alto`, `Hombres` (hardcode en padre) |
| Opciones del dropdown | Fijas en el padre; sin endpoint auxiliar por ahora |
| Columna `team` | Bind a `caseTeam` / `victim_case_team` |
| Combinación con `agentId` | **No permitida** — dropdown oculto en Mis casos; backend ignora equipo si hay scope agente |
| Combinación con chip / riesgo / search | Sí, vía parámetros aditivos |
