━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario selecciona un filtro
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en un chip de filtro o cambio en un dropdown de filtro
               dentro de FilterBar → FilterChips

INPUT: {
  filterKey:    key del filtro seleccionado   → interacción del usuario
                valores: 'casos_nuevos' | 'riesgo' | 'equipo' | 'persona_asignada'
  filterValue:  valor seleccionado            → del dropdown (vacío para chips booleanos)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar si el filtro ya está activo

  SI activeFilter.key === filterKey Y activeFilter.value === filterValue:
    → No hacer nada (mismo filtro seleccionado)
    → TERMINAR ejecución

  SI NO:
    → CONTINÚA PASO 2


PASO 2 — Actualizar estado del filtro activo

  activeFilter  = { key: filterKey, value: filterValue }
  searchText    = ""       // resetear búsqueda al cambiar de filtro
  loading       = true
  loadError     = null
  cases         = []       // limpiar tabla mientras carga


PASO 3 — Consultar backend con el nuevo filtro

  GET /api/v1/cases/list?filter_key={filterKey}&filter_value={filterValue}&sort={sortBy}&order={sortOrder}

  SI respuesta no ok (status != 2xx):
    → loading   = false
    → loadError = data.error || 'Error al aplicar el filtro'
    → Mostrar ErrorState
    → TERMINAR ejecución

  SI respuesta ok:
    → cases   = data.cases
    → loading = false
    → CONTINÚA PASO 4


PASO 4 — Recalcular filteredCases

  filteredCases = cases
  // searchText está vacío (reseteado en PASO 2), no hay filtro de búsqueda que aplicar

  SI filteredCases.length === 0:
    → Mostrar EmptyState "No se encontraron casos"

  SI filteredCases.length > 0:
    → Renderizar tabla con los nuevos casos

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| ¿Se puede tener más de un filtro activo a la vez (multi-filtro)?                 | PASO 1, 2, 3  |
| Al seleccionar 'casos_nuevos' chip, ¿se puede deseleccionar para volver a todos? | PASO 1        |
| ¿El filtro por persona_asignada usa agent_id (icode) o nombre del agente?        | PASO 3        |
| ¿Los dropdowns de riesgo/equipo cargan sus opciones del backend o son hardcoded? | Interfaz      |
