# flow-E05 — Cuando presiona "Gestionar"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Gestionar"
   Tipo: User Interaction
   Función: openGestionarModal(oficio)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  oficio:  objeto del oficio seleccionado  → fila donde el usuario presionó el botón
}

PASO 1 — Guardar el oficio seleccionado
  selectedOficio = oficio

PASO 2 — Determinar qué modal abrir según el estado del oficio

SEGÚN oficio.status:
  CASO 'por_proyectar':       → activeModal = 'proyectar'        → CONTINÚA
  CASO 'para_revisar':        → activeModal = 'revisar'          → CONTINÚA
  CASO 'en_correccion':       → activeModal = 'corregir'         → CONTINÚA
  CASO 'aprobacion_juridica': → activeModal = 'aprobar'          → CONTINÚA
  CASO 'para_radicar':        → activeModal = 'radicar'          → CONTINÚA
  CASO 'radicado':            → activeModal = 'registrar_respuesta' → CONTINÚA
  DEFAULT:                    → activeModal = null (modal no se muestra) → CONTINÚA

PASO 3 — Pre-cargar el formulario con datos existentes del oficio
  modalForm.nivel                = oficio.nivel            || ''
  modalForm.kofaxPath            = oficio.urlKofax         || ''
  modalForm.prioridad            = oficio.letterPriority   || 'normal'
  modalForm.entidad              = oficio.entidad          || oficio.barrierOrg || ''
  modalForm.asunto               = oficio.asuntoRadicado   || ''
  modalForm.correoEntidad        = oficio.correoEntidad    || ''
  modalForm.numeroRadicado       = oficio.numeroRadicado   || ''
  modalForm.fechaRespuesta       = ''
  modalForm.correoRemitente      = oficio.correo           || ''
  modalForm.asuntoRespuesta      = ''
  modalForm.respuestaRecibidaPor = ''
  modalForm.reasonCorrection     = ''

PASO 4 — Limpiar errores y mostrar el modal
  saveError = null
  showModal = true
  → Vue evalúa activeModal y muestra el modal correspondiente
```
