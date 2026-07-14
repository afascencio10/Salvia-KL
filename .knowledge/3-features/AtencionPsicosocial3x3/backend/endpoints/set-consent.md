---
okf_version: "1.0"
type: API_Endpoint
title: "Registrar decisión de Consentimiento Informado"
description: "Persiste la decisión del Consentimiento Informado (acepta/no acepta) sobre un intento de contacto exitoso. Si no acepta, marca el proceso para cierre psicosocial (formulario de cierre, futuro)."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, consent]
resource: "PATCH /api/v1/contact-attempts/{id}/consent"
related_models:
  - 2-data-dictionary/contact-attempts-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
last_updated: "2026-07-09"
---

# Registrar decisión de consentimiento (`PATCH /api/v1/contact-attempts/{id}/consent`)

> Nota de ruta: el prefijo es `/api/v1/contact-attempts/...` (no `/api/v1/psychosocial/...`) para evitar el conflicto de Gin entre el parámetro `:psicosocialId` y un segmento estático hermano bajo `/psychosocial/`.

Fija `consent_given` sobre un `contact_attempts` **exitoso** (`was_answered = true`).

- **`true` (acepta)**: habilita el paso de agendamiento de sesión (ver [schedule-session](./schedule-session.md)).
- **`false` (no acepta)**: dispara el **cierre del proceso psicosocial**. El cierre efectivo requiere el **formulario de cierre** (🔴 futuro); mientras no exista, el endpoint devuelve `requiresClosureForm=true` y el frontend abre ese formulario (dependencia externa).

## Parámetros de Ruta

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `id` | `uuid` | ID del `contact_attempts` exitoso. |

## Contrato de Petición (Request) — JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["consent_given"],
  "additionalProperties": false,
  "properties": {
    "consent_given": { "type": "boolean" }
  }
}
```

### Ejemplo de Request
```json
{ "consent_given": false }
```

## Contrato de Respuesta (Response) — JSON Schema

### `200 OK`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["id", "consentGiven", "requiresClosureForm"],
  "properties": {
    "id": { "type": "string" },
    "consentGiven": { "type": "boolean" },
    "requiresClosureForm": {
      "type": "boolean",
      "description": "true cuando consent_given=false: el frontend debe abrir el formulario de cierre (🔴 futuro)."
    },
    "canScheduleSession": {
      "type": "boolean",
      "description": "true cuando consent_given=true."
    }
  }
}
```

### Ejemplo de Response (200)
```json
{ "id": "b0e1...", "consentGiven": false, "requiresClosureForm": true, "canScheduleSession": false }
```

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Decisión registrada | Si `true` abrir agendamiento; si `false` abrir formulario de cierre |
| `400` | `VALIDATION_FAILED` | Input inválido o intento no exitoso (`was_answered=false`) | Mostrar error |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Mostrar mensaje de error |
| `404` | `RESOURCE_NOT_FOUND` | `id` de intento no existe | Mostrar error |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: JWT válido en cookie HttpOnly.
- **Autorización**: Roles `ps`, `ts`. Permiso RBAC: `set_psychosocial_consent`.

## Dependencia 🔴 PENDIENTE

El **cierre efectivo** por "no acepta" depende del **formulario de cierre psicosocial** (spec futuro). Este endpoint solo registra la decisión y señala `requiresClosureForm`.
