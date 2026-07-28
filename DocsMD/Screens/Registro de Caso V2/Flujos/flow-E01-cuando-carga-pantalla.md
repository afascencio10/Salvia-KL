━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Función: mounted()  ·  bloque middleware del gate de autorización
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  contactId:  UUID opcional del contacto previo   → path param :id (si viene de "Registro de Contacto Víctima")
  formId:     UUID fijo del formulario            → hardcoded "0a24ab30-3cfc-4861-b74d-65d21524bc00"
}


PASO 1 — Resolver el gate de autorización (fuera de dinamic-form)

  authorizationAnswer = null

  SI authorizationAnswer.code != 'y':
    → Mostrar SeccionAutorizacion (texto legal Ley 1581 de 2012 + select Sí/No)
    → NO montar <dinamic-form> todavía
    → TERMINAR ejecución (se re-ejecuta este paso cada vez que cambia authorizationAnswer)

  SI authorizationAnswer.code == 'y':
    → CONTINÚA PASO 2


PASO 2 — Resolver submissionId (si viene de un contacto previo)

  SI contactId existe:
    GET /api/v1/victim-contacts/{contactId}   // o el endpoint que ya usa la pantalla actual
    → SI el contacto ya tiene un form_submission asociado a este formId:
        submissionId = ese id  (continuar un intento existente)
      SI NO:
        submissionId = null  (formulario nuevo, sin prellenado — ver GAP)

  SI contactId no existe:
    submissionId = null


PASO 3 — Cargar catálogo de enums y geografía en paralelo (sin bloquear)

  Promise.all([
    GET /api/v1/enums?categories=victim_case_form2_*,yes_no,docType2,aggressor_occupation
      → formState.enums = { [categoria]: [{label, value}] }
    GET /api/v1/locations/departments
      → formState.departments = [{ label, value }]
  ])
  → Sin loader propio, sin bloquear el montaje de dinamic-form
  → Si alguna falla: se loggea advertencia, las opciones quedan vacías []


PASO 4 — Inicializar formState con los flags calculados

  formState = {
    enums:            {},   // se puebla en PASO 3
    departments:      [],   // se puebla en PASO 3
    geo: {
      residencia: { cities: [], towns: [] },
      hechos:     { cities: [], towns: [] },
      atencion:   { cities: [], towns: [] },
    },
    wasPartner:                '',
    partnerKnown:               '',
    hasRisk:                    '',
    riskBadgeText:               '',
    hasWorkplaceScope:          '',
    hasSelectedViolenceTypes:   '',
    violenceSubtypes:           [],
  }


PASO 5 — Montar DinamicForm

  DinamicForm recibe:
    :form-id       = "0a24ab30-3cfc-4861-b74d-65d21524bc00"
    :submission-id = submissionId       (null si es nuevo)
    :can-edit      = true               (Registro de Caso siempre es editable al crear)
    :form-state    = formState

  → El componente dispara su propio E-01 interno (carga la estructura +
    respuestas previas si submissionId existe) — ver dinamic-form-events.md

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                          | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| Endpoint real para cargar el catálogo de enums (¿existe ya uno genérico, o hay que crearlo?) — hoy se inyectan server-side en el template, no vía API | PASO 3 |
| Si viene de un contacto previo, cómo se prellenan las 128 respuestas del submission a partir de `victim_contact_form1`/`form2` (mapeo campo a campo pendiente) | PASO 2 |
| docType2 / aggressor_occupation no están en victim_case_form2_enums (ver G-01 de form-registro-caso-v2-data.md) — el endpoint de PASO 3 necesita saber de dónde traerlas | PASO 3 |
