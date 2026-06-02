━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando filtra por nivel de riesgo
   Tipo: User Interaction
   Código: E-10  (antes parte de E-02)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio de valor en el DropdownFilter **"Nivel de riesgo"**
               (`filter.type === 'dropdown'`, `filter.key === 'riesgo'`)
               Handler: `setDropdownFilter('riesgo', value)`

INPUT: {
  filterValue:   valor seleccionado en el `<select>`
                 opciones UI: 'bajo' | 'moderado' | 'alto' | 'extremo'
                 valor vacío "" → resetear al defaultFilter
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si el filtro ya está activo

  SI activeDropdown.key === 'riesgo' Y activeDropdown.value === filterValue:
    → No hacer nada
    → TERMINAR ejecución

  SI filterValue === "":
    → activeDropdown = null (solo si key === 'riesgo')
    → CONTINÚA PASO 2

  SI NO:
    → activeDropdown = { key: 'riesgo', value: filterValue }
    → CONTINÚA PASO 2

  // activeDropdown es independiente de agentId y activeChipKey


PASO 2 — Resetear búsqueda y paginación

  searchText   = ""
  currentPage  = 1
  loading      = true
  loadError    = null
  cases        = []


PASO 3 — Consultar backend

  GET /api/v1/cases/list
    ?dropdown_filter_key=riesgo
    &dropdown_filter_value={filterValue}   // bajo | moderado | alto | extremo
    &filter_key=persona_asignada           // si prop agentId (Mis casos)
    &filter_value={agentId}
    &chip_filter=casos_nuevos              // si chip activo (E-09)
    &search=...                            // si búsqueda activa (E-03)
    &sort={sortBy}
    &order={sortOrder}
    &page=1
    &page_size={pageSize}

  // dropdown_filter_* es aditivo: combina con agentId, chip, search y paginación

  SI respuesta no ok:
    → loadError = data.error || 'Error al aplicar el filtro'
    → TERMINAR ejecución

  SI respuesta ok:
    → cases, totalCases, loading = false
    → CONTINÚA PASO 4


PASO 4 — Actualizar vista

  filteredCases = cases
  → Renderizar tabla o EmptyState

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/cases/list  (filter_key = riesgo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

**Fuente de dato:** `salvia.victim_case_form2.victim_case_form2_risk_level`  
**DAO:** `VictimCaseForm2DAO` — campo `VictimCaseForm2RiskLevel`  
**JOIN requerido (ya en query base E-01):**

```sql
LEFT JOIN salvia.victim_case_form2 vf2
       ON vf2.victim_case_form2_victim_case = vc.victim_case_id
```

**Escala numérica (HU-027):**

| `filter_value` (UI) | `victim_case_form2_risk_level` | Etiqueta   |
|-----------------------|--------------------------------|------------|
| `bajo`                | 1                              | Bajo       |
| `moderado`            | 2                              | Moderado   |
| `alto`                | 3                              | Alto       |
| `extremo`             | 4                              | Extremo    |

PASO 5 — Agregar cláusula WHERE

  CASO dropdown_filter_key = 'riesgo' (o filter_key legacy = 'riesgo'):

    WHERE vf2.victim_case_form2_risk_level = {nivel_entero}

  // El backend traduce filter_value → entero (1-4) antes de ejecutar la query
  // Casos sin form2 quedan excluidos al filtrar por un nivel concreto


PASO 6 — Respuesta

  `riskStatus` en cada ítem: slug bajo/moderado/alto/extremo; sin form2 → `desconocido`.
  Badge UI: etiqueta "Desconocido".


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  Decisiones aplicadas (E-10 implementado)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Decisión | Resolución |
|----------|------------|
| Opciones del dropdown | Fijas en el padre (bajo/moderado/alto/extremo, escala 1-4) |
| Combinación con otros filtros | `dropdown_filter_key/value` aditivo con agentId, chip, search |
| Casos sin form2 | `riskStatus = 'desconocido'`, badge "Desconocido" en tabla |
