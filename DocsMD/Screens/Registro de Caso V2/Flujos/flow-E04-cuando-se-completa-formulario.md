━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando se completa el formulario (todas las secciones respondidas)
   Tipo: Backend/Scheduled
   Función: processVictimCaseSubmission(submissionId, actorId)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: mismo mecanismo que processFollowUpSubmission — goroutine
                lanzada dentro de SaveSection cuando allAnswered == true,
                vía el dispatch existente OnEndFormSubmission (SEGÚN formID,
                form_service.go:1919-1923):

                  func (s *formService) OnEndFormSubmission(...) error {
                      switch formID {
                      case seguimientoFormID:
                          return s.processFollowUpSubmission(...)
                      case registroCasoFormID:              // NUEVO case
                          return s.processVictimCaseSubmission(...)
                      }
                  }

Archivo:       nuevo — propuesto src/salvia/service/victim_case_form_service.go
                (misma función mencionada en flow-E03, ahora en su etapa final)

A diferencia del diseño original de este flujo, el `victim_case` **ya existe**
(creado en Borrador por flow-E03 al guardar la Sección 1). Este evento NO
crea el caso — lo completa y lo activa.

INPUT: {
  submissionId:  UUID del FormSubmission   → propagado desde saveSection
  actorId:       i_code del usuario        → propagado desde saveSection
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — processVictimCaseSubmission
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Cargar el victim_case en Borrador vinculado a este submission

  caseId = DB.form_submission.FindVictimCaseID({ submissionId })

  SI no existe (no se creó el draft en flow-E03 — no debería pasar):
    → Retornar error (aborta — inconsistencia grave, se loggea con severidad alta)

  DB.victim_case.FindByICode({ iCode: caseId })
  SI victim_case_status != "bo":
    → Loggear "el caso ya fue activado o no está en borrador — se omite (idempotente)"
    → TERMINAR ejecución  // evita reprocesar si el goroutine se reintenta


PASO 2 — Construir answerMap completo (las 128 respuestas)

  DB.answers.FindDirectBySubmissionID({ submissionId })
  → answerMap = { [questionId]: value }


PASO 3 — Validar completitud

  (idéntico a lo ya documentado — recorre las preguntas required=true y
  visibles según las visibility_condition ya evaluadas)

  SI falta alguna respuesta requerida y visible:
    → Loggear advertencia con la lista de questionIds faltantes
    → TERMINAR ejecución (defensivo — dinamic-form ya no debería permitir
      llegar aquí sin todo lo requerido respondido)


PASO 4 — Calcular wasPartner y riskScore/riskLevel

  wasPartner = answerMap[qRelationshipAggressor] ∈ {'pi', 'ex'}
  riskScore, riskLevel = getRiskScore(mapearATamizajeForm2(answerMap, wasPartner))

  → Mismo gap ya señalado: fórmula duplicada entre este paso y el bloque 2
    del frontend (flow-E02) — solo para el banner en vivo, no afecta el
    resultado final que siempre se recalcula aquí como fuente de verdad.


PASO 5 — Proyectar las 128 respuestas sobre victim_case_form2 (UPDATE completo)

  salvia_daos.SetVictimCase({
    VictimCaseICode:     caseId,             // UPDATE, no INSERT — el caso ya existe
    VictimCaseNames:     answerMap[qNames],
    VictimCaseLastNames: answerMap[qLastNames],
    VictimCaseDocType:   answerMap[qDocType],
    VictimCaseDocNumber: answerMap[qDocNumber],
    VictimCaseForm2: {
      ... mapeo 1:1 completo de las 128 respuestas a VictimCaseForm2DTO,
          igual que la Sección 1 ya escrita en flow-E03, más las 7 secciones
          restantes ahora disponibles ...
      RiskScore: riskScore,
      RiskLevel: riskLevel,
    },
  })

  SI falla:
    → Retornar error (aborta — el caso queda en "bo", se puede reintentar
      manualmente o por un reintento del goroutine)


PASO 6 — Generar usuario para la víctima

  newPassword = generarPasswordAleatoria()
  hashedPassword, err = bcrypt.GenerateFromPassword([]byte(newPassword), 10)
  SI err:
    → Retornar error (aborta — sin usuario no se activa el caso)

  → DB.general_user.Create({
        login:    generarLoginDesdeDocumento(answerMap[qDocNumber]),
        password: hashedPassword,
    })
  SI falla:
    → Retornar error (aborta)
  → newUserID = general_user.id


PASO 7 — Activar el caso

  DB.victim_case.UpdateStatus({ iCode: caseId, status: "ra" })
  // "ra" = "enrutado aprobado" — mismo estado que SetVictimCase asigna hoy
  // de forma atómica al crear un caso como rol "op". No se necesita un
  // código "Activo" nuevo — "bo" (Borrador) es el único código nuevo.


PASO 8 — Generar calendario de seguimientos

  calendarInput = { RiskLevel: int(riskLevel) }
  followUps, calErr = FollowUpSvc.GenerateOrRecalculate(ctx, caseId, calendarInput)
  SI calErr:
    → Loggear advertencia WARN (no aborta — el caso ya está activo)


PASO 9 — Asignar equipo y agente

  team = resolverTeamPorRiskLevel(riskLevel)
  agentID = calcularAgente(...)  // Borda/dense-rank, igual que reasignarCaso en Seguimiento
  DB.victim_case_light.UpdateTeamAndAgent({ iCode: caseId, team, agentID })
  SI error:
    → Loggear advertencia WARN (no aborta)


PASO 10 — Registrar evento en el timeline

  DB.case_timeline_events.Create({
    case_id:  caseId,
    category: "General",
    type:     "Caso Activado",
    icon:     "folder-plus",
    color:    "#22c55e",
    event_user_id: actorId,
    date:     now(),
  })
  SI falla:
    → Loggear advertencia WARN (no aborta)


PASO 11 — Guardar el resultado para que el frontend lo recupere

  DB.form_submission_result.Upsert({   // tabla/columna nueva — ver GAPS
    submissionId,
    caseId,
    newUserLogin: login,
    newUserPass:  newPassword,          // texto plano, uso único — ver GAPS de seguridad
    readyAt:      now(),
  })

  → Ver flow-E05-cuando-emite-form-completed.md para cómo el frontend
    recupera esto vía GET.

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                          | Paso afectado |
|------------------------------------------------------------------------------------------------|---------------|
| Tabla/columna para almacenar el resultado (newUser/caseId) hasta que el frontend lo pida — necesita expiración/limpieza (contraseña en texto plano de un solo uso) | PASO 11 |
| Mapeo campo-por-campo completo de las 128 respuestas → VictimCaseForm2DTO (aquí se resume, falta el mapeo exhaustivo 1:1) | PASO 5 |
| Rollback si falla PASO 6 después de actualizar form2 en PASO 5 (form2 completo pero sin usuario, caso sigue en "bo") | PASO 5-6 |
| Confirmar que "ra" es semánticamente correcto para "caso activo recién completado" vs. crear un código nuevo dedicado | PASO 7 |
