━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre el modal de reasignación
   Tipo: User Interaction
   Código: E-08
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por la barra "↺ Reasignar N seleccionada(s)" (E-07).

Exige puedeReasignar (roles sm y sv, permiso reassign_remisiones_medidas).

"Asignar agente" no abre este modal. Busca la tarea ToDo asignar_agente_medidas
de esa remisión y abre el modal de la tarea. Si no hay tarea pendiente, muestra
"Esta remisión no tiene una asignación pendiente."


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Qué queda seleccionado

  Desde la barra:
    → selectedIds se queda como está
    → si está vacío, la barra no se pinta y este camino no corre

PASO 2 — Abrir

  reasignarAgenteId = ""
  reasignarAbierto  = true

  Título: "Reasignar remisiones"
  Texto: "{N} remisión(es) seleccionada(s). Elige el agente destino."
  El botón "Reasignar" permanece deshabilitado mientras el select esté vacío.

PASO 3 — Catálogo

  Si E-04 todavía no cargó la lista:

  GET /api/v1/emergency-measure-referrals/agentes

  Misma respuesta que E-04: { items: [{ id, nombre }] }.
  El catálogo es role_code = am (M-01, paso 5).

  Este evento no escribe assigned_agent_id ni status.

  → FIN EJECUCIÓN ✓
