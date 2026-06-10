━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma la reasignación (guardado)
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en ReassignBtn ("Reasignar") en el footer del modal
               (.rcm-btn--primary)

Precondiciones:
  selectedAgentIcode no vacío
  cases.length > 0
  loadingAgents === false
  agentsError === null
  saving === false

INPUT (frontend): {
  cases:            casos del modal          → this.cases (Array<CaseListItem>)
  agentIcode:       agente seleccionado      → selectedAgent.icode
  agentFullName:    nombre del agente        → selectedAgent.fullName (solo UI / descripción)
}

INPUT (backend — sesión): {
  supervisorIcode:  usuario logueado         → session.UserICode
  supervisorName:   nombre del supervisor    → session.Names + " " + session.LastNames
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — reasignar-casos-modal.js
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar estado

  SI !canConfirmReassign:
    → No hacer nada
    → TERMINAR ejecución


PASO 2 — Confirmar con el usuario (opcional)

  → Mostrar diálogo de confirmación (Swal o nativo)
    "¿Confirma reasignar {cases.length} caso(s) a {selectedAgent.fullName}?"
  SI el usuario cancela:
    → TERMINAR ejecución


PASO 3 — Marcar estado de guardado

  saving = true
  Deshabilitar Cancelar, cerrar (✕) y select


PASO 4 — Llamar al backend

  POST /api/v1/cases/reasignar-bulk
  Body: {
    case_icodes:  cases.map(c => c.i_code),
    agent_icode:  selectedAgentIcode
  }

  // ActorID / ActorName se resuelven en el servidor desde la sesión (no enviar desde el cliente)


PASO 5 — Procesar respuesta

  SI respuesta no ok:
    → saving = false
    → Mostrar mensaje de error en el modal (agentsError o saveError)
    → TERMINAR ejecución

  SI respuesta ok:
    → saving = false
    → emit('reassigned', {
        cases:     cases,
        agent:     selectedAgent,
        updated:   result.updated   // cantidad de casos procesados
      })
    → close()   // limpia estado del modal (M-04)
    → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE (list_cases.html)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  <reasignar-casos-modal
    ref="reasignarModal"
    @closed="onReasignarModalClosed"
    @reassigned="onReasignarCompletado"
  />

  onReasignarCompletado(payload) {
    → Recargar tabla: this.$refs.casosComponent.reload()
      (método público a exponer en casos-component — llama fetchCases internamente)
    → La selección de checkboxes se limpia al recargar (E-14)
  }


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/cases/reasignar-bulk
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Archivos sugeridos:
  src/salvia/controller/cases_reassign_controller.go   *(nuevo)*
  src/salvia/service/cases_reassign_service.go         *(nuevo)*
  src/internal/repository/cases_reassign_repository.go *(nuevo, o extender repos existentes)*

Request body:
  {
    "case_icodes": ["uuid-1", "uuid-2"],
    "agent_icode":  "uuid-agente"
  }

Response 200:
  {
    "ok": true,
    "updated": 2,
    "follow_ups_updated": 5
  }

Response 4xx/5xx:
  { "error": "mensaje" }


PASO B1 — Autenticación y validación

  → Leer sesión del supervisor (ActorID, ActorName)
  → Validar case_icodes no vacío y agent_icode presente
  → Verificar que el agente existe y está activo (general_user_status = 'e')
  → Obtener team y nombre del agente desde general_user + general_user_profile


PASO B2 — Transacción por lote

  Iniciar transacción DB
  Para cada caseIcode en case_icodes:

    ── B2.1 — Leer caso actual ──
    Tabla: salvia.victim_case
    Campos: victim_case_i_code, agent_id, victim_case_team (VictimCaseTeam)

    SI el caso no existe → skipped++; continuar con el siguiente (no abortar el lote)

    ── B2.2 — Actualizar victim_case ──
    Tabla: salvia.victim_case
    Modelo: models.VictimCase

    SIEMPRE actualizar:
      agent_id = agent_icode seleccionado

    SI victim_case_team está vacío o NULL:
      victim_case_team = general_user_team del agente seleccionado

    SI victim_case_team ya tiene valor:
      NO modificar victim_case_team (conservar equipo del caso)

    SQL de referencia:
      UPDATE salvia.victim_case
      SET agent_id = $agentIcode,
          victim_case_team = CASE
            WHEN COALESCE(TRIM(victim_case_team), '') = '' THEN $agentTeam
            ELSE victim_case_team
          END
      WHERE victim_case_i_code = $caseIcode

    ── B2.3 — Actualizar seguimientos no realizados ──
    Tabla: salvia.follow_up_v2
    Modelo: models.FollowUpV2

    Criterio: status NOT IN ('REALIZADO', 'CERRADO')
    Incluye: PENDIENTE, VENCIDO, REPROGRAMADO, etc.
    Excluye: REALIZADO y CERRADO

    UPDATE salvia.follow_up_v2
    SET agent_id = $agentIcode,
        team = $agentTeam,
        updated_at = NOW()
    WHERE case_id = $caseIcode
      AND status NOT IN ('REALIZADO', 'CERRADO')
      AND deleted_at IS NULL

    ── B2.4 — Crear evento en timeline (uno por caso) ──
    Tabla: salvia.case_timeline_event
    Modelo: models.CaseTimelineEvent

    Insertar:
      CaseID      = caseIcode
      Category    = models.TimelineCategorySeguimientos     // "Seguimientos"
      Type        = models.TimelineTypeReasignacionCaso     // "Reasignación de Caso"
      EventType   = models.TimelineEventReasignacion          // legacy "REASIGNACION"
      Icon        = models.TimelineIconReasignacion          // "arrows-rotate"
      Color       = models.TimelineColorBlue                 // o TimelineColorOrange
      Date        = time.Now()
      CreatedAt   = time.Now()
      ActorID     = supervisorIcode    (sesión — supervisor que reasigna)
      ActorName   = supervisorName     (sesión)
      EventUserID = supervisorIcode
      Description = "Caso reasignado a {agentFullName} por {supervisorName}"

  Commit transacción
  SI cualquier paso falla → Rollback completo


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  TABLAS Y MODELOS INVOLUCRADOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Tabla | Modelo | Operación |
|---|---|---|
| salvia.victim_case | VictimCase | UPDATE agent_id; UPDATE victim_case_team solo si vacío |
| salvia.follow_up_v2 | FollowUpV2 | UPDATE agent_id WHERE status <> REALIZADO |
| salvia.case_timeline_event | CaseTimelineEvent | INSERT uno por caso |
| security.general_user | — | READ team del agente |
| security.general_user_profile | — | READ nombre del agente |

Nota: la lista de casos (cases_list_repository) obtiene ownerNames desde vc.agent_id
JOIN general_user, por lo que actualizar agent_id en victim_case refleja el cambio
en la tabla sin tocar rel_case_owner_victim_case.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  DIFERENCIAS VS ReasignarCaso EXISTENTE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

El endpoint actual POST /api/v1/casos/:id/reasignar (case_detail_service.reasignarCasoInternal)
tiene comportamiento distinto. M-05 define reglas nuevas para reasignación masiva:

| Aspecto | ReasignarCaso actual | M-05 (nuevo) |
|---|---|---|
| Alcance | Un caso | Varios casos (bulk) |
| victim_case_team | Siempre sobrescribe con team del agente | Solo si el caso no tiene team |
| follow_up_v2 | Solo status = PENDIENTE | status <> REALIZADO |
| Timeline | No crea evento | Crea CaseTimelineEvent por caso |
| Actor timeline | No aplica | Supervisor logueado (sesión) |


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  DIAGRAMA DE FLUJO
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```
[ReassignBtn] → confirmReassign() [M-05]
       │
       ▼
POST /api/v1/cases/reasignar-bulk
       │
       ├── Por cada caseIcode ─────────────────────────────┐
       │    1. UPDATE victim_case (agent_id + team cond.)  │
       │    2. UPDATE follow_up_v2 (status <> REALIZADO)   │
       │    3. INSERT case_timeline_event                  │
       └───────────────────────────────────────────────────┘
       │
       ▼
emit('reassigned') → padre recarga casos-component
```


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                           | Paso afectado |
|---------------------------------------------------------------|---------------|
| follow_up_v2.team al reasignar                                  | Implementado — team del agente |
| victim_case_owner_description                                   | No se actualiza |
| Caso inexistente en el lote                                     | Se omite (skipped) |
| Diálogo de confirmación antes del POST                          | Swal confirm |
| Validación de rol supervisor                                    | No — solo sesión válida |
| reload() en casos-component                                     | Implementado |
