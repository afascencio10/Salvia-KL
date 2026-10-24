━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario filtra por agente asignado
   Tipo: User Interaction
   Código: E-04
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: select "Agente asignado"
Solo vista supervisor.


INPUT: {
  agent_id: icode del agente, "sin-asignar", o vacío
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Catálogo

  Al montar el supervisor:

  GET /api/v1/emergency-measure-referrals/agentes

  {
    "items": [
      { "id": "general_user_i_code", "nombre": "Nombres Apellidos" }
    ]
  }

  El rol que alimenta esta lista está pendiente (M-01, paso 5).
  La primera opción del select es "Agente asignado" (vacío).
  La segunda es "Sin asignar" (value sin-asignar).
  Después, un option por item.

PASO 2 — Volver a pedir E-01

  vacío        → no se envía agent_id de filtro
                 (no confundir con el agent_id de alcance de la vista agente)
  sin-asignar  → &agent_id=sin-asignar
  un agente    → &agent_id={icode}

  Cards y alcance no cambian.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  agent_id = sin-asignar
    → AND assigned_agent_id IS NULL

  agent_id = un icode
    → AND assigned_agent_id = {icode}

  Este parámetro, en vista supervisor, es el filtro de UI.
  En vista agente, agent_id es el recorte (E-01) y este select no existe.

  → FIN EJECUCIÓN ✓
