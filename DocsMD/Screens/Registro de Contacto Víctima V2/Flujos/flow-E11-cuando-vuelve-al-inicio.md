━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Volver al Inicio"
   Tipo: User Interaction
   Función: resetForm() — botón del paso 'confirmacion'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}

PASO 1 — Reiniciar el estado del wizard
  form.value = blankForm()
  numRadicado.value = ''
  step.value = 'auth'

Permite radicar un segundo reporte en la misma sesión de navegador sin
recargar la página. No hay llamada a backend.
