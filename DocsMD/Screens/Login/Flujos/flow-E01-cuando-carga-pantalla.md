# flow-E01 — Cuando carga la pantalla

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando carga la pantalla
   Tipo: Lifecycle
   Función: mounted()
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  windowTitle:  título de la ventana del navegador  → inyectado por Go (.windowTitle)
  resetToken:   token de restablecimiento            → inyectado por Go (.resetPswd)
  captchaID:    ID del CAPTCHA generado              → inyectado por Go (.CaptchaID)
  nav_rules:    reglas de navegación por rol         → inyectado por Go (.nav_rules)
}

PASO 1 — Establecer el título de la pestaña del navegador
  document.title = windowTitle

PASO 2 — Mostrar el contenedor principal de la app
  document.getElementById('app').style.display = 'block'

PASO 3 — Registrar las reglas de navegación como propiedad global de Vue
  home.config.globalProperties.enums = this.enums

PASO 4 — Verificar si hay un token de restablecimiento activo

SI resetToken.length > 4:
  → Mostrar overlay de restablecer contraseña
     document.getElementById('resetOverlay').style.display = 'flex'
  → El usuario ve el formulario de nueva contraseña directamente

SI resetToken.length <= 4:
  → No se muestra ningún overlay
  → El usuario ve el formulario de login normal

PASO 5 — Vue renderiza las URLs del CAPTCHA reactivamente
  captchaImgURL   = /seguridad/login/captcha/{captchaID}/img
  captchaAudioURL = /seguridad/login/captcha/{captchaID}/audio
  → El navegador carga la imagen y el audio del CAPTCHA desde el servidor
```
