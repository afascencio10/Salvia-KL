---
okf_version: "1.0"
type: Error_Catalog
title: "Catálogo de Errores Estándar"
description: "Taxonomía unificada de códigos de error del sistema con mensajes, acciones del cliente, y mapeo a HTTP status codes."
owner: "@backend-squad"
status: active
tags: [errors, api, standards, backend, frontend]
last_updated: "2026-06-23"
---

# Catálogo de Errores Estándar

Todos los endpoints de la API deben usar estos códigos de error. El frontend consume el `code` para decidir qué mensaje mostrar al usuario.

## Formato de Error Response

```json
{
  "error": {
    "code": "DOMAIN_ERROR_NAME",
    "message": "Mensaje legible para el usuario (nunca exponer detalles internos).",
    "status": 401,
    "details": {},
    "request_id": "req_abc123def456"
  }
}
```

## Reglas

1. **`code`**: Siempre en formato `DOMAIN_ERROR_NAME` (UPPER_SNAKE_CASE).
2. **`message`**: Orientado al usuario final. **Nunca** incluir stack traces, nombres de tablas, o queries SQL.
3. **`details`**: Objeto opcional con información adicional (ej. campos con errores de validación).
4. **`request_id`**: UUID generado por el middleware para trazabilidad en logs.

## Errores de Autenticación (`AUTH_*`)

| Código | HTTP | Mensaje (Usuario) | Cuándo se usa |
| :--- | :---: | :--- | :--- |
| `AUTH_INVALID_INPUT` | 400 | "El formato del email o contraseña no es válido." | Request body no pasa validación de schema |
| `AUTH_INVALID_CREDENTIALS` | 401 | "Las credenciales son incorrectas." | Email no existe O password incorrecto |
| `AUTH_TOKEN_EXPIRED` | 401 | "Tu sesión ha expirado. Inicia sesión nuevamente." | JWT `exp` pasó |
| `AUTH_TOKEN_INVALID` | 401 | "Token de autenticación inválido." | JWT con firma inválida o malformado |
| `AUTH_TOKEN_REVOKED` | 401 | "Tu sesión fue cerrada. Inicia sesión nuevamente." | Token en blacklist (logout previo) |
| `AUTH_RATE_LIMITED` | 429 | "Demasiados intentos. Intenta de nuevo en {minutes} minutos." | Rate limiter activado |
| `AUTH_ACCOUNT_LOCKED` | 403 | "Tu cuenta ha sido bloqueada por seguridad. Contacta soporte." | 5+ intentos fallidos consecutivos |
| `AUTH_INSUFFICIENT_ROLE` | 403 | "No tienes permisos para realizar esta acción." | RBAC: rol insuficiente |

> [!WARNING]
> `AUTH_INVALID_CREDENTIALS` se usa tanto cuando el email no existe como cuando el password es incorrecto. **Nunca** diferenciar estos casos para evitar enumeración de usuarios.

## Errores de Validación (`VALIDATION_*`)

| Código | HTTP | Mensaje (Usuario) | Cuándo se usa |
| :--- | :--- | :--- | :--- |
| `VALIDATION_FAILED` | 400 | "Algunos campos contienen errores." | Input del usuario no pasa schema |
| `VALIDATION_REQUIRED_FIELD` | 400 | "El campo '{field}' es obligatorio." | Campo requerido ausente |
| `VALIDATION_INVALID_FORMAT` | 400 | "El campo '{field}' tiene un formato inválido." | Formato incorrecto (email, fecha, etc.) |
| `VALIDATION_TOO_LONG` | 400 | "El campo '{field}' excede el largo máximo de {max} caracteres." | Excede maxLength del schema |

### Details para errores de validación

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Algunos campos contienen errores.",
    "status": 400,
    "details": {
      "fields": [
        { "field": "email", "code": "VALIDATION_INVALID_FORMAT", "message": "Formato de email inválido." },
        { "field": "password", "code": "VALIDATION_TOO_SHORT", "message": "La contraseña debe tener al menos 8 caracteres." }
      ]
    }
  }
}
```

## Errores de Recursos (`RESOURCE_*`)

| Código | HTTP | Mensaje (Usuario) | Cuándo se usa |
| :--- | :--- | :--- | :--- |
| `RESOURCE_NOT_FOUND` | 404 | "El recurso solicitado no fue encontrado." | Entidad no existe en DB |
| `RESOURCE_ALREADY_EXISTS` | 409 | "Ya existe un recurso con estos datos." | Violación de UNIQUE constraint |
| `RESOURCE_CONFLICT` | 409 | "El recurso fue modificado por otro usuario. Recarga la página." | Optimistic locking conflict |
| `RESOURCE_GONE` | 410 | "Este recurso ya no está disponible." | Entidad eliminada permanentemente |

## Errores de Sistema (`SYSTEM_*`)

| Código | HTTP | Mensaje (Usuario) | Cuándo se usa |
| :--- | :--- | :--- | :--- |
| `SYSTEM_INTERNAL_ERROR` | 500 | "Ocurrió un error interno. Intenta de nuevo más tarde." | Error inesperado no manejado |
| `SYSTEM_SERVICE_UNAVAILABLE` | 503 | "El servicio está temporalmente no disponible." | Dependencia caída (DB, Redis, etc.) |
| `SYSTEM_TIMEOUT` | 504 | "La operación tardó demasiado. Intenta de nuevo." | Timeout en operación upstream |

> [!CAUTION]
> Los errores `SYSTEM_*` **siempre** deben loggearse con stack trace completo en el server. El usuario solo ve el mensaje genérico.

## Mapeo Frontend — Cómo consumir errores

```text
// api/errorHandler.ts
function handleApiError(error: ApiError): void {
  switch (error.code) {
    case 'AUTH_INVALID_CREDENTIALS':
      showBanner({ type: 'error', message: error.message });
      break;
    case 'AUTH_TOKEN_EXPIRED':
      router.push('/login');
      break;
    case 'AUTH_RATE_LIMITED':
      showBanner({ type: 'warning', message: error.message });
      disableForm(error.details.retry_after_seconds);
      break;
    case 'VALIDATION_FAILED':
      setFieldErrors(error.details.fields);
      break;
    default:
      showBanner({ type: 'error', message: 'Ocurrió un error inesperado.' });
      logger.error('Unhandled API error', { error });
  }
}
```
