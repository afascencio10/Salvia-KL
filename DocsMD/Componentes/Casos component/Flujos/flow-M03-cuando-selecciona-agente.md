━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona un agente
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en el <select> de "Nueva persona asignada"
               (.rcm-select)

Precondición: agents cargados (M-02 completado sin error)
               loadingAgents === false

INPUT: {
  agentIcode:   valor del select   → selectedAgentIcode (string)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reasignar-casos-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Leer valor del select

  agentIcode = event.target.value   // o v-model selectedAgentIcode


PASO 2 — Actualizar estado

  SI agentIcode === '':
    → selectedAgentIcode = ''
    → selectedAgent      = null
    → Botón "Reasignar" queda :disabled

  SI agentIcode !== '':
    → selectedAgentIcode = agentIcode
    → selectedAgent      = agents.find(a => a.icode === agentIcode) || null
    → Botón "Reasignar" queda habilitado
      (la acción de guardado se implementará en M-05)


PASO 3 — Validación opcional en UI

  SI selectedAgent es null tras el find:
    → selectedAgentIcode = ''
    → agentsError = 'El agente seleccionado no es válido'
    → TERMINAR ejecución

  // No se emite evento hacia el padre en este flujo.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ESTADO RESULTANTE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Campo | Valor tras selección |
|---|---|
| `selectedAgentIcode` | `icode` del agente elegido |
| `selectedAgent` | `{ icode, fullName, team }` |
| Botón Reasignar | Habilitado si `selectedAgentIcode` no vacío |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                           | Paso afectado |
|---------------------------------------------------------------|---------------|
| ¿Mostrar preview del agente seleccionado bajo el select?      | Interfaz      |
| ¿Advertir si el agente elegido ya es el asignado en algún caso? | UX / M-05  |
