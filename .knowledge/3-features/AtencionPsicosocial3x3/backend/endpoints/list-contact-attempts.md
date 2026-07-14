---
okf_version: "1.0"
type: API_Endpoint
title: "Listar intentos de contacto psicosocial (historial del card)"
description: "Devuelve el historial de intentos de contacto de un proceso psicosocial más los contadores y las fechas de último/próximo intento. Alimenta el card 'Intentos de contacto 3x3'."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, contact-attempts]
resource: "GET /api/v1/psychosocial/{psicosocialId}/contact-attempts"
related_models:
  - 2-data-dictionary/contact-attempts-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
  - "src/internal/repository/contact_attempt_repository.go"
last_updated: "2026-07-09"
---

# Listar intentos de contacto (`GET /api/v1/psychosocial/{psicosocialId}/contact-attempts`)

Devuelve el historial de un proceso ordenado **descendente por `attempt_at` (más nuevo → más viejo)**, más los contadores y las fechas de último/próximo intento. El `sequenceNumber` permanece **cronológico** (1 = intento más antiguo) aunque el arreglo se entregue del más nuevo al más viejo. Es la fuente de datos del card **`psychosocial-contact-card`**, que es autosuficiente al montarse en cualquier pantalla.

## Parámetros de Ruta

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `psicosocialId` | `uuid` | ID del proceso `psychosocial_support`. |

## Contrato de Respuesta (Response) — JSON Schema

### `200 OK`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["attempts", "counters", "lastAttemptAt", "nextContactAttemptAt"],
  "properties": {
    "attempts": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "sequenceNumber", "wasAnswered", "attemptAt"],
        "properties": {
          "id": { "type": "string" },
          "sequenceNumber": { "type": "integer", "description": "Nº de intento (1-based) para el título 'Intento #N'." },
          "wasAnswered": { "type": "boolean" },
          "note": { "type": ["string", "null"] },
          "consentGiven": { "type": ["boolean", "null"] },
          "attemptAt": { "type": "string", "format": "date-time" }
        }
      }
    },
    "counters": {
      "type": "object",
      "required": ["dailyFailedCount", "totalCount", "distinctDaysCount", "dailyThresholdReached", "maxAttemptsReached", "closureEligible"],
      "properties": {
        "dailyFailedCount": { "type": "integer" },
        "totalCount": { "type": "integer" },
        "distinctDaysCount": { "type": "integer" },
        "dailyThresholdReached": { "type": "boolean" },
        "maxAttemptsReached": { "type": "boolean" },
        "closureEligible": { "type": "boolean" }
      }
    },
    "lastAttemptAt": { "type": ["string", "null"], "format": "date-time", "description": "attempt_at del intento más reciente (tile 'Último intento realizado')." },
    "nextContactAttemptAt": { "type": ["string", "null"], "format": "date-time", "description": "psychosocial_support.next_contact_attempt_at (tile 'Próximo intento' y fila 'PROGRAMADO')." },
    "processStatus": { "type": "string" }
  }
}
```

### Ejemplo de Response (200)
```json
{
  "attempts": [
    { "id": "a2", "sequenceNumber": 2, "wasAnswered": true, "note": "Sí contestó, pide que lo llamen otro día.", "consentGiven": true, "attemptAt": "2026-07-09T18:46:00-05:00" },
    { "id": "a1", "sequenceNumber": 1, "wasAnswered": false, "note": "No contestó, celular apagado.", "consentGiven": null, "attemptAt": "2026-07-09T10:15:00-05:00" }
  ],
  "counters": { "dailyFailedCount": 1, "totalCount": 2, "distinctDaysCount": 1, "dailyThresholdReached": false, "maxAttemptsReached": false, "closureEligible": false },
  "lastAttemptAt": "2026-07-09T18:46:00-05:00",
  "nextContactAttemptAt": "2026-07-10T10:00:00-05:00",
  "processStatus": "en_gestion"
}
```

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Historial devuelto | Renderizar card |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Ocultar card |
| `404` | `RESOURCE_NOT_FOUND` | `psicosocialId` no existe | Mostrar error |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: JWT válido en cookie HttpOnly.
- **Autorización**: Roles `ps`, `ts`. Permiso RBAC: `get_psychosocial_contact_attempts`.
- **PII**: `note` puede contener contexto sensible (Nivel 2). No loguear el cuerpo completo.
