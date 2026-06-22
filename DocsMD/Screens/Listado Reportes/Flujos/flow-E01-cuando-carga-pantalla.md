━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  rol:     sv | ro   → sesión del usuario
  filter:  (vacío)    → default fcv en frontend
  page:    0
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — VictimCaseFacade.VictimCaseGET
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Resolver template según rol

  case "sv" → tplName = "get_victim_cases_sv"
  case "ro" → tplName = "get_victim_cases_ro"

PASO 2 — Cargar reportes iniciales (sin id en URL)

  switch filter:
    case "", "fcv":
      GetVictimContactsWithoutVictimCase(page, "v", ...)
    case "fci":
      GetVictimContactsWithoutVictimCase(page, "i", ...)

  // Solo contactos sin victim_case asociado + status indicado

PASO 3 — Preparar menuToolsContactTable

  MenuTools[lang][rol]["menu_tool_get_victim_contacts"]
  → sv: tras cambio en Menu.go incluirá Ver + Crear caso
  → ro: ya incluye Ver + Crear caso

PASO 4 — Renderizar HTML con datos embebidos

  victimContacts: JSON.parse({{.victimContacts}}).data
  currentFilter:  "fcv"
  currentPage:      0
  filterNames:      ""      // NUEVO — inicializar vacío
  filterLastNames:  ""      // NUEVO
  filterPhone:      ""      // NUEVO

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  document.title = windowTitle
  Muestra tab Recontacto activo
  Tabla victimContacts con paginación si numPages > 1

  // Bloques de casos comentados — no se renderizan

  → FIN EJECUCIÓN ✓
