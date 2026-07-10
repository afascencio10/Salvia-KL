---
okf_version: "1.0"
type: API_Endpoint
title: "Iniciar formulario de cierre psicosocial"
description: "Crea una form_submission del formulario de Cierre de proceso psicosocial, preselecciona (editable) el Motivo de cierre según el disparador (no consentimiento / imposibilidad de contacto) y devuelve el submission_id para renderizar dinamic-form."
owner: "@backend-squad"
status: active
tags: [backend, api, psicosocial, 3x3, closure, dynamic-form]
resource: "POST /api/v1/psychosocial/{psicosocialId}/init-closure-form"
related_models:
  - 2-data-dictionary/closure-form-model.md
code_refs:
  - "src/salvia/controller/psychosocial_contact_controller.go"
  - "src/salvia/service/psychosocial_3x3_service.go"
last_updated: "2026-07-10"
---

# Iniciar formulario de cierre (`POST /api/v1/psychosocial/{psicosocialId}/init-closure-form`)

Prepara el formulario dinámico de **Cierre de proceso psicosocial** (`form_id = fcc7dc8d-835b-4559-9283-2ea8b8e6092b`) para un proceso, análogo al `init-closure-form` de seguimientos. Crea la `form_submission`, **preselecciona** el Motivo de cierre según el disparador  y lo devuelve; el frontend lo muestra preseleccionado pero **editable** (dinamic-form no bloquea por pregunta).

## Parámetros

| Parámetro | Ubicación | Tipo | Descripción |
| :--- | :--- | :--- | :--- |
| `psicosocialId` | ruta | `uuid` | Proceso `psychosocial_support`. |
| `reason` | query | `enum` | Disparador: `no_consentimiento` \| `imposibilidad_contacto_3x3`. Determina el Motivo preseleccionado (editable). |

## Contrato de Respuesta — JSON Schema

### `201 Created`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["submission_id", "form_id", "lockedMotivo"],
  "properties": {
    "submission_id": { "type": "string" },
    "form_id": { "type": "string" },
    "lockedMotivo": {
      "type": "string",
      "description": "Valor del Motivo de cierre preseleccionado (editable; coincide con reason).",
      "enum": ["no_consentimiento", "imposibilidad_contacto_3x3"]
    }
  }
}
```

### Ejemplo (201)
```json
{
  "submission_id": "b71c...",
  "form_id": "fcc7dc8d-835b-4559-9283-2ea8b8e6092b",
  "lockedMotivo": "imposibilidad_contacto_3x3"
}
```

## Comportamiento

1. Crea `salvia.form_submission (form_id = fcc7dc8d…)`.
2. Pre-responde la pregunta **Motivo de cierre** (`1c6d97ac…`) con la opción del `reason` (`no_consentimiento` → `9c25c5f0…`; `imposibilidad_contacto_3x3` → `d2c6065a…`).
3. Devuelve `submission_id`. El frontend renderiza `<dinamic-form :form-id :submission-id :form-state>` con el Motivo preseleccionado (editable).

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `201` | — | Submission creada | Abrir modal con `<dinamic-form>` |
| `400` | `VALIDATION_FAILED` | `reason` inválido/ausente | Mostrar error |
| `401` | `AUTH_TOKEN_INVALID` | No autenticado | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | Rol distinto de `ps`/`ts` | Mostrar mensaje |
| `404` | `RESOURCE_NOT_FOUND` | `psicosocialId` no existe | Mostrar error |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error interno | Mostrar error genérico |

## Seguridad

- **Autenticación**: JWT válido en cookie HttpOnly.
- **Autorización**: Roles `ps`, `ts`. Permiso RBAC: `init_psychosocial_closure_form`.
