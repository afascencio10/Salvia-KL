━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando alterna el switch "Asignar en dupla"
   Tipo: User Interaction
   ID: RRM-02
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: cambio en ToggleSwitch `asignarEnDupla`

Precondiciones:
  visible === true
  saving === false


PASO 1 — Actualizar modo

  asignarEnDupla = nuevo valor del switch


PASO 2 — Limpiar selección previa

  selectedProfessionalId = ''
  selectedDuplaId = ''
  optionsError = null


PASO 3 — Recargar opciones del select

  → Disparar RRM-03 según modo:
    SI asignarEnDupla === false → label "Profesional", optgroups ps/ts
    SI asignarEnDupla === true  → label "Dupla", lista enriquecida


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  UI — Cambio de etiqueta del select
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Switch | Label select | Placeholder |
|---|---|---|
| OFF | Profesional | Seleccionar profesional... |
| ON | Dupla | Seleccionar dupla... |
