━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se completa el formulario (todas las secciones respondidas)
   Tipo: Backend/Síncrono (no goroutine)
   Función: Activate(submissionId, actorId)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: el hook existente OnEndFormSubmission (SEGÚN formID), que ya
                disparaba processFollowUpSubmission para Seguimiento —
                Registro de Caso agrega su propio `case`:

                  func (s *formService) OnEndFormSubmission(ctx, formID, submissionID, actorID string) error {
                      switch formID {
                      case SeguimientoFormID:
                          return s.processFollowUpSubmission(...)
                      case service.RegistroCasoFormID:
                          return s.victimCaseFormSvc.Activate(ctx, submissionID, actorID)
                      }
                  }

                A diferencia del hook OnSectionUpdate (flow-E03), este SÍ es
                condicional: solo dispara cuando isAnswered es true para
                TODAS las secciones (mismo mecanismo ya usado por
                Seguimiento para decidir la redirección).

Archivo:       src/salvia/service/victim_case_form_service.go (Activate)

El `victim_case` ya existe en estado 'bo' (creado por UpdateCaseDraft desde
el primer saveSection — ver flow-E03) — este evento no lo crea, lo confirma
y lo activa.

INPUT: {
  submissionId:  UUID del FormSubmission   → propagado desde saveSection
  actorId:       i_code del usuario        → propagado desde saveSection
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — Activate
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PASO 1 — Cargar el victim_case vinculado a este submission

  iCode, found = formRepo.FindICodeBySubmissionID(ctx, submissionID)
  SI !found:
    → Error — no debería pasar: UpdateCaseDraft (flow-E03) ya lo crea desde
      el primer saveSection, con o sin datos dummy.


PASO 2 — Leer TODAS las respuestas y recalcular riskScore/riskLevel

  answers = formRepo.BuildAnswersByFieldKey(ctx, RegistroCasoFormID, submissionID)
  wasPartner = resolveWasPartner(answers)
  riskScore, riskLevel = computeRiskScore(answers, wasPartner)


PASO 3 — Re-proyectar sobre victim_case_form2 (defensivo/idempotente)

  formRepo.UpsertForm2FromAnswers(ctx, iCode, answers, riskScore, riskLevel)
  // Ya debería estar al día (UpdateCaseDraft lo dejó así en esta misma
  // request, justo antes de que isAnswered diera true en todas las
  // secciones) — se repite aquí solo como defensa, es idempotente.


PASO 4 — Activar el caso (bo → ra), refrescando datos de Sección 1 si cambiaron

  formRepo.MarkActive(ctx, iCode, answers)
  // UPDATE victim_case: status='ra' ("ra" = mismo código "enrutado
  // aprobado" ya usado por SetVictimCase hoy — no se creó un código nuevo
  // para "Activo"). Solo sobrescribe nombres/apellidos/doc/municipio de
  // atención (FieldKey(9,3)) si la respuesta real NO está vacía — nunca
  // pisa un valor bueno con un dummy en este paso.


PASO 5 — Generar calendario de seguimientos y asignar equipo/agente

  followUps, err = followUpV2Svc.GenerateOrRecalculate(ctx, iCode, { RiskLevel: riskLevel })
  SI err:
    → Loggear WARN (no aborta — el caso ya está activo)
  SI ok Y len(followUps) > 0:
    → team, agentID = followUps[0].Team, followUps[0].AgentID
    → caseRepo.UpdateTeamAndAgent(ctx, iCode, team, agentID)


PASO 6 — Registrar evento en el timeline

  caseTimelineRepo.Create({
    CaseID: iCode, Category: "General", Type: "Caso Activado",
    Icon: "folder-plus", Color: "#22c55e", EventUserID: actorID, Date: now(),
  })
  SI falla:
    → Loggear WARN (no aborta)

  → FIN EJECUCIÓN ✓ — ver flow-E05-cuando-emite-form-completed.md para
    cómo el frontend recupera caseId/login/password vía GET
    /api/v1/victim-case-forms/:submissionId/result (VictimCaseFormService.GetResult).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ GAPS del diseño original — YA RESUELTOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Gap original | Cómo se resolvió |
|---|---|
| Tabla/columna para el resultado (newUser/caseId) con expiración/limpieza | Se descartó la expiración — `victim_case.victim_case_new_user_credentials` (jsonb) guarda login+password en texto plano de forma PERMANENTE (decisión explícita del negocio: nunca se limpia). GetResult las lee de ahí, no de una tabla aparte. |
| Mapeo campo-por-campo de las 128 respuestas → VictimCaseForm2DTO | Resuelto vía la tabla estática `victimCaseFormFields` (victim_case_form_fields.go) — FieldKey → columna + tipo de conversión (KindText/Int/Date/Time/Timestamp/Enum1/EnumN/BoolEnum/Scale), aplicada genéricamente por UpsertForm2FromAnswers. |
| Rollback si falla generar el usuario después de actualizar form2 | Ya no aplica — el usuario/login/password se crean en flow-E03 (createDraftShell), no aquí. Activate no crea usuario nuevo. |
| ¿"ra" es semánticamente correcto para "Activo"? | Confirmado que sí — se reutiliza tal cual, sin crear un código dedicado. |
