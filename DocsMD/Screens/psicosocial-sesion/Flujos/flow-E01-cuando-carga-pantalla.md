━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  psicosocialId:  UUID de la remisión psicosocial  → route param :psicosocial_id
  userICode:      i_code del profesional           → plantilla Go ({{ .userICode }})
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — registrar_sesion.html
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Inicializar la app Vue

  Vue.createApp monta sobre #psicosocial-app
  → Se ejecuta mounted() → llama a loadPsicosocial()


PASO 2 — Llamar al backend para cargar la sesión

  GET /api/v1/psicosocial-support/{psicosocialId}/load?agent_id={userICode}

  SI respuesta no ok (status != 2xx):
    → this.loadError = data.error || 'Error al cargar la sesión psicosocial'
    → TERMINAR ejecución   // DinamicForm no se monta; se muestra ErrorAlert

  SI respuesta ok:
    → data = JSON parseado { formId, submissionId, victimInfo, psicosocialState, formState }
    → CONTINÚA PASO 3


PASO 3 — Poblar estado de la pantalla

  formId           = data.formId           // UUID del formulario seleccionado por el backend
  submissionId     = data.submissionId     // UUID del form_submission activo
  victimInfo       = data.victimInfo
  psicosocialState = data.psicosocialState  // { yaHizoPrimerContacto, yaHizoPrimeraAtencion, sessionCount, status }
  formState        = data.formState         // {} — objeto vacío; skipContact ya no se usa

  // psicosocialState se usa para:
  //   - Calcular el label del badge "Sesión X de 6" en la PsicosocialCard
  //   - Mostrar el FormTypeBadge con el nombre del formulario cargado
  //   - Derivar el texto del CompletedOverlay al terminar


PASO 4 — Calcular canEdit

  SI psicosocialState.status === 'cerrado':
    → canEdit = false
    → TERMINAR cálculo   // La remisión está cerrada: solo lectura

  SI cualquier otro status:
    → canEdit = true


PASO 5 — DinamicForm se monta con los props calculados

  DinamicForm recibe:
    :form-id       = formId        (UUID devuelto por el backend — varía por sesión)
    :submission-id = submissionId  (del team_contact activo)
    :can-edit      = canEdit       (calculado en PASO 4)
    :form-state    = formState     (objeto reactivo — actualmente vacío, reservado para extensiones futuras)

  Nota: La visibilidad de S2 en el Form de Primer Contacto (Primera Atención) y de S3 en el
  Form de Cierre es gestionada internamente por DinamicForm a través de visibility_conditions
  en la base de datos. No se necesita ningún flag externo en formState.

  → El componente renderiza el formulario en modo lectura o edición según canEdit

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — GET /api/v1/psicosocial-support/:id/load
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


INPUT: {
  psicosocialId:  UUID de la remisión  → path param :id
  agentId:        i_code del usuario   → query param agent_id
}


PASO 6 — Cargar la remisión psicosocial

  DB.psychosocial_support.FindByID({ psicosocialId })

  SI no existe:
    → 404 { error: "remisión psicosocial no encontrada" }
    → TERMINAR

  SI ps.ProfessionalID != agentId  Y  ps.DuplaID no contiene al agente:
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

  SI ps.YaHizoPrimeraAtencion == true  Y  ps.SessionCount >= 1  Y  ps.SessionCount < 3:
    → formKey = SEGUIMIENTO
    → formId  = FORM_ID_SEGUIMIENTO

  SI ps.YaHizoPrimeraAtencion == true  Y  ps.SessionCount >= 3:
    → formKey = CIERRE
    → formId  = FORM_ID_CIERRE
    // S3 (Cierre) permanece oculta hasta que el profesional marque
    // "Cerrar remisión = Sí" en S2; DinamicForm gestiona esa visibilidad.


PASO 8 — Resolver el team_contact activo

  Buscar un team_contact pendiente para este psicosocial_id:
    DB.team_contact.FindOne({
      psicosocial_id   = psicosocialId,
      is_completed     = false,
      is_psico_session = true,
      deleted_at IS NULL
    })

  SI existe un team_contact pendiente:
    → tc = registro encontrado

  SI no existe:
    → Crear nuevo team_contact:
         case_id          = ps.CaseID
         psicosocial_id   = psicosocialId
         professional_id  = agentId
         dupla_id         = ps.DuplaID
         is_completed     = false
         is_psico_session = true
         created_at       = NOW()
    → tc = nuevo registro


PASO 9 — Resolver o crear el FormSubmission

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

  200 {
    formId:      "{UUID del formulario seleccionado en PASO 7}",
    submissionId: "{UUID del form_submission}",
    victimInfo: {
      names, lastNames, phone, genderIdentity, townName, riskLevel, caseICode
    },
    psicosocialState: {
      yaHizoPrimerContacto:  ps.YaHizoPrimerContacto,
      yaHizoPrimeraAtencion: ps.YaHizoPrimeraAtencion,
      sessionCount:          ps.SessionCount,
      status:                ps.Status
    },
    formState: {}
    // formState ya no incluye skipContact. La visibilidad de secciones
    // es gestionada completamente por DinamicForm a través de
    // salvia.visibility_condition.
  }

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CONSTANTES DE FORMULARIO (backend Go)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Generadas al ejecutar seed_psicosocial.sql. Los UUIDs se capturan del RAISE NOTICE y
se fijan como constantes en un archivo Go (ej: src/internal/constants/psicosocial_forms.go):

  const (
    FORM_ID_PRIMER_CONTACTO   = "UUID-del-seed"
    FORM_ID_PRIMERA_ATENCION  = "UUID-del-seed"
    FORM_ID_SEGUIMIENTO       = "UUID-del-seed"
    FORM_ID_CIERRE            = "UUID-del-seed"
  )


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                           | Paso afectado |
|-----------------------------------------------------------------------------------------------|---------------|
| Confirmar si un profesional de dupla puede cargar la sesión aunque no sea el `professional_id` directo | PASO 6 |
| ¿Se muestra un loader mientras se decide el formulario? (formId llega en la misma respuesta) | PASO 2-3      |
| Definir la ruta de "Ver remisión" al completar la sesión                                      | INTERFAZ      |
| ¿Puede el profesional rellenar el Form de Cierre completando solo S1-S2 (Seguimiento) con "Cerrar remisión = No" y volver luego para cerrar? | PASO 7 |
