---
okf_version: "1.0"
type: API_Endpoint
title: "Registrar intento de contacto psicosocial"
description: "Registra un intento de contacto (exitoso o fallido) de un proceso de Atención Psicosocial en la tabla contact_attempts. En contacto exitoso pasa el proceso a en_gestion. Devuelve los contadores para gobernar el flujo 3x3 (umbral diario, tope 50, elegibilidad de cierre)."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, contact-attempts]
resource: "POST /api/v1/psychosocial/{psicosocialId}/contact-attempts"
related_models:
  - 2-data-dictionary/contact-attempts-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
  - "src/internal/repository/contact_attempt_repository.go"
last_updated: "2026-07-09"
---

# Registrar intento de contacto (`POST /api/v1/psychosocial/{psicosocialId}/contact-attempts`)

Registra un intento de contacto de un proceso psicosocial. Sirve tanto para "Sí contestó" (`was_answered=true`) como para "No contestó" (`was_answered=false`).

- **Éxito (`was_answered=true`)**: crea la fila, deja evento en el timeline y transiciona `psychosocial_support.status` a `en_gestion`. El `consent_given` se fija después vía [set-consent](./set-consent.md). El frontend abre luego el modal de Consentimiento.
- **Fallo (`was_answered=false`)**: crea la fila y el proceso se mantiene `abierto`. El frontend cierra el modal (no hay segundo paso).
- **`note`** (opcional en ambos casos): nota libre capturada en el modal "¿Contestó?" (ej. "No contestó, celular apagado" / "Sí contestó, pide que lo llamen otro día").
- **Limpieza del próximo intento**: si el proceso tenía un `next_contact_attempt_at` cuya fecha es **igual o anterior** al día del intento registrado (`next_contact_attempt_at::date <= attempt_at::date`), se limpia a `NULL` — el intento planeado ya se ejecutó (sin importar si contestó o no; cubre que el agente lo hiciera el mismo día o días después de lo programado). Solo se conserva si lo programado es a **futuro** respecto al intento.

## Parámetros de Ruta

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `psicosocialId` | `uuid` | ID del proceso `psychosocial_support`. |

## Contrato de Petición (Request) — JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["was_answered"],
  "additionalProperties": false,
  "properties": {
    "was_answered": { "type": "boolean" },
    "note": { "type": ["string", "null"], "maxLength": 2000, "description": "Nota libre (opcional) capturada en el modal ¿Contestó?; aplica a exitosos y fallidos." },
    "attempt_at": {
      "type": ["string", "null"],
      "format": "date-time",
      "description": "Fecha/hora del intento (editable). Si es null, el servidor usa now()."
    }
  }
}
```

### Ejemplo de Request (fallido)
```json
{
  "was_answered": false,
  "note": "No contestó, celular apagado.",
  "attempt_at": "2026-07-09T10:15:00-05:00"
}
```

### Ejemplo de Request (exitoso)
```json
{ "was_answered": true, "note": "Sí contestó, pide que lo llamen otro día.", "attempt_at": null }
```

## Contrato de Respuesta (Response) — JSON Schema

### `201 Created`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["attempt", "counters"],
  "properties": {
    "attempt": {
      "type": "object",
      "required": ["id", "psicosocialId", "wasAnswered", "attemptAt"],
      "properties": {
        "id": { "type": "string" },
        "psicosocialId": { "type": "string" },
        "caseId": { "type": "string" },
        "wasAnswered": { "type": "boolean" },
        "note": { "type": ["string", "null"] },
        "consentGiven": { "type": ["boolean", "null"] },
        "attemptAt": { "type": "string", "format": "date-time" }
      }
    },
    "counters": {
      "type": "object",
      "required": ["dailyFailedCount", "totalCount", "distinctDaysCount", "dailyThresholdReached", "maxAttemptsReached", "closureEligible"],
      "properties": {
        "dailyFailedCount": { "type": "integer", "description": "Intentos fallidos hoy (attempt_at::date = hoy)." },
        "totalCount": { "type": "integer", "description": "Intentos totales del proceso." },
        "distinctDaysCount": { "type": "integer", "description": "Días distintos con al menos un intento." },
        "dailyThresholdReached": { "type": "boolean", "description": "true si dailyFailedCount >= 3." },
        "maxAttemptsReached": { "type": "boolean", "description": "true si totalCount >= 50." },
        "closureEligible": { "type": "boolean", "description": "true si totalCount >= 9 y distinctDaysCount >= 3." }
      }
    },
    "processStatus": { "type": "string", "description": "Estado resultante del proceso (p.ej. en_gestion)." }
  }
}
```

### Ejemplo de Response (201, fallido)
```json
{
  "attempt": {
    "id": "b0e1...",
    "psicosocialId": "9a2f...",
    "caseId": "c-123",
    "wasAnswered": false,
    "note": "No contestó, celular apagado.",
    "consentGiven": null,
    "attemptAt": "2026-07-09T10:15:00-05:00"
  },
  "counters": {
    "dailyFailedCount": 3,
    "totalCount": 3,
    "distinctDaysCount": 1,
    "dailyThresholdReached": true,
    "maxAttemptsReached": false,
    "closureEligible": false
  },
  "processStatus": "abierto"
}
```

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `201` | — | Intento registrado | Actualizar contadores; si `dailyThresholdReached` abrir modal de acciones; si exitoso abrir Consentimiento |
| `400` | `VALIDATION_FAILED` | Input inválido (falta `was_answered`, `attempt_at` mal formada) | Mostrar errores de campo |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Mostrar mensaje de error |
| `404` | `RESOURCE_NOT_FOUND` | `psicosocialId` no existe | Mostrar error |
| `409` | `RESOURCE_CONFLICT` | Tope de 50 intentos alcanzado (regla de negocio) | Deshabilitar "Añadir intento" |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: Requiere JWT válido en cookie HttpOnly.
- **Autorización**: Roles permitidos: `ps`, `ts`. Permiso RBAC: `register_psychosocial_contact_attempt`.
- **PII**: `note` puede contener contexto sensible (Nivel 2). No loguear el cuerpo completo.
