━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de módulo en el sidebar
   Tipo: User Interaction
   Función: cambiarModulo(moduloDestino) → set_victim_case_v3.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  moduloDestino: 'CI' | 'DP' | 'CA' | 'HV' | 'AG' | 'TZ' | 'GS'   → click en ModuloItem del sidebar
}

PASO 1 — Verificar si el módulo está bloqueado

SI moduloDestino en ['TZ', 'GS'] Y !triadaCompleta:
  → No navega
  → Muestra tooltip "Completa Víctima + Hechos + Agresor para desbloquear"
  → TERMINAR ejecución

SI NO:
  → CONTINÚA FLUJO GENERAL

PASO 2 — Cambiar de módulo
  moduloActual = moduloDestino
  Vue re-evalúa moduloActual
  → El `<dinamic-form>` renderiza la sección correspondiente a moduloDestino
  → El sidebar resalta el ModuloItem activo

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno — `triadaCompleta` se calcula en el sub-flujo "Recalcular habilitación secuencial" del evento E11, este evento solo lo consulta.
