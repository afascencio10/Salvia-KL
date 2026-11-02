━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma «Registro de caso»
   Tipo: User Interaction / Backend write           (antes M-10 «Apertura de caso», diferido; REESCRITO en WP-1 T7, 2026-11-02)
   Estado: ✅ Implementado y probado (ver QA/M10-M12-registro-de-caso-qa-cases.md)
   Funciones: seleccionarPendingClose('registro') · convertirEnRegistro()   (gestion-reporte-modal.js)
              · ReportConversionController.Convert → ReportConversionService.ConvertToDraft
              · ReportConversionRepository.CheckConvertible / ConvertReportToDraft
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: { icode ← reporte abierto en el modal ; sesión ← cookie (UserICode, CurrentRole) }

FRONTEND
PASO 1 — Visibilidad: el botón «Registro de caso» solo existe en step 'postContacto' con canManage y canConvert (Víctima y Tercero).
         SI el tipo es Familiar u Hombre: aviso «Por ahora no hay procesos de atención para este tipo de reportes» (M-01 PASO 7) → TERMINAR ejecución
PASO 2 — Clic en «Registro de caso» → pendingClose = 'registro' → ConfirmBox:
         «Se creará un registro de caso en borrador con los datos del reporte (nombres, apellidos, teléfono, hechos y ajuste razonable).
          El reporte quedará en "Borrador iniciado" hasta que el caso se complete.»   [Cancelar] [Crear registro de caso]
PASO 3 — «Crear registro de caso»: saving = true ; actionError = null
         POST /api/v1/victim-cases-v3/desde-reporte/:icode        // sin body

BACKEND
PASO 4 — Controller: sesión → 401 ; permisos manage_victim_contacts y manage_registro_v3 (op y sv) → 401
PASO 5 — Service.ConvertToDraft → repo.CheckConvertible (falla rápido, sin generar credenciales):
         SI no existe → 404
         SI el caso de Registro V3 vinculado más reciente está en te/cd, o el reporte ya tiene motivo de cierre o está invalidado → 409 «el reporte ya está cerrado»
         SI hay un borrador vinculado (bo/in) → 409 { error: «el reporte ya tiene un borrador de caso», caseIcode }
         SI hay un caso vinculado de otro flujo → 409 «el reporte ya tiene un caso vinculado»
         SI el tipo no es Víctima ni Tercero (o es de V1, sin tipo V2) → 409 «Por ahora no hay procesos de atención para este tipo de reportes»
         SI no tiene al menos un intento exitoso → 409 «aún no hay contacto exitoso»
PASO 6 — Repository.ConvertReportToDraft: UNA transacción (db.Transaction) con el reporte bloqueado `FOR UPDATE` (un doble clic no crea dos
         borradores: la segunda llamada espera, vuelve a validar y recibe el 409 con el caseIcode) :
   a. Crear el borrador con los mismos pasos de CreateDraft (se invoca `CreateDraft` dentro de un SAVEPOINT de la transacción, sin modificarlo):
        general_user_profile y general_user (login y clave aleatorios, rol 'us') · victim_case en 'bo' con agent_id VACÍO
        (el agente del equipo de seguimiento se asigna después) · victim_case_form2 con valores temporales · teléfono (CI_4) si el reporte lo trae
   b. Vincular: UPDATE victim_case SET victim_case_victim_contact = victim_contact_id del reporte, nombres y apellidos de la persona afectada (CI_2);
      el documento queda VACÍO (no el temporal cc/0000000000: ver el GAP de Registro V3 flow-E19)
   c. Prellenar y marcar como respondidas (markKeysAnswered): CI_2 nombres y apellidos · CI_4 teléfono · HV_1 hechos · CI_8 ajuste razonable (copia de las opciones del reporte)
   d. Hito «Borrador de caso iniciado» en el timeline (VictimContactID y CaseID = el caso nuevo; actor = quien convierte)
   e. Commit. El reporte NO se cierra: queda «Borrador iniciado» (estado derivado del caso vinculado) hasta que el caso se complete.
        SI cualquier paso falla → Rollback → 500 «error al crear el registro» (no queda caso, usuario, vínculo ni evento)
PASO 7 — Después del commit: StoreCredentials(caso, login, clave) como hace CreateDraft. SI falla: WARN en el log ; el caso queda creado
PASO 8 — Responder 200 { "caseIcode": "<uuid>", "redirect": "/salvia/registro-v3?caso=<uuid>" }

FRONTEND (continuación)
PASO 9 — SI 200: window.location.assign(redirect)  → Registro V3 E19
         SI 409 con caseIcode: draftCaseIcode = caseIcode ; step = 'borradorIniciado'  → M-12
         SI otro 409 o 500: actionError con el mensaje ; saving = false ; el modal conserva el paso

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| **GAP-6 en simple:** antes se pensó preguntar «¿Remitir a Atención Especializada? Sí/No» al abrir el caso, solo para Hombre o Familiar. Nadie confirmó que esa pregunta se necesite ni dónde se guardaría. Como hoy esos tipos no se convierten en caso, no se hace | PASO 2 |
| CI_8 depende de CI_7: no se marca CI_7 (decisión D4) | PASO 6.c |
| El tipo y número de documento quedan vacíos hasta que el agente complete CI_5 y CI_6 | PASO 6.a |
| El estado final del reporte es DERIVADO del caso (ningún código escribe `ap`); si se prefiere guardar el cierre al completarse el caso, WP-3 tendría que avisar desde `Finalize` (decisión D8) | PASO 6.e |
| No hay forma de descartar el borrador desde el reporte (fuera de alcance); mientras exista, el reporte no admite más intentos ni cierres (409) | PASO 6.e |
