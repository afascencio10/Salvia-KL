---
okf_version: "1.0"
type: API_Endpoint
title: "Descargar reporte consolidado de seguimientos"
description: "Genera y devuelve un Excel (.xlsx) con el historial consolidado de seguimientos de los casos creados en un rango de fechas."
owner: "@backend-squad"
status: active
tags: [backend, api, reportes]
resource: "POST /api/v1/reportes/seguimientos-consolidado"
dependencies:
  - Feature: 3-features/ReporteConsolidadoSeguimientos/index.md
  - Service: 3-features/ReporteConsolidadoSeguimientos/backend/services/reporte-consolidado-service.md
last_updated: "2026-07-03"
code_refs:
  - src/salvia/facades/ConsolidatedReportFacade.go
  - src/salvia/controller/report_controller.go
  - src/salvia/facades/ConsolidatedReportFacade_test.go
  - src/salvia/config/Menu.go
---

# Descargar reporte consolidado de seguimientos (`POST /api/v1/reportes/seguimientos-consolidado`)

Recibe un rango de fechas y devuelve un archivo **Excel** (`.xlsx`) como descarga (no JSON). El rango filtra los casos por su **fecha de creación** (`victim_case_creation_date`), en zona horaria `America/Bogota`.

## Contrato de Petición (Request) — JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["start_date", "end_date"],
  "additionalProperties": false,
  "properties": {
    "start_date": {
      "type": "string",
      "format": "date",
      "description": "Fecha inicial del rango (inclusive), formato YYYY-MM-DD."
    },
    "end_date": {
      "type": "string",
      "format": "date",
      "description": "Fecha final del rango (inclusive), formato YYYY-MM-DD."
    }
  }
}
```

### Reglas de validación
- `start_date` y `end_date` son obligatorias y deben ser fechas válidas `YYYY-MM-DD`.
- `start_date` <= `end_date`.
- La ventana no puede exceder **366 días** (máximo 1 año).
- El filtro se aplica como: `victim_case_creation_date::date BETWEEN start_date AND end_date` (en `America/Bogota`).

### Ejemplo de Request
```json
{ "start_date": "2026-01-01", "end_date": "2026-06-30" }
```

## Contrato de Respuesta (Response)

### `200 OK` — archivo binario (no JSON)

| Header | Valor |
| :--- | :--- |
| `Content-Type` | `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` |
| `Content-Disposition` | `attachment; filename="reporte-casos_<start_date>_<end_date>.xlsx"` |

El cuerpo es el `.xlsx` generado por el [servicio de generación](/.knowledge/3-features/ReporteConsolidadoSeguimientos/backend/services/reporte-consolidado-service.md) con las hojas: **Casos**, **Registro**, **Seguimientos**, **Respuestas de formulario** y **Timeline**.

Si no se encuentran casos creados en el rango de fechas seleccionado, el servidor **no genera el Excel vacío**; en su lugar, devuelve un código `400 Bad Request` con el código de error `VALIDATION_FAILED` y un JSON explicativo, evitando la descarga inútil.

> [!NOTE]
> El reporte se construye consultando la BD **en el momento del request** (sin caché): refleja la información más reciente.

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Excel generado y descargado | Guardar el archivo |
| `400` | `VALIDATION_FAILED` | Rango de fechas vacío, formato inválido, `start_date > end_date`, rango > 366 días, o **cero casos encontrados** | Mostrar el error en el modal |
| `401` | `AUTH_TOKEN_INVALID` | Sesión no válida / expirada | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | El rol no tiene el permiso `report_followups_consolidated` | Ocultar el botón / mostrar mensaje |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error al consultar datos o construir el Excel | Mostrar error genérico |

Los códigos se mapean al [Catálogo de Errores](/.knowledge/2-data-dictionary/error-catalog.md). En error, la respuesta es **JSON** (no archivo).

## Seguridad

- **Autenticación**: sesión server-side por cookie firmada (ver [ADR-001](/.knowledge/1-architecture/adrs/adr-001-stack-go-gin-vue-postgres.md)). **No** es JWT.
- **Autorización**: se verifica en el facade con
  `utils.CheckPermission(salvia_config.PermissionsByRole, "report_followups_consolidated", s.CurrentRole, c)`.
  Roles permitidos: **`ad` (Administrador)** y **`sv` (Supervisor)**. Ver [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md).
- **Visibilidad de datos**: el reporte incluye solo los casos que el usuario puede ver según su rol/alcance (misma regla que la Lista de Casos). Contiene PII de víctimas (Nivel sensible, ver [data-classification](/.knowledge/6-security/data-classification.md)); no debe loguearse el contenido.
- **Registro en el ruteo**: el endpoint está registrado en `ReportController.RegisterRoutes` (`src/salvia/controller/report_controller.go:19`) bajo el grupo `/api/v1`, e instanciado en `src/main.go:193`. Implementado y activo en `develop` (verificado: `PermissionsByRole["report_followups_consolidated"]` en `Menu.go` ya contiene `{ad: true, sv: true}`, coincidiendo con esta sección).
