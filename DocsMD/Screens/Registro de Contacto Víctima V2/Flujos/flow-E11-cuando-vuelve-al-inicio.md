━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Volver al Inicio"
   Tipo: User Interaction
   Función: resetForm() — botón del paso 'confirmacion'
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}

PASO 1 — Reiniciar el estado del wizard
  this.form = this.blankForm(this.form.captchaID)   // conserva el MISMO captchaID
  this.radicado = ''
  this.step = 'auth'
  (fechaActual y horaActual no se limpian: se sobrescriben en el próximo envío exitoso)

No hay llamada a backend. Permite iniciar otro reporte sin recargar la página.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| **L2** — el captchaID reutilizado ya fue consumido por el envío anterior (`captcha.Verify` lo borra): el segundo envío en la misma página recibe 400 de captcha y recién ahí un id nuevo; la imagen del captcha mostrará un id inexistente hasta entonces. Por lectura de código, no verificado en el navegador. Solución natural: pedir un captchaID nuevo al reiniciar | PASO 1 |
