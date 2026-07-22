---
okf_version: "1.0"
type: API_Endpoint
title: "Descargar reporte consolidado de contactos"
description: "Genera y devuelve un Excel (.xlsx) de una sola hoja con todos los reportes (victim_contact) creados en un rango de fechas, incluyendo la información completa de Form1 o Form2."
owner: "@backend-squad"
status: active
tags: [backend, api, reportes, contactos]
resource: "POST /api/v1/reportes/contactos-consolidado"
dependencies:
  - Feature: 3-features/ReporteConsolidadoContactos/index.md
  - Service: 3-features/ReporteConsolidadoContactos/backend/services/reporte-contactos-service.md
last_updated: "2026-07-17"
code_refs:
  - src/salvia/facades/ContactsReportFacade.go
  - src/salvia/facades/ContactsReportFacade_test.go
  - src/salvia/controller/report_controller.go
  - src/salvia/config/Menu.go
---

# Descargar reporte consolidado de contactos (`POST /api/v1/reportes/contactos-consolidado`)

Recibe un rango de fechas y devuelve un archivo **Excel** (`.xlsx`) como descarga (no JSON). El rango filtra los reportes por su **fecha de creación** (`victim_contact_creation_date`), en zona horaria `America/Bogota`. Sigue el mismo patrón que `POST /api/v1/reportes/seguimientos-consolidado` ([spec](/.knowledge/3-features/ReporteConsolidadoSeguimientos/backend/endpoints/reporte-consolidado-seguimientos-post.md)).

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
- El filtro se aplica como: `victim_contact_creation_date::date BETWEEN start_date AND end_date` (en `America/Bogota`).
- Se incluyen los reportes en **todos los estados** (nuevos, gestionados, cerrados/descartados).

### Ejemplo de Request
```json
{ "start_date": "2026-01-01", "end_date": "2026-06-30" }
```

## Contrato de Respuesta (Response)

### `200 OK` — archivo binario (no JSON)

| Header | Valor |
| :--- | :--- |
| `Content-Type` | `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` |
| `Content-Disposition` | `attachment; filename="reporte-contactos_<start_date>_<end_date>.xlsx"` |

El cuerpo es el `.xlsx` generado por el [servicio de generación](/.knowledge/3-features/ReporteConsolidadoContactos/backend/services/reporte-contactos-service.md): **una sola hoja** (`Reportes`) con una fila por reporte y todas las columnas definidas en el spec del servicio.

Si no se encuentran reportes creados en el rango de fechas seleccionado, el servidor **no genera ningún archivo**; devuelve `400 Bad Request` con el código `VALIDATION_FAILED` y un JSON explicativo, y el modal muestra "No hay reportes en el rango seleccionado" (decisión de la entrevista 2026-07-17).

> [!NOTE]
> El reporte se construye consultando la BD **en el momento del request** (sin caché): refleja la información más reciente.

## Códigos de Respuesta

| Status | Código de Error | Descripción | Acción del Cliente |
| :---: | :--- | :--- | :--- |
| `200` | — | Excel generado y descargado | Guardar el archivo |
| `400` | `VALIDATION_FAILED` | Rango vacío, formato inválido, `start_date > end_date`, rango > 366 días, o **cero reportes encontrados** | Mostrar el error en el modal |
| `401` | `AUTH_TOKEN_INVALID` | Sesión no válida / expirada | Redirigir a login |
| `403` | `AUTH_INSUFFICIENT_ROLE` | El rol no tiene el permiso `report_contacts_consolidated` | Ocultar el botón / mostrar mensaje |
| `500` | `SYSTEM_INTERNAL_ERROR` | Error al consultar datos o construir el Excel | Mostrar error genérico |

Los códigos se mapean al [Catálogo de Errores](/.knowledge/2-data-dictionary/error-catalog.md). En error, la respuesta es **JSON** (no archivo).

## Seguridad

- **Autenticación**: sesión server-side por cookie firmada (ver [ADR-001](/.knowledge/1-architecture/adrs/adr-001-stack-go-gin-vue-postgres.md)). **No** es JWT.
- **Autorización**: se verifica en el facade con
  `utils.CheckPermission(salvia_config.PermissionsByRole, "report_contacts_consolidated", s.CurrentRole, c)`.
  Roles permitidos: **`ad` (Administrador)** y **`sv` (Supervisor)**. Ver [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md). El permiso es **nuevo** y debe agregarse a `PermissionsByRole` en `src/salvia/config/Menu.go`.
- **Datos sensibles**: el archivo contiene PII de nivel sensible (documento, teléfono, dirección, descripción de hechos — ver [data-classification](/.knowledge/6-security/data-classification.md)). **No** debe loguearse el contenido del archivo ni los datos de las filas; solo metadatos del request (usuario, rango, conteo).
- **Auditoría**: sin registro persistente de descargas (decisión de la entrevista 2026-07-17, consistente con el reporte de seguimientos).
- **Registro en el ruteo**: la ruta se agrega en `ReportController.RegisterRoutes` (`src/salvia/controller/report_controller.go`) bajo el grupo `/api/v1`, junto a `/reportes/seguimientos-consolidado`.
