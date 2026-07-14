---
okf_version: "1.0"
type: API_Endpoint
title: "Completar cierre psicosocial"
description: "Al finalizar el formulario de Cierre de proceso psicosocial, transiciona el proceso a cerrado o en_devolucion según el Motivo de cierre registrado en la submission."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, closure]
resource: "POST /api/v1/psychosocial/{psicosocialId}/close"
related_models:
  - 2-data-dictionary/closure-form-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
last_updated: "2026-07-10"
---

# Completar cierre (`POST /api/v1/psychosocial/{psicosocialId}/close`)

Se invoca cuando el formulario de cierre se completó (`form-completed` del `dinamic-form`). Lee el **Motivo de cierre** de la submission y transiciona el estado del proceso psicosocial. El cierre afecta **solo** al proceso psicosocial (no al caso de víctima).

## Parámetros de Ruta

| Parámetro | Tipo | Descripción |
| :--- | :--- | :--- |
| `psicosocialId` | `uuid` | Proceso `psychosocial_support`. |

## Contrato de Petición — JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["submission_id"],
  "additionalProperties": false,
  "properties": {
    "submission_id": { "type": "string" }
  }
}
```

## Contrato de Respuesta — JSON Schema

### `200 OK`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["psicosocialId", "status", "motivo"],
  "properties": {
    "psicosocialId": { "type": "string" },
    "status": { "type": "string", "enum": ["cerrado", "en_devolucion"] },
    "motivo": { "type": "string" }
  }
}
```

### Ejemplo (200)
```json
{ "psicosocialId": "9a2f...", "status": "cerrado", "motivo": "imposibilidad_contacto_3x3" }
```

## Mapeo Motivo → Estado (RN)

| Motivo (`value`) | Estado resultante |
| :--- | :--- |
| `imposibilidad_contacto_3x3` | `cerrado` |
| `cumplimiento_objetivos` | `cerrado` |
| `cumplimiento_esquema` | `cerrado` |
| `no_consentimiento` | `en_devolucion` |
| `desistimiento_proceso` | `en_devolucion` |

> Actualmente solo se disparan `no_consentimiento` (→ `en_devolucion`) e `imposibilidad_contacto_3x3` (→ `cerrado`). Los demás motivos quedan mapeados para cuando existan sus disparadores.

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Proceso cerrado/devuelto | Cerrar modal, refrescar card, quitar de la lista |
| `400` | `VALIDATION_FAILED` | `submission_id` ausente o sin Motivo respondido | Mostrar error |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Mostrar mensaje |
| `404` | `RESOURCE_NOT_FOUND` | `psicosocialId` o submission no existe | Mostrar error |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: JWT válido en cookie HttpOnly.
- **Autorización**: Roles `ps`, `ts`. Permiso RBAC: `close_psychosocial_process`.
