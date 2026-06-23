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
      • nivel:      required
      • entidad:    required
      • kofaxPath:  required
    SI validación con errores:
      → alert() con el campo faltante
      → TERMINAR ejecución
    SI validación ok:
      → payload.nivel    = modalForm.nivel
      → payload.entidad  = modalForm.entidad
      → payload.urlKofax = modalForm.kofaxPath
      → payload.priority = modalForm.prioridad || 'normal'

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

PATCH /api/v1/entity-letters/{selectedOficio.id}/action

// payload:
{
  "action":  acción ejecutada           → determinada en Sub-flujo
  "userId":  currentUserId              → sesión
  ...campos adicionales según acción
}

→ resultado: objeto entity_letter actualizado o error

PASO 5 — Manejar respuesta

SI status === 200:
  → Ocultar overlay de éxito después de 1.2s
  → Buscar el oficio en this.oficios por ID
  → Actualizar en la lista (sin recargar):
      status    = response.state
      canManage = canManageForRole(response.state)
  → Vue re-evalúa reactivamente: badge de pendientes, filtros, tabla
  → Cerrar modal → Ver flujo: Cuando cancela el modal (E06)
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
| ¿Qué campos exactos devuelve el PATCH en la respuesta?       | PASO 5        |
| Las validaciones usan alert() — ¿se migrarán a errores       | Sub-flujo     |
| inline para ser consistentes con el resto del sistema?       |               |
| El modal 'aprobar' llama action='radicar' igual que          | PASO 2        |
| 'para_radicar' — ¿son la misma transición en el backend?     |               |
```
