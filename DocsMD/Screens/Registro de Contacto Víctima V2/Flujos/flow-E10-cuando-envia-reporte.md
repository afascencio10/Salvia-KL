━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Enviar Reporte"
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — 2026-09-27 (ver QA/E10-cuando-envia-reporte-qa-cases.md) · MD sincronizado con el código 2026-10-31
   Funciones: isFormValid (computed) · submitForm() · buildPayload() (frontend) · VictimContactPOST_Public_V2() (facade)
              · SetVictimContactV2() (controller) · SetVictimContact() / SetVictimContactForm2()
              (DAO existentes de V1, sin modificar) · SetVictimContactForm2V2() (DAO nuevo)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tipoReporte:           form.tipoReporte              → paso 'perfil' (E-04)
  ajusteRazonable:       form.ajusteRazonable           → paso 'perfil' (E-05)
  autorizacion:          form.autorizacion              → paso 'auth' (E-02)
  anonimato:             form.anonimato                 → paso 'formulario', solo Tercero (E-08)
  nombreReportante:      form.nombreReportante          → input, solo si mostrarNombreReportante
  apellidoReportante:    form.apellidoReportante        → input, solo si mostrarNombreReportante
  correoReportante:      form.correoReportante          → input, opcional
  parentesco:            form.parentesco                → select, solo Familiar
  telefonoReportante:    form.telefonoReportante        → input, solo Familiar/Hombre
  rolCuidado:            form.rolCuidado                → select, solo Tercero/Familiar
  victimaConoce:         form.victimaConoce             → select, solo Tercero
  disposicionReportante: form.disposicionReportante     → select, solo Familiar/Hombre
  horaReportante:        form.horaReportante            → input time, si disposicionReportante.code === 'y'
  nombreVictima:         form.nombreVictima             → input, si tipoReporte.code !== 'ho'
  apellidoVictima:       form.apellidoVictima           → input, si tipoReporte.code !== 'ho'
  telefonoVictima:       form.telefonoVictima           → input, solo Tercero ('r') / Víctima ('v')
  disposicionLlamada:    form.disposicionLlamada        → select, solo Tercero ('r') / Víctima ('v')
  horaContacto:          form.horaContacto              → input time, si disposicionLlamada.code === 'y'
  descripcionHechos:     form.descripcionHechos         → textarea, todos los tipos
  captchaID:             form.captchaID                 → generado por el backend en E-01 (embebido en el HTML)
  captchaSolution:       form.captchaSolution            → input, paso 'formulario' (BloqueCaptcha)
}

Nota: en la implementación final, todos los campos que referencian un enum
(tipoReporte, ajusteRazonable, autorizacion, anonimato, parentesco, rolCuidado,
victimaConoce, disposicionReportante, disposicionLlamada) se guardan en `form`
como el OBJETO completo `{icode, code, name}` devuelto por el backend en
E-01/E-07 (no como string) — así el payload solo necesita reenviar `{icode}`
y las comparaciones usan `.code` (ej. `form.disposicionLlamada.code === 'y'`).

PASO 1 — Habilitar el botón según tipoReporte (isFormValid, cliente)
  El botón "Enviar Reporte" lleva :disabled="!isFormValid"; no hay validación al hacer clic ni errores por campo.
  isFormValid === true cuando:
    • tipoReporte y ajusteRazonable:  respondidos
    • descripcionHechos:   ≥ 10 caracteres (trim)
    • captchaSolution:     no vacío
    • SI tipoReporte.code === 'v' || 'r':
        - nombreVictima, apellidoVictima, telefonoVictima, disposicionLlamada: respondidos
        - SI disposicionLlamada.code === 'y': horaContacto respondida
    • SI tipoReporte.code === 'r' (Tercero):
        - anonimato respondido
        - SI anonimato.code === 'pr': nombreReportante y apellidoReportante respondidos
        - rolCuidado y victimaConoce respondidos
    • SI tipoReporte.code === 'fm' || 'ho':
        - nombreReportante, apellidoReportante, telefonoReportante, disposicionReportante: respondidos
        - SI disposicionReportante.code === 'y': horaReportante respondida
    • SI tipoReporte.code === 'fm' (Familiar), además:
        - nombreVictima, apellidoVictima, parentesco, rolCuidado: respondidos
  El teléfono NO se valida por formato (solo "no vacío"); al armar el payload se limpia con replace(/\D/g,'').

  SI isFormValid === false:
    → El botón permanece deshabilitado (sin mensaje de qué falta)
    → TERMINAR ejecución
  SI isFormValid === true:
    → CONTINÚA FLUJO GENERAL

PASO 2 — Armar el payload para el backend (buildPayload(), implementación final)
  victimContact: {
    // Decisión 7A: sin víctima nombrada (Hombre, code 'ho'), Names/LastNames van vacíos —
    // el agente los completa tras el recontacto. Van como '' (no null) y así se
    // guardan: en BD queda la cadena vacía, no NULL (verificado en una fila real 'ho').
    names:      tipoReporte.code === 'ho' ? '' : form.nombreVictima,   // NO se filtra por anonimato: lo digitado viaja igual
    lastNames:  tipoReporte.code === 'ho' ? '' : form.apellidoVictima,
    captchaID:       form.captchaID,
    captchaSolution: form.captchaSolution,

    // Columnas EXISTENTES de victim_contact_form2 — sin cambio de esquema.
    // reporterNames también es una columna EXISTENTE: se reutiliza para el
    // nombre (no el apellido) del reportante — corrección sobre la planeación
    // original, ver related-tables.md.
    form2: {
      reportType:          form.tipoReporte,              // {icode} de la clave victim_case_form2_report_type_v2 — decisión 3A
      adjustmentsGBV:      [ form.ajusteRazonable ],       // arreglo de 1 elemento
      reporterNames:       form.nombreReportante || '',
      reporterPhone:       form.telefonoReportante ? parseInt(...) : 0,
      hasCareRole:         form.rolCuidado,
      victimAwareOfReport: form.victimaConoce,
      victimColPhone:      form.telefonoVictima ? parseInt(...) : 0,
      // Decisión 6B: este par queda EXCLUSIVO de la víctima (Tercero/Victima).
      // Para Familiar/Hombre, victim_col_phone y best_contact_time son
      // NOT NULL en la tabla existente (no se puede alterar sin tocar V1) —
      // el backend los fuerza a 0 / vacío en esos casos (ver PASO 3).
      willReceiveCall:     form.disposicionLlamada,
      bestContactTime:     form.horaContacto ? (form.horaContacto + ':00') : '00:00:00',   // "HH:MM:SS"
      factsDescription:    form.descripcionHechos,
    },

    // Tabla de extensión NUEVA victim_contact_form2_v2 — decisión 2A.
    // authorizationAnswer vive AQUÍ, no en form2: victim_contact_form2 no
    // tiene esa columna (V1 la valida pero nunca la persiste — hallazgo de
    // la implementación, ver related-tables.md).
    form2V2: {
      authorizationAnswer:     form.autorizacion,
      anonimato:               form.anonimato || null,                          // decisión 5B
      parentesco:              form.parentesco || null,
      reporterLastNames:       form.apellidoReportante || '',
      reporterEmail:           form.correoReportante || '',
      // Decisión 6B: par EXCLUSIVO del reportante (Familiar/Hombre)
      reporterWillReceiveCall: form.disposicionReportante || null,
      reporterBestContactTime: form.horaReportante ? (form.horaReportante + ':00') : null,   // "HH:MM:SS" — ver L1: el DTO no lo deserializa
    }
  }

PASO 3 — Enviar el reporte
  POST /salvia/public-v2
  {
    "victimContact": <payload del PASO 2>
  }

  Backend (SetVictimContactV2):
    1. Decodifica el JSON en DOS pasadas independientes sobre los mismos bytes
       crudos — una para VictimContactDTO (raíz + form2), otra para el wrapper
       de form2V2. (Embeber ambos en un solo struct rompía la decodificación
       de form2V2 — bug encontrado en implementación, ver related-tables.md.)
       utils.JSONToStruct descarta los errores de decodificación: un campo que
       no se pueda parsear queda en su valor cero (ver L1).
    2. Verifica captcha (bypass solo si User-Agent = "flutter-client", igual
       que V1).
    3. authorizationAnswer debe resolver a una fila de enum con code 'y'.
    4. reportType debe resolver a una fila de la clave `_v2`; adjustmentsGBV
       debe traer ≥ 1 opción y todas deben resolver (error de getAndVerify…).
    5. factsDescription: len < 10 → error "Este campo es obligatorio".
    6. SEGÚN code del tipo de reporte:
         'v' | 'r' → names, lastNames, victimColPhone ≠ 0, willReceiveCall resuelto
                     · 'r' además: anonimato resuelto; si es 'pr': reporterNames y
                     reporterLastNames; hasCareRole y victimAwareOfReport resueltos
         'fm' | 'ho' → reporterNames, reporterLastNames, reporterPhone ≠ 0,
                     reporterWillReceiveCall resuelto (si es 'y': reporterBestContactTime
                     no puede ser cero)
                     · 'fm' además: names, lastNames, parentesco y hasCareRole
                     · 'ho': names = lastNames = '' (decisión 7A)
                     · ambos: victimColPhone se fuerza a 0
         otro código → error en reportType
    7. SI hay errores acumulados: collectedErrors["default"]["captchaID"] = captcha.New()
       y responde 400 (el captchaID se regenera con CUALQUIER 400, no solo el de captcha).

  Sistema.responde({
    codigo:  200 | 400 | 500,
    cuerpo:  {"icode": "..."} | {"default":{...}, "form2":{...}, "form2V2":{...}}
  })

En el callback de newEntity (siempre): form.captchaSolution = ""

SI código === 400 (captcha inválido, enum no encontrado en la clave `_v2`,
  campo requerido faltante según tipoReporte):
  → errors = respuesta  (el alert #alert muestra errors.default.global_msg si existe; el watcher de
    errors arma un resumen "El formulario tiene N errores. <etiqueta>: <mensaje>…" en la región
    aria-live #global-errors-live, resolviendo cada etiqueta por data-entity/data-field/data-label)
  → SI errors.default.captchaID existe: form.captchaID = ese valor (el backend lo envía en todo 400)
  → TERMINAR ejecución (el usuario permanece en 'formulario'; no se marcan los campos en rojo)

SI código es distinto de 200 y 400:
  → alert("Se produjo un error interno, por favor vuelva a intentarlo")
  → TERMINAR ejecución

SI código === 200:
  → Dentro de una sola transacción: SetVictimContact (a) → SetVictimContactForm2
    (b, con FK a (a)) → setVictimContactEnumsMultiple (adjustmentsGBV) →
    SetVictimContactForm2V2 (c, con FK a (b)) — las 3 funciones DAO existentes
    de V1 se invocan sin modificarlas
  → VictimContactICode = utils.GetUUID() generado al insertar (a) — es el
    "número de radicado", mismo mecanismo que V1. A diferencia de V1 (que
    responde con cuerpo vacío), V2 sí devuelve {"icode": "..."} en el body
  → PASO 4 — Guardar el resultado en el cliente
    this.radicado = respuesta.icode
    this.fechaActual = new Date().toLocaleDateString('es-CO')
    this.horaActual  = new Date().toLocaleTimeString('es-CO')
    this.step = 'confirmacion'
    Se oculta explícitamente #successOverlay (bug #6, ver related-tables.md) —
    el popup de éxito compartido no se auto-oculta y bloqueaba el botón
    "Volver al Inicio" de la confirmación
  → CONTINÚA FLUJO GENERAL

Casos de prueba, edge cases y verificación en BD real: ver
QA/E10-cuando-envia-reporte-qa-cases.md.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| **L1** — `reporterBestContactTime` viaja como "HH:MM:SS" pero `VictimContactForm2V2DTO` (`time.Time`, sin `UnmarshalJSON`) no lo acepta: `JSONToStruct` ignora el error, el campo queda en cero y el backend responde 400 `form2V2.reporterBestContactTime` cuando Familiar/Hombre eligen disposición "Sí". Con "No" funciona. Observado por API el 2026-10-13; no verificado por el wizard | PASO 2 y PASO 3 (validación 6) |
| **L4** — el 400 siempre regenera el captchaID, aunque el captcha haya sido correcto | PASO 3 (validación 7) |
| El cliente no marca en rojo los campos con error devuelto por el backend (solo alert global + región aria-live) | PASO 3 (400) |

