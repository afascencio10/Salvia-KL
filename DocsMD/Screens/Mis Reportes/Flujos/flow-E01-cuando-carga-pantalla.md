━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: navegación del usuario a GET /salvia/mis-reportes

Precondición: sesión activa con rol 'op'

INPUT: {
  sessionID:        cookie de sesión  → sessions.Default(ctx).Get("userData")
  currentUserId:    s.UserICode       → obtenido de la sesión, no de la URL (se ignora cualquier ?filter_value= / ?agent=)
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — MyReportsGET (salvia/facades/VictimContactFacade.go)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar sesión y rol

  obtener s desde utils.GetCommonSession(sessionID)
  CheckPermission(PermissionsByRole, "get_my_reports", s.CurrentRole)
  SI s.CurrentRole != 'op':
    → HTTP 401 (CheckPermission responde 401, no 403; también para sv)
    → TERMINAR ejecución


PASO 2 — Renderizar template

  RenderTemplate(c, ..., "reportes/", "get_my_reportes", ..., {
    "currentUserId": s.UserICode,   // se pasa como :agent-id del componente
    "currentUser", "currentRole", "menu", "locale", "lang", "windowTitle": "Mis Reportes"
  })


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — get_my_reportes.html (mount)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 3 — Montar casos-component con scope fijo

  <casos-component
    ref="casosComponent"
    list-endpoint="/api/v1/victim-contacts/list"
    :agent-id="currentUserId"                  // fuerza filter_key=persona_asignada&filter_value={agentId}
    :default-filter="{}"
    :columns="reportesColumns"                 // sin 'assigned_agent' (redundante)
    :available-filters="reportesAvailableFilters"   // estado_reporte + busqueda, sin 'asignacion'
    :sort-options="reportesSortOptions"
    item-label-plural="reportes"
    :buttons="[{ id: 'ver', label: 'Ver Reporte' }]"
    :reasignacion="false"
    @action-clicked="onReporteActionClicked"
  />

  // Dispara el E-01 de casos-component contra
  // /api/v1/victim-contacts/list?filter_key=persona_asignada&filter_value={icode}

  // Columna «Fecha Reporte» (REP-12): cada fila trae creationDateCO ("AAAA-MM-DD", día de creación en hora de
  // Bogotá) y la celda lo muestra como DD/MM/AAAA sin new Date(): no cambia de día por la zona del navegador.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                  | Resolución |
|-----------------------------------------------------------------------|------------|
| ¿`persona_asignada` en el backend de reportes se compara contra `victim_contact_agent_id`? | Sí: `WHERE victim_contact_agent_id = filter_value` (caso `MR-E01-C7`) |
| ¿El scope se puede quitar desde la UI o la URL? | No: sale de la sesión (`MR-E01-C4`) y sobrevive a filtros, búsqueda y orden (`MR-E01-C5`, `MR-E01-C6`) |
| La columna «Fecha Reporte» salía «—» (H11) porque el listado no traía `creationDateCO`, y debía verse el mismo día de creación en Bogotá | Resuelto 2026-11-02 (WP-2 · REP-12): el endpoint trae `creationDateCO`; casos `MR-E01-C10` y `MR-E01-C11` |
