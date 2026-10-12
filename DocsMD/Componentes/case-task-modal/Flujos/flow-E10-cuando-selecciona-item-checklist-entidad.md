━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando selecciona ítem del checklist entidad
   Tipo: User Interaction
   Función: toggleObligacion(itemId)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  itemId:  uuid de entity_case_obligation  → click en el chip
  form:    { items, seleccion }            → estado del componente
}


PASO 1 — Alternar selección

  SI form.seleccion[itemId] no existe:
    form.seleccion[itemId] = { selected: false, descripcion: "" }

  form.seleccion[itemId].selected = !form.seleccion[itemId].selected

  SI selected === false:
    form.seleccion[itemId].descripcion = ""


PASO 2 — Vue re-evalúa reactivamente

  SI selected === true:
    → Chip con fondo púrpura + check
    → Textarea visible debajo (placeholder de descripción, maxlength 500)

  SI selected === false:
    → Chip con borde gris
    → Textarea oculto

  formularioValido:
    ≥1 ítem selected Y cada selected tiene descripcion.trim() !== ""
    → habilita / deshabilita "Enviar solicitudes"

→ FIN EJECUCIÓN ✓
