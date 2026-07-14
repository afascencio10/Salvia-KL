---
okf_version: "1.0"
type: Service_Logic
title: "Servicio: Generación del Excel consolidado de seguimientos"
description: "Consulta los casos creados en un rango de fechas y arma un libro Excel multi-hoja (Casos, Seguimientos, Respuestas de formulario, Timeline) con excelize."
owner: "@backend-squad"
status: active
tags: [backend, service, reportes]
dependencies:
  - Endpoint: 3-features/ReporteConsolidadoSeguimientos/backend/endpoints/reporte-consolidado-seguimientos-post.md
last_updated: "2026-07-08"
code_refs:
  - src/salvia/service/report_service.go
  - src/internal/repository/report_repository.go
  - src/salvia/service/followups_report_excel.go
  - src/salvia/service/report_service_test.go
---

# Servicio: Generación del Excel consolidado de seguimientos

## Responsabilidades

1. Recibir el rango `[start_date, end_date]` ya validado y el rol/alcance del solicitante.
2. Consultar los **casos creados** en el rango (por `victim_case_creation_date`, `America/Bogota`), respetando la visibilidad del rol.
3. Para esos casos, cargar en lote: seguimientos (`FollowUpV2`), respuestas de formulario (`FormSubmission`→`Answer`) y eventos de timeline (`CaseTimelineEvent`).
4. Construir un libro **Excel** con `github.com/xuri/excelize/v2` y devolverlo como stream (`WriteTo`) para que el endpoint lo escriba en la respuesta.

## Estructura del libro (multi-hoja)

| Hoja | Grano | Columnas (mínimo) |
| :--- | :--- | :--- |
| `Casos` | 1 fila por caso | `case_icode`, fecha de creación, datos de víctima (documento, nombre), datos de hechos (de `VictimCaseForm1`/`Form2`), preguntas de riesgo pareja y no pareja, estado |
| `Seguimientos` | 1 fila por seguimiento | `case_icode`, `follow_up_id`, `sequence_number`, `scheduled_date`, `completed_at`, `status`, `team`, `risk_status`, `summary` |
| `Respuestas de formulario` | 1 fila por respuesta (o por seguimiento con columnas por pregunta) | `case_icode`, `follow_up_id`, `form_submission_id`, pregunta, respuesta |
| `Timeline` | 1 fila por evento | `case_icode`, `date`, `category`, `type`, `description`, `actor_name` |

- La **clave de relación** entre hojas es el `case_icode` (`victim_case_i_code`), que es también el `case_id` de `FollowUpV2` y `CaseTimelineEvent`.
- Cada seguimiento queda **diferenciado** por su fila (`follow_up_id` + `sequence_number`).
- La primera fila de cada hoja es de **encabezados**; el número de columnas es estable aunque no haya datos.

## Interfaz Pública (Contract)

```go
// En ReportService (src/salvia/service/report_service.go)
type FollowUpsReportRange struct {
    StartDate time.Time
    EndDate   time.Time
}

// Devuelve el .xlsx en memoria/stream listo para escribir en la respuesta HTTP.
GenerateConsolidatedFollowUpsReport(
    ctx context.Context,
    r FollowUpsReportRange,
    scope ReportScope, // rol + filtros de visibilidad del solicitante
) (*excelize.File, error)
```

## Errores que Lanza

| Error | Código Catálogo | Cuándo |
| :--- | :--- | :--- |
| `ErrInvalidRange` | `VALIDATION_FAILED` | Rango inválido que llegue sin validar (defensa) |
| `ErrReportQuery` | `SYSTEM_INTERNAL_ERROR` | Falla al consultar casos/seguimientos/timeline |
| `ErrReportBuild` | `SYSTEM_INTERNAL_ERROR` | Falla al construir/serializar el `.xlsx` |

## Rendimiento y consideraciones

- Consultas **en lote por lista de `case_icode`** (evitar N+1): un query por hoja filtrando `case_id IN (...)`.
- La cota de **≤ 366 días** acota el volumen; aun así, con muchos casos el archivo puede ser grande. Para esta primera versión la generación es **síncrona** (el request espera el archivo).
- Si en producción el volumen resulta alto, se evaluará generación **asíncrona** (job + descarga diferida) en una iteración posterior — **fuera de alcance** de este spec.
- No se cachea: cada llamada refleja el estado actual de la BD.

## Dependencias

- `github.com/xuri/excelize/v2` (ya en `go.mod`).
- `report_repository.go` (GORM) para las consultas en lote.
