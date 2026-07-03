━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando abre el modal de reasignación
   Tipo: Lifecycle / User Interaction
   ID: RRM-01
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: padre llama `open(remisiones)` tras recibir E-15

Precondiciones:
  remisiones.length > 0
  Ninguna en status cerrado (garantizado por E-14)

INPUT: {
  remisiones: Array<PsychosocialListItem>  → payload de E-15
}


PASO 1 — Validar payload

  SI remisiones vacío o no es array:
    → No abrir modal
    → TERMINAR


PASO 2 — Inicializar estado

  visible = true
  remisiones = copia del array recibido
  asignarEnDupla = false
  selectedProfessionalId = ''
  selectedDuplaId = ''
  professionalGroups = []
  duplaOptions = []
  optionsError = null
  saveError = null
  saving = false


PASO 3 — Cargar opciones iniciales

  → Disparar RRM-03 (modo profesional — switch OFF por defecto)


PASO 4 — Focus accesibilidad

  → Enfocar primer control interactivo del modal (switch o select)


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  RESPONSABILIDAD DEL PADRE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

```html
<reasignar-remisiones-modal
  ref="reasignarModal"
  @closed="onReasignarModalClosed"
  @reassigned="onReasignarCompletado"
/>
```

```javascript
onReasignarRemisiones: function(payload) {
  this.$refs.reasignarModal.open(payload.remisiones);
}
```
