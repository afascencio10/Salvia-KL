━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  psicosocialId:  UUID de la remisión psicosocial  → route param :psicosocial_id
  userICode:      i_code del profesional           → plantilla Go ({{ .userICode }})
  contactId:      UUID opcional del team_contact   → query ?contactId= (desde "Ver sesión")
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — registrar_sesion.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar la app Vue

  Vue.createApp monta sobre #psicosocial-app
  → Se ejecuta mounted() → llama a loadPsicosocial()


PASO 2 — Llamar al backend para cargar la sesión

  Mientras la petición está en curso: pageLoading = true → se muestra un loader de
  página (spinner) en lugar de la tarjeta de víctima y del formulario. Implementado.

  contactId = URLSearchParams.get('contactId')   // desde "Ver sesión" en detalle remisión

  GET /api/v1/psychosocial-support/{psicosocialId}/load?agent_id={userICode}
      [&contact_id={contactId}]   // opcional — ver PASO 8

  Nota: el endpoint real usa el prefijo "psychosocial-support" (igual que el resto
  de la API: /detail, /contacts, etc.), no "psicosocial-support". El mock inicial
  del frontend tenía este typo y ya fue corregido.

  SI respuesta no ok (status != 2xx):
    → this.loadError = data.error || 'Error al cargar la sesión psicosocial'
    → TERMINAR ejecución   // DinamicForm no se monta; se muestra ErrorAlert

  SI respuesta ok:
    → data = JSON parseado { formId, formType, submissionId, victimInfo,
                              psicosocialState, formState, canEdit, teamContactId, isCompleted }
    → CONTINÚA PASO 3

  Al finalizar (éxito o error): pageLoading = false.


PASO 3 — Poblar estado de la pantalla

  formId           = data.formId           // UUID del formulario seleccionado por el backend
  submissionId     = data.submissionId     // UUID del form_submission del team_contact resuelto
  victimInfo       = data.victimInfo
  psicosocialState = data.psicosocialState  // { yaHizoPrimerContacto, yaHizoPrimeraAtencion, sessionCount, status }
  formState        = data.formState         // incluye currentBarriers, etc.
  canEdit          = data.canEdit           // lo calcula el backend (PASO 4 / 11)


PASO 4 — canEdit (fuente de verdad: backend)

  El frontend NO recalcula canEdit. Usa `data.canEdit` del backend:

  canEdit = false SI:
    - team_contact.is_completed == true   (vista "Ver sesión" de sesión ya guardada)
    - O psicosocial_support.status == 'cerrado'
  canEdit = true en caso contrario (contacto pendiente / en curso)

  Entrada desde detalle remisión:
    - "Ver sesión" → /salvia/psicosocial/registrar/{id}?contactId={tc.id}
    - "Registrar sesión" → /salvia/psicosocial/registrar/{id}  (sin contactId → pendiente)


PASO 5 — DinamicForm se monta con los props calculados

  DinamicForm recibe:
    :form-id       = formId        (UUID devuelto por el backend — varía por sesión)
    :submission-id = submissionId  (del team_contact resuelto — con answers si ya se guardó)
    :can-edit      = canEdit       (desde data.canEdit del backend)
    :form-state    = formState

  Nota: La visibilidad de S2 en el Form de Primer Contacto (Primera Atención) y de S3 en el
  Form de Cierre es gestionada internamente por DinamicForm a través de visibility_conditions
  en la base de datos. No se necesita ningún flag externo en formState.

  → El componente renderiza el formulario en modo lectura o edición según canEdit

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/psychosocial-support/:id/load
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  IMPLEMENTADO en PsychosocialDetailController.LoadSession (mismo controller que
  /detail y /contacts) → PsychosocialDetailService.LoadSession.


INPUT: {
  psicosocialId:  UUID de la remisión  → path param :id
  agentId:        i_code del usuario   → query param agent_id
  contactId:      UUID opcional        → query param contact_id  (team_contact a reabrir)
}


PASO 6 — Cargar la remisión psicosocial

  DB.psychosocial_support.FindByID({ psicosocialId })

  SI no existe:
    → 404 { error: "remisión psicosocial no encontrada" }
    → TERMINAR

  Control de acceso (GAP resuelto — ver decisión más abajo):
  SI ps.ProfessionalID == agentId:
    → acceso permitido (profesional directo)
  SINO SI ps.DuplaID existe Y agentId ∈ { dupla.PsychologistID, dupla.SocialWorkerID }:
    → acceso permitido (miembro de la dupla asignada — psicóloga o trabajador social)
  SINO:
    → 403 { error: "sesión no asignada a este profesional" }
    → TERMINAR

  SI ps.Status === 'cerrado':
    → Continúa (se carga en modo solo lectura; canEdit = false en frontend)


PASO 7 — Seleccionar el formulario según el estado del proceso

  Evaluar en orden (primera condición que se cumpla gana):

  SI ps.YaHizoPrimerContacto == false:
    → formKey = PRIMER_CONTACTO
    → formId  = FORM_ID_PRIMER_CONTACTO   (constante Go del seed)
    // El formulario tiene S2 oculta por defecto; DinamicForm la muestra si
    // el profesional marca "Continuar Primera Atención = Sí" durante la sesión.

  SI ps.YaHizoPrimerContacto == true  Y  ps.YaHizoPrimeraAtencion == false:
    → formKey = PRIMERA_ATENCION
    → formId  = FORM_ID_PRIMERA_ATENCION
    // Solo aplica al Escenario B (segunda llamada independiente).
    // Ya no existe redirección desde PC a PA: si el profesional marcó
    // "Continuar = Sí" en el PC, el backend ya habrá seteado
    // ya_hizo_primera_atencion = true y este bloque no se ejecutará.

  SI ps.YaHizoPrimeraAtencion == true  Y  ps.SessionCount < 3:
    → formKey = ATENCION_PSICOSOCIAL   (antes "SEGUIMIENTO" — ver rebranding Jul 2026)
    → formId  = FORM_ID_ATENCION_PSICOSOCIAL

  SI ps.YaHizoPrimeraAtencion == true  Y  ps.SessionCount >= 3:
    → formKey = CIERRE
    → formId  = FORM_ID_CIERRE
    // S3 (Cierre) permanece oculta hasta que el profesional marque
    // "Cerrar remisión = Sí" en S2; DinamicForm gestiona esa visibilidad.

  IMPLEMENTADO como selectPsicosocialForm(ps) en psychosocial_detail_service.go.


PASO 8 — Resolver el team_contact activo

  ── Rama A: contact_id presente (flujo "Ver sesión") ──

  DB.team_contact.FindOne({ id = contactId, psicosocial_id = psicosocialId })

  SI no existe o no pertenece a la remisión:
    → 404 { error: "contacto de sesión no encontrado" }
    → TERMINAR

  → tc = registro encontrado
  → NO crear team_contact nuevo
  → formId / formKey: reutilizar tc.FormID (si vacío, calcular PASO 7 solo para respuesta;
    no reescribir el contacto completado de forma agresiva)
  → SI tc.FormSubmissionID vacío:
    → 400/404 { error: "la sesión no tiene formulario asociado" }
    → TERMINAR
  → submissionId = tc.FormSubmissionID  (contiene las answers ya guardadas)
  → canEdit = false SI tc.IsCompleted O ps.Status == 'cerrado'; else true

  ── Rama B: sin contact_id (flujo "Registrar sesión" / carga normal) ──

  Buscar un team_contact pendiente para este psicosocial_id:
    DB.team_contact.FindOne({
      psicosocial_id   = psicosocialId,
      is_completed     = false,
      is_psico_session = true,
      deleted_at IS NULL
    })

  SI existe un team_contact pendiente:
    → tc = registro encontrado
    → SI tc.FormID ya está fijado: se reutiliza SIEMPRE (no se vuelve a evaluar el
      PASO 7), para que la sesión no cambie de formulario a mitad de camino si el
      estado del proceso cambia mientras el contacto sigue pendiente.
    → SI tc.FormID está vacío (contacto creado antes de esta migración): se calcula
      con el PASO 7 y se persiste ahora.

  SI no existe:
    → Se calcula formKey/formId con el PASO 7 y se crea un nuevo team_contact:
         case_id          = ps.CaseID
         psicosocial_id   = psicosocialId
         is_completed     = false
         is_psico_session = true
         form_id          = formId
         session_type     = formKey
         created_at       = NOW()

    → Asignación profesional/dupla (nunca se guardan ambos campos):
         SI ps.DuplaID existe:      tc.dupla_id = ps.DuplaID
         SINO SI ps.ProfessionalID: tc.professional_id = ps.ProfessionalID
         SINO:                      tc.professional_id = agentId


PASO 9 — Resolver o crear el FormSubmission

  SI Rama A (contact_id): ya resuelto en PASO 8 — no crear submission nuevo.

  SI Rama B:
    SI tc.FormSubmissionID existe:
      → submissionId = tc.FormSubmissionID
    SI tc.FormSubmissionID está vacío:
      → Crear nueva FormSubmission en DB ({ formId })
      → submissionId = nuevo ID
      → Actualizar tc.FormSubmissionID = submissionId en DB


PASO 10 — Cargar información de la víctima

  DB.victim_cases.FindByICode({ caseId: ps.CaseID })
  → victimInfo = {
      Names, LastNames, Phone,
      GenderIdentity,           // traducido al español
      TownName,
      RiskLevel,                // COALESCE(victim_case_form2_risk_level, 0) — entero 1-4
      CaseICode
    }


PASO 11 — Responder al frontend

  canEdit = !(tc.IsCompleted || ps.Status == 'cerrado')

  200 {
    formId:         "{UUID del formulario del team_contact}",
    formType:       "{PRIMER_CONTACTO | PRIMERA_ATENCION | ATENCION_PSICOSOCIAL | CIERRE}",
    submissionId:   "{UUID del form_submission}",
    teamContactId:  "{UUID del team_contact resuelto}",
    isCompleted:    tc.IsCompleted,
    canEdit:        canEdit,
    victimInfo: { ... },
    psicosocialState: {
      yaHizoPrimerContacto, yaHizoPrimeraAtencion, sessionCount, status
    },
    formState: { currentBarriers: [...] }
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONSTANTES DE FORMULARIO (backend Go)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

IMPLEMENTADO en src/internal/constants/psicosocial_forms.go, con los UUIDs reales
capturados tras ejecutar seed_psicosocial.sql contra Supabase el 2026-07-14:

  const (
    FormIDPrimerContacto      = "439b57e6-07ea-4da4-9721-8ed28c6ca43f"
    FormIDPrimeraAtencion     = "501ab3d7-8382-4447-96a6-f463152693c6"
    FormIDAtencionPsicosocial = "7a7b61b3-7bc3-4088-a74d-0974e54a3563"
    FormIDCierre              = "c31026f8-7ce2-49b9-8990-04b44ed4513c"
  )

El archivo también expone `PsicosocialFormKeyByID` / `PsicosocialFormIDByKey` para
convertir entre el UUID y la clave legible (PRIMER_CONTACTO, etc.), usado tanto en
la selección del formulario (PASO 7) como al reutilizar el `form_id` ya fijado en
un team_contact existente.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅  GAPS — Resueltos (Jul 2026)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                           | Paso afectado | Resolución |
|-----------------------------------------------------------------------------------------------|---------------|------------|
| Confirmar si un profesional de dupla puede cargar la sesión aunque no sea el `professional_id` directo | PASO 6 | Sí puede — se valida que agentId sea `dupla.psychologist_id` o `dupla.social_worker_id`. Pero un `team_contact` nunca guarda `professional_id` Y `dupla_id` a la vez: refleja el modo de asignación del `psychosocial_support` padre. |
| ¿Se muestra un loader mientras se decide el formulario? (formId llega en la misma respuesta) | PASO 2-3 | Sí — `pageLoading` en el frontend muestra un spinner de página completa hasta que la respuesta de `/load` llega (éxito o error). |
| Definir la ruta de "Ver remisión" al completar la sesión                                      | INTERFAZ | Ya existe: `/salvia/remision-psicosocial/:id` (pantalla "Detalle de Remisión"). Se usa como botón del overlay de completado, pero la lógica de guardado/transición de estado (incrementar `session_count`, actualizar `status`, etc.) queda pendiente para el evento de "guardar sesión" — no implementada en este evento E-01. |
| ¿Puede el profesional rellenar el Form de Cierre completando solo S1-S4 (Contacto + Atención Psicosocial, antes "Seguimiento"; tras los ajustes de Barreras y rebranding Jul 2026) con "Cerrar remisión = No" y volver luego para cerrar? | PASO 7 | Sí — el formulario de Cierre siempre carga completo (S1-S5); la sección S5 "Cierre" permanece oculta vía `visibility_condition` hasta que el profesional marque "Cerrar remisión = Sí" en S4. Esto ya está en el seed; la lógica de qué pasa al guardar con "No" queda para el evento de guardar. |

Nota: el `form_id` de un `team_contact` se fija una sola vez (la primera carga) y se
reutiliza en cargas posteriores mientras el contacto siga pendiente, evitando que la
sesión "salte" de formulario si `ya_hizo_primera_atencion` o `session_count` cambian
mientras el profesional todavía no ha completado esa sesión.
