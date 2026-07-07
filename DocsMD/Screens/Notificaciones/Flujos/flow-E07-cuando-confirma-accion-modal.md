# flow-E07 — Cuando confirma acción en modal

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando confirma acción en modal
   Tipo: User Interaction
   Función: submitModal(action)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  action:         string con la acción a ejecutar    → botón presionado en el modal
                  valores posibles: 'proyectar' | 'revisar' | 'corregir' |
                  'por_corregir' | 'radicar' | 'registrar_respuesta'
  selectedOficio: oficio actualmente en el modal     → state del componente
  modalForm:      campos del formulario del modal    → inputs llenados por el usuario
  currentUserId:  ID del usuario en sesión           → window.NotifConfig
}

PASO 1 — Verificar que hay un oficio seleccionado

SI selectedOficio es null:
  → TERMINAR ejecución (guard clause)

SI selectedOficio existe:
  → CONTINÚA FLUJO GENERAL

PASO 2 — Construir payload base
  payload = {
    action: action,          // acción a ejecutar
    userId: currentUserId    // ID del usuario en sesión
  }

┌─────────────────────────────────────────────────────────────────┐
│  SUB-FLUJO: Validación y campos adicionales según acción        │
└─────────────────────────────────────────────────────────────────┘

  SI action === 'proyectar':
    Validar formulario:
      • departmentId:      required
      • cityId:            required
      • townId:            required
      • entityBranchId:    required (integer PK de entity_branch, o 'otra')
      • entityName:        required si entityBranchId === 'otra'
      • officialDependency: required
      • subject:           required
      • kofaxPath:         required
    SI validación con errores:
      → alert() con el campo faltante
      → TERMINAR ejecución
    SI validación ok:
      → payload.departmentId       = modalForm.departmentId
      → payload.cityId             = modalForm.cityId
      → payload.townId             = modalForm.townId
      → payload.entityBranchId     = modalForm.entityBranchId !== 'otra'
                                       ? modalForm.entityBranchId   // integer (PK FK a entity_branch)
                                       : null
      → payload.entityName         = nombre legible de la sede o texto libre si 'otra'
      → payload.officialDependency = modalForm.officialDependency
      → payload.subject            = modalForm.subject
      → payload.urlKofax           = modalForm.kofaxPath
      → payload.priority           = modalForm.prioridad || 'normal'

  SI action === 'por_corregir':
    Validar formulario:
      • reasonCorrection: required
    SI validación con errores:
      → alert()
      → TERMINAR ejecución
    SI validación ok:
      → payload.reasonCorrection = modalForm.reasonCorrection

  SI action === 'radicar':
    Validar formulario:
      • asunto:         required
      • correoEntidad:  required
      • numeroRadicado: required
    SI validación con errores:
      → alert() con el campo faltante
      → TERMINAR ejecución
    SI validación ok:
      → payload.asuntoRadicado = modalForm.asunto
      → payload.correoEntidad  = modalForm.correoEntidad
      → payload.numeroRadicado = modalForm.numeroRadicado

  SI action === 'registrar_respuesta':
    Validar formulario:
      • fechaRespuesta:       required
      • correoRemitente:      required
      • asuntoRespuesta:      required
      • respuestaRecibidaPor: required
    SI validación con errores:
      → alert() con el campo faltante
      → TERMINAR ejecución
    SI validación ok:
      → payload.responseDate     = modalForm.fechaRespuesta
      → payload.correoRemitente  = modalForm.correoRemitente
      → payload.asuntoRespuesta  = modalForm.asuntoRespuesta
      → payload.responseReviewBy = modalForm.respuestaRecibidaPor

  SI action === 'revisar' o 'corregir':
    → No requiere campos adicionales
    → payload solo contiene { action, userId }

  → FIN SUB-FLUJO → CONTINÚA FLUJO GENERAL

PASO 3 — Activar estado de guardado
  isSaving  = true
  saveError = null

PASO 4 — Enviar acción al API

PUT /api/v1/entity-letters/{selectedOficio.id}/action

// payload:
{
  "action":  acción ejecutada           → determinada en Sub-flujo
  "userId":  currentUserId              → sesión
  ...campos adicionales según acción
}

// Payload completo para action === 'proyectar':
{
  "action":             "proyectar",
  "userId":             string (UUID del usuario),
  "departmentId":       string,
  "cityId":             string,
  "townId":             string (DIVIPOLA town_code, ej. "11001000"),
  "entityBranchId":     integer | null  (PK de salvia.entity_branch; null si eligió "otra"),
  "entityName":         string  (nombre legible de la sede o texto libre),
  "officialDependency": string,
  "subject":            string,
  "urlKofax":           string,
  "priority":           "normal" | "alta"
}

→ resultado: objeto entity_letter actualizado o error

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  BACKEND — entity_letter_service.go → PerformAction()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

BACK 1 — Cargar el oficio

  DB.entity_letters.FindByID({ id })
  SI no existe → retornar 404

BACK 2 — Ejecutar lógica según action

  SEGÚN action:

  ┌── 'proyectar'  (por_proyectar → para_revisar) ──────────────┐

    Validar estado actual = 'por_proyectar'  (422 si no coincide)
    Validar requeridos: departmentId, cityId, townId, entityName,
                        officialDependency, subject, urlKofax

    → DB.entity_letters.UpdateFields({
          department_id:       departmentId,
          city_id:             cityId,
          town_id:             townId,
          entidad:             entityName,
          official_dependency: officialDependency,
          subject:             subject,
          url_kofax:           urlKofax,
          entity_branch_id:    entityBranchId  (null si eligió 'otra'),
          priority:            priority || 'normal',
          register_by:         userId,
          state:               'para_revisar',
      })

    → completarCaseTask(id)  [fire-and-forget]
        • FindTodoByEntityLetterID(id)
          → UpdateFields({ status: 'Done', completed_at: now() })

    → registrarEventoOficio('para_revisar')  [fire-and-forget]
        • CaseTimelineEvent { type: "Oficio Para Revisar", color: blue }

  └──────────────────────────────────────────────────────────────┘

  ┌── 'revisar'  (para_revisar → aprobacion_juridica) ──────────┐

    Validar estado actual = 'para_revisar'  (422 si no)
    → DB.entity_letters.UpdateFields({ review_by: userId, state: 'aprobacion_juridica' })
    → registrarEventoOficio('aprobacion_juridica')  [fire-and-forget]
        • CaseTimelineEvent { type: "Oficio en Aprobación Jurídica", color: purple }

  └──────────────────────────────────────────────────────────────┘

  ┌── 'por_corregir'  (para_revisar → en_correccion) ───────────┐

    Validar estado actual = 'para_revisar'  (422 si no)
    Validar reasonCorrection requerido

    → DB.entity_letters.UpdateFields({
          reason_correction: reasonCorrection,
          state:             'en_correccion',
      })

    → crearCaseTaskCorreccion()  [fire-and-forget]
        • CaseTask {
              type:             'Corregir oficio',
              assigned_user_id: letter.AgentID,
              description:      "Corregir oficio — Razón: {reasonCorrection}",
              status:           'ToDo',
              entity_letter_id: id,
          }

    → registrarEventoOficio('en_correccion')  [fire-and-forget]
        • CaseTimelineEvent { type: "Oficio En Corrección",
                              description: "... — Razón: {reason}", color: orange }

  └──────────────────────────────────────────────────────────────┘

  ┌── 'corregir'  (en_correccion → para_revisar) ───────────────┐

    Validar estado actual = 'en_correccion'  (422 si no)
    → DB.entity_letters.UpdateFields({ state: 'para_revisar' })

    → completarCaseTask(id)  [fire-and-forget]
        • Busca CaseTask ToDo de corrección con entity_letter_id = id
        • UpdateFields({ status: 'Done', completed_at: now() })

    → registrarEventoOficio('para_revisar')  [fire-and-forget]
        • CaseTimelineEvent { type: "Oficio Para Revisar", color: blue }

  └──────────────────────────────────────────────────────────────┘

  ┌── 'radicar'  (aprobacion_juridica → radicado) ──────────────┐

    Validar estado actual = 'aprobacion_juridica'  (422 si no)
    Validar requeridos: asuntoRadicado, correoEntidad, numeroRadicado

    → DB.entity_letters.UpdateFields({
          asunto_radicado:  asuntoRadicado,
          correo_entidad:   correoEntidad,
          numero_radicado:  numeroRadicado,
          radicado_by:      userId,
          state:            'radicado',
      })

    → registrarEventoOficio('radicado')  [fire-and-forget]
        • CaseTimelineEvent { type: "Oficio Radicado", color: green }

  └──────────────────────────────────────────────────────────────┘

  ┌── 'registrar_respuesta'  (radicado → respondido) ───────────┐

    Validar estado actual = 'radicado'  (422 si no)
    Validar requeridos: responseDate, correoRemitente, asuntoRespuesta, responseReviewBy
    Parsear responseDate: "YYYY-MM-DD" → time.Time

    → DB.entity_letters.UpdateFields({
          correo_remitente:    correoRemitente,
          asunto_respuesta:    asuntoRespuesta,
          response_review_by:  responseReviewBy,
          response_date:       time.Time parseado,
          state:               'respondido',
      })

    → registrarEventoOficio('respondido')  [fire-and-forget]
        • CaseTimelineEvent { type: "Oficio Respondido", color: green }

  └──────────────────────────────────────────────────────────────┘

BACK 3 — Retornar el oficio actualizado

  DB.entity_letters.FindByID({ id }) → retornar objeto completo

PASO 5 — Manejar respuesta

SI status === 200:
  → Ocultar overlay de éxito después de 1.2s
  → Cerrar modal → Ver flujo: Cuando cancela el modal (E06)
  → Recargar página actual desde BD: loadOficios()
      (actualiza tabla, total, pendingCount y paginador)
  → isSaving = false

SI status === 401:
  → Redirigir a /static/landing.html
  → TERMINAR ejecución

SI status === 422:
  → saveError = response.error || 'Transición de estado no permitida.'
  → isSaving = false
  → El modal permanece abierto mostrando el error

SI otro status:
  → saveError = response.error || 'No se pudo guardar el oficio. Intenta de nuevo.'
  → isSaving = false
  → El modal permanece abierto mostrando el error

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                          | Paso afectado |
|--------------------------------------------------------------|---------------|
| Las validaciones usan alert() — ¿se migrarán a errores       | Sub-flujo     |
| inline para ser consistentes con el resto del sistema?       |               |
| El modal 'aprobar' llama action='radicar' igual que          | PASO 2        |
| 'para_radicar' — ¿son la misma transición en el backend?     |               |
```
