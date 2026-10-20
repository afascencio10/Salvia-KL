━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma la reasignación
   Tipo: User Interaction / Backend write
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en "Reasignar" dentro de reasignar-casos-modal
               (evento M-05 del modal, extendido — ver
               reasignar-casos-modal-events.md)

INPUT: {
  reports:        reportes a reasignar      → cases (del modal)
  agentICode:     agente integral elegido   → selectedAgentIcode
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — POST /api/v1/victim-contacts/reasignar-bulk (nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar sesión, permiso y payload (igual a CasesReassignController)

  obtener sess desde la sesión (supervisorICode, supervisorName)
  SI no hay sesión: → HTTP 401
  SI sess.CurrentRole != 'sv' (permiso reassign_victim_contacts): → HTTP 401   // op, ro y demás roles quedan fuera
  SI body.report_icodes vacío O body.agent_icode vacío (o el JSON no trae los campos):
    → HTTP 400
    → TERMINAR ejecución


PASO 2 — Resolver agente destino

  agente = buscar general_user por agent_icode (mismo FindAgentForReassign)
  SI agente no existe O agente.status != 'e':
    → HTTP 400 "agente no encontrado / inactivo"
    → TERMINAR ejecución


PASO 3 — Por cada report_icode (sin transacción global, igual que CasesReassignService)

  icodes = report_icodes sin duplicados ni vacíos
  PARA CADA icode EN icodes:
    SI victim_contact con ese icode no existe:
      skipped++
      CONTINUAR con el siguiente
    SI NO:
      UPDATE victim_contact SET
        victim_contact_agent_id = agente.icode,
        victim_contact_team     = agente.team        // '' si el agente no tiene equipo
      WHERE victim_contact_i_code = icode
      INSERT salvia.app_log { key: icode, message: 'Reasignación manual de reporte (bulk): agente actualizado',
                              data: { agentICode, agentName, team, actorICode, actorName } }
                           // de mejor esfuerzo: si el log falla solo se avisa en el log del servidor
      updated++

  // Cada UPDATE se aplica por separado: si la BD falla a mitad, los reportes anteriores YA quedaron
  // reasignados y la respuesta es 400 con el mensaje del error (el controller responde 400 ante
  // cualquier error del servicio). Reasignar un reporte ya cerrado se permite: no hay regla que lo impida.

  // A diferencia de UpdateVictimCaseAgent (que solo pisa el equipo si estaba
  // vacío), aquí SIEMPRE se pisa victim_contact_team con el equipo del nuevo
  // agente — un reporte no tiene equipo propio previo que proteger (GAP-3).


PASO 4 — Responder

  { ok: true, updated, skipped }

  → FIN EJECUCIÓN BACKEND ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — list_reportes.html (padre)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 5 — onReasignarCompletado()

  this.$refs.casosComponent.reload()   // limpia selectedCases y vuelve a pedir la página actual
  // (mismo método reload() ya pendiente/documentado en reasignar-casos-modal-interface.md)

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Resolución |
|------------------------------------------------------------------------------------|------------|
| ¿Se registra un evento de timeline/log para reportes, análogo a case_timeline_event? | Sí, sin timeline visual: una fila de `salvia.app_log` por reporte reasignado (`AppLogRepository`, ver PASO 3) |
| ¿Reasignar un reporte YA asignado a otro agente requiere confirmación extra? | No: el supervisor puede reasignar libremente (igual que casos). El modal ya pide una confirmación general ("¿Confirmar reasignación?") |
| ¿Atomicidad de la reasignación masiva? | No es atómica (igual que la de casos): ver el comentario del PASO 3. Anotado como hallazgo en el plan de QA |
