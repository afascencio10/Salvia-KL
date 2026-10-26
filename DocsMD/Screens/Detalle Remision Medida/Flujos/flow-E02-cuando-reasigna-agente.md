━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando reasigna el agente de medidas
   Tipo: User Interaction / Backend write
   Código: E-02
   Reutiliza: E-10 del componente (POST reasignar)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Pantalla: Detalle de Remisión — Medida de Emergencia
Visible solo si currentRole es sm o sv.


INPUT: {
  remisionId: uuid de la remisión abierta
  agentId:    general_user_i_code elegido en el select
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Abrir

  El botón "Reasignar agente de medidas" muestra el select y
  "Confirmar" / "Cancelar".
  El select no incluye al agente que ya está asignado.

  GET /api/v1/emergency-measure-referrals/agentes
  El catálogo es el de la bandeja: usuarios con role_code am.
  Ese GET exige get_remisiones_medidas. sm y sv lo tienen.

PASO 2 — Cancelar

  Oculta el select. No llama al API.

PASO 3 — Confirmar

  "Confirmar" queda deshabilitado si no hay agente.
  POST /api/v1/emergency-measure-referrals/reasignar

  {
    "ids": ["{remisionId}"],
    "agentId": "{agentId}"
  }

  Es el mismo cuerpo de E-10. Una sola id.

PASO 4 — Si responde 200

  Cierra el formulario y repite la lectura de E-01.
  La cabecera y el campo "Agente asignado" muestran al nuevo agente.
  E-10 deja la remisión en status asignada. El badge pasa a "Asignada".

PASO 5 — Si responde 4xx o 5xx

  El formulario sigue abierto.
  Se muestra el mensaje de error. El detalle no se recarga.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  El POST no es nuevo. Lo documenta E-10 del componente.

  Exige reassign_remisiones_medidas (sm y sv).
  Rechaza 400 si agentId no es un agente am activo, o si el body es inválido.
  Rechaza 404 si el id no existe.
  Actualiza assigned_agent_id y pone status = asignada.
  Vale para cualquier estado previo, igual que en la bandeja.

  em, sc, cr y am no ven el bloque. Si llaman al POST, el permiso responde 401.

  → FIN EJECUCIÓN ✓
