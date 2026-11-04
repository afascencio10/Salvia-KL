━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Estado: ✅ Implementado y probado — 2026-09-27 (MD sincronizado con el código 2026-10-31)
   Funciones: VictimContactPOST_GET_Public_V2() (Go) → data() y mounted() del app Vue `appV2` de set_victim_contact_v2.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  Accept:  cabecera de la petición   → navegador (HTML) o cliente móvil/QA (`application/json`)
  sesión:  cookie `userData` (opcional)   → 🆕 [WP-1 T8] navegador; sin sesión el flujo es el de siempre
}
(GET simple, sin :id ni querystring; no exige sesión ni permisos)

PASO 1 — El backend genera un captchaID nuevo
  GET /salvia/public-v2/nuevo → VictimContactPOST_GET_Public_V2()
  captcha.New() → captchaID
  SetHeaderNoCache(c)

SI la cabecera Accept es "application/json":
  → Responde 200 con los catálogos y el captchaID, sin renderizar HTML:
    { "yesNo": [...], "reportType": [...], "adjustmentsGBV": [...], "anonimato": [...], "parentesco": [...], "captchaID": "..." }
    // cada ítem: {icode, code, name}; reportType sale de la clave victim_case_form2_report_type_v2
  → TERMINAR ejecución (lo usan la app Flutter y qa-salvia: crearReporteQA)

SI NO:
  → CONTINÚA FLUJO GENERAL

PASO 2 — 🆕 [T8] Leer la sesión opcional del reportante (readOptionalReporter)
  sessionID = sessions.Default(c).Get("userData"); GetCommonSession(sessionID)
  SI no hay cookie, la sesión venció o no existe → reporter = nil (anónimo, sin error)
  SI NO → reporter = {UserICode, RoleCode = CurrentRole, DisplayName = nombres + apellidos (o el login),
                      RoleLabel y ShowBanner solo si el rol es en, op o et}

PASO 3 — Renderizar el wizard
  RenderTemplate("set_victim_contact_v2", {
    windowTitle:    Locale["sp"]["set_victim_contact_window_title"],
    yes_no, reportTypeV2, adjustmentsGBV, anonimato, parentesco,   // catálogos desde el caché de enums
    CaptchaID:      captchaID,
    formPath:       FormPaths["sp"]["VictimContactPOST_Public_V2"],  // "/salvia/public-v2"
    reporterName, reporterRole, reporterRoleLabel:  los de reporter solo si ShowBanner; si no ""   // 🆕 [T8]
    sinPasoUsuario: reporter ≠ nil Y NO ShowBanner   // 🆕 [T8] sesión de otro rol (sv, ad, ro…): arranca en 'auth'
  })
  → El captchaID viaja embebido en el HTML ({{.CaptchaID}}); NO es una llamada AJAX aparte.

PASO 4 — Inicializar el wizard (cliente)
  data():    pasoInicial = sinPasoUsuario ? 'auth' : 'usuario'   // ✏️ [T8] antes siempre 'auth'
             step = pasoInicial, errors = {}, form = blankForm(CaptchaID), radicado = '', copiado = false, fechaActual = '', horaActual = ''
  El paso 'usuario' muestra «Está reportando como …» (si reporterName ≠ "") o la pregunta «¿Tiene usuario?» — ver E-12.
  mounted(): document.title = windowTitle; #app (el <body> del layout) pasa de display:none a block

El captcha (imagen + audio + input) se muestra recién en el paso 'formulario', pero el id ya está disponible desde la carga.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno.
