━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Función: mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  barrierICode:  icode de la barrera a mostrar  → inyectado por Go: {{.barrierICode}}
  currentUser:   nombre del usuario en sesión   → inyectado por Go: {{.currentUser}}
  currentRole:   rol del usuario                → inyectado por Go: {{.currentRole}}
  currentUserId: icode del usuario              → inyectado por Go: {{.currentUserId}}
}

PASO 1 — Mostrar el contenedor principal
  document.getElementById('app').style.display = 'block'
  (el div #app se renderiza oculto mientras Go carga el template)

PASO 2 — Renderizar con datos mock
  Vue monta la pantalla con los datos hardcoded en data():
    barrera.victimName    = 'María González'
    barrera.age           = 32
    barrera.location      = 'San Salvador, San Salvador'
    barrera.priority      = 'Alto'
    barrera.sector        = 'Justicia'
    barrera.org           = 'Policía Nacional'
    barrera.description   = 'Negación de recibir denuncia formal'
    barrera.status        = 'Por articular'
    barrera.identifiedAt  = '15 may. 2026'
    barrera.active        = true
    barrera.pendingTasks  = []
    barrera.completedTasks = []

  Tab inicial: activeTab = 'info'
  La pantalla muestra la sección de Información General con los datos mock.
  El tab de Timeline muestra el estado vacío.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                                       | Paso afectado |
|---------------------------------------------------------------------------|---------------|
| Endpoint para cargar datos reales de la barrera por barrierICode          | PASO 2        |
| Shape del payload de respuesta (campos de barrier_v2 + barrier_update)    | PASO 2        |
| Endpoint para actualizar el estado activo/inactivo del toggle             | E03           |
| Cómo se poblará pendingTasks y completedTasks (¿desde case_task?)         | PASO 2        |
| Endpoint para cargar eventos del timeline de la barrera                   | PASO 2 (tab timeline) |
