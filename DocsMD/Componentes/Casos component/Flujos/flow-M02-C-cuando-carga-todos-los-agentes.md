━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga todos los agentes (contingencia cross-team)
   Tipo: Lifecycle / Backend read
   ID alternativo: M-02-C
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: M-01 al abrir el modal, cuando
               REASSIGN_CROSS_TEAM_CONTINGENCY === true

Precondición: contingencia activa (ver
               contingencia-reasignacion-cross-team.md)

INPUT: {
  (ninguno — no se filtra por equipo del caso)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reasignar-casos-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Preparar estado de carga

  loadingAgents = true
  agentsError   = null
  agents          = []
  selectedAgentIcode = ''
  selectedAgent      = null


PASO 2 — Consultar backend (sin filtro team)

  GET /api/v1/equipo-operadores?role=ro

  // Sin query param team → GetOperadoresByTeam omite filtro general_user_team
  // role=ro vía AGENT_ROLE_CODES (mismo criterio que M-02 original)

  // LÓGICA ORIGINAL COMENTADA (M-02):
  // GET /api/v1/equipo-operadores?team={resolvedTeam}&role=ro


PASO 3 — Procesar respuesta

  SI respuesta no ok:
    → loadingAgents = false
    → agentsError   = 'Error al cargar los agentes'
    → TERMINAR ejecución

  SI respuesta ok:
    → agents = data (array)
    → Ordenar alfabéticamente por fullName
    → loadingAgents = false

  SI agents.length === 0:
    → agentsError = 'No hay agentes disponibles'
    → El <select> queda deshabilitado; botón Reasignar deshabilitado

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Consulta de agentes (sin team)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Endpoint: GET /api/v1/equipo-operadores
Handler:  case_detail_controller.GetOperadoresByTeam

Query:
  role=ro   (default si no se envía)

Cuando team está vacío, SQL efectivo:

  WHERE r.role_code IN ('ro')
    AND gu.general_user_status = 'e'
  -- SIN: AND gu.general_user_team IN (...)

Campos retornados: icode, fullName, team


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  UI RECOMENDADA
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Mostrar agent.team en cada <option> del select para distinguir equipos.

Ocultar o adaptar TeamBanner (.rcm-team-banner) — ver interfaz del modal.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  REVERSIÓN
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

REASSIGN_CROSS_TEAM_CONTINGENCY = false
→ Restaurar M-02 (fetchAgentsByTeam con resolvedTeam)
→ Comentar o no invocar fetchAllAgents / M-02-C
