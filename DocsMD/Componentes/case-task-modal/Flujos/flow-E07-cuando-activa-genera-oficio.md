━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando activa "¿Genera oficio?"
   Tipo: User Interaction
   Función: onToggleGeneraOficio()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Aplica a: formulario 'gestion_llamada' únicamente.

INPUT: {
  form.generaOficio:  valor nuevo del switch  → Switch "¿Genera oficio?"
                      true | false
}


PASO 1 — Actualizar estado del switch

  form.generaOficio = nuevoValor


SI form.generaOficio === false:
  → form.asunto    = ""
  → form.rutaKofax = ""
  → Vue oculta reactivamente los campos Asunto y Ruta al archivo
    ([v-if form.generaOficio] → false)
  → CONTINÚA FLUJO GENERAL

SI form.generaOficio === true:
  → Vue muestra reactivamente los campos Asunto y Ruta al archivo
  → Ambos quedan vacíos y requeridos para la validación del formulario
  → CONTINÚA FLUJO GENERAL


PASO 2 — Vue re-evalúa formularioValido

  La computed property `formularioValido` incluye la condición:
  SI generaOficio: asunto requerido Y rutaKofax requerido

  → BtnConfirmar se habilita o deshabilita según el nuevo estado de validación

→ FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
(sin gaps — lógica completamente definida)
