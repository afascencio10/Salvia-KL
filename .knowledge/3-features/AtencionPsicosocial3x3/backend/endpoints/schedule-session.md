---
okf_version: "1.0"
type: API_Endpoint
title: "Agendar sesión de Atención Psicosocial"
description: "Crea la sesión terapéutica en team_contact tras aceptar el Consentimiento Informado. Soporta sesión inmediata (redirige al formulario de atención) o agendada con fecha/hora flexible."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, session, team-contact]
resource: "POST /api/v1/psychosocial/{psicosocialId}/sessions"
related_models:
  - 2-data-dictionary/team-contact-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
last_updated: "2026-07-09"
---

# Agendar sesión (`POST /api/v1/psychosocial/{psicosocialId}/sessions`)

Crea una fila en `salvia.team_contact` (`is_psico_session = true`, `is_completed = false`) para la sesión de atención. Solo válido si el Consentimiento Informado fue aceptado.

- **`immediate = true`**: `scheduled_date = now()`; la respuesta incluye `redirectUrl` al formulario de atención (🔴 futuro).
- **`immediate = false`**: usa `scheduled_date` + `scheduled_time` provistos (fecha/hora flexible); sin redirección.

## Parámetros de Ruta

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `psicosocialId` | `uuid` | ID del proceso `psychosocial_support`. |

## Contrato de Petición (Request) — JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["immediate"],
  "additionalProperties": false,
  "properties": {
    "immediate": { "type": "boolean" },
    "scheduled_date": { "type": ["string", "null"], "format": "date" },
    "scheduled_time": { "type": ["string", "null"], "pattern": "^([01]\\d|2[0-3]):[0-5]\\d(:[0-5]\\d)?$" }
  },
  "allOf": [
    {
      "if": { "properties": { "immediate": { "const": false } } },
      "then": { "required": ["scheduled_date", "scheduled_time"] }
    }
  ]
}
```

### Ejemplo de Request (agendada)
```json
{ "immediate": false, "scheduled_date": "2026-07-14", "scheduled_time": "15:30" }
```

### Ejemplo de Request (inmediata)
```json
{ "immediate": true }
```

## Contrato de Respuesta (Response) — JSON Schema

### `201 Created`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["session"],
  "properties": {
    "session": {
      "type": "object",
      "required": ["id", "psicosocialId", "isPsicoSession", "isCompleted"],
      "properties": {
        "id": { "type": "string" },
        "psicosocialId": { "type": "string" },
        "caseId": { "type": "string" },
        "scheduledDate": { "type": ["string", "null"], "format": "date-time" },
        "scheduledTime": { "type": ["string", "null"] },
        "isPsicoSession": { "type": "boolean" },
        "isCompleted": { "type": "boolean" },
        "formSubmissionId": { "type": ["string", "null"] }
      }
    },
    "redirectUrl": {
      "type": ["string", "null"],
      "description": "Presente solo si immediate=true. Ruta al formulario de atención (🔴 futuro)."
    }
  }
}
```

### Ejemplo de Response (201, inmediata)
```json
{
  "session": {
    "id": "t-77",
    "psicosocialId": "9a2f...",
    "caseId": "c-123",
    "scheduledDate": "2026-07-09T10:20:00-05:00",
    "scheduledTime": null,
    "isPsicoSession": true,
    "isCompleted": false,
    "formSubmissionId": null
  },
  "redirectUrl": "/salvia/psicosocial/sesion/t-77"
}
```

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `201` | — | Sesión creada | Si `redirectUrl` presente, navegar; si no, cerrar modal con éxito |
| `400` | `VALIDATION_FAILED` | Input inválido (falta fecha/hora en agendada; formato de hora inválido) | Mostrar errores de campo |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Mostrar mensaje de error |
| `404` | `RESOURCE_NOT_FOUND` | `psicosocialId` no existe | Mostrar error |
| `409` | `RESOURCE_CONFLICT` | Consentimiento no aceptado para este proceso | Reabrir Consentimiento |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: JWT válido en cookie HttpOnly.
- **Autorización**: Roles `ps`, `ts`. Permiso RBAC: `schedule_psychosocial_session`.

## Dependencia 🔴 PENDIENTE

El destino de `redirectUrl` (formulario de atención) y el llenado de `form_submission_id` dependen del **formulario de atención** (spec futuro).
