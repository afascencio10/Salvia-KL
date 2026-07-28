━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se guarda la primera sección
   Tipo: Backend/Scheduled
   Función: createDraftVictimCase(submissionID, actorID)  [NUEVA]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: el mismo endpoint POST /api/v1/forms/saveSection que usa
                dinamic-form — en el punto EXACTO donde hoy se crea un
                FormSubmission nuevo (form_service.go:1744-1749, dentro de
                SaveSection):

                  if submissionID == "" {
                      fs := &models.FormSubmission{FormID: input.FormID}
                      s.submissionRepo.Create(ctx, fs)
                      submissionID = fs.ID
                  }

                Se agrega un dispatch SEGÚN input.FormID justo después de
                crear el submission — mismo patrón estructural que
                OnEndFormSubmission (SEGÚN formID) pero en un punto distinto
                del ciclo de vida, y de forma SÍNCRONA (no goroutine): el
                caseId resultante debe existir antes de responder al
                frontend, para poder devolverlo si la pantalla lo necesita
                más adelante.

Archivo:       nuevo — propuesto src/salvia/service/victim_case_form_service.go
Constante:     registroCasoFormID = "0a24ab30-3cfc-4861-b74d-65d21524bc00"

INPUT: {
  submissionId:    UUID recién creado                → fs.ID (PASO 1 de SaveSection)
  actorId:         i_code del usuario                 → propagado desde el controller
  directAnswers:   respuestas de la sección 1 recién enviadas → input.DirectAnswers
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — createDraftVictimCase (llamado desde SaveSection, síncrono)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Verificar que aplica (dispatch SEGÚN formID)

  SI input.FormID != registroCasoFormID:
    → No hacer nada — comportamiento genérico de saveSection sin cambios
    → TERMINAR

  SI input.FormID == registroCasoFormID:
    → CONTINÚA PASO 2


PASO 2 — Verificar idempotencia

  DB.form_submission.FindVictimCaseID({ submissionId })
  SI ya tiene un victim_case_id asociado (reintento del mismo submission):
    → Loggear "draft ya existe — se omite creación"
    → TERMINAR ejecución (no crea un segundo caso)


PASO 3 — Extraer los campos mínimos disponibles de la Sección 1

  answerMap = { [questionId]: value }  // solo las de la sección 1 recién guardada

  names       = answerMap[qNames]
  lastNames   = answerMap[qLastNames]
  docType     = answerMap[qDocType]
  docNumber   = answerMap[qDocNumber]

  SI falta cualquiera de estos 4 campos (deberían venir siempre, son
  required en Sección 1, pero se valida defensivamente):
    → Loggear advertencia y TERMINAR sin crear el draft
    → El caso se creará en el SIGUIENTE saveSection donde ya estén completos
      (no debería pasar en la práctica, dinamic-form ya exige estos campos
      required antes de dejar guardar la sección)


PASO 4 — Crear victim_case + victim_case_form2 en estado Borrador

  salvia_daos.SetVictimCase({
    VictimCaseNames:     names,
    VictimCaseLastNames: lastNames,
    VictimCaseDocType:   docType,
    VictimCaseDocNumber: docNumber,
    VictimCaseStatus:    "bo",   // NUEVO código — ver GAPS de form-registro-caso-v2-data.md (G-11)
    VictimCaseForm2: {
      // Resto de campos de Sección 1 ya disponibles (teléfono, residencia,
      // accesibilidad) se mapean aquí también; el resto del form2 queda
      // con sus defaults hasta que se completen más secciones.
      ... mapeo de los campos restantes de Sección 1 disponibles en answerMap ...
    },
  })

  SI falla:
    → Retornar error (aborta — el saveSection completo falla, el usuario
      ve el error de guardado en la Sección 1)
  → newCaseID = victim_case.i_code


PASO 5 — Vincular el form_submission al caso recién creado

  DB.form_submission.UpdateVictimCaseID({ submissionId, caseId: newCaseID })
  // Columna nueva necesaria en form_submission — ver GAPS


PASO 6 — Continuar el flujo normal de SaveSection

  → El resto de SaveSection sigue igual (guardar answers, responder
    formStructure/formSubmission/currentSection al frontend)
  → El frontend NO necesita hacer nada especial con newCaseID en este
    punto — solo se usa internamente hasta que el formulario se completa

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                          | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| Columna nueva `form_submission.victim_case_id` (o tabla de relación) — no existe hoy | PASO 2, 5 |
| Agregar código "bo" a VICTIM_CASE_STATUS y decidir su label exacto ("Borrador") | PASO 4 |
| Si las secciones 2 en adelante deben re-sincronizar victim_case_form2 en cada saveSection, o solo se escribe completo al final (E-04) — este flujo asume que SOLO Sección 1 se escribe aquí, y el resto se escribe de una vez al completar | General |
| Qué pasa si el usuario nunca vuelve a este formulario — el caso queda en "bo" indefinidamente; ¿necesita un job de limpieza/expiración? | General |
