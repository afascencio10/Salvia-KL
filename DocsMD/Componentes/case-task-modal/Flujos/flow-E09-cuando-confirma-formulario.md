━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma el formulario
   Tipo: User Interaction
   Función: confirmar()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  tarea:          objeto case_task cargado en E02     → estado del componente
  form:           campos llenados por el usuario      → estado del componente
  currentUserId:  icode del usuario en sesión         → window config del componente
}


┌─────────────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Validación según tarea.type                         │
└─────────────────────────────────────────────────────────────────┘

  CASO 'gestion_llamada':
    Validar formulario:
      • departamentoId:  required
      • ciudadId:        required
      • municipioId:     required
      • entidadId:       required
      • entidadNombre:   required si entidadId === 'otra'
      • funcionario:     required
      • descripcion:     opcional
      • asunto:          required si form.generaOficio === true
      • rutaKofax:       required si form.generaOficio === true

  CASO 'proyectar_oficio':
    Validar formulario:
      • departamentoId:  required
      • ciudadId:        required
      • municipioId:     required
      • entidadId:       required
      • entidadNombre:   required si entidadId === 'otra'
      • funcionario:     required
      • asunto:          required
      • rutaKofax:       required

  CASO 'comite_caso':
    Validar formulario:
      • decisiones:      required, min 1 opción seleccionada
      • nivelMecanismo:  required si 'mecanismo_articulador' ∈ form.decisiones

  SI validación con errores:
    → saveError = mensaje descriptivo del campo faltante
    → TERMINAR ejecución  // modal permanece abierto

  SI validación ok:
    → saveError = null
    → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


┌─────────────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Construir formData según tarea.type                 │
└─────────────────────────────────────────────────────────────────┘

  CASO 'gestion_llamada':
    formData = {
      departamentoId:  form.departamentoId,
      ciudadId:        form.ciudadId,
      municipioId:     form.municipioId,
      entidadId:       form.entidadId !== 'otra'
                         ? form.entidadId        // integer PK de entity_branch
                         : null,
      entidadNombre:   form.entidadId !== 'otra'
                         ? nombre de entidades.find(e => e.id == form.entidadId)
                         : form.entidadNombre,   // texto libre
      funcionario:     form.funcionario,
      descripcion:     form.descripcion || null,
      generaOficio:    form.generaOficio,
      asunto:          form.generaOficio ? form.asunto    : null,
      rutaKofax:       form.generaOficio ? form.rutaKofax : null,
    }

  CASO 'proyectar_oficio':
    formData = {
      departamentoId:  form.departamentoId,
      ciudadId:        form.ciudadId,
      municipioId:     form.municipioId,
      entidadId:       form.entidadId !== 'otra'
                         ? form.entidadId
                         : null,
      entidadNombre:   form.entidadId !== 'otra'
                         ? nombre de entidades.find(e => e.id == form.entidadId)
                         : form.entidadNombre,
      funcionario:     form.funcionario,
      asunto:          form.asunto,
      rutaKofax:       form.rutaKofax,
    }

  CASO 'comite_caso':
    formData = {
      decisiones:                   form.decisiones,
      observacionesOficio:          form.decisiones.includes('oficio')
                                      ? form.observacionesOficio || null
                                      : null,
      observacionesRecomendaciones: form.decisiones.includes('recomendaciones_agente')
                                      ? form.observacionesRecomendaciones || null
                                      : null,
      nivelMecanismo:               form.decisiones.includes('mecanismo_articulador')
                                      ? form.nivelMecanismo
                                      : null,
      observacionesMecanismo:       form.decisiones.includes('mecanismo_articulador')
                                      ? form.observacionesMecanismo || null
                                      : null,
    }

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


PASO 1 — Activar estado de guardado

  guardando = true
  saveError = null


PASO 2 — Enviar al API

PUT /api/v1/tasks/{tarea.id}/complete

// payload:
{
  "userId":   currentUserId,   // icode del usuario en sesión
  "formData": formData         // objeto construido en el sub-flujo anterior
}

→ resultado: objeto case_task actualizado o error


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — CaseTaskService.CompletarTarea()
  PUT /api/v1/tasks/:taskId/complete
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

BACK 1 — Cargar la tarea

  DB.case_tasks.FindByID({ taskId })
  SI no existe → retornar 404


BACK 2 — Persistir resultado del formulario y marcar como Done

  DB.case_tasks.UpdateFields({
      form_data:    JSON(formData),
      status:       'Done',
      completed_at: now(),
  })


BACK 3 — Ejecutar efectos de lado según tarea.type

  ┌── 'proyectar_oficio' ────────────────────────────────────────┐

    → EntityLetterService.PerformAction("proyectar", tarea.entity_letter_id, {
          userId:             userId,
          departmentId:       formData.departamentoId,
          cityId:             formData.ciudadId,
          townId:             formData.municipioId,
          entityBranchId:     formData.entidadId  (null si 'otra'),
          entityName:         formData.entidadNombre,
          officialDependency: formData.funcionario,
          subject:            formData.asunto,
          urlKofax:           formData.rutaKofax,
      })

    PerformAction se encarga de:
      • Validar entity_letter.state = 'por_proyectar'
      • UpdateFields en entity_letter (campos + state = 'para_revisar')
      • completarCaseTask() → no encontrará tareas ToDo (ya marcada Done en BACK 2) → noop
      • registrarEventoOficio() → CaseTimelineEvent { type: "Oficio Para Revisar", color: blue }

    SI error en PerformAction:
      → Retornar 422 con el mensaje de PerformAction
        (la tarea ya quedó Done — se informa el problema de estado del oficio)

  └──────────────────────────────────────────────────────────────┘

  ┌── 'gestion_llamada' ─────────────────────────────────────────┐

    SI formData.generaOficio === true:
      → DB.entity_letters.Create({
            barrier_id:          tarea.barrier_id,
            case_id:             tarea.case_id,
            agent_id:            userId,
            state:               'para_revisar',   // ya proyectado, sin pasar por por_proyectar
            department_id:       formData.departamentoId,
            city_id:             formData.ciudadId,
            town_id:             formData.municipioId,
            entidad:             formData.entidadNombre,
            entity_branch_id:    formData.entidadId  (null si 'otra'),
            official_dependency: formData.funcionario,
            subject:             formData.asunto,
            url_kofax:           formData.rutaKofax,
            register_by:         userId,
        })
      SI falla el insert:
        → Log WARN (no aborta — la tarea ya está Done)

    SI formData.generaOficio === false:
      → No crear entity_letter

    → registrarEventoTarea()  [fire-and-forget, función nueva]
        • CaseTimelineEvent {
              category:    'Barreras',
              type:        'Gestión de Llamada',
              case_id:     tarea.case_id,
              barrier_id:  tarea.barrier_id,
          }

  └──────────────────────────────────────────────────────────────┘

  ┌── 'comite_caso' ─────────────────────────────────────────────┐

    SI formData.decisiones.includes('activar_enlace'):
      → DB.case_tasks.UpdateFields({ enlace_activado: true })
        // ⚠️ GAP: confirmar si este campo vive en case_task o en victim_case

    SI formData.decisiones.includes('oficio'):
      → agentId = DB.victim_cases.FindAgentIDByICode(tarea.case_id)

      → newLetter = DB.entity_letters.Create({
            barrier_id: tarea.barrier_id,
            case_id:    tarea.case_id,
            agent_id:   agentId,
            state:      'por_proyectar',
        })

      → DB.case_tasks.Create({
            category:         'Barreras',
            type:             'proyectar_oficio',
            description:      'Proyectar oficio — Decisión de comité',
            status:           'ToDo',
            assigned_user_id: agentId,
            case_id:          tarea.case_id,
            barrier_id:       tarea.barrier_id,
            entity_letter_id: newLetter.id,
        })

      SI falla algún insert:
        → Log WARN (no aborta)

    → registrarEventoTarea()  [fire-and-forget, función nueva]
        • CaseTimelineEvent {
              category:   'General',
              type:       'Decisiones del Comité',
              case_id:    tarea.case_id,
              barrier_id: tarea.barrier_id,
          }

  └──────────────────────────────────────────────────────────────┘


BACK 4 — Retornar la tarea actualizada

  DB.case_tasks.FindByID({ taskId }) → retornar objeto completo


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS adicionales (backend)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                             | Paso afectado |
|-----------------------------------------------------------------|---------------|
| ¿enlace_activado vive en case_task o en victim_case?            | BACK 3 comité |
| ¿registrarEventoTarea() escribe en case_timeline_event igual    | BACK 3 todos  |
| que registrarEventoOficio()? ¿Mismo modelo, campos distintos?   |               |
| ¿town_id en DB es DIVIPOLA town_code (string) o PK integer?     | BACK 3 llamada|


PASO 3 — Manejar respuesta

  SI status === 200:
    → guardando = false
    → cancelar()                         // cierra el modal y limpia el estado
    → emit('completed', tareaActualizada) // notifica al padre
    → TERMINAR ejecución

  SI status === 401:
    → Redirigir a /static/landing.html
    → TERMINAR ejecución

  SI status === 422:
    → saveError  = response.error || 'Transición no permitida.'
    → guardando  = false
    → Modal permanece abierto mostrando saveError

  SI otro status:
    → saveError  = response.error || 'No se pudo completar la tarea. Intenta de nuevo.'
    → guardando  = false
    → Modal permanece abierto mostrando saveError

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                             | Paso afectado |
|-----------------------------------------------------------------|---------------|
| Shape exacto del objeto tareaActualizada devuelto por el API    | PASO 3        |
| ¿El backend devuelve la tarea completa o solo { id, status }?   | PASO 3        |
| ¿El municipioId en formData es el town_code (string DIVIPOLA)   | Sub-flujo     |
| o el id integer de la tabla town?                               | formData      |
