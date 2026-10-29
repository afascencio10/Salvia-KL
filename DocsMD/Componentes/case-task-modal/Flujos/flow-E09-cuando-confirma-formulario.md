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

  CASO 'decidir_via_riesgo_inminente':
    formularioValido si form.via es "A" o "B" y form.agentId no está vacío.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete
    con formData { via, agentId }.
    En una transacción: asigna el agente a la remisión,
    Vía A deja el estado en asignada y crea intentar_contacto_3x3 para ese agente,
    Vía B deja el estado en riesgo_inminente_por_escalar_cruz_roja y crea
    escalar_cruz_roja con assigned_role sc,
    y marca esta tarea Done.
    Después, timeline del caso tipo "Vía de atención decidida":
    "Vía A — asignación al agente" o "Vía B — escalar a Cruz Roja".
    Un fallo del timeline no revierte la decisión.

  CASO 'escalar_cruz_roja':
    formularioValido es true.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete.
    La remisión pasa a riesgo_inminente_por_entregar_medida.
    Se crea la tarea entregar_medida_inmediata con assigned_role cr
    y descripción "Entregar medida inmediata (sin caracterización previa)".
    El timeline queda "Escalamiento a Cruz Roja confirmado":
    "Supervisión del Convenio confirmó el escalamiento a Cruz Roja".
    Un fallo del timeline no revierte la confirmación.

  CASO 'entregar_medida_inmediata':
    formularioValido exige al menos una medida. El comentario puede ir vacío.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete
    con formData { medidas, comentario }.
    Las medidas elegidas pasan a ENTREGADA.
    La remisión pasa a riesgo_inminente_por_confirmar_seguridad.
    Se crea confirmar_lugar_seguro para el agente ya asignado a la remisión,
    con descripción "Confirmar que la persona está en lugar seguro".
    Esta tarea queda Done.
    El timeline queda "Medida entregada por Cruz Roja" y nombra las medidas
    elegidas. Si hay comentario, lo agrega. El texto indica que queda en
    espera de confirmación del agente.
    Un fallo del timeline no revierte la entrega.

  CASO 'confirmar_lugar_seguro':
    formularioValido es true.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete
    con formData { confirmado: true }.
    La remisión pasa a asignada.
    Se crea intentar_contacto_3x3 para el mismo agente,
    con descripción "Intentar contacto (3×3)".
    Esta tarea queda Done.
    El timeline queda "Entrega de medida confirmada":
    "El agente de medidas confirmó la entrega de la medida."
    Un fallo del timeline no revierte la confirmación.

  CASO 'revision_equipo_tecnico':
    formularioValido exige concepto y uno de los tres resultados.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete
    con formData { concepto, resultado }.
    Si el resultado es aprobado_para_caracterizacion_me:
      la remisión pasa de pendiente_de_revision a por_asignar,
      se crea asignar_agente_medidas con assigned_role sm
      y el timeline queda "Enviada a asignación de agente".
    Si el resultado es devuelto_para_ajuste o devuelto_para_acompanamiento:
      la remisión pasa a devuelta_para_justificacion o devuelta_a_seguimiento,
      se crea ajustar_remision_devuelta para submitted_by
      y el timeline "Devolución del Equipo Técnico" incluye el concepto.
    Se guardan el concepto, el resultado, la fecha y el usuario de la sesión.
    Esta tarea queda Done.
    Un fallo del timeline no revierte la revisión.
    Si ya hay vueltas, el modal muestra el hilo de GET /case-tasks/:id/chat.

  CASO 'ajustar_remision_devuelta':
    El modal muestra el hilo, el último concepto del equipo y el ajuste.
    Sin la casilla, formularioValido exige el ajuste.
    confirmar() reenvía: la remisión vuelve a pendiente_de_revision
    y se crea otra revision_equipo_tecnico para el rol em.
    El timeline incluye el texto del ajuste.
    Con la casilla, el botón descarta: la remisión pasa a descartada,
    las medidas ACTIVE de esa remisión pasan a CLOSED
    y no nace otra revisión.
    El timeline queda "Remisión descartada".
    Un fallo del timeline no revierte el reenvío ni el descarte.

  CASO 'asignar_agente_medidas':
    El modal muestra el select de Agente de Medidas.
    Reasignar se habilita con un agente.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete
    con formData { agentId }.
    La remisión pasa de por_asignar a asignada y guarda el agente.
    Se crea intentar_contacto_3x3 para ese agente.
    Esta tarea queda Done.
    El timeline queda "Agente de medidas asignado":
    "Se asignó el agente de medidas."
    Un fallo del timeline no revierte la asignación.
    Si la remisión ya no está por asignar, la operación se rechaza.

  CASO 'intentar_contacto_3x3':
    formularioValido exige Sí o No. La nota puede ir vacía.
    confirmar() llama a PUT /api/v1/case-tasks/:id/complete
    con formData { contesto, nota }.
    La remisión tiene que estar en asignada.
    Si contestó, se crea entrevista_caracterizacion para el mismo agente
    y el timeline queda "Contacto con la víctima".
    Si no contestó y van menos de tres, se crea otro intentar_contacto_3x3
    y el timeline queda "Intento de contacto fallido" con el número.
    A la tercera, la remisión pasa a cerrada_por_no_contacto,
    las medidas ACTIVE pasan a CLOSED y no nace otra tarea.
    El timeline queda "Remisión cerrada por no contacto".
    Esta tarea queda Done.
    Un fallo del timeline no revierte el intento.

  CASO 'entrevista_caracterizacion':
    El botón está habilitado.
    confirmar() cierra el modal y no llama al API.

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

  CASO 'Corregir oficio':
    Validar formulario:
      • Sin campos que validar — solo requiere cargandoOficio === false

  CASO 'justificar_remision':
    Validar formulario:
      • respuesta: required (trim)

  CASO 'solicitud_entidad':
    Validar formulario:
      • ítems seleccionados: min 1
      • descripcion: required (trim) en cada ítem seleccionado
      • cargandoObligaciones === false

  CASO 'respuesta_entidad':
    Validar formulario:
      • estadoInstitucional: required (iniciado | en_proceso | concluido)
      • radicado: opcional
      • ampliacion: opcional

  SI validación con errores:
    → saveError = mensaje descriptivo del campo faltante
    → TERMINAR ejecución  // modal permanece abierto

  SI validación ok:
    → saveError = null
    → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


┌─────────────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Construir formData según tarea.type                 │
│  (incluye resolución de labels legibles para lectura posterior) │
└─────────────────────────────────────────────────────────────────┘

  lookupLabel(lista, valor) = lista.find(o => o.value == valor)?.label ?? ''
    // usado sobre this.departamentos / this.ciudades / this.municipios

  CASO 'gestion_llamada':
    formData = {
      departamentoId:      form.departamentoId,
      departamentoNombre:  lookupLabel(departamentos, form.departamentoId),
      ciudadId:            form.ciudadId,
      ciudadNombre:        lookupLabel(ciudades, form.ciudadId),
      municipioId:         form.municipioId,
      municipioNombre:     lookupLabel(municipios, form.municipioId),
      entidadId:           form.entidadId !== 'otra'
                             ? form.entidadId        // integer PK de entity_branch
                             : null,
      entidadNombre:       form.entidadId !== 'otra'
                             ? nombre de entidades.find(e => e.id == form.entidadId)
                             : form.entidadNombre,   // texto libre
      funcionario:         form.funcionario,
      descripcion:         form.descripcion || null,
      generaOficio:        form.generaOficio,
      asunto:              form.generaOficio ? form.asunto    : null,
      rutaKofax:           form.generaOficio ? form.rutaKofax : null,
    }

  CASO 'proyectar_oficio':
    formData = {
      departamentoId:      form.departamentoId,
      departamentoNombre:  lookupLabel(departamentos, form.departamentoId),
      ciudadId:            form.ciudadId,
      ciudadNombre:        lookupLabel(ciudades, form.ciudadId),
      municipioId:         form.municipioId,
      municipioNombre:     lookupLabel(municipios, form.municipioId),
      entidadId:           form.entidadId !== 'otra'
                             ? form.entidadId
                             : null,
      entidadNombre:       form.entidadId !== 'otra'
                             ? nombre de entidades.find(e => e.id == form.entidadId)
                             : form.entidadNombre,
      funcionario:         form.funcionario,
      asunto:              form.asunto,
      rutaKofax:           form.rutaKofax,
    }

  CASO 'comite_caso':
    DECISION_LABELS = {
      activar_enlace: 'Activar Enlace', oficio: 'Generar Oficio',
      recomendaciones_agente: 'Recomendaciones al Agente',
      mecanismo_articulador: 'Mecanismo Articulador',
    }
    NIVEL_LABELS = { municipal: 'Municipal', departamental: 'Departamental', nacional: 'Nacional' }

    formData = {
      decisiones:                   form.decisiones,
      decisionesTexto:              form.decisiones.map(d => DECISION_LABELS[d]),
      observacionesOficio:          form.decisiones.includes('oficio')
                                      ? form.observacionesOficio || null
                                      : null,
      observacionesRecomendaciones: form.decisiones.includes('recomendaciones_agente')
                                      ? form.observacionesRecomendaciones || null
                                      : null,
      nivelMecanismo:               form.decisiones.includes('mecanismo_articulador')
                                      ? form.nivelMecanismo
                                      : null,
      nivelMecanismoTexto:          form.decisiones.includes('mecanismo_articulador')
                                      ? NIVEL_LABELS[form.nivelMecanismo]
                                      : null,
      observacionesMecanismo:       form.decisiones.includes('mecanismo_articulador')
                                      ? form.observacionesMecanismo || null
                                      : null,
    }

  CASO 'Corregir oficio':
    formData = {}   // objeto vacío, sin campos

  CASO 'justificar_remision':
    formData = { respuesta: form.respuesta.trim() }

  CASO 'solicitud_entidad':
    formData = {
      entityName: entityName,
      solicitudes: form.items
        .filter(item => form.seleccion[item.id].selected)
        .map(item => ({
          obligationId: item.id,
          label:        item.label,
          descripcion:  form.seleccion[item.id].descripcion.trim(),
        })),
    }

  CASO 'respuesta_entidad':
    formData = {
      solicitudTaskId, obligationId, label, descripcionSolicitud,  // heredados de tarea.formData
      radicado:                   form.radicado.trim() || null,
      ampliacion:                 form.ampliacion.trim() || null,
      estadoInstitucional:        form.estadoInstitucional,
      estadoInstitucionalTexto:   Iniciado | En proceso | Concluido,
    }

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL


PASO 1 — Activar estado de guardado

  guardando = true
  saveError = null


PASO 2 — Enviar al API

PUT /api/v1/case-tasks/{tarea.id}/complete

// payload:
{
  "userId":   currentUserId,   // icode del usuario en sesión
  "formData": formData         // objeto construido en el sub-flujo anterior
}

→ resultado: objeto case_task actualizado o error


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — CaseTaskService.CompleteWithFormData()
  PUT /api/v1/case-tasks/:id/complete
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Confirmado contra el código real (`case_task_service.go` / `case_task_controller.go`).
Este endpoint existió desde el commit `34784f5` ("feat: Modal para manejar cada tipo
de tarea"), fue borrado accidentalmente un día después por un merge de otro
desarrollador (`cfd1739`), y se restauró en esta sesión con la misma lógica.

BACK 1 — Cargar la tarea y parsear formData

  DB.case_tasks.FindByID({ id })
  SI no existe → retornar 404 "tarea no encontrada"

  fd = JSON.parse(formData)   // map[string]interface{}


BACK 2 — Guardar form_data y marcar Done (según tipo — no es uniforme)

  SEGÚN task.Type:

    CASO 'proyectar_oficio':
      → Solo guarda form_data (status se marca Done más abajo, vía PerformAction → completarCaseTask())
        DB.case_tasks.UpdateFields({ form_data: JSON(formData) })

    CASO 'Corregir oficio':
      → No guarda form_data (formData siempre es {} — nada que persistir)
        Status se marca Done más abajo, vía PerformAction → completarCaseTask()

    CASO 'gestion_llamada' / 'comite_caso' / 'justificar_remision':
      → Guarda form_data Y marca Done en la misma escritura
        DB.case_tasks.UpdateFields({
            form_data:    JSON(formData),
            status:       'Done',
            completed_at: now(),
        })

    CASO 'solicitud_entidad':
      → NO marca Done todavía. Primero valida solicitudes y resuelve el usuario
        único de la sede. Si falla, 400 y la tarea queda ToDo.
      → Si ok: transacción — Done + form_data + N case_task respuesta_entidad
      → Tras la transacción: INSERT case_timeline_event "Solicitud a entidad"
        (WARN si falla; no revierte Done)

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

    PerformAction se encarga de (dentro de entity_letter_service.go):
      • Validar entity_letter.state === 'por_proyectar' + todos los campos requeridos
        (departmentId, cityId, townId, entityName, officialDependency, subject, urlKofax)
        SI falla cualquier validación → retorna error INMEDIATAMENTE, sin tocar la BD
      • UpdateFields en entity_letter (campos + state = 'para_revisar')
      • completarCaseTask() → busca la CaseTask 'ToDo' vinculada a este entity_letter_id
        (que es justo ESTA tarea, todavía ToDo — BACK 2 solo guardó form_data,
        no cambió status) → la marca Done + completed_at = now()
      • actualizarBarreraEnGestion() → si barrier_v2.status === 'OPEN', pasa a 'En Gestion'
      • registrarEventoOficio() → CaseTimelineEvent { type: "Oficio Para Revisar", color: azul }

    SI error en PerformAction (estado inválido o campo faltante):
      → Retornar 422 (o 400 si es campo faltante) con el mensaje de PerformAction
      → ⚠️ La tarea NO llegó a completarse — se queda en 'ToDo' porque
        completarCaseTask() nunca se ejecuta (la función retorna antes,
        en la validación). El único cambio persistido en BACK 2 fue el form_data.

  └──────────────────────────────────────────────────────────────┘

  ┌── 'gestion_llamada' ─────────────────────────────────────────┐

    ⚡ Todo este bloque corre en goroutine fire-and-forget
    (sideEffectsGestionLlamada) — el request ya respondió 200 antes
    de que esto termine. Los errores solo se loguean (WARN), nunca
    se propagan al cliente.

    SI formData.generaOficio === true:
      → DB.entity_letters.Create({
            barrier_id:          tarea.barrier_id,
            case_id:             tarea.case_id,
            agent_id:            userId,
            register_by:         userId,
            state:               'para_revisar',   // ya proyectado, sin pasar por por_proyectar
            department_id:       fd.departamentoId,
            city_id:             fd.ciudadId,
            town_id:             fd.municipioId,
            entidad:             fd.entidadNombre,
            entity_branch_id:    fd.entidadId  (int64, solo si vino como número — 'otra' no castea)
            official_dependency: fd.funcionario,
            subject:             fd.asunto,
            url_kofax:           fd.rutaKofax,
        })
      SI falla el insert:
        → Log WARN (no aborta — la tarea ya está Done)

    SI formData.generaOficio === false:
      → No crear entity_letter

    → CaseTimelineEvent {
          category:    'Barreras',
          type:        'Gestión de Llamada',
          eventType:   'GESTION_LLAMADA',
          icon:        'phone',
          color:       azul (#3b82f6),
          case_id:     tarea.case_id,
          barrier_id:  tarea.barrier_id,
          event_user_id: userId,
          task_id:     tarea.id,
      }
      SI falla el insert: → Log WARN (no aborta)

  └──────────────────────────────────────────────────────────────┘

  ┌── 'comite_caso' ─────────────────────────────────────────────┐

    ⚡ Todo este bloque corre en goroutine fire-and-forget
    (sideEffectsComiteCaso) — mismo patrón que gestion_llamada.

    PARA CADA decision EN fd.decisiones:

      SEGÚN decision:

        CASO 'activar_enlace':
          ✅ Resuelto: enlace_activado vive en barrier_v2, NO en case_task
          ni en victim_case (campo agregado en el mismo commit original:
          BarrierV2.EnlaceActivado, columna enlace_activado).
          SI tarea.barrier_id existe:
            → DB.barrier_v2.UpdateFields(tarea.barrier_id, { enlace_activado: true })
            SI falla: → Log WARN (no aborta)

        CASO 'oficio':
          → newLetter = DB.entity_letters.Create({
                barrier_id: tarea.barrier_id,
                case_id:    tarea.case_id,
                agent_id:   userId,   // el mismo usuario que completa la tarea, no un lookup separado
                state:      'por_proyectar',
            })
          SI falla el insert: → Log WARN → continuar con la siguiente decisión (no crea la task)

          → DB.case_tasks.Create({
                category:         'Barreras',
                type:             'proyectar_oficio',
                description:      'Proyectar oficio — Decisión de comité',
                status:           'ToDo',
                assigned_user_id: userId,
                case_id:          tarea.case_id,
                barrier_id:       tarea.barrier_id,
                entity_letter_id: newLetter.id,
            })
          SI falla: → Log WARN (no aborta)

        (recomendaciones_agente / mecanismo_articulador no tienen efecto de lado
        propio más allá de quedar guardadas en form_data — solo informan al humano)

    → CaseTimelineEvent {
          category:   'General',
          type:       'Decisiones del Comité',
          eventType:  'COMITE_CASO',
          icon:       'users',
          color:      morado (#7c3aed),
          case_id:    tarea.case_id,
          barrier_id: tarea.barrier_id,
          event_user_id: userId,
          task_id:    tarea.id,
      }
      SI falla el insert: → Log WARN (no aborta)

  └──────────────────────────────────────────────────────────────┘

  ┌── 'Corregir oficio' ─────────────────────────────────────────┐

    ✅ Confirmado contra el código real — ya no es una inferencia.

    SI tarea.entity_letter_id es null:
      → Retornar error "Corregir oficio sin entity_letter_id" (400)
      → TERMINAR

    → EntityLetterService.PerformAction(ctx, tarea.entity_letter_id, {
          action: "corregir",
          userId: userId,
      })

    PerformAction se encarga de (dentro de entity_letter_service.go):
      • Validar entity_letter.state === 'en_correccion'
        SI falla → retorna error INMEDIATAMENTE, sin tocar la BD
      • UpdateFields en entity_letter → state = 'para_revisar'
      • completarCaseTask() → busca la CaseTask 'ToDo' vinculada a este entity_letter_id
        (que es justo ESTA tarea — BACK 2 no hizo nada para 'Corregir oficio',
        la tarea sigue ToDo hasta este punto) → la marca Done + completed_at = now()
      • registrarEventoOficio() → CaseTimelineEvent { type: "Oficio Para Revisar", color: azul }

    SI error en PerformAction (ej. entity_letter.state !== 'en_correccion'):
      → Retornar 422 con el mensaje de PerformAction
      → ⚠️ La tarea NO se completa — se queda en 'ToDo' porque completarCaseTask()
        nunca se ejecuta (la función retorna antes, en la validación).
        No hay ningún cambio persistido en BACK 2 para este tipo (form_data
        no se guarda), así que en un error la tarea queda exactamente como estaba.

  └──────────────────────────────────────────────────────────────┘


  ┌── 'solicitud_entidad' ────────────────────────────────────────┐

    ⚡ Síncrono — si falla, el cliente recibe 400 y la tarea original
    permanece ToDo.

    1. Parsear formData.solicitudes[] (obligationId, label, descripcion)
       SI length < 1 o alguna descripcion vacía → 400
    2. Cargar entity_case.entity_branch_id desde tarea.entity_case_id
    3. SELECT general_user_i_code FROM security.general_user
         WHERE entity_branch_id = sede
       SI 0 usuarios → 400 "no hay usuario de entidad asociado a esta sede"
       SI >1 usuario → 400 "hay más de un usuario de entidad para esta sede"
    4. Transacción:
       • Marcar solicitud_entidad Done + form_data
       • Crear N case_task type=respuesta_entidad, category=Entidades,
         assigned_user_id = ese usuario, entity_case_id heredado,
         form_data = { solicitudTaskId, obligationId, label, descripcionSolicitud }
    5. Tras éxito de la transacción: INSERT case_timeline_event
         category = Entidades
         type     = "Solicitud a entidad"
         icon     = folder
         color    = #5106A7
         description = "Solicitudes enviadas a {entityName} (N ítems)"
         event_user_id = userId del complete
       SI el INSERT falla → WARN; la tarea permanece Done
    6. No cambia entity_case_obligation.status

  └──────────────────────────────────────────────────────────────┘


  ┌── 'respuesta_entidad' ────────────────────────────────────────┐

    Síncrono. Valida estadoInstitucional ∈ {iniciado, en_proceso, concluido}.
    SI inválido → 400, la tarea permanece ToDo.

    SI ok:
      DB.case_tasks.UpdateFields({
          form_data:    JSON(formData),   // incluye contexto heredado + respuesta
          status:       'Done',
          completed_at: now(),
      })
      INSERT case_timeline_event
        category = Entidades
        type     = "Respuesta institucional"
        icon     = reply
        color    = #15803d
        description = "Respuesta sobre: {label} ({estado legible})"
        event_user_id = userId del complete
      SI el INSERT falla → WARN; la tarea permanece Done
      No cambia entity_case_obligation.status. No crea más tareas.

  └──────────────────────────────────────────────────────────────┘


BACK 4 — Retornar la tarea actualizada

  DB.case_tasks.FindByID({ taskId }) → retornar objeto completo


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ GAPS resueltos (ya no aplican)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Lo que se creía un GAP                                          | Resolución |
|-------------------------------------------------------------------|-----------|
| ¿enlace_activado vive en case_task o en victim_case?              | Vive en **barrier_v2** (`BarrierV2.EnlaceActivado`, columna `enlace_activado`) |
| ¿registrarEventoTarea() escribe igual que registrarEventoOficio()?| No existe una función compartida — cada tipo construye su `CaseTimelineEvent` inline dentro de `sideEffectsGestionLlamada` / `sideEffectsComiteCaso`, con `EventType`/`Icon`/`Color` propios |
| Mismatch PUT vs POST / formData vs result                        | Era una regresión real: el commit `34784f5` (autor original) sí implementaba `PUT .../complete` + `CompleteWithFormData`; el commit `cfd1739` (otro dev, un día después) lo borró por error de merge. **Restaurado en esta sesión.** |
| Rama 'Corregir oficio' — ¿inferencia o código real?               | Era código real del commit original, no una inferencia — confirmado y restaurado tal cual |

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información aún pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                             | Paso afectado |
|-----------------------------------------------------------------|---------------|
| ¿town_id en formData/entity_letter es DIVIPOLA town_code (string)| BACK 3 llamada|
| o el id integer de la tabla town? El modelo lo tipa como *string,|               |
| pero no se confirmó contra datos reales de esta sesión.         |               |


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

CASO 'articulacion_territorial':
  Validar: fechaContacto, compromiso; fechaSeguimiento si requiereSeguimiento
  POST /api/v1/case-tasks/articulacion-territorial
  Body: { barrierId, taskId si existe, formData }
  Backend:
    SI barrera no asignada al usuario → 403
    SI taskId: completa ToDo articulacion_territorial
    SI NO taskId: crea Done description="Gestión propia"
    SI requiereSeguimiento: crea ToDo compromiso_con_enlace
    SI barrera OPEN → En Gestion
  SI 200 → emit completed + cerrar

CASO 'gestion_mecanismo':
  Validar según mecanismoCtx.scope:
    municipal:     abordadoEnSesion + reporteGestion
    departamental: fechaProxima + detalleProxima + reporteGestion
    nacional:      reporteGestion + asistenciasFecha + asistenciasTexto
  POST /api/v1/case-tasks/gestion-mecanismo
  Body: { mechanismId, formData }
  Backend:
    SI mecanismo no asignado al usuario → 403
    SI no active → 409
    Siempre INSERT case_task Done (nunca PUT complete)
    result = resumen concatenado del formulario
  SI 200 → emit completed + cerrar

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                             | Paso afectado |
|-----------------------------------------------------------------|---------------|
| ¿El municipioId en formData es el town_code (string DIVIPOLA)   | Sub-flujo     |
| o el id integer de la tabla town?                               | formData      |

✅ Resuelto: `tareaActualizada` es el objeto `CaseTask` completo — `CompleteWithFormData`
retorna `s.repo.FindByID(ctx, id)` al final, y el controller hace `ctx.JSON(http.StatusOK, task)`.
