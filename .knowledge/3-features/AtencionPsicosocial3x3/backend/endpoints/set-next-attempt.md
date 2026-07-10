---
okf_version: "1.0"
type: API_Endpoint
title: "Fijar fecha/hora del próximo intento de contacto (3a)"
description: "Agendamiento manual flexible del punto 3a: la profesional fija la fecha/hora exacta del siguiente intento de contacto. Persiste en psychosocial_support.next_contact_attempt_at. No usa reagendamiento automático ni follow_up_v2."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, scheduling]
resource: "PUT /api/v1/psychosocial/{psicosocialId}/next-attempt"
related_models:
  - 2-data-dictionary/contact-attempts-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
  - "src/internal/repository/psychosocial_support_repository.go"
last_updated: "2026-07-09"
---

# Fijar próximo intento (`PUT /api/v1/psychosocial/{psicosocialId}/next-attempt`)

Corresponde a la acción **3a** del modal de acciones. A diferencia de Seguimientos (posponer automático +1/+2/+3 días), aquí la profesional elige **fecha y hora exactas** del próximo intento de contacto, según sus turnos rotativos y la disponibilidad de la ciudadana.

Persiste el valor en **`psychosocial_support.next_contact_attempt_at`** (columna nueva). **No** modifica `follow_up_v2` ni crea sesión en `team_contact`.

## Parámetros de Ruta

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `psicosocialId` | `uuid` | ID del proceso `psychosocial_support`. |

## Contrato de Petición (Request) — JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["next_contact_attempt_at"],
  "additionalProperties": false,
  "properties": {
    "next_contact_attempt_at": {
      "type": "string",
      "format": "date-time",
      "description": "Fecha y hora exactas del próximo intento de contacto."
    }
  }
}
```

### Ejemplo de Request
```json
{ "next_contact_attempt_at": "2026-07-10T14:00:00-05:00" }
```

## Contrato de Respuesta (Response) — JSON Schema

### `200 OK`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["psicosocialId", "nextContactAttemptAt"],
  "properties": {
    "psicosocialId": { "type": "string" },
    "nextContactAttemptAt": { "type": "string", "format": "date-time" }
  }
}
```

### Ejemplo de Response (200)
```json
{ "psicosocialId": "9a2f...", "nextContactAttemptAt": "2026-07-10T14:00:00-05:00" }
```

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Próximo intento agendado | Cerrar modal de acciones y refrescar |
| `400` | `VALIDATION_FAILED` | Fecha/hora inválida o ausente | Mostrar error de campo |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Mostrar mensaje de error |
| `404` | `RESOURCE_NOT_FOUND` | `psicosocialId` no existe | Mostrar error |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: JWT válido en cookie HttpOnly.
- **Autorización**: Roles `ps`, `ts`. Permiso RBAC: `set_psychosocial_next_attempt`.

## Nota de diseño

Se decidió una columna dedicada `next_contact_attempt_at` en `psychosocial_support` (en vez de reutilizar `follow_up_v2.scheduled_date`) porque el flujo 3x3 psicosocial es autónomo del motor de Seguimientos.
