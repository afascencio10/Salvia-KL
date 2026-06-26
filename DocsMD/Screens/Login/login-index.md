# Login — Index

## Navegación

| Documento | Link |
|---|---|
| Interfaz | [login-interface.md](./login-interface.md) |

---

## Resumen de eventos

| # | Evento | Tipo | Flujo |
|---|---|---|---|
| E01 | Cuando carga la pantalla | Lifecycle | [📄 flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md) |
| E02 | Cuando presiona "Ingresar" | User Interaction | [📄 flow-E02-cuando-presiona-ingresar.md](./Flujos/flow-E02-cuando-presiona-ingresar.md) |
| E03 | Cuando presiona "¿Olvidé mi contraseña?" | User Interaction | — |
| E04 | Cuando presiona "Enviar" en recuperar contraseña | User Interaction | [📄 flow-E04-cuando-envia-recuperar-contrasena.md](./Flujos/flow-E04-cuando-envia-recuperar-contrasena.md) |
| E05 | Cuando presiona "Volver" en recuperar contraseña | User Interaction | — |
| E06 | Cuando presiona "Guardar" en restablecer contraseña | User Interaction | [📄 flow-E06-cuando-guarda-restablecer-contrasena.md](./Flujos/flow-E06-cuando-guarda-restablecer-contrasena.md) |
| E07 | Cuando presiona "Volver" en restablecer contraseña | User Interaction | — |

---

## Inventario de eventos

---

**Nombre del evento:** Cuando carga la pantalla
**Tipo:** Lifecycle
**Descripción:** Se ejecuta al montar el componente Vue. Establece el título de la ventana, muestra el contenedor principal y verifica si existe un token de restablecimiento en la URL. Si el token tiene más de 4 caracteres, muestra automáticamente el overlay de restablecer contraseña.
**Requerido:** Sí

📄 [Ver flujo → flow-E01-cuando-carga-pantalla.md](./Flujos/flow-E01-cuando-carga-pantalla.md)

---

**Nombre del evento:** Cuando presiona "Ingresar"
**Tipo:** User Interaction
**Descripción:** El usuario presiona el botón "Ingresar" o pulsa Enter dentro del formulario. Envía usuario, contraseña y solución del CAPTCHA al endpoint de login. Si es exitoso, muestra el overlay de bienvenida y redirige según las reglas de navegación del rol. Si falla, muestra errores y renueva el CAPTCHA.
**Requerido:** Sí

📄 [Ver flujo → flow-E02-cuando-presiona-ingresar.md](./Flujos/flow-E02-cuando-presiona-ingresar.md)

---

**Nombre del evento:** Cuando presiona "¿Olvidé mi contraseña?"
**Tipo:** User Interaction
**Descripción:** El usuario presiona el link de contraseña olvidada. Muestra el overlay de recuperar contraseña (`forgotOverlay`).
**Requerido:** Sí

---

**Nombre del evento:** Cuando presiona "Enviar" en recuperar contraseña
**Tipo:** User Interaction
**Descripción:** El usuario ingresa su nombre de usuario y presiona "Enviar". Llama al endpoint de forgot para generar y enviar el correo de restablecimiento. Si es exitoso, cierra el overlay de recuperar contraseña y muestra el overlay de confirmación de correo enviado. Si falla, muestra errores en el formulario.
**Requerido:** Sí

📄 [Ver flujo → flow-E04-cuando-envia-recuperar-contrasena.md](./Flujos/flow-E04-cuando-envia-recuperar-contrasena.md)

---

**Nombre del evento:** Cuando presiona "Volver" en recuperar contraseña
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Volver" en el overlay de recuperar contraseña o en el overlay de confirmación. Oculta ambos overlays (`forgotOverlay` y `mailOverlay`) y regresa al formulario de login.
**Requerido:** Sí

---

**Nombre del evento:** Cuando presiona "Guardar" en restablecer contraseña
**Tipo:** User Interaction
**Descripción:** El usuario ingresa la nueva contraseña y su confirmación, luego presiona "Guardar". Envía el token de reset y la nueva contraseña al endpoint de reset. Si es exitoso, cierra el overlay y espera 3 segundos antes de ocultar el overlay de éxito. Si falla, muestra errores en el formulario.
**Requerido:** Sí

📄 [Ver flujo → flow-E06-cuando-guarda-restablecer-contrasena.md](./Flujos/flow-E06-cuando-guarda-restablecer-contrasena.md)

---

**Nombre del evento:** Cuando presiona "Volver" en restablecer contraseña
**Tipo:** User Interaction
**Descripción:** El usuario presiona "Volver" en el overlay de restablecer contraseña. Oculta el overlay de reset (`resetOverlay`) y regresa al formulario de login.
**Requerido:** Sí

---

## Checklist de validación

- [x] ¿Cada acción del usuario puede ser manejada?
- [x] ¿La carga inicial de datos está cubierta?
- [x] ¿Todos los envíos de formulario están incluidos?
- [ ] ¿Hay actualizaciones en tiempo real? — No aplica en esta pantalla
