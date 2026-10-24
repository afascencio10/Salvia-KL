━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre el modal de reasignación
   Tipo: User Interaction
   Código: E-08
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  - la barra "↺ Reasignar N seleccionada(s)" (E-07)
  - el botón "Asignar agente" de una fila en estado por_asignar

Las dos exigen puedeReasignar (roles sm y sv, permiso reassign_remisiones_medidas).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Qué queda seleccionado

  Desde la barra:
    → selectedIds se queda como está
    → si está vacío, la barra no se pinta y este camino no corre

  Desde "Asignar agente":
    → selectedIds = [id de esa fila]
    → cualquier otra fila que estuviera marcada deja de estarlo

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
