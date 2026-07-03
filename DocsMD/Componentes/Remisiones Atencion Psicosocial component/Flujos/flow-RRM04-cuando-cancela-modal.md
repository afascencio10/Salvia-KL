━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando cancela o cierra el modal
   Tipo: User Interaction
   ID: RRM-04
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por:
  - Botón "Cancelar"
  - Clic en backdrop (.rrm-backdrop)
  - Tecla Escape (opcional)

Precondiciones:
  saving === false   // bloquear cierre durante guardado


PASO 1 — Cerrar UI

  visible = false


PASO 2 — Limpiar estado interno

  remisiones = []
  asignarEnDupla = false
  selectedProfessionalId = ''
  selectedDuplaId = ''
  professionalGroups = []
  duplaOptions = []
  loadingOptions = false
  optionsError = null
  saveError = null


PASO 3 — Notificar al padre

  emit('closed', {})


NOTA: No recargar `remisiones-psicosocial-component`. La selección de checkboxes permanece hasta que el padre recargue tras RRM-05.
