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

IDs de preguntas relevantes (preguntas directas del form):
  qEquipos                 = "e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc"  // Derivaciones a equipos (multi-select)
  qMedidasEmergencia       = "1a36260c-33a4-4ebd-bffb-e387d7964b96"  // Medidas de emergencia (multi-select)
  qCriteriosPsico          = "71c42c4a-f640-47ad-b2c1-5d4c18480449"  // Criterios remisión psicosocial
  qCriteriosHombres        = "f7edf4fc-d1cd-4591-a358-31566c806365"  // Criterios remisión hombres
  qCriteriosEstabilizacion = "28accaa6-99dc-4ec4-967b-f26045ad707c"  // Criterios remisión estabilización
  qServiciosDiscapacidad   = "47b122b1-0151-4b58-a007-d4afb722c1b9"  // Servicios equipo discapacidad
  qNuevosHechosViolencia   = "69ccecbe-8fd5-44a4-901a-17084ea7134d"  // ¿Se registraron nuevos hechos? (boolean)
  qDescripcionHechos       = "f8453544-2a8c-461d-9986-302ee4719492"  // Descripción de hechos (text)
  qFechaHechos             = "7373aeab-b6d8-44e4-b1ce-4d40196f17bc"  // Fecha en que ocurrieron (date)
  qCierraCaso              = "08950a38-3db3-4dc7-852c-3b06b4b1ed72"  // ¿Realiza cierre del caso? (boolean)

IDs de repeater groups:
  rgBarreras              = "5fd3ecdc-2e5f-4b31-97ef-8a994580586a"  // Sección 2 — Nuevas barreras
  rgSeguimientoBarreras   = "b536f16c-67b3-4370-810d-7cc8c9d5463e"  // Sección 3 — Seguimiento a barreras activas

IDs de preguntas dentro del repeater de barreras (rgBarreras):
  qBarreraSector           = "f19378b6-55c5-4fdf-b765-7ebcc3978741" // Q1  Sector (dropdown: salud|justicia|proteccion)
  qBarreraSalud            = "5fc1f2af-cc30-41f0-aa31-4731e5cb674c" // Q2  Barreras Salud (multiple)
  qInstitucionSalud        = "e88fb2bf-7196-4bea-91e6-fe78ccdedac1" // Q3  Institución Salud (multiple)
  qOtraBarreraSalud        = "bebf6e6c-0200-4b53-886c-c01f591eaa62" // Q4  Otra barrera Salud (text)
  qBarreraJusticia         = "2bec977e-97c7-42d7-a00a-b536af8038eb" // Q5  Barreras Justicia (multiple)
  qInstitucionJusticia     = "4b4997fa-67a0-4479-852d-df9e7bfb2b3e" // Q6  Institución Justicia (multiple)
  qOtraBarreraJusticia     = "96f0c507-64c1-446b-b07d-83d5ad1d7172" // Q7  Otra barrera Justicia (text)
  qBarreraProteccion       = "66c9fc1e-9b5e-4ad4-999f-7aeb483d84dc" // Q8  Barreras Protección (multiple)
  qInstitucionProteccion   = "78474c82-61b9-4a0c-beca-439274813c03" // Q9  Institución Protección (multiple)
  qOtraBarreraProteccion   = "68f6bf06-a6a8-443f-9a1e-001148d4eb45" // Q10 Otra barrera Protección (text)
  qBarreraOtraInstitucion  = "a1573bc3-28f6-485e-8d8b-b3567db42ae3" // Q11 Nombre institución (text)
  qBarreraDepartamento     = "31c7f8ba-880e-4c9a-89f0-1a6a1e43b7b9" // Q12 Departamento (dropdown)
  qBarreraCiudad           = "c6f2d54a-3c61-4ae7-b2bf-50a884031fa4" // Q13 Ciudad (dropdown)
  qBarreraMunicipio        = "8225d03f-8de9-4ff0-9345-67b5e71bf02d" // Q14 Municipio (dropdown)
  qEstructuralInstitucional= "64754fb5-04a8-43e4-ac86-d60f0a84010b" // Q15 (multiple)
  qEstructuralEconomico    = "94dfc417-b59f-4afe-b63b-c6839a8dfecc" // Q16 (multiple)
  qEstructuralTerritorial  = "cf162686-0327-4b28-9639-e4d947272483" // Q17 (multiple)
  qEstructuralDiferencial  = "81275638-2837-4ecc-8686-675dfeab4f32" // Q18 (multiple)
  qBarreraFecha            = "4592d85f-8c11-4785-ad32-06aa810cc491" // Q19 Fecha (date)
  qBarreraFuncionario      = "33c3961e-4252-4b02-b491-615b11d6bf56" // Q20 Funcionario (text)
  qBarreraDescripcion      = "d2be610f-44eb-4526-b4ec-86eaaaba08c8" // Q21 Descripción (text)
  qBarreraGestion          = "572ad72a-8174-4ff3-9c56-5c8c65ac63ac" // Q22 Gestión (multiple, CSV)

IDs de preguntas dentro del repeater de seguimiento (rgSeguimientoBarreras):
  qSBPersiste              = "4325a514-332a-42a9-9a89-a1d249b87c71" // Q2 ¿Persiste la barrera? (boolean)
  qSBRespuestaInstitucional= "a234a906-7b02-4a43-b18d-b802d430d272" // Q3 Respuesta institucional (single)
  qSBGestion               = "c785a994-3a38-4337-bd96-67329144affa" // Q4 Gestión realizada (multiple)
  qSBActuaciones           = "c808b590-128c-43cf-af94-c191a4bc31ac" // Q5 Actuaciones (text)
  qSBCierra                = "7442394d-256a-45f4-a63a-ca1164e865c2" // Q6 ¿Cierra la barrera? (boolean)
  qSBMotivoCierre          = "d2d3cee8-b5d3-422f-8ae0-e5a40effc983" // Q7 Motivo de cierre (single)

IDs de preguntas — Sección 1 Valoración del Riesgo (usadas en reasignarCaso):
  qProtectores = "a0fdcf67-b05b-4d92-9c19-14753a32bbe3"  // Factores protectores (multiple)
  qRiesgos     = "ec5bb242-6f86-4c64-8b9f-afabd5a51878"  // Factores de riesgo (multiple)
  qExtremo     = "65f2d582-a39d-4c93-9519-1a20efbebb03"  // Factores de riesgo extremo (multiple)
  qConfirmHigh = "df7a0293-e2ca-49e4-88a8-66a70519a9e5"  // ¿Confirmar reasignación a riesgo alto? (boolean)
  qConfirmLow  = "7ec8d66d-7015-470a-a596-edebe574b5c2"  // ¿Confirmar reasignación a riesgo bajo? (boolean)


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


PASO 3 — Crear barreras nuevas (rgBarreras)

  DB.repeater_entries.FindBySubmissionIDAndGroupIDs({ submissionId, groupIds: [rgBarreras] })
  → Por cada entry del repeater:

      3.1 Leer answers de la entry → entryMap

          sector = entryMap[qBarreraSector]   // dropdown: "salud" | "justicia" | "proteccion"

          → DB.barrier_v2.Create({
                case_id:         fu.case_id,
                follow_up_id:    fu.id,
                created_by_id:   actorId,
                status:          "OPEN",
                sector:          sector,
                specific_barriers:     entryMap[qBarrera{sector}],       // preguntas Q2/Q5/Q8 según sector
                specific_institutions: entryMap[qInstitucion{sector}],   // preguntas Q3/Q6/Q9 según sector
                other_barrier_desc:    entryMap[qOtraBarrera{sector}],   // preguntas Q4/Q7/Q10 según sector
                institution_name:      entryMap[qBarreraOtraInstitucion],
                department_id:         entryMap[qBarreraDepartamento],
                city_id:               entryMap[qBarreraCiudad],
                town_id:               entryMap[qBarreraMunicipio],
                structural_institutional: entryMap[qEstructuralInstitucional],
                structural_economic:      entryMap[qEstructuralEconomico],
                structural_territorial:   entryMap[qEstructuralTerritorial],
                structural_differential:  entryMap[qEstructuralDiferencial],
                barrier_date:          entryMap[qBarreraFecha],
                official_dependency:   entryMap[qBarreraFuncionario],
                description:           entryMap[qBarreraDescripcion],
                management_actions:    entryMap[qBarreraGestion],
            })
          SI falla el insert:
            → Retornar error (aborta la ejecución)
          → barrierCount++
          → newBarrierID = barrier_v2.id

      3.2 Crear tareas (y oficios si aplica) por gestión de barrera (Q22)

          gestionRaw = entryMap[qBarreraGestion]   // CSV de opciones seleccionadas

          SI gestionRaw está vacío:
            → No crear tareas ni oficios → CONTINÚA al siguiente entry

          Opciones y su comportamiento:
            "orientacion_llamada"                → OMITIR (no genera ningún registro)
            "gestion_llamada"                    → Solo case_task  (type="gestion_llamada")
            "alerta_barreras"                    → Solo case_task  (type="comite_caso")
            "activacion_ruta_interinstitucional" → entity_letter + case_task (type="proyectar_oficio")
            "articulacion_institucional"         → entity_letter + case_task (type="proyectar_oficio")
            "escalamiento_organismo_control"     → entity_letter + case_task (type="proyectar_oficio")

          Por cada gVal en splitCSV(gestionRaw):

            SI gVal == "orientacion_llamada":
              → OMITIR → CONTINÚA

            SI gVal ∈ { "activacion_ruta_interinstitucional", "articulacion_institucional",
                        "escalamiento_organismo_control" }:
              → DB.entity_letter.Create({
                    barrier_id: newBarrierID,
                    case_id:    fu.case_id,
                    state:      "por_proyectar",
                    agent_id:   actorId,
                })
              SI falla el insert:
                → Loggear advertencia WARN (no aborta)
                → newLetterID = nil
              SI ok:
                → newLetterID = entity_letter.id
              → taskType = "proyectar_oficio"

            SI gVal == "alerta_barreras":
              → taskType = "comite_caso"

            SI gVal == "gestion_llamada":
              → taskType = "gestion_llamada"

            // Descripción enriquecida: label + sector + institución
            label = gestionLabels[gVal] + " (Barreras) | Sector: {sector} | Institución: {institution}"

            → DB.case_task.Create({
                  category:         "Barreras",
                  type:             taskType,
                  description:      label,
                  assigned_user_id: actorId,
                  status:           "ToDo",
                  case_id:          fu.case_id,
                  follow_up_id:     fu.id,
                  barrier_id:       newBarrierID,
                  entity_letter_id: newLetterID  // nil si no aplica oficio
              })
            SI falla el insert:
              → Loggear advertencia WARN (no aborta)


PASO 3b — Seguimiento a barreras activas (rgSeguimientoBarreras)

  DB.repeater_entries.FindBySubmissionIDAndGroupIDs({ submissionId, groupIds: [rgSeguimientoBarreras] })
  → Obtener activeIDs desde fu.ActiveBarrierIDs (CSV de IDs de barreras activas del follow-up)

  Por cada entry (indexada por posición):

      sbMap = answers de la entry (indexado por questionId)
      barrierID = activeIDs[idx]  // relación por posición, no por ID explícito

      SI sbMap[qSBCierra] == "true" Y barrierID != "":
        → DB.barrier_v2.UpdateStatus({ id: barrierID, status: "MANAGED" })
        SI error:
          → Log advertencia (no aborta)

      → DB.case_timeline_events.Create({
            case_id:     fu.case_id,
            follow_up_id: fu.id,
            category:    "Barreras",
            type:        "Seguimiento a Barrera",
            icon:        "shield-halved",
            color:       "#6366f1",
            description: buildBarrierFollowUpSummary(qSBPersiste, qSBRespuestaInstitucional, qSBActuaciones),
            event_user_id: actorId,
            date:        now(),
        })
      SI error:
        → Log advertencia (no aborta)


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
    SI falla el insert:
      → Retornar error (aborta la ejecución)
    → Log: "✅ derivacion creada -> atencion_psico (id=...)"
    → remisionCount++

    → Crear tarea de validación de la remisión (equipo psicosocial confirma si procede)
      DB.case_task.Create({
        category:                "Psicosocial",
        type:                    "validar_remision",
        description:             "Validar que remision a psicosocial es valida",
        assigned_user_id:        "",                        // sin asignar
        status:                  "ToDo",
        case_id:                 fu.case_id,
        follow_up_id:            fu.id,
        psychosocial_support_id: ps.id                       // id de la remisión recién creada
      })
      SI falla el insert:
        → Loggear advertencia WARN (no aborta — la remisión psicosocial ya se creó)

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "atencion_hombres"                                     │
  └─────────────────────────────────────────────────────────────┘

    criteriosHombresVal = answerMap[qCriteriosHombres]

    SI criteriosHombresVal NO contiene "criterio_hombres":
      → Log: "derivacion atencion_hombres NO cumple criterio — omitida"
      → OMITIR (break)

    → DB.men_team_remision.Create({ case_id, follow_up_id, type: "derivacion", status: "ACTIVE" })
    SI falla el insert:
      → Retornar error (aborta la ejecución)
    → Log: "✅ derivacion creada -> atencion_hombres (id=...)"
    → remisionCount++

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "discapacidad"                                         │
  └─────────────────────────────────────────────────────────────┘

    serviciosVal = answerMap[qServiciosDiscapacidad]

    SI serviciosVal está vacío:
      → Log: "derivacion discapacidad sin servicios seleccionados — omitida"
      → OMITIR (break)

    Por cada servicio en splitCSV(serviciosVal):
      → DB.discapacidad_remision.Create({ case_id, follow_up_id, service, status: "ACTIVE" })
      SI falla el insert:
        → Retornar error (aborta la ejecución)
      → Log: "✅ derivacion creada -> discapacidad servicio=X (id=...)"
      → remisionCount++

    Servicios posibles: apoyo_lsc | enfoque_discapacidad
    Se crea 1 registro por servicio seleccionado.

  ┌─────────────────────────────────────────────────────────────┐
  │ CASO: "estabilizacion"                                       │
  └─────────────────────────────────────────────────────────────┘

    criteriosEstVal = answerMap[qCriteriosEstabilizacion]

    SI criteriosEstVal está vacío:
      → Log: "derivacion estabilizacion sin criterios seleccionados — omitida"
      → OMITIR (break)

    → DB.economic_stabilization.Create({ case_id, follow_up_id, type: "derivacion", status: "ACTIVE" })
    SI falla el insert:
      → Retornar error (aborta la ejecución)
    → Log: "✅ derivacion creada -> estabilizacion (id=...)"
    → remisionCount++


PASO 5 — Procesar medidas de emergencia

  medidasVal = answerMap[qMedidasEmergencia]

  SI medidasVal está vacío → saltar este paso

  Por cada medida en splitCSV(medidasVal):
    → DB.emergency_measure.Create({ case_id, follow_up_id, type: medida, status: "ACTIVE" })
    SI falla el insert:
      → Retornar error (aborta la ejecución)
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
    follow_up_id:  fu.id,
    category:      "Seguimientos",
    type:          "Seguimiento Ejecutado",
    icon:          "calendar-check",
    color:         "#22c55e",
    description:   description,
    event_user_id: actorId,
    actor_name:    actorName,
    date:          now(),
  })

  SI falla el insert:
    → Loggear advertencia WARN (no aborta — el seguimiento ya fue marcado REALIZADO)


PASO 7b — Registrar nuevos hechos de violencia (Sección 1)

  SI answerMap[qNuevosHechosViolencia] == "true":

    descripcion = answerMap[qDescripcionHechos]
    fechaHechos = parsear answerMap[qFechaHechos] (formato "2006-01-02")
                  SI no parseable → usar now()

    → DB.case_timeline_events.Create({
          case_id:       fu.case_id,
          follow_up_id:  fu.id,
          category:      "General",
          type:          "Hechos del caso",
          icon:          icon_hechos_caso,
          color:         "#f87171",
          description:   descripcion,
          event_user_id: actorId,
          date:          fechaHechos,
      })
    SI falla el insert:
      → Loggear advertencia WARN (no aborta)


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
      → Loggear advertencia WARN
      → Continuar (no aborta el flujo principal)

  SI answerMap[qCierraCaso] != "true":
    → No se cierra el caso


PASO 9 — Reasignación de caso

  → Llamar reasignarCaso(ctx, fu, answerMap, actorId)

  ┌──────────────────────────────────────────────────────────────┐
  │  SUB-FLUJO: reasignarCaso                                    │
  └──────────────────────────────────────────────────────────────┘

  S1. Leer nivel de riesgo actual del caso
      DB.victim_case_form2.FindRiskLevelByICode({ iCode: fu.case_id })
      → currentLevel = 1 | 2 | 3 | 4
      SI error:
        → Loggear advertencia → TERMINAR sub-flujo (no aborta el flujo principal)

  S2. Leer y filtrar respuestas de factores
      protectores = splitCSV(answerMap[qProtectores]).filter(v != "ninguno")
      riesgos     = splitCSV(answerMap[qRiesgos]).filter(v != "ninguno")
      extremos    = splitCSV(answerMap[qExtremo]).filter(v != "ninguno")
      confirmHigh = answerMap[qConfirmHigh] == "true"
      confirmLow  = answerMap[qConfirmLow]  == "true"
      // "ninguno" se excluye antes de contar — no activa ninguna lógica

  S3. Determinar nuevo nivel según respuestas y nivel actual

      SI currentLevel ∈ [1, 2] (riesgo bajo):
        SI len(extremos) > 0:
          → newLevel = 4  // reasignación automática sin confirmación
        SI confirmHigh == true Y len(riesgos) >= 4:
          → newLevel = 3
        SI NINGUNA condición:
          → newLevel = 0  // sin reasignación

      SI currentLevel >= 3 (riesgo alto):
        SI confirmLow == true Y len(extremos) == 0 Y len(protectores) >= 3:
          → newLevel = 2
        SI NINGUNA condición:
          → newLevel = 0  // sin reasignación

      SI newLevel == 0:
        → Loggear decisión con detalle (extremos, riesgos, protectores, confirms)
        → TERMINAR sub-flujo

  S4. Actualizar nivel de riesgo en victim_case_form2
      DB.victim_case_form2.UpdateRiskLevelByICode({ iCode: fu.case_id, newLevel })
      SI error:
        → Loggear advertencia → TERMINAR sub-flujo

  S5. Reasignar calendario o agente según nuevo nivel

      → FollowUpV2Svc.ReasignarCalendario(ctx, fu.case_id, newLevel)

      ┌─────────────────────────────────────────────────────────┐
      │  Lógica interna de ReasignarCalendario                  │
      └─────────────────────────────────────────────────────────┘

      SI newLevel >= 3 (alto/extremo):
        team = "Riesgo alto"
        → Calcular nuevas fechas según riskMatrix[newLevel]
        → calcularAgente(ctx, fechas, team)  // Borda/dense-rank
        → DB.follow_up_v2.DeletePendingByCaseID({ caseID })
        → DB.follow_up_v2.BulkCreate(nuevos seguimientos)

      SI newLevel == 2 (riesgo bajo destino):
        team = "Riesgo bajo"
        → Obtener fechas de los PENDIENTE existentes
        → calcularAgente(ctx, fechas, team)  // Borda/dense-rank
        → DB.follow_up_v2.UpdateAgentForPendingByCaseID({ caseID, agentID })
        // No se borra ni regenera el calendario

      SI error en ReasignarCalendario:
        → Loggear advertencia → TERMINAR sub-flujo

  S6. Actualizar equipo y agente en victim_case
      team = "Riesgo alto" si newLevel >= 3, "Riesgo bajo" si newLevel == 2
      newAgentID = AgentID del primer follow_up_v2 PENDIENTE del caso
      DB.victim_case.UpdateTeamAndAgent({ iCode: fu.case_id, team, agentID: newAgentID })
      SI error:
        → Loggear advertencia WARN (no aborta — la reasignación ya ocurrió)

  S7. Registrar evento en el timeline
      DB.case_timeline_events.Create({
        case_id:       fu.case_id,
        follow_up_id:  fu.id,
        category:      "General",
        type:          "Reasignación de Caso",
        icon:          "arrows-rotate",
        color:         "#f59e0b",
        event_user_id: actorId,
        date:          now(),
      })
      SI falla el insert:
        → Loggear advertencia WARN (no aborta — la reasignación ya ocurrió)

      → FIN SUB-FLUJO reasignarCaso

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
| buildBarrierFollowUpSummary no está documentada — genera el texto del evento de timeline   | PASO 3b       |
| El type "validar_remision" de la nueva case_task no tiene manejo en `case-task-modal.js` ni `case-task-history.js` — hoy no se puede completar ni ver su detalle desde la UI | PASO 4 |
