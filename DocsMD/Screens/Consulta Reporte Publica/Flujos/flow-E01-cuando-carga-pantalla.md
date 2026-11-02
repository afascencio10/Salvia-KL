━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Estado: ✅ Implementado y probado — WP-1 T4 (2026-11-02)
   Funciones: ReportStatusGET_Public() (Go) → data() y mounted() del app Vue `appConsulta`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { GET /salvia/public-consulta  → navegador (sin sesión, sin parámetros; la URI contiene /public: AuthMiddleware la deja pasar) }

PASO 1 — SetHeaderNoCache(c)
PASO 2 — Renderizar la plantilla
  RenderTemplate("get_report_status", {
    windowTitle:   Locale["sp"]["report_status_window_title"],     // «Consultar mi reporte»
    lang: "sp", locale: Locale,
    apiPath:       "/api/v1/public/consulta-reporte",               // POST de CR-E03
    newReportPath: "/salvia/public-v2/nuevo"                        // enlace del pie
  })
PASO 3 — Inicializar la pantalla (cliente)
  data():    vista = 'consulta', radicado = '', consultando = false, bloqueado = false, mensajeError = '', resultado = null
  mounted(): document.title = windowTitle ; #app pasa de display:none a block ; foco en el campo radicado

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno.
