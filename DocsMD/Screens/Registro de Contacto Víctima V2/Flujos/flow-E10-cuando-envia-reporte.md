━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Enviar Reporte"
   Tipo: User Interaction
   Estado: ✅ Implementado y probado — 2026-09-27 (ver QA/E10-cuando-envia-reporte-qa-cases.md)
   Funciones: submitForm() (frontend) · VictimContactPOST_Public_V2() (facade)
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
  nombreVictima:         form.nombreVictima             → input, si tipoReporte !== 'Hombre'
  apellidoVictima:       form.apellidoVictima           → input, si tipoReporte !== 'Hombre'
  telefonoVictima:       form.telefonoVictima           → input, solo Tercero/Victima
  disposicionLlamada:    form.disposicionLlamada        → select, solo Tercero/Victima
  horaContacto:          form.horaContacto              → input time, si disposicionLlamada.code === 'y'
  descripcionHechos:     form.descripcionHechos         → textarea, todos los tipos
  captchaID:             form.captchaID                 → obtenido en E-01
  captchaSolution:       form.captchaSolution            → input, paso 'formulario' (BloqueCaptcha)
}

Nota: en la implementación final, todos los campos que referencian un enum
(tipoReporte, ajusteRazonable, autorizacion, anonimato, parentesco, rolCuidado,
victimaConoce, disposicionReportante, disposicionLlamada) se guardan en `form`
como el OBJETO completo `{icode, code, name}` devuelto por el backend en
E-01/E-07 (no como string) — así el payload solo necesita reenviar `{icode}`
y las comparaciones usan `.code` (ej. `form.disposicionLlamada.code === 'y'`).

PASO 1 — Validar el formulario según tipoReporte (isFormValid, cliente)
  Validar formulario:
    • autorizacion:        ya validado en E-02 (guard: código === 'y')
    • tipoReporte:         required
    • ajusteRazonable:     required
    • descripcionHechos:   required, min_chars: 10
    • captchaSolution:     required
    • SI tipoReporte.code === 'r' (Tercero):
        - anonimato:       required
        - SI anonimato.code === 'pr': nombreReportante y apellidoReportante required
        - rolCuidado:      required
        - victimaConoce:   required
    • SI tipoReporte.code === 'fm' (Familiar):
        - nombreReportante, apellidoReportante: required
        - parentesco:      required
        - rolCuidado:      required
        - disposicionReportante: required
        - SI disposicionReportante.code === 'y': horaReportante required
    • SI tipoReporte.code === 'ho' (Hombre):
        - nombreReportante, apellidoReportante: required
        - telefonoReportante: required, formato teléfono
        - disposicionReportante: required
        - SI disposicionReportante.code === 'y': horaReportante required
    • SI tipoReporte.code === 'r' || tipoReporte.code === 'v':
        - nombreVictima, apellidoVictima: required
        - telefonoVictima: required, formato teléfono
        - disposicionLlamada: required
        - SI disposicionLlamada.code === 'y': horaContacto required

  SI validación con errores:
    → Mostrar error por campo (mismo patrón visual que el resto del wizard)
    → TERMINAR ejecución (el usuario permanece en el paso 'formulario')

  SI validación ok:
    → CONTINÚA FLUJO GENERAL

PASO 2 — Armar el payload para el backend (buildPayload(), implementación final)
  victimContact: {
    // Decisión 7A: sin víctima nombrada (Hombre), Names/LastNames van vacíos —
    // el agente los completa tras el recontacto. Van como '' (no null): las
    // columnas de victim_contact son VARCHAR nullable a nivel de esquema, pero
    // el DTO Go es un string plano, así que el valor vacío es lo que se envía.
    names:      tipoReporte.code === 'ho' ? '' : form.nombreVictima,
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
      bestContactTime:     form.horaContacto ? (form.horaContacto + ':00') : '00:00:00',
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
      reporterBestContactTime: form.horaReportante ? (form.horaReportante + ':00') : null,
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
    2. Verifica captcha (bypass solo si User-Agent = "flutter-client", igual
       que V1).
    3. Resuelve y valida los enums recibidos (reportType contra la clave `_v2`,
       adjustmentsGBV, authorizationAnswer === 'y', y según tipoReporte:
       anonimato/parentesco/hasCareRole/victimAwareOfReport/
       reporterWillReceiveCall) contra el caché local de enums.
    4. SEGÚN tipoReporte:
         'v' | 'r' → valida nombre/apellido/teléfono/disposición de la víctima
         'fm'      → valida datos del reportante + nombre de la víctima (sin
                     teléfono/disposición) + parentesco; fuerza victimColPhone = 0
         'ho'      → valida solo datos del reportante; limpia names/lastNames
                     a '' (decisión 7A); fuerza victimColPhone = 0

  Sistema.responde({
    codigo:  200 | 400 | 500,
    cuerpo:  {"icode": "..."} | {"default":{...}, "form2":{...}, "form2V2":{...}}
  })

SI código === 400 (captcha inválido, enum no encontrado en la clave `_v2`,
  campo requerido faltante según tipoReporte):
  → Mostrar los errores devueltos, mapeados a los campos del wizard por
    data-entity/data-field
  → SI el error es de captcha: regenerar captchaID (el backend devuelve uno
    nuevo en errors.default.captchaID, mismo patrón que V1)
  → TERMINAR ejecución (el usuario permanece en 'formulario')

SI código === 200:
  → Dentro de una sola transacción: SetVictimContact (a) → SetVictimContactForm2
    (b, con FK a (a)) → setVictimContactEnumsMultiple (adjustmentsGBV) →
    SetVictimContactForm2V2 (c, con FK a (b)) — las 3 funciones DAO existentes
    de V1 se invocan sin modificarlas
  → VictimContactICode = utils.GetUUID() generado al insertar (a) — es el
    "número de radicado", mismo mecanismo que V1. A diferencia de V1 (que
    responde con cuerpo vacío), V2 sí devuelve {"icode": "..."} en el body
  → PASO 4 — Guardar el resultado en el cliente
    radicado.value = respuesta.icode
    fechaActual.value = hoy formateada
    horaActual.value  = hora actual formateada
    step.value = 'confirmacion'
    Se oculta explícitamente #successOverlay (bug #6, ver related-tables.md) —
    el popup de éxito compartido no se auto-oculta y bloqueaba el botón
    "Volver al Inicio" de la confirmación
  → CONTINÚA FLUJO GENERAL

Casos de prueba, edge cases y verificación en BD real: ver
QA/E10-cuando-envia-reporte-qa-cases.md.
