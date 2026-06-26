# Login — Interfaz

Pantalla de autenticación de SOG Salvia. Permite a los usuarios ingresar con usuario, contraseña y CAPTCHA. Incluye un flujo de recuperación de contraseña por correo y un flujo de restablecimiento de contraseña mediante token (cuando el usuario llega desde un enlace de email). Al autenticarse correctamente, muestra un overlay de bienvenida y redirige al usuario a su ruta principal según su rol.

La pantalla tiene cuatro estados distintos visibles mediante overlays: formulario de login (default), recuperar contraseña, restablecer contraseña y confirmación de correo enviado.

---

## Archivos relevantes

| Archivo | Descripción |
|---|---|
| `src/frontend/html/security/login.html` | Template principal: formulario de login, overlays y lógica Vue |
| `src/frontend/js/login.js` | Scripts de apoyo para el login |

---

## Árbol de interfaz

```
Pantalla: Login
│
├── Layout de dos columnas
│   ├── Columna izquierda — Imagen de fondo decorativa (.login-bg)
│   └── Columna derecha — Formulario de login
│       │
│       ├── [v-if: errors.default.global_msg] Alerta de error global
│       │   └── Mensaje de error dismissible
│       │
│       ├── Logo de Salvia
│       ├── Título: "Ingresar a mi Salvia"
│       │
│       ├── Campo "Usuario" (tipo text)
│       │   └── [si error] Mensaje de error bajo el campo
│       │
│       ├── Campo "Contraseña" (tipo password)
│       │   └── [si error] Mensaje de error bajo el campo
│       │
│       ├── Bloque CAPTCHA
│       │   ├── Imagen del CAPTCHA (:src → /seguridad/login/captcha/:id/img)
│       │   ├── Audio del CAPTCHA (:src → /seguridad/login/captcha/:id/audio)
│       │   └── Campo "Solución CAPTCHA" (tipo text)
│       │       └── [si error] Mensaje de error bajo el campo
│       │
│       ├── Botón "Ingresar" (@click → submit('login'))
│       └── Link "¿Olvidé mi contraseña?" (@click → submit('openForgot'))
│
├── Overlay: Video de bienvenida [id: videoOverlay, display:none por defecto]
│   └── Contenedor vacío (reservado para video de bienvenida por rol)
│
├── Overlay: Recuperar contraseña [id: forgotOverlay, display:none por defecto]
│   ├── Título: "Recuperar contraseña"
│   ├── Texto explicativo
│   ├── Campo "Usuario" (tipo text)
│   └── Botones: [Volver → submit('closeForgot')] [Enviar → submit('forgot')]
│
├── Overlay: Restablecer contraseña [id: resetOverlay, display:none por defecto]
│   │   Se muestra automáticamente al cargar si resetToken.length > 4
│   ├── Título: "Restablecer contraseña"
│   ├── Texto explicativo
│   ├── Campo "Nueva contraseña" (tipo password)
│   ├── Campo "Repetir contraseña" (tipo password)
│   └── Botones: [Volver → submit('closeReset')] [Guardar → submit('reset')]
│
└── Overlay: Correo enviado [id: mailOverlay, display:none por defecto]
    ├── Mensaje de confirmación de envío de correo
    └── Botón: [Volver → submit('closeForgot')]
```
