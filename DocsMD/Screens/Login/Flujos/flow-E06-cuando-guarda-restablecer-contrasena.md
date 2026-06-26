# flow-E06 — Cuando presiona "Guardar" en restablecer contraseña

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Guardar" en restablecer contraseña
   Tipo: User Interaction
   Función: submit('reset')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  data.user.pass:       nueva contraseña ingresada       → input del overlay resetOverlay
  data.user.pass2:      confirmación de contraseña        → input del overlay resetOverlay
  data.user.resetToken: token de restablecimiento         → inyectado por Go (.resetPswd)
}

PASO 1 — Enviar nueva contraseña al API

POST /seguridad/login/reset

// payload:
{
  "user": {
    "pass":       data.user.pass,         // nueva contraseña
    "pass2":      data.user.pass2,        // confirmación de contraseña
    "resetToken": data.user.resetToken    // token recibido por correo
  }
}

→ resultado: confirmación de restablecimiento o errores de validación

PASO 2 — Manejar respuesta

SI status === 200:
  → Cerrar overlay de restablecer contraseña
     document.getElementById('resetOverlay').style.display = 'none'
  → Esperar 3 segundos y ocultar overlay de éxito
     setTimeout(() => successOverlay.style.display = 'none', 3000)
  → El usuario queda en el formulario de login para ingresar con su nueva contraseña

SI status === 400:
  → Asignar errores: this.errors = response
  → Vue muestra los errores por campo en el formulario del overlay

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| ¿El backend valida que pass === pass2 o lo hace el frontend?     | PASO 1        |
| ¿El token tiene expiración? ¿Qué responde el backend si venció? | PASO 1        |
```
