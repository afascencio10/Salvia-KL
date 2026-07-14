---
okf_version: "1.0"
type: Test_Plan
title: "Plan de Pruebas: Reporte consolidado de seguimientos"
description: "Matriz de casos de prueba (API y UI) para la descarga del Excel consolidado de seguimientos por rango de fechas."
owner: "@qa-squad"
status: draft
tags: [qa, testing, reportes]
dependencies:
  - Feature: 3-features/ReporteConsolidadoSeguimientos/index.md
last_updated: "2026-07-02"
---

# Estrategia de Testing (Reporte consolidado de seguimientos)

Los casos derivan del [PRD](/.knowledge/0-product/prd-reporte-consolidado-seguimientos.md) y del [contrato del endpoint](/.knowledge/3-features/ReporteConsolidadoSeguimientos/backend/endpoints/reporte-consolidado-seguimientos-post.md). Herramienta oficial: **Playwright** (API + E2E).

## Pruebas de API (Backend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Script |
| :--- | :--- | :--- | :--- |
| `TC-RPT-API-01` | `sv` o `ad` con rango válido con casos | `200`; `Content-Type` xlsx; `Content-Disposition` attachment; libro con 4 hojas | `tests/api/` |
| `TC-RPT-API-02` | Rango válido **sin casos** creados | `200`; xlsx válido con encabezados y sin filas de datos | `tests/api/` |
| `TC-RPT-API-03` | Falta `start_date` o `end_date` | `400 VALIDATION_FAILED` (JSON) | `tests/api/` |
| `TC-RPT-API-04` | `start_date > end_date` | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-RPT-API-05` | Rango > 366 días | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-RPT-API-06` | Formato de fecha inválido (`2026/13/40`) | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-RPT-API-07` | Rol sin permiso (ej. `op`, `us`) | `403 AUTH_INSUFFICIENT_ROLE` | `tests/api/` |
| `TC-RPT-API-08` | Sin sesión / sesión expirada | `401 AUTH_TOKEN_INVALID` | `tests/api/` |
| `TC-RPT-API-09` | Un caso con varios seguimientos | Cada seguimiento es una fila diferenciada en la hoja `Seguimientos` (`follow_up_id` únicos) | `tests/api/` |
| `TC-RPT-API-10` | Filtro por fecha de creación | Solo aparecen casos con `victim_case_creation_date` dentro del rango (`America/Bogota`) | `tests/api/` |
| `TC-RPT-API-11` | Dato registrado justo antes de la descarga | El xlsx refleja el valor más reciente (sin caché) | `tests/api/` |

## Pruebas E2E / UI (Frontend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Script |
| :--- | :--- | :--- | :--- |
| `TC-RPT-UI-01` | `sv`/`ad` ve el botón "Descargar reporte consolidado" en Lista de Casos | Botón visible | `tests/e2e/` |
| `TC-RPT-UI-02` | Rol sin permiso no ve el botón | Botón ausente | `tests/e2e/` |
| `TC-RPT-UI-03` | Abrir modal, elegir rango válido, confirmar | Se descarga un archivo `.xlsx` con el nombre esperado | `tests/e2e/` |
| `TC-RPT-UI-04` | Rango > 1 año en el modal | Error visible en el DOM; no se dispara la petición | `tests/e2e/` |
| `TC-RPT-UI-05` | Backend responde error | Mensaje de error visible en el modal; no se descarga | `tests/e2e/` |

## Pruebas Unitarias (Go)

- Validación del rango (obligatorio, `desde<=hasta`, ≤366 días) en el facade/servicio.
- Construcción del libro: nombres de hojas, encabezados y que cada seguimiento genere exactamente una fila.
- Relación por `case_icode` consistente entre hojas.

## Evidencia de ejecución

### Pruebas Unitarias (Go)

* **Caso: Pruebas unitarias de lógica y construcción de Excel**
  * **Comando**: `go test -v ./salvia/service/report_service_test.go ./salvia/service/report_service.go ./salvia/service/followups_report_excel.go`
  * **Fecha**: 2026-07-02
  * **Resultado**: `PASS`
  * **Salida**:
    ```text
    === RUN   TestReportService_GenerateConsolidatedFollowUpsReport_Empty
    --- PASS: TestReportService_GenerateConsolidatedFollowUpsReport_Empty (0.00s)
    === RUN   TestReportService_GenerateConsolidatedFollowUpsReport_WithData
    --- PASS: TestReportService_GenerateConsolidatedFollowUpsReport_WithData (0.00s)
    === RUN   TestReportService_GenerateConsolidatedFollowUpsReport_Error
    --- PASS: TestReportService_GenerateConsolidatedFollowUpsReport_Error (0.00s)
    PASS
    ok  	command-line-arguments	0.569s
    ```

* **Caso: Pruebas unitarias de validación y binding en Facade**
  * **Comando**: `go test -v ./salvia/facades/ReporteConsolidadoFacade_test.go ./salvia/facades/ReporteConsolidadoFacade.go`
  * **Fecha**: 2026-07-02
  * **Resultado**: `PASS`
  * **Salida**:
    ```text
    === RUN   TestDownloadConsolidatedReport_CompileAndVerify
    --- PASS: TestDownloadConsolidatedReport_CompileAndVerify (0.00s)
    PASS
    ok  	command-line-arguments	0.513s
    ```

### Validación de Specs OKF

* **Caso: Validación global de referencias y sintaxis OKF**
  * **Comando**: `make validate`
  * **Fecha**: 2026-07-02
  * **Resultado**: `PASS`
  * **Salida**:
    ```text
    python3 .scripts/validate_okf.py --check-refs
    ======================================================================
      OKF Validator — Open Knowledge Format
      Mode: STANDARD | Refs: ON
    ======================================================================
      ✅ .knowledge/0-product/prd-reporte-consolidado-seguimientos.md (Product_Requirement)
      ✅ .knowledge/3-features/ReporteConsolidadoSeguimientos/backend/endpoints/reporte-consolidado-seguimientos-post.md (API_Endpoint)
      ✅ .knowledge/3-features/ReporteConsolidadoSeguimientos/backend/services/reporte-consolidado-service.md (Service_Logic)
      ✅ .knowledge/3-features/ReporteConsolidadoSeguimientos/frontend/components/BotonReporteConsolidado.md (UI_Component)
      ✅ .knowledge/3-features/ReporteConsolidadoSeguimientos/index.md (Project)
      ✅ .knowledge/3-features/ReporteConsolidadoSeguimientos/testing/reporte-consolidado-test-plan.md (Test_Plan)
      ...
    ──────────────────────────────────────────────────────────────────────
      Archivos: 25 | Errores: 0 | Warnings: 0
    ──────────────────────────────────────────────────────────────────────
      🚀 Validación OKF exitosa.
    ```
