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

IDs de preguntas relevantes:
  qEquipos             = "e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc"  // ¿Cuáles equipos? (multi-select)
  qMedidasEmergencia   = "1a36260c-33a4-4ebd-bffb-e387d7964b96"  // Medidas de emergencia (multi-select)
  qCriteriosPsico      = "71c42c4a-f640-47ad-b2c1-5d4c18480449"  // Criterios remisión psicosocial
  qCriteriosHombres    = "f7edf4fc-d1cd-4591-a358-31566c806365"  // Criterios remisión hombres
  qCriteriosEstab      = "28accaa6-99dc-4ec4-967b-f26045ad707c"  // Criterios remisión estabilización
  qServiciosDiscap     = "47b122b1-0151-4b58-a007-d4afb722c1b9"  // Servicios equipo discapacidad
  qCierraCaso          = "08950a38-3db3-4dc7-852c-3b06b4b1ed72"  // ¿Realiza cierre del caso?


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — processFollowUpSubmission
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Cargar el FollowUpV2

  DB.follow_up_v2.FindByFormSubmissionID({ submissionId })

  SI no existe:
    → Retornar error
    → TERMINAR ejecución

  SI fu.status == "REALIZADO":
    → Resolver nombre del actor (agentLightRepo.FindByICode)
    → DB.case_timeline_events.Create({ type: "Seguimiento Editado", icon: "pospuesto", color: "teal" })
    → TERMINAR ejecución  // idempotente — no reprocesa


PASO 2 — Construir answerMap

  DB.answers.FindDirectBySubmissionID({ submissionId })
  → answerMap = { [questionId]: value }   // mapa plano de todas las respuestas directas


PASO 3 — Procesar barreras (BarrierV2)

  DB.repeater_entries.FindBySubmissionIDAndGroupIDs({ submissionId, groupId: rgBarreras })
  → Por cada entry:
      Leer answers de la entry → entryMap
      sectorExplicit = entryMap[qBarreraSector]  // dropdown de sector
      Por cada pregunta de barreras (salud / justicia / proteccion):
        SI valor no vacío:
          sector = sectorExplicit || sectorFallback
          Por cada opción del CSV:
            → DB.barrier_v2.Create({ case_id, follow_up_id, sector, description, status: "OPEN" })
            → barrierCount++


PASO 4 — Procesar derivaciones a equipos

  equiposVal = answerMap[qEquipos]

  SI equiposVal está vacío → saltar este paso

  Por cada equipo en splitCSV(equiposVal):

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "atencion_psico"                                       │
  └─────────────────────────────────────────────────────────────┘

    REGLA DE EXCLUSIÓN:
    SI equiposVal contiene "medidas_emergencia":
      → Log: "atencion_psico omitida — incompatible con medidas_emergencia"
      → OMITIR (break)

    VALIDAR CRITERIOS:
    criteriosVal = answerMap[qCriteriosPsico]
    tieneCriterioObligatorio = criteriosVal contiene "criterio_obligatorio"

    Puntaje por criterio:
      conducta_suicida=3, interseccionalidad=2, sin_ruta=1,
      condiciones_territoriales=1, sin_acceso_psico=1, naturalizacion_vbg=1

    totalPuntos = suma de puntos de los criterios seleccionados

    SI NOT tieneCriterioObligatorio OR totalPuntos < 3:
      → Log: "derivacion psicosocial NO cumple criterios (obligatorio=X, puntos=N)"
      → OMITIR (break)

    → DB.psychosocial_support.Create({ case_id, follow_up_id, type: "derivacion", status: "ACTIVE" })
    → Log: "✅ derivacion creada -> atencion_psico (id=...)"
    → remisionCount++

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "atencion_hombres"                                     │
  └─────────────────────────────────────────────────────────────┘

    criteriosHombresVal = answerMap[qCriteriosHombres]

    SI criteriosHombresVal NO contiene "criterio_hombres":
      → Log: "derivacion atencion_hombres NO cumple criterio — omitida"
      → OMITIR (break)

    → DB.men_team_remision.Create({ case_id, follow_up_id, type: "derivacion", status: "ACTIVE" })
    → Log: "✅ derivacion creada -> atencion_hombres (id=...)"
    → remisionCount++

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "discapacidad"                                         │
  └─────────────────────────────────────────────────────────────┘

    serviciosVal = answerMap[qServiciosDiscap]

    SI serviciosVal está vacío:
      → Log: "derivacion discapacidad sin servicios seleccionados — omitida"
      → OMITIR (break)

    Por cada servicio en splitCSV(serviciosVal):
      → DB.discapacidad_remision.Create({ case_id, follow_up_id, service, status: "ACTIVE" })
      → Log: "✅ derivacion creada -> discapacidad servicio=X (id=...)"
      → remisionCount++

    Servicios posibles: apoyo_lsc | enfoque_discapacidad
    Se crea 1 registro por servicio seleccionado.

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "estabilizacion"                                       │
  └─────────────────────────────────────────────────────────────┘

    criteriosEstVal = answerMap[qCriteriosEstab]

    SI criteriosEstVal está vacío:
      → Log: "derivacion estabilizacion sin criterios seleccionados — omitida"
      → OMITIR (break)

    → DB.economic_stabilization.Create({ case_id, follow_up_id, type: "derivacion", status: "ACTIVE" })
    → Log: "✅ derivacion creada -> estabilizacion (id=...)"
    → remisionCount++


PASO 5 — Procesar medidas de emergencia

  medidasVal = answerMap[qMedidasEmergencia]

  SI medidasVal está vacío → saltar este paso

  Por cada medida en splitCSV(medidasVal):
    → DB.emergency_measure.Create({ case_id, follow_up_id, type: medida, status: "ACTIVE" })
    → Log: "✅ medida emergencia creada tipo=X (id=...)"
    → remisionCount++

  Valores posibles: alojamiento | transporte | alimentacion | vestuario | apoyo_psico | otras_me


PASO 6 — Marcar follow-up como REALIZADO

  DB.follow_up_v2.UpdateStatus({ id: fu.id, status: "REALIZADO" })


PASO 7 — Registrar evento en el timeline

  Resolver nombre del actor (agentLightRepo.FindByICode)

  description = "Seguimiento ejecutado. Se identificaron N barreras y se realizaron M remisiones a Equipos Salvia"

  DB.case_timeline_events.Create({
    case_id:       fu.case_id,
    category:      "Seguimientos",
    type:          "Seguimiento Ejecutado",
    icon:          "calendar-check",
    color:         "#22c55e",
    description:   description,
    event_user_id: actorId,
    actor_name:    actorName,
    follow_up_id:  fu.id,
    date:          now(),
  })

  SI falla el insert:
    → Loggear advertencia (no aborta — el seguimiento ya fue marcado REALIZADO)


PASO 8 — Cierre del caso (Sección 5)

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

      S3. DB.case_timeline_events.Create({
            case_id:  fu.case_id,
            category: "General",
            type:     "Cierre de Caso",
            icon:     "circle-xmark",
            color:    "red",
            date:     now(),
          })
          SI falla el insert:
            → Loggear advertencia (el cierre ya ocurrió — no se revierte)

      → FIN SUB-FLUJO

    SI CerrarCaso retorna error:
      → Loggear advertencia
      → Continuar (no aborta el flujo principal)

  SI answerMap[qCierraCaso] != "true":
    → No se cierra el caso

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                        | Paso afectado |
|--------------------------------------------------------------------------------------------|---------------|
| No hay rollback si el goroutine falla a mitad de ejecución (ej. falla al crear timeline)   | General       |
| La lógica de puntaje de criterios psicosociales está duplicada en frontend y backend       | PASO 4        |
| Agregar un criterio nuevo a qCriteriosPsico requiere cambio en código (no solo en BD)      | PASO 4        |
| salvia_dignidad aparece como opción en qEquipos pero no tiene lógica en el switch          | PASO 4        |
