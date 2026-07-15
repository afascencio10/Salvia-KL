---
okf_version: "1.0"
type: API_Endpoint
title: "Autenticación de Usuario General"
description: "Valida las credenciales del usuario, crea la sesión del lado del servidor y retorna el mapa de navegación y menús según el rol."
owner: "@security-guild"
status: active
tags: [backend, api, auth]
resource: "POST /seguridad/login"
code_refs:
  - src/security/facades/GeneralUserFacade.go
last_updated: "2026-07-02"
---

# Autenticación de Usuario General (`POST /seguridad/login`)

## Descripción
Este endpoint autentica a los usuarios en el sistema. Valida las credenciales contra la tabla `security.general_user`, compara la contraseña con bcrypt (costo 10), y valida la solución del captcha. 
Al autenticar correctamente, crea una sesión (`userData` uuid) y retorna el mapa de navegación (`nav_rules`) del rol activo asignado al usuario.

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Login exitoso y redirección configurada | Redirigir a la ruta asignada en `nav_rules` |
| `400` | `login_fail` / `security_general_user_captcha_fail` | Credenciales incorrectas o error de captcha | Mostrar el error y reiniciar captcha |
| `403` | — | Sesión denegada (ej. usuario inactivo o sin roles válidos) | Mostrar error |
| `500` | — | Error interno del servidor | Mostrar mensaje de error |

## Redirección por Roles (Navigation Rules)

El endpoint resuelve la redirección de los roles activos de la siguiente manera:

| Rol | Código | Redirección por Defecto |
| :--- | :---: | :--- |
| Administrador | `ad` | `/seguridad/usuarios` (o ruta de configuración) |
| Supervisor | `sv` | `/salvia/casos` |
| Operador | `op` | `/salvia/casos` |
| Operador de riesgo | `ro` | `/salvia/casos` |
| Operador territorial departamental | `do` | `/salvia/casos` |
| Operador territorial nacional | `no` | `/salvia/casos` |
| Entidad | `et` | `/salvia/casos` |
| Usuario externo | `us` | `/salvia/casos` |
| Psicólogo | `ps` | `/salvia/casos` |
| Trabajador social | `ts` | `/salvia/casos` |
| Agente de notificaciones | `an` | `/salvia/notificaciones` |
| Enlace territorial | `en` | `/salvia/mis-barreras` |
