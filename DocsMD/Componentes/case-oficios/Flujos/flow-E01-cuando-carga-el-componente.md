━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga el componente
   Tipo: Lifecycle
   Función: mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  caseId:  ID del caso (victim_case)   → prop recibido del padre
}


PASO 1 — Inicializar estado

  cargando         = true
  error            = null
  oficios          = []
  buscador         = ''
  temaActivo       = null   // 'barrera' | null (los otros 3 temas deshabilitados — GAP)
  estadoActivo     = null   // null = "Todos"
  oficioSeleccionado = null


PASO 2 — Consultar los oficios del caso

GET /api/v1/entity-letters?caseId={caseId}

→ resultado: array de EntityLetter (sin relaciones de caso/barrera enriquecidas)


SI status === 401:
  → Redirigir a /static/landing.html
  → TERMINAR ejecución

SI error o status !== 200:
  → error    = "No se pudieron cargar los oficios de este caso."
  → cargando = false
  → TERMINAR ejecución  // el componente renderiza ErrorMsg + BtnReintentar

SI ok:
  → oficios  = resultado
  → cargando = false
  → CONTINÚA FLUJO GENERAL


PASO 3 — Vue re-evalúa reactivamente

  → oficiosFiltrados (computed) se recalcula sobre `oficios` aplicando
    buscador / temaActivo / estadoActivo — todos en null/'' al inicio,
    así que la primera renderización muestra todos los oficios sin filtrar

  SI oficiosFiltrados.length === 0:
    → Renderiza EmptyState "No hay oficios que coincidan con los filtros."
      (en la carga inicial, esto significa que el caso no tiene oficios)

  SI NO:
    → Renderiza una OficioCard por cada elemento

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                    | Paso afectado |
|--------------------------------------------------------------------------|---------------|
| Los chips "Medidas de Emergencia" / "Apoyo Psicosocial" /                | PASO 1 / 3    |
| "Estabilización Económica" no tienen columna en entity_letter todavía   |               |
| — requiere migración de schema (ver case-oficios-interface.md).        |               |
| ¿El buscador debe tener debounce? Con fetch único y filtro client-side  | PASO 3        |
| probablemente no haga falta, pero depende del volumen real de oficios  |               |
| por caso.                                                               |               |
