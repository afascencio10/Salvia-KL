# flow-E02 — Cuando presiona "Ingresar"

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Ingresar"
   Tipo: User Interaction
   Función: submit('login')
   Disparadores: click en botón "Ingresar" | keyup.enter en el formulario
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  data.user.login:           usuario ingresado          → input del formulario
  data.user.pass:            contraseña ingresada       → input del formulario
  data.user.captchaSolution: solución del CAPTCHA       → input del formulario
  data.user.captchaID:       ID del CAPTCHA activo      → inyectado por Go / renovado en error
}

PASO 1 — Enviar credenciales al API

POST /seguridad/login

// payload:
{
  "user": {
    "login":           data.user.login,
    "pass":            data.user.pass,
    "captchaSolution": data.user.captchaSolution,
    "captchaID":       data.user.captchaID
  }
}

→ resultado: reglas de navegación por rol o errores de validación

PASO 2 — Manejar respuesta

SI status === 200:
  → Mostrar overlay de bienvenida
     document.getElementById('videoOverlay').style.display = 'flex'
  → Guardar reglas de navegación recibidas
     enums.nav_rules = response
  → Esperar 1ms y redirigir al usuario

┌─────────────────────────────────────────────────────┐
│  SUB-FLUJO: Determinar ruta de redirección          │
└─────────────────────────────────────────────────────┘

  SI nav_rules["default"] existe:
    → location.assign(nav_rules["default"])

  SI NO, si nav_rules["login"] existe:
    → location.assign(nav_rules["login"])

  → FIN SUB-FLUJO → TERMINAR ejecución

SI status === 400:
  → Ocultar overlay de éxito
     document.getElementById('successOverlay').style.display = 'none'
  → Asignar errores: this.errors = response
  → Mostrar mensaje global de error en el bloque de alerta
     document.getElementById("fail_msg").innerHTML = errors.default.global_msg

  SI response.default.captchaID existe (el CAPTCHA fue incorrecto):
    → Limpiar solución ingresada: data.user.captchaSolution = ''
    → Actualizar ID del CAPTCHA: data.user.captchaID = response.default.captchaID
    → Renovar imagen del CAPTCHA:
       document.getElementById("captchaIMG").src =
         /seguridad/login/captcha/{nuevo captchaID}/img
    → Crear nuevo objeto de audio del CAPTCHA:
       data.audio = new Audio(/seguridad/login/captcha/{nuevo captchaID}/audio)
    → Vue re-evalúa captchaAudioURL reactivamente con el nuevo ID

SI otro status:
  → Ocultar overlay de éxito
  → alert("Se produjo un error interno, verifique que el usuario tenga un rol asignado.")

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                         | Paso afectado |
|-------------------------------------------------------------|---------------|
| ¿Qué estructura exacta tiene nav_rules en la respuesta 200? | PASO 2        |
| ¿Existe validación de campos vacíos en el frontend          | PASO 1        |
| antes de hacer el POST, o la hace solo el backend?          |               |
```
