━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga agentes por equipo
   Tipo: Lifecycle / Backend read
   Modo: NORMAL (comentar al activar contingencia — ver M-02-C)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

> **Contingencia:** cuando `REASSIGN_CROSS_TEAM_CONTINGENCY === true`, este flujo
> no se ejecuta. Usar [flow-M02-C-cuando-carga-todos-los-agentes.md](./flow-M02-C-cuando-carga-todos-los-agentes.md).
> Plan: [contingencia-reasignacion-cross-team.md](../contingencia-reasignacion-cross-team.md).

Disparado por: M-01 al abrir el modal, una vez resuelto resolvedTeam
               (solo si contingencia desactivada)

INPUT: {
  team:   equipo resuelto   → resolvedTeam (string)
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


PASO 2 — Consultar backend

  GET /api/v1/equipo-operadores?team={team}&role=ro

  // role=ro por defecto; para habilitar 'op' en el futuro: &role=ro&role=op
  // Endpoint: case_detail_controller.GetOperadoresByTeam
  // Normaliza variantes "Riesgo alto" / "Riesgo Alto" vía teamVariantsForQuery


PASO 3 — Procesar respuesta

  SI respuesta no ok:
    → loadingAgents = false
    → agentsError   = 'Error al cargar los agentes del equipo'
    → TERMINAR ejecución

  SI respuesta ok:
    → agents = data (array)
    → Ordenar alfabéticamente por fullName (si el backend no lo hace)
    → loadingAgents = false

  SI agents.length === 0:
    → agentsError = 'No hay agentes disponibles para el equipo "' + team + '"'
    → El <select> queda deshabilitado; botón Reasignar deshabilitado

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Consulta de agentes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Tablas:
  security.general_user          (u / gu)
  security.general_user_profile  (gup / p)
  security.rel_role_general_user (rr)
  security.role                  (r)

Criterios:
  r.role_code              = 'ro'
  gu.general_user_status   = 'e'        // activo
  gu.general_user_team     = {team}     // equipo resuelto del caso

Campos retornados al frontend:
  icode     ← gu.general_user_i_code
  fullName  ← gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names
  team      ← gu.general_user_team

SQL de referencia (case_detail_controller.GetOperadoresByTeam):

  SELECT DISTINCT gu.general_user_i_code,
         gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names AS full_name,
         gu.general_user_team,
         gup.general_user_profile_names,
         gup.general_user_profile_last_names
  FROM security.general_user gu
  JOIN security.general_user_profile gup
    ON gup.general_user_profile_id = gu.general_user_general_user_profile
  JOIN security.rel_role_general_user rr ON rr.general_user_id = gu.general_user_id
  JOIN security.role r ON r.role_id = rr.role_id
  WHERE r.role_code = 'ro'
    AND gu.general_user_status = 'e'
    AND gu.general_user_team = ?
  ORDER BY full_name ASC


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  REFERENCIA GeneralUserDAO
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

GeneralUserDAO.go define la entidad GeneralUser con:
  - GeneralUserTeam        → columna general_user_team
  - GeneralUserProfile     → relación con general_user_profile (nombres del agente)

La consulta de agentes por equipo puede implementarse:
  (A) Reutilizando GET /api/v1/equipo-operadores (recomendado — ya existe)
  (B) Nuevo método en GeneralUserDAO / AgentLightRepository.FindAllByRoleAndTeam
      expuesto como endpoint dedicado


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                           | Paso afectado |
|---------------------------------------------------------------|---------------|
| ¿Incluir también rol 'op' además de 'ro'?                     | Implementado — `AGENT_ROLE_CODES` / `?role=` |
| ¿Excluir al agente actualmente asignado del listado?           | No — incluido en listado |
| ¿Normalizar "Riesgo Alto" vs "Riesgo alto" en la consulta?    | Implementado — `teamVariantsForQuery` |
