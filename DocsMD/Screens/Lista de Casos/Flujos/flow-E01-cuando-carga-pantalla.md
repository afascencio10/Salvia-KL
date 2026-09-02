━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  currentRole:            rol de sesión
  assignedDepartmentId:   CommonSession.AssignedDepartmentID (puede ser "")
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — ListCasesGET
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Permiso HTML

  CheckPermission(..., "get_list_cases", s.CurrentRole)
  SI denegado → 401 y TERMINAR

PASO 2 — Renderizar list_cases.html

  currentRole, currentUser, assignedDepartmentId, menu


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — list_cases.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 3 — Resolver variante por rol

  SI currentRole === 'en' AND assignedDepartmentId vacío:
    → enMissingDepartment = true
    → Mostrar alert de departamento faltante
    → NO montar casos-component
    → TERMINAR

  SI currentRole === 'en' AND assignedDepartmentId presente:
    → Mostrar subtítulo "Casos del departamento asignado al usuario"
    → CONTINÚA PASO 4

  SI otro rol:
    → CONTINÚA PASO 4


PASO 4 — Montar casos-component

  → E-01 del componente: GET /api/v1/cases/list
  → Si sesión en, el API aplica AssignedDepartmentID de sesión
     (victim_case_victim_town_code → town → city.department_id)
  → Si API 403, casos-component muestra loadError

  → FIN EJECUCIÓN ✓
