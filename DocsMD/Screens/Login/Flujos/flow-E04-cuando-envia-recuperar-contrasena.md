# flow-E04 — Cuando presiona "Enviar" en recuperar contraseña

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona "Enviar" en recuperar contraseña
   Tipo: User Interaction
   Función: submit('forgot')
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  data.user.login:  nombre de usuario ingresado  → input del overlay forgotOverlay
}

PASO 1 — Enviar solicitud de recuperación al API

POST /seguridad/login/forgot

// payload:
{
  "user": {
    "login": data.user.login    // usuario que solicita recuperación
  }
}

→ resultado: confirmación de envío o errores de validación

PASO 2 — Manejar respuesta

SI status === 200:
  → Ocultar overlay de éxito
     document.getElementById('successOverlay').style.display = 'none'
  → Cerrar overlay de recuperar contraseña
     document.getElementById('forgotOverlay').style.display = 'none'
  → Mostrar overlay de confirmación de correo enviado
     document.getElementById('mailOverlay').style.display = 'flex'
  → El usuario ve la confirmación de que el correo fue enviado

SI status === 400:
  → Asignar errores: this.errors = response
  → Vue muestra los errores por campo en el formulario del overlay

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
| Variable / decisión                                              | Paso afectado |
|------------------------------------------------------------------|---------------|
| ¿El backend valida que el usuario exista antes de enviar        | PASO 1        |
| el correo, o siempre responde 200 por seguridad?                |               |
```
