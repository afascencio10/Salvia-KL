━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Estado: ✅ Implementado y probado — 2026-09-27 (MD sincronizado con el código 2026-10-31)
   Funciones: VictimContactPOST_GET_Public_V2() (Go) → data() y mounted() del app Vue `appV2` de set_victim_contact_v2.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  Accept:  cabecera de la petición   → navegador (HTML) o cliente móvil/QA (`application/json`)
}
(GET simple, sin :id ni querystring; sin sesión ni permisos)

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

PASO 2 — Renderizar el wizard
  RenderTemplate("set_victim_contact_v2", {
    windowTitle:    Locale["sp"]["set_victim_contact_window_title"],
    yes_no, reportTypeV2, adjustmentsGBV, anonimato, parentesco,   // catálogos desde el caché de enums
    CaptchaID:      captchaID,
    formPath:       FormPaths["sp"]["VictimContactPOST_Public_V2"]  // "/salvia/public-v2"
  })
  → El captchaID viaja embebido en el HTML ({{.CaptchaID}}); NO es una llamada AJAX aparte.

PASO 3 — Inicializar el wizard (cliente)
  data():    step = 'auth', errors = {}, form = blankForm(CaptchaID), radicado = '', fechaActual = '', horaActual = ''
  mounted(): document.title = windowTitle; #app (el <body> del layout) pasa de display:none a block

El captcha (imagen + audio + input) se muestra recién en el paso 'formulario', pero el id ya está disponible desde la carga.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Ninguno.
