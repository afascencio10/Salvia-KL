━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cancela o cierra el modal
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  (A) clic en CancelBtn ("Cancelar")
  (B) clic en CloseBtn ("✕") del encabezado
  (C) clic en ModalBackdrop (.rcm-backdrop) — click fuera del diálogo

INPUT: ninguno


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reasignar-casos-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar que no hay guardado en curso

  SI saving === true:
    → No cerrar (evitar interrumpir M-05 futuro)
    → TERMINAR ejecución


PASO 2 — Cerrar y limpiar estado

  visible            = false
  cases              = []
  resolvedTeam       = ''
  agents             = []
  selectedAgentIcode = ''
  selectedAgent      = null
  loadingAgents      = false
  agentsError        = null
  saving             = false


PASO 3 — Emitir evento hacia el padre

  this.$emit('closed', {})

  // El padre NO recarga la tabla en este flujo.
  // La selección en casos-component permanece intacta.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  @closed="onReasignarModalClosed"

  onReasignarModalClosed() {
    // Opcional: no hacer nada
    // La tabla conserva los checkboxes marcados
  }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                           | Paso afectado |
|---------------------------------------------------------------|---------------|
| ¿Limpiar selectedCases en casos-component al cancelar?        | UX / padre    |
| ¿Confirmar cierre si ya hay agente seleccionado?              | UX            |
