━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma la reasignación
   Tipo: User Interaction / Backend write
   Código: E-10
   Prerequisito: M-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: botón "Reasignar" del modal, habilitado solo con un agente elegido


INPUT: {
  ids:     selectedIds          → uuids de emergency_measure_referral
  agentId: reasignarAgenteId    → general_user_i_code
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Guardar

  POST /api/v1/emergency-measure-referrals/reasignar
  Content-Type: application/json

  {
    "ids": ["uuid", "uuid"],
    "agentId": "general_user_i_code"
  }

PASO 2 — Si responde 200

  selectedIds = []
  cerrar el modal (E-09, ya sin selección)
  repetir E-01 con los filtros de UI que sigan activos

  La fila queda con el agente nuevo y el badge "Asignada".
  La card Asignadas del supervisor sube porque el alcance se recalcula.
  Si el estado anterior era por_asignar, esa card baja.


PASO 3 — Si responde 4xx o 5xx

  El modal sigue abierto.
  selectedIds no se vacía.
  Se muestra el mensaje de error. El listado no se recarga.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Para cada id de ids, en una transacción:

    UPDATE salvia.emergency_measure_referral
       SET assigned_agent_id = {agentId},
           status            = 'asignada',
           updated_at        = now()
     WHERE id = {id}
       AND deleted_at IS NULL

  Vale para cualquier estado previo. La barra del supervisor puede incluir
  filas que no estaban en por_asignar; todas pasan a asignada.
  No modifica la revisión técnica ni las medidas hijas.

  Rechaza con 400 si ids está vacío, si agentId está vacío o si agentId
  no pertenece al catálogo de GET .../agentes.
  Rechaza con 404 si algún id no existe.

  Respuesta 200:

  { "updated": [ "uuid", "uuid" ] }

  → FIN EJECUCIÓN ✓
