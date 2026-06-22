━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cambia de tab o de página
   Tipo: User Interaction
   Código: E-03
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  TAB — goToTab(tabName)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Validar tab permitido

  PERMITIDOS:
    recontact_v → currentFilter = "fcv"  (status v)
    recontact_i → currentFilter = "fci"  (status i)

  COMENTADOS (no invocables desde UI):
    routing, routedToApprove, expired, completed, issues

PASO 2 — Resetear paginación

  currentTab  = tabName
  currentFilter = fcv | fci
  numPages    = 1
  goToPage(0)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  PAGINACIÓN — goToPage(page)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 3 — Validar rango

  SI page < 0 || page >= numPages → TERMINAR

PASO 4 — Actualizar estado

  currentPage = page
  victimContacts = []

PASO 5 — Consultar backend

  url = buildContactListUrl()
  // Incluye currentFilter, currentPage y query params de filtro activos

  getEntity(url, callback)

  // REPORTES-ONLY: rama victimCases comentada
  /*
  if (currentFilter != "fcv" && currentFilter != "fci") {
    getEntity(salviaVictimCaseGETFormPath + ...);
  }
  */

PASO 6 — Actualizar vista

  scrollTo(0, 0)
  Renderizar tabla o empty state

  → FIN EJECUCIÓN ✓
