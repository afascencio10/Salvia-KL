━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el backend procesa el seguimiento completado
   Tipo: Backend/Scheduled
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: E-04 de dinamic-form — goroutine lanzada cuando allAnswered == true
Archivo:       src/salvia/service/form_service.go → processFollowUpSubmission()

INPUT: {
  submissionId:  UUID del FormSubmission   → propagado desde el form component
  actorId:       i_code del usuario        → propagado desde el form component
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — processFollowUpSubmission
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Cargar el FollowUpV2

  DB.follow_ups.FindByFormSubmissionID({ submissionId })

  SI no existe:
    → Loggear error
    → TERMINAR ejecución


PASO 2 — Construir answerMap

  DB.answers.FindDirectBySubmissionID({ submissionId })
  → answerMap = { [questionId]: value }   // mapa plano de todas las respuestas directas


PASO 3 — Marcar el follow-up como REALIZADO

  DB.follow_ups.Update({
    id:           fu.id,
    status:       "REALIZADO",
    completed_at: now(),
  })


PASO 4 — Registrar evento de timeline

  DB.case_timeline_events.Create({
    case_id:       fu.case_id,
    category:      "General",
    type:          "Seguimiento Realizado",
    icon:          "clipboard-check",
    color:         "green",
    description:   "Seguimiento completado por el profesional",
    event_user_id: actorId,
    date:          now(),
  })


PASO 5 — Generar intentos del calendario de seguimiento

  SI el follow-up aún no tiene intentos generados:
    → Crear los registros de FollowUpAttempt en DB según la configuración
      del calendario (fechas programadas de contacto)


PASO 6 — Procesar barreras (BarrierV2)

  A partir de las respuestas de la sección de barreras en answerMap:
  → Crear o actualizar registros BarrierV2 asociados al caso


PASO 7 — Procesar medidas complementarias

  Según las respuestas en answerMap:
  → EmergencyMeasure:        crear/actualizar si aplica
  → PsychosocialSupport:     crear/actualizar si aplica
  → EconomicStabilization:   crear/actualizar si aplica


PASO 8 — Cierre del caso (si el profesional marcó cierre en Sección 5)

  qCierraCaso = "08950a38-3db3-4dc7-852c-3b06b4b1ed72"

  SI answerMap[qCierraCaso] == "true":

    → Llamar CasoCierreService.CerrarCaso({
          CaseICode:               fu.case_id,
          Motivo:                  answerMap["95fb963e-99de-4a1d-a170-8e30845d1f7d"],
          OtroMotivo:              answerMap["e259ff16-d049-41d1-b92a-b0f09ee6fa91"],
          Descripcion:             answerMap["50ab05f3-95d9-42e4-bfed-6abfcd0a8050"],
          AccionesInstitucionales: answerMap["0e7b61c8-9401-4429-81ed-61c8fc97808e"] == "true",
          ActorID:                 actorId,
      })

    ┌──────────────────────────────────────────────────────────────┐
    │  SUB-FLUJO: CasoCierreService.CerrarCaso                     │
    └──────────────────────────────────────────────────────────────┘

      S1. DB.victim_cases.FindByICode({ iCode: fu.case_id })
          SI caso ya tiene status == "cd":
            → Loggear "ya está cerrado — se omite"
            → TERMINAR sub-flujo (idempotente)

      S2. DB.victim_cases.UpdateStatus({ iCode: fu.case_id, status: "cd" })
          → victim_case_status = "cd"

      S3. DB.case_timeline_events.Create({
            case_id:       fu.case_id,
            category:      "General",
            type:          "Cierre de Caso",
            icon:          "circle-xmark",
            color:         "red",
            description:   buildCierreDescription(Motivo, OtroMotivo, Descripcion),
            event_user_id: actorId,
            date:          now(),
          })
          SI falla el insert:
            → Loggear advertencia (el cierre ya ocurrió — no se revierte)

      → FIN SUB-FLUJO

    SI CerrarCaso retorna error:
      → Loggear advertencia
      → Continúar (no aborta el resto del procesamiento)

  SI answerMap[qCierraCaso] != "true":
    → No se cierra el caso

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                              | Paso afectado |
|----------------------------------------------------------------------------------|---------------|
| Lógica interna exacta de generación de intentos del calendario                   | PASO 5        |
| Criterios para crear/actualizar vs. ignorar BarrierV2 según answerMap            | PASO 6        |
| Criterios para EmergencyMeasure, PsychosocialSupport, EconomicStabilization      | PASO 7        |
| Qué pasa si el goroutine falla a mitad de ejecución (no hay rollback)            | General       |
