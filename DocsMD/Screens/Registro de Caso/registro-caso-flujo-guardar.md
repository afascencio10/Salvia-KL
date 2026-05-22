# Registro de Caso — Flujo Logico: Guardar

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona boton Guardar
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  data.victimCase:             objeto completo del caso          → estado Vue (formulario)
  data.victimCase.form2:       sub-objeto con todos los campos   → estado Vue (formulario)
  home.config.globalProperties.latitude:   coordenada GPS        → navigator.geolocation (obtenida en mounted)
  home.config.globalProperties.longitude:  coordenada GPS        → navigator.geolocation (obtenida en mounted)
  action:                      "save"                            → parametro hardcodeado del boton
}

═══════════════════════════════════════════
 FRONTEND (Vue)
═══════════════════════════════════════════

PASO 1 — Ejecutar submit('save')
  El usuario hace click en el boton "Guardar" (btn-success).
  Se invoca el metodo submit con action = "save".

PASO 2 — Verificar accion
  SI action == "back":
    → Redirigir a nav_rules["default"] con location.assign
    → TERMINAR ejecucion

  SI action == "save":
    → CONTINUA FLUJO GENERAL

PASO 3 — Adjuntar coordenadas GPS al caso
  data.victimCase.livingLatitude  = home.config.globalProperties.latitude
  data.victimCase.livingLongitude = home.config.globalProperties.longitude

PASO 4 — Enviar POST al backend
  newEntity(salviaFormPath, {victimCase: data.victimCase}, callback)

  [POST] /salvia/casos
  // o [POST] /salvia/casos/:id  (si viene de un contacto previo)

  payload:
  {
    "victimCase": data.victimCase    // origen: estado Vue completo
  }

  → Esperar respuesta del backend (el callback maneja status)

PASO 5 — Procesar respuesta del backend
  Ocultar successOverlay (display = 'none')

  SEGUN status:
    CASO 200:
      → data.victimCase.newUser = response  (contiene login y password)
      → Mostrar finishedOverlay (display = 'flex') con credenciales
      → TERMINAR ejecucion frontend

    CASO 400:
      → this.errors = response  (errores campo por campo)
      → Vue reactivamente muestra errores inline en cada campo
      → El watcher de errors construye lista para aria-live (accesibilidad)
      → TERMINAR ejecucion frontend

    DEFAULT:
      → Mostrar alert("Se produjo un error interno, por favor vuelva a intentarlo")
      → TERMINAR ejecucion frontend


═══════════════════════════════════════════
 BACKEND (Go — Facade + Controller)
═══════════════════════════════════════════

PASO 6 — Verificar sesion y permisos (VictimCasePOST)
  Obtener sesion del contexto Gin:
    sessionID = session.Get("userData")
    s = utils.GetCommonSession(sessionID)

  Verificar permiso "set_victim_case" para s.CurrentRole

  SI no tiene permiso:
    → Retornar 403 Forbidden
    → TERMINAR

  SI tiene permiso:
    → Leer body de la request
    → Llamar salvia_ctrl.SetVictimCase(body, session, dbConfig, dbServerConfig)
    → CONTINUA FLUJO GENERAL

PASO 7 — Parsear JSON de entrada (SetVictimCase)
  Parsear body JSON a VictimCaseRequest usando GetDTOMap y JSONToStruct
  Asignar defaults:
    SetVictimCaseDefaults(victimCase, SQL_INSERT, session)
    SetVictimCaseForm2Defaults(victimCase.form2, SQL_INSERT)

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Validacion del formulario   │
└─────────────────────────────────────────┘

  S1. Validar campos nivel victimCase (required):
      • names:       required, string
      • lastNames:   required, string
      • docType:     required, enum
      • docNumber:   required, string

  S2. Validar campos nivel form2 (~80 campos):
      Campos required (true): phone, factsDescription, factsDate, factsStartTime,
        factsTownCode, factsZone, factsAddress, scenarioViolence, reportedPreviously,
        recurrenceAggression, numAgressors, proximityPrincipalAggressor,
        relationshipWithPresumedAggressor, economicallyDependent, birthDate,
        physicalMentalSensoryDifficulties, nationality, migrationCondition,
        genderIdentity, sexualOrientation, assignedSexAtBirth, ethnicAffiliation,
        campesinoRecognition, maritalStatus, lastEducationLevel, occupation,
        incomeGenerationMethod, housingTenancyForm, housingStratum, currentlyPregnant,
        residenceTownCode, residenceAddress, residenceZone, personWithDisability,
        requireLanguageInterpreter, violenceMotivatedByGender, attentionWasAppropriate,
        aggressorOccupation, allowsEasyReport, salivaManagementExplanation,
        6 preguntas comunes tamizaje
      Campos opcionales (false): identityName, aggressorNames, aggressorDocType,
        aggressorDocNumber, aggressorAddress, aggressorPhone, aggressorGenderIdentity,
        supportContact*, activitiesUnable*, preguntas especificas tamizaje,
        specifiedNationality, indigenousPeople, employmentRelationship, approxStartAsp,
        workplaceSectorOccurrence

  S3. Verificar enums unicos
      getAndVerifyVictimCaseEnums: cada campo enum single se valida contra
      la tabla victim_case_form2_enums para confirmar que el codigo existe

  S4. Verificar enums multiples
      getAndVerifyVictimCaseEnumsMultiple: cada campo enum multiple se valida
      (typeViolenceExperienced, subtypeViolenceExperienced, scopeOfViolence,
       adjustmentsGBV, speciallyProtectedPopulation, actionPlan, hasDependents,
       whoReportTo, law1996, aspMode, reasonASP)

  S5. Validar campos condicionales:
      SI incomeGenerationMethod.code == "pr":
        → approxStartAsp es required (validar fecha no vacia)
        SI approxStartAsp vacio:
          → Agregar error
        → CONTINUA

      SI NO:
        → No validar approxStartAsp
        → CONTINUA

      SI requireLanguageInterpreter.code == "y":
        → languageAssistance es required (validar longitud min/max)
        SI languageAssistance vacio o fuera de rango:
          → Agregar error
        → CONTINUA

      SI NO:
        → No validar languageAssistance
        → CONTINUA

  S6. Validar fechas no futuras:
      SI factsDate > hoy:
        → Agregar error "fecha futura"
      SI birthDate > hoy:
        → Agregar error "fecha futura"

  S7. Evaluar errores acumulados
      SI collectedErrors tiene errores:
        → Retornar 400 con JSON de errores por campo
        → TERMINAR

      SI NO:
        → CONTINUA FLUJO GENERAL

  → FIN SUB-FLUJO → CONTINUA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Calculo de riesgo (backend) │
└─────────────────────────────────────────┘

  S1. Calcular riskScore y riskLevel en el backend con getRiskScore:
      partnerValidated = (relationshipWithPresumedAggressor.code == "pi" || "ex")
      riskScore = suma de respuestas "y" en las 6 comunes

      SI partnerValidated:
        → Sumar 18 preguntas de pareja
        → Nivel: 0-4=1(bajo), 5-8=2(moderado), 9-15=3(alto), 16-24=4(extremo)

      SI NO (no pareja):
        → Sumar 14 preguntas de no-pareja
        → Nivel: 0-2=1(bajo), 3-5=2(moderado), 6-8=3(alto), 9-20=4(extremo)

  S2. Comparar riskScore del cliente vs backend
      SI riskScore_cliente != riskScore_backend:
        → Retornar 400 con error "victim_case_form2_risk_score_client_error"
        → TERMINAR (proteccion contra manipulacion del cliente)

      SI coinciden:
        → Asignar riskLevel calculado al form2
        → CONTINUA FLUJO GENERAL

  → FIN SUB-FLUJO → CONTINUA FLUJO GENERAL

PASO 8 — Buscar contacto previo (si aplica)
  SI victimContact.iCode != "":
    → Consultar tabla victim_contact por iCode
    SI existe:
      → Asignar victimContactId al caso
    SI NO:
      → Ignorar (el caso se crea sin contacto previo)
    → CONTINUA FLUJO GENERAL

  SI NO (iCode vacio):
    → CONTINUA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Gestion de usuario          │
└─────────────────────────────────────────┘

  S1. Buscar perfil de usuario por docType + docNumber
      Consultar tabla general_user_profile:
        filter: { docType: victimCase.docType, docNumber: victimCase.docNumber }

  S2. Evaluar si el usuario existe

      SI perfil existe:
        → Buscar general_user asociado al perfil
        → usuario_existente = true

      SI NO existe:
        → usuario_existente = false

  S3. Iniciar transaccion de BD
      db.StartTransaction(connData)
      SI error:
        → Retornar 500
        → TERMINAR

  S4. Generar password aleatorio
      newPassword = getRandomPassword(4)  // 4 caracteres aleatorios

  S5. Crear o actualizar usuario

      SI usuario_existente == false:
        → Crear GeneralUserProfile con: docNumber, docType, gender, names, lastNames, nick
        → INSERT en general_user_profile
        SI error:
          → ROLLBACK → Retornar 500 → TERMINAR

        → Crear GeneralUser con: lang, status="e", login=getRandomUserName(6), password=bcrypt(newPassword)
        → INSERT en general_user
        SI error:
          → ROLLBACK → Retornar 500 → TERMINAR

        → Obtener rol "us" de la tabla role
        → INSERT en rel_role_general_user (asociar usuario con rol "us")
        SI error:
          → ROLLBACK → Retornar 500 → TERMINAR

      SI usuario_existente == true:
        SI usuario tiene rol "us":
          → Resetear password: bcrypt(newPassword)
          → UPDATE general_user (nueva password)
          SI error:
            → ROLLBACK → Retornar 500 → TERMINAR

        SI NO tiene rol "us":
          → No cambiar password
          → newPassword = "Usa tu contraseña actual" (mensaje traducido)

  S6. Asignar usuario al caso
      victimCase.generalUser = user.iCode
      victimCase.newUser = user  (incluye login y password para la respuesta)

  → FIN SUB-FLUJO → CONTINUA FLUJO GENERAL

PASO 9 — Determinar dueño del caso (CaseOwner)
  Consultar tabla case_owner:
    filter: { generalUser: session.userICode }
  SI error:
    → ROLLBACK → Retornar 500 con "error cargando propietario"
    → TERMINAR

  SI session.currentRole == "op":
    → victimCase.approvedBy = owner  (auto-aprobacion)
  → CONTINUA FLUJO GENERAL

PASO 10 — Resolver lugar de atencion
  SI victimCase.townCode == "" (no se eligio lugar de atencion):
    → victimCase.townCode = form2.residenceTownCode  (usar residencia como fallback)

  SI NO:
    → Mantener townCode seleccionado
  → CONTINUA FLUJO GENERAL

PASO 11 — Insertar caso en BD
  INSERT en victim_case:
    salvia_daos.SetVictimCase(victimCase)
  SI error:
    → ROLLBACK → Retornar 500 → TERMINAR
  → CONTINUA FLUJO GENERAL

PASO 12 — Insertar formulario form2
  form2.victimCase = victimCase  (asociar form2 al caso recien creado)
  INSERT en victim_case_form2:
    salvia_daos.SetVictimCaseForm2(form2)
  SI error:
    → ROLLBACK → Retornar 500 → TERMINAR
  → CONTINUA FLUJO GENERAL

PASO 13 — Insertar enums multiples del form2
  setVictimCaseEnumsMultiple: inserta las relaciones many-to-many entre form2
  y cada grupo de enums seleccionados (typeViolenceExperienced, scopeOfViolence,
  adjustmentsGBV, speciallyProtectedPopulation, actionPlan, hasDependents, etc.)
  INSERT multiples en tablas rel_victim_case_form2_*
  SI error:
    → ROLLBACK → Retornar 500 → TERMINAR
  → CONTINUA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Asignacion de ruta          │
│  (entity branches / momentos)           │
└─────────────────────────────────────────┘

  ↺ PARA CADA momentCode EN victimCase.entityBranches:
    ↺ PARA CADA sector EN sectors:
      ↺ PARA CADA branchICode EN entities:
        SI branchICode != "":
          S1. Obtener EntityBranch por iCode
              GetEntityBranchByICode(branchICode)
              SI status != 200:
                → ROLLBACK → Retornar status → TERMINAR

          S2. Crear Moment
              moment = {
                code:           momentCode,           // origen: clave del mapa entityBranches
                entityBranch:   branch,               // origen: paso S1
                approvalSource: MOMENT_APPROVAL_SOURCE["MANUAL"],
                victimCase:     victimCase             // origen: caso recien creado
              }

          S3. Insertar Moment en BD
              INSERT en moment:
                salvia_daos.SetMoment(moment)
              SI error:
                → ROLLBACK → Retornar 500 → TERMINAR

  → FIN SUB-FLUJO → CONTINUA FLUJO GENERAL

PASO 14 — Asociar dueño al caso
  Crear relacion case_owner ↔ victim_case:
    rel = {
      caseOwner:  owner.id,          // origen: PASO 9
      victimCase: victimCase.id       // origen: PASO 11
    }
  INSERT en rel_case_owner_victim_case
  SI error:
    → ROLLBACK → Retornar 500 → TERMINAR
  → CONTINUA FLUJO GENERAL

PASO 15 — Incrementar contador de casos del dueño
  owner.numCases = owner.numCases + 1
  UPDATE case_owner
  SI error:
    → ROLLBACK → Retornar 500 → TERMINAR
  → CONTINUA FLUJO GENERAL

PASO 16 — Co-autoria para rol "et"
  SI session.currentRole == "et":
    → Buscar todos los CaseOwners con rol "op" activos
      GetCaseOwnersByActiveUsersAndRoleCode("op")
    SI error o lista vacia:
      → ROLLBACK → Retornar 500 → TERMINAR

    → Crear segunda relacion case_owner ↔ victim_case con opOwners[0]
    INSERT en rel_case_owner_victim_case
    SI error:
      → ROLLBACK → Retornar 500 → TERMINAR

  SI NO (no es "et"):
    → No hacer nada
  → CONTINUA FLUJO GENERAL

PASO 17 — Actualizar columna informativa de funcionarios
  UpdateVictimCaseOwnersAndRolesByVictimCaseId(victimCase.id)
  (Actualiza una columna desnormalizada con los nombres y roles de todos los owners)
  SI error:
    → ROLLBACK → Retornar 500 → TERMINAR
  → CONTINUA FLUJO GENERAL

PASO 18 — Commit de la transaccion
  db.CommitTransaction(connData)
  SI error:
    → Retornar 500 → TERMINAR

  // Post-commit: se lanzan 2 goroutines independientes (no bloquean la respuesta):
  //   1. SUB-FLUJO: Generar calendario de seguimientos (FollowUpSvc)
  //   2. SUB-FLUJO: Evento de hechos del caso (CaseTimelineRepo)
  → CONTINUA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Generar calendario de       │
│  seguimientos (asincrono — goroutine)   │
└─────────────────────────────────────────┘

  NOTA: Se ejecuta DESPUES del commit en una goroutine separada.
  Si falla, NO deshace la creacion del caso.

  SI FollowUpSvc != nil:
    → Lanzar goroutine:

    S1. Construir input del calendario
        calendarInput = {
          riskLevel: riskLevel,    // origen: SUB-FLUJO Calculo de riesgo
          agentID:   "",           // vacio — se asignara automaticamente en S3
          team:      ""            // se define por nivel de riesgo en S2
        }

    S2. Determinar equipo por nivel de riesgo (GenerateOrRecalculate)
        SI riskLevel >= 3 (alto o extremo):
          → input.Team = "Riesgo alto"
        SI NO (bajo o moderado):
          → input.Team = "Riesgo bajo"

    S3. Calcular fechas programadas (computeScheduledDates)
        ↺ PARA CADA offset EN riskMatrix[riskLevel]:
          SI offset == 0 Y riskLevel == 4 (Extremo S1):
            → fecha = ahora + 4 horas
          SI NO:
            → fecha = hoy + offset dias

        riskMatrix:
          Extremo (4): offsets [0, 1, 2, 3, 15, 30]  → 6 seguimientos
          Alto (3):    offsets [1, 3, 15, 30]          → 4 seguimientos
          Moderado (2): offsets [2, 15, 30, 45]        → 4 seguimientos
          Bajo (1):    offsets [5, 15, 30, 60]          → 4 seguimientos

        → scheduledDates = lista de fechas calculadas

  ┌─────────────────────────────────────────┐
  │  SUB-FLUJO: Auto-asignacion de agente   │
  │  (calcularAgente — Borda / dense-rank)  │
  └─────────────────────────────────────────┘

    S3.1. Obtener agentes del equipo
          agentRepo.FindAllByRoleAndTeam("ro", input.Team)
          → agents = lista de agentes con rol "ro" del equipo
          SI agents vacio:
            → Retornar error "no hay agentes ro en el equipo"
            → AgentID queda vacio (NULL en BD, asignacion diferida)
            → CONTINUA FLUJO GENERAL

    S3.2. Obtener matriz de carga en una sola query (getDateMatrix)
          repo.FindWorkloadByDates(team, scheduledDates)
          → fullMatrix = map[fechaStr] → map[agentID] → count
          Ejemplo:
            "2026-05-22": { "agent-A": 3, "agent-B": 1 }
            "2026-05-23": { "agent-A": 2, "agent-B": 4 }

    S3.3. Obtener carga global por agente (tiebreaker final)
          repo.FindGlobalWorkloadByTeam(team)
          → globalLoad = map[agentID] → total_seguimientos_pendientes

    S3.4. Calcular dense-rank por fecha
          ↺ PARA CADA fecha EN scheduledDates:
            → Obtener carga de cada agente en esa fecha (0 si no tiene)
            → Ordenar cargas de menor a mayor
            → Asignar posicion dense-rank (menor carga = posicion 1)
            → Sumar posicion al acumulado de cada agente (agentPosSum)
            → Sumar carga al acumulado por fechas (agentDateLoad)

          Ejemplo con 2 fechas y 3 agentes:
            Fecha 1: A=3 (pos 2), B=1 (pos 1), C=3 (pos 2)
            Fecha 2: A=2 (pos 1), B=4 (pos 2), C=2 (pos 1)
            posSum:  A=3,          B=3,          C=3

    S3.5. Seleccionar agente con criterios en cascada
          avg_pos(agente) = agentPosSum / n_fechas

          Criterio 1: menor avg_pos
          SI empate:
            → Criterio 2: menor agentDateLoad (carga sumada en las fechas del caso)
          SI sigue empate:
            → Criterio 3: menor globalLoad (carga total pendiente del equipo)

          → bestID = agente seleccionado
          → input.AgentID = bestID

    SI calcularAgente retorna error:
      → Log warning "AutoAsignacion: {error} — se usara el agentID del input"
      → AgentID queda vacio → se persiste como NULL (asignacion diferida por supervisor)
      → CONTINUA FLUJO GENERAL

    SI calcularAgente retorna exito:
      → input.AgentID = agente seleccionado
      → CONTINUA FLUJO GENERAL

  → FIN SUB-FLUJO AUTO-ASIGNACION

    S4. Construir seguimientos (buildFollowUps)
        ↺ PARA CADA i EN scheduledDates:
          followUp = {
            caseID:       victimCase.iCode,
            agentID:      input.AgentID (o NULL si vacio),
            riskLevel:    riskLevelStr,
            scheduledDate: scheduledDates[i],
            order:        i + 1 (S1, S2, S3...),
            status:       "pending",
            team:         input.Team
          }
        → INSERT bulk en follow_up_v2 (dentro de transaccion GORM)

    S5. Registrar eventos en timeline (no bloquea si falla)
        registrarEventosCreacion(caseID, agentID, riskLevel, followUps, now)

    S6. Manejar resultado
        SI error:
          → Log warning: "Error generando calendario para caso {iCode}"
        SI exito:
          → Log info: "Calendario generado para caso {iCode} (risk_level={n})"

  SI FollowUpSvc == nil:
    → No generar calendario (servicio no inicializado)

  → FIN SUB-FLUJO → CONTINUA FLUJO GENERAL

┌─────────────────────────────────────────┐
│  SUB-FLUJO: Evento de hechos del caso   │
│  (asincrono — goroutine)                │
└─────────────────────────────────────────┘

  NOTA: Se ejecuta DESPUES del commit en una goroutine separada.
  Si falla, se loguea pero NO deshace la creacion del caso.

  SI CaseTimelineRepo != nil:
    → Lanzar goroutine:

    S1. Extraer subtipos de violencia
        ↺ PARA CADA st EN form2.subtypeViolenceExperienced:
          SI st.name != "":
            → Agregar st.name a subtypeNames

    S2. Construir descripcion
        SI len(subtypeNames) > 0:
          → description = "[subtypeNames unido por ", "] — [factsDescription]"
        SI NO:
          → description = factsDescription

    S3. Insertar evento en timeline
        INSERT en case_timeline_event:
        {
          caseID:      victimCase.iCode,          // origen: caso recien creado
          eventType:   TimelineEventRegistro,      // legacy — compatibilidad
          category:    TimelineCategoryGeneral,
          type:        TimelineTypeHechosCaso,     // "Hechos del caso"
          icon:        TimelineIconHechosCaso,     // "file-lines"
          color:       TimelineColorLightRed,      // "#f87171"
          description: description,               // origen: S2
          eventUserID: session.userICode,          // origen: sesion activa
          date:        form2.factsDate,            // fecha de los hechos, NO now
          createdAt:   now
        }

        SI error:
          → Log warning: "timeline hechos: error insertando evento para caso {iCode}"
        SI exito:
          → Log info: "timeline hechos: evento creado para caso {iCode} (subtipos=N)"

  SI CaseTimelineRepo == nil:
    → No registrar evento (repositorio no inicializado)

  → FIN SUB-FLUJO → CONTINUA FLUJO GENERAL

PASO 19 — Retornar respuesta exitosa
  Retornar HTTP 200 con JSON:
  {
    "login":    user.login,        // origen: getRandomUserName(6) o existente
    "pass":     newPassword         // origen: getRandomPassword(4) o mensaje "usa tu password"
  }
  → TERMINAR ejecucion backend


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Informacion pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decision                                      | Paso afectado          |
|----------------------------------------------------------|------------------------|
| getRandomPassword(4): logica exacta de generacion        | SUB-FLUJO Usuario S4   |
| getRandomUserName(6): logica exacta de generacion        | SUB-FLUJO Usuario S5   |
| No hay envio de SMS con credenciales (boton existe en UI pero no hay logica) | PASO 19 |
```
