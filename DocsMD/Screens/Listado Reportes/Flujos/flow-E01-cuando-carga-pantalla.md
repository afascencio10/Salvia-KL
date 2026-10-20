━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: navegación del usuario a GET /salvia/reportes

Precondición: sesión activa con rol 'op' o 'sv'

INPUT: {
  sessionID:  cookie de sesión  → sessions.Default(ctx).Get("userData")
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — ListReportsGET (salvia/facades/VictimContactFacade.go)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Validar sesión y rol

  obtener s desde utils.GetCommonSession(sessionID)
  CheckPermission(PermissionsByRole, "get_list_reports", s.CurrentRole)
  SI s.CurrentRole NOT IN {'op', 'sv'}:
    → HTTP 401 (CheckPermission responde 401, no 403 — igual que el resto de pantallas protegidas por Menu.go)
    → TERMINAR ejecución


PASO 2 — Renderizar template

  RenderTemplate(c, ..., "reportes/", "list_reportes", ..., {
    "currentRole": s.CurrentRole,   // el template deriva canReasignarReportes = (currentRole === 'sv')
    "currentUser", "menu", "locale", "lang", "windowTitle": "Reportes"
  })

  // El HTML no trae reportes — los carga casos-component en su propio E-01
  // contra GET /api/v1/victim-contacts/list (ver GAP-1/2/3 en related-tables.md)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — list_reportes.html (mount)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 3 — Montar casos-component con configuración de reportes

  <casos-component
    ref="casosComponent"
    list-endpoint="/api/v1/victim-contacts/list"
    :default-filter="{}"                       // sin scope: ve todos los reportes
    :available-filters="reportesAvailableFilters"
    :columns="reportesColumns"
    :sort-options="reportesSortOptions"        // solo Fecha de Reporte ASC/DESC
    item-label-plural="reportes"               // textos del contador, del botón y del vacío
    :buttons="[{ id: 'ver', label: 'Ver Reporte' }]"
    :reasignacion="canReasignarReportes"       // currentRole === 'sv'
    @action-clicked="onReporteActionClicked"
    @reasignar-casos="onReasignarReportes"
  />

  // Dispara el E-01 propio de casos-component (ya documentado) contra el
  // endpoint nuevo. Esta pantalla no repite esa lógica.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — resueltos tras la implementación
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                  | Resolución |
|-----------------------------------------------------------------------|------------|
| Backend de GET /api/v1/victim-contacts/list (filtros, columnas)       | Implementado — ver flow-E02, E03, E07, E08 y related-tables.md |
| ¿`op` ve el botón "Reasignar Reportes" oculto o el prop reasignacion directamente en false? | `reasignacion` va en false por rol (`canReasignarReportes`): para `op` no existen ni las casillas ni el botón, y `reasignar-casos-modal` ni se monta |
