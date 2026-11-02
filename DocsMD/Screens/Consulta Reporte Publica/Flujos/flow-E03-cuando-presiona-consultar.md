━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona «Consultar»
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — WP-1 T4 (2026-11-02) (ver QA/CR-E03-consultar-qa-cases.md)
   Funciones: consultar() · ReportStatusController.Lookup() · ReportStatusService.Lookup() · ReportStatusRepository.FindByRadicado()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { radicado: campo normalizado (CR-E02) → input (también Enter) ; ip: c.ClientIP() → servidor }

FRONTEND
PASO 1 — SI !radicadoValido o consultando o bloqueado → TERMINAR ejecución
         consultando = true ; mensajeError = ''
PASO 2 — POST /api/v1/public/consulta-reporte   { "radicado": "K7QM-4R2X" }

BACKEND (público: sin sesión, sin captcha)
PASO 3 — Cabeceras no-cache ; leer el body (máx. 1 KiB; 400 «Solicitud inválida.» si es JSON inválido)
PASO 4 — Limitador por IP (en memoria)
  SI la IP acumula ≥ 5 fallos en los últimos 15 min:
    → 429 { "mensaje": "Demasiados intentos. Intente de nuevo en N minutos.", "reintentarEnSegundos": S } + cabecera Retry-After
    → TERMINAR ejecución                       // no consulta la BD
PASO 5 — Normalizar: trim, mayúsculas, sin espacios ni guiones (incluidos los guiones largos), luego «XXXX-XXXX» ;
         formatoValido = ^[A-HJ-NP-Z2-9]{4}-[A-HJ-NP-Z2-9]{4}$
PASO 6 — SI formatoValido: fila = FindByRadicado(canónico)   // SELECT parametrizado; une el caso vinculado más reciente de Registro V3 (bo/in/te/cd)
         SI NO: fila = nil
PASO 7 — SI fila = nil (formato inválido o inexistente): registrar un fallo para la IP →
         404 { "mensaje": "No encontramos un reporte con ese radicado. Verifique que esté bien escrito." }   // mismo cuerpo en ambos casos
PASO 8 — Derivar el estado público (ver Listado Reportes/related-tables.md, «Estado derivado del reporte»)
  cerrado    (motivo de cierre, estado legado 'i', o caso vinculado en te/cd)
  en_gestion (sin cierre y con agente, con algún intento o con borrador bo/in)
  recibido   (el resto)
  → 200 { radicado, estado, estadoTexto, descripcion, fechaRadicacion: "YYYY-MM-DD" }
  // Nada más. Un acierto no suma fallos. Un error de BD responde 500 «No fue posible consultar…» y tampoco suma fallos.

FRONTEND (continuación)
PASO 9 — consultando = false ; SEGÚN el status:
  200: resultado = respuesta ; vista = 'resultado'
  404: mensajeError = respuesta.mensaje
  429: mensajeError = respuesta.mensaje ; bloqueado = true ; setTimeout(reintentarEnSegundos) → bloqueado = false
  otro: mensajeError = "No fue posible consultar en este momento. Intente de nuevo."

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| La IP real detrás de un proxy depende de los proxies de confianza de Gin (con la configuración por defecto un cliente podría falsear X-Forwarded-For) | PASO 4 |
| El limitador en memoria se reinicia con el servicio y no se comparte entre réplicas | PASO 4 |
| «Asignado sin intentos» se cuenta como «en gestión» (decisión de lenguaje para la persona) | PASO 8 |
