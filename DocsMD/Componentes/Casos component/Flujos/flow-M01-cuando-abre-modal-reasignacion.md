━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre el modal de reasignación
   Tipo: Lifecycle / User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: el padre llama `this.$refs.reasignarModal.open(cases)`
               tras recibir @reasignar-casos de casos-component (E-15)

INPUT: {
  cases:   casos seleccionados   → payload.cases de E-15 (Array<CaseListItem>)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reasignar-casos-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar entrada

  SI cases vacío o no es array:
    → No abrir modal
    → TERMINAR ejecución


PASO 2 — Inicializar estado del modal

  visible            = true
  cases              = copia superficial de cases
  selectedAgentIcode = ''
  selectedAgent      = null
  agents             = []
  agentsError        = null
  loadingAgents      = false
  saving             = false


PASO 3 — Resolver equipo (resolvedTeam)

  // ── CONTINGENCIA (REASSIGN_CROSS_TEAM_CONTINGENCY === true) ──
  SI crossTeamContingency:
    → resolvedTeam = ''   // o calcular solo para banner informativo (opcional)
    → NO validar caseTeam / riesgo
    → fetchAllAgents()    // → M-02-C
    → FIN EJECUCIÓN ✓

  // ── MODO NORMAL (lógica original — comentar en código al activar contingencia) ──
  referenceCase = cases[0]

  SI referenceCase.caseTeam no está vacío:
    → resolvedTeam = referenceCase.caseTeam.trim()

  SI NO (caseTeam vacío):
    → risk = referenceCase.riskStatus (normalizado a minúsculas)
    → SI risk ∈ ['bajo', 'moderado', 'medio']:
        resolvedTeam = 'Riesgo bajo'
    → SI risk ∈ ['alto', 'extremo']:
        resolvedTeam = 'Riesgo alto'
    → SI NO se puede derivar:
        agentsError = 'No se pudo determinar el equipo del caso seleccionado.'
        → Mostrar modal con error; NO llamar M-02
        → TERMINAR ejecución


PASO 4 — Disparar carga de agentes (modo normal)

  → fetchAgentsByTeam(resolvedTeam)   // → M-02

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE (list_cases.html)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  <reasignar-casos-modal ref="reasignarModal" @closed="onReasignarModalClosed" />

  onReasignarCasos({ cases }) {
    this.$refs.reasignarModal.open(cases);
  }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                           | Paso afectado |
|---------------------------------------------------------------|---------------|
| ¿Bloquear foco en el modal (trap focus) para accesibilidad?   | Interfaz      |
| ¿Mostrar i_code de cada caso en la lista del modal?           | Interfaz      |
| Contingencia cross-team                                       | [contingencia-reasignacion-cross-team.md](../contingencia-reasignacion-cross-team.md) |
