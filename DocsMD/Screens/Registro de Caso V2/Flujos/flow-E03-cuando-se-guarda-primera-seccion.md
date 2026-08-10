━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se guarda CUALQUIER sección (no solo la primera)
   Tipo: Backend/Síncrono
   Hook genérico: OnSectionUpdate(formID, submissionID, actorID)
   Función específica: UpdateCaseDraft(submissionID, actorID)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️ Este flujo cambió de diseño respecto a la versión original de este
documento: en vez de un hook que solo dispara UNA VEZ al guardar la
Sección 1, se implementó un hook GENÉRICO — OnSectionUpdate — análogo a
OnEndFormSubmission pero que se dispara después de CADA saveSection, para
cualquier formulario. Registro de Caso es simplemente el primer `case` de
su dispatch. UpdateCaseDraft (la función específica de este formulario)
también cambió de diseño: en vez de posponer la creación de victim_case
hasta que Sección 1 esté completa, la crea/actualiza SIEMPRE, usando
valores dummy para lo que aún falte — nunca se pospone.

Disparado por: SaveSection (form_service.go), incondicionalmente, justo
                después de persistir las respuestas de la sección (no solo
                al crear el submission):

                  if err := s.OnSectionUpdate(ctx, input.FormID, submissionID, input.ActorID); err != nil {
                      log.Printf("[saveSection] OnSectionUpdate error ...")
                  }

                OnSectionUpdate (form_service.go) hace dispatch por formID:

                  func (s *formService) OnSectionUpdate(ctx, formID, submissionID, actorID string) error {
                      switch formID {
                      case service.RegistroCasoFormID:
                          return s.victimCaseFormSvc.UpdateCaseDraft(ctx, submissionID, actorID)
                      }
                      return nil
                  }

                Un error de OnSectionUpdate se loggea pero NO aborta
                saveSection — el guardado de la sección ya se completó; el
                draft del caso se pondrá al día en el siguiente saveSection.

Archivo:       src/salvia/service/victim_case_form_service.go (UpdateCaseDraft, createDraftShell)
                src/internal/repository/victim_case_form_repository.go (CreateDraft, UpdateDraftCoreFields, UpsertForm2FromAnswers)
Constante:     RegistroCasoFormID = "0a24ab30-3cfc-4861-b74d-65d21524bc00"


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — UpdateCaseDraft (llamada en CADA saveSection de este formulario)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Leer TODAS las respuestas disponibles del submission

  answers = formRepo.BuildAnswersByFieldKey(ctx, RegistroCasoFormID, submissionID)
  // indexado por FieldKey ("S{sección}Q{orden}"), no por questionID — evita
  // ambigüedad entre preguntas de Tamizaje pareja/no-pareja con texto idéntico.


PASO 2 — ¿Ya existe un victim_case para este submission?

  iCode, found = formRepo.FindICodeBySubmissionID(ctx, submissionID)
  // busca por victim_case.victim_case_form_submission_id (columna nueva,
  // ver form-registro-caso-v2-data.md) — no por una tabla de relación aparte.


PASO 3a — SI NO existe: createDraftShell (primera vez, cualquier sección)

  Para cada uno de los 5 campos obligatorios de Sección 1
  (names, lastNames, docType, docNumber, residenceTown):
    valor = answers[FieldKey(1, n)] SI no está vacío, SI NO dummy:
      names/lastNames → "Pendiente"
      docType         → "cc"
      docNumber       → "0000000000"
      residenceTown   → "00000000"

  → Genera login/password aleatorios, hashea con bcrypt
  → formRepo.CreateDraft(...) crea:
      security.general_user_profile + general_user + rel_role_general_user (rol "us")
      salvia.victim_case en status 'bo', con victim_case_form_submission_id = submissionID
  → formRepo.StoreCredentials(...) guarda login/password EN TEXTO PLANO en
    victim_case.victim_case_new_user_credentials (jsonb, nunca se limpia —
    ver GAPS de seguridad ya resueltos como decisión consciente del negocio)
  → iCode = recién creado

  Nunca retorna vacío / nunca pospone — a diferencia del diseño original.


PASO 3b — SI YA existe: UpdateDraftCoreFields (sincronizar, cada llamada)

  formRepo.UpdateDraftCoreFields(ctx, iCode,
    coalesce(answers[names],     "Pendiente"),
    coalesce(answers[lastNames], "Pendiente"),
    coalesce(answers[docType],   "cc"),
    coalesce(answers[docNumber], "0000000000"),
    coalesce(answers[residenceTown], "00000000"),
  )
  // UPDATE directo sobre victim_case — los valores reales sobrescriben el
  // dummy en cuanto llegan, sin esperar a que el formulario se complete.


PASO 4 — Calcular wasPartner y riskScore/riskLevel (con lo disponible hasta ahora)

  wasPartner = resolveWasPartner(answers)  // FieldKey(4,3) ∈ {'pi','ex'}
  riskScore, riskLevel = computeRiskScore(answers, wasPartner)
  // misma fórmula que el tamizaje ya usado hoy — se recalcula en cada
  // llamada con lo que haya, y de nuevo (como fuente de verdad) en Activate.


PASO 5 — Proyectar TODAS las respuestas disponibles sobre victim_case_form2

  formRepo.UpsertForm2FromAnswers(ctx, iCode, answers, riskScore, riskLevel)

  → Crea la fila si no existe, o la actualiza si ya existe (UPDATE +
    reinserción idempotente de relaciones multi-valor).
  → Las 55 columnas NOT NULL de victim_case_form2 (todas enum1/bool_enum)
    reciben el i_code "No" (categoría yes_no) como valor dummy si la
    pregunta correspondiente aún no tiene respuesta — así el INSERT nunca
    falla por NOT NULL y nunca se pospone. Se sobrescribe con la respuesta
    real en la siguiente sección que la responda.
  → Ya NO se captura ni se trata SQLSTATE 23502 como "incompleto, reintentar" —
    con los dummies, esa violación ya no debería ocurrir; si ocurre, es un
    bug real (columna NOT NULL sin mapear en victim_case_form_fields.go).

  → FIN EJECUCIÓN ✓ (el resto de SaveSection sigue igual — responde
    formStructure/formSubmission/currentSection al frontend sin cambios)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ GAPS del diseño original — YA RESUELTOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Gap original | Cómo se resolvió |
|---|---|
| Columna nueva `form_submission.victim_case_id` | Se descartó — en su lugar, `victim_case.victim_case_form_submission_id` (patrón análogo a follow_up_v2, pero como columna propia de victim_case, no de form_submission). |
| Código "bo" en VICTIM_CASE_STATUS | Agregado ("borrador") en config/Enums.go. |
| ¿Se re-sincroniza form2 en cada saveSection o solo al final? | Se re-sincroniza en CADA saveSection (UpsertForm2FromAnswers), no solo al completar. |
| Campos mínimos faltantes → ¿posponer? | NO — valores dummy, nunca se pospone (instrucción explícita del negocio). |
| Caso abandonado en "bo" indefinidamente | Sin cambios — no se implementó job de expiración; sigue siendo un gap abierto si se necesita a futuro. |
