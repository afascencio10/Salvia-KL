━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla de Registro
   Tipo: Lifecycle
   Función: mounted() → set_victim_case_v3.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  icode:   código de un Caso existente, opcional   → query param de la URL (ej. link desde otra pantalla)
  agentId: id del Agente Integral autenticado        → sesión
}

PASO 1 — Inicializar estado de la pantalla
  stage = 'habeasData'
  registro = null
  caso = null
  moduloActual = null

SI icode viene en la URL:
  → PASO 2a — Cargar el Caso existente
    VictimCaseService.getByICode({ icode })
    → Si existe: registro = perfil del caso, caso = datos del caso, stage = 'formulario', moduloActual = último módulo con progreso
    → Si no existe o no autorizado: mostrar error "Caso no encontrado" y mantener stage = 'habeasData'
  → CONTINÚA FLUJO GENERAL

SI NO viene icode:
  → stage se mantiene en 'habeasData'
  → CONTINÚA FLUJO GENERAL

PASO 3 — Renderizar
  Angular/Vue re-evalúa stage
  → Si stage === 'habeasData': muestra HabeasDataGate
  → Si stage === 'formulario': muestra FormularioIngresoRegistro directamente en moduloActual

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                    | Paso afectado |
|----------------------------------------------------------|---------------|
| Quién puede generar un link directo con `icode` (¿solo Expediente de víctima, que está diferido?) | PASO 2a |
