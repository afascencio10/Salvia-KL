━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Volver al Inicio"
   Tipo: User Interaction
   Función: resetForm() · refreshCaptcha() — botón del paso 'confirmacion'
   Estado: ✅ Actualizado en WP-1 (T2, 2026-11-02): ya no reutiliza el captchaID consumido
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {}

PASO 1 — Reiniciar el estado del wizard
  this.form = this.blankForm(this.form.captchaID)   // provisional; el PASO 2 lo reemplaza
  this.radicado = '' ; this.copiado = false ; this.copiaError = false
  this.step = this.pasoInicial                      // 'usuario' o 'auth' (T8)
  (fechaActual y horaActual no se limpian: se sobrescriben en el próximo envío exitoso)

PASO 2 — 🆕 [T2] refreshCaptcha(): pedir un captchaID nuevo
  GET <ruta de la página>   Accept: application/json   // E-01 en modo JSON, sin tocar el backend
  SI 200: form.captchaID = respuesta.captchaID ; form.captchaSolution = ''
  SI falla: console.warn y se conserva el id anterior (el primer envío falla y el 400 trae uno nuevo, como antes)

Permite iniciar otro reporte sin recargar la página y con un captcha vigente (el anterior lo consume `captcha.VerifyString`).

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión | Paso afectado |
|---|---|
| ✅ **L2 resuelta (WP-1 T2)** — el id reutilizado ya estaba consumido; ahora `refreshCaptcha()` pide uno nuevo | PASO 2 |
