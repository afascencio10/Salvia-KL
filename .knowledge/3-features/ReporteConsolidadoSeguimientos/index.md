---
okf_version: "1.0"
type: Project
title: "Feature: Reporte consolidado de seguimientos"
description: "Exportación a Excel del historial consolidado de seguimientos (detalle de caso, seguimientos, respuestas de formulario y timeline) de los casos creados en un rango de fechas, desde Lista de Casos."
owner: "@salvia-squad"
status: draft
tags: [feature, reportes]
dependencies:
  - PRD: 0-product/prd-reporte-consolidado-seguimientos.md
last_updated: "2026-07-02"
---

# Feature: Reporte consolidado de seguimientos

## Descripción

Un botón **"Descargar reporte consolidado"** en la pantalla **Lista de Casos** abre un modal donde el usuario elige un **rango de fechas** (máx. 1 año). El backend genera un **Excel multi-hoja** con, para cada caso **creado** en ese rango: el detalle del caso, sus seguimientos, las respuestas del formulario de cada seguimiento y el timeline de eventos. Solo disponible para **Supervisor** y **Administrador**.

Reemplaza el script SQL manual que dejó de funcionar tras el cambio de infraestructura (ver [PRD](/.knowledge/0-product/prd-reporte-consolidado-seguimientos.md)).

## Documentos Relacionados

### Producto
- [PRD: Reporte consolidado de seguimientos](/.knowledge/0-product/prd-reporte-consolidado-seguimientos.md)

### Backend — Endpoints
| Endpoint | Spec | Estado |
| :--- | :--- | :---: |
| `POST /api/v1/reportes/seguimientos-consolidado` | [reporte-consolidado-seguimientos-post.md](/.knowledge/3-features/ReporteConsolidadoSeguimientos/backend/endpoints/reporte-consolidado-seguimientos-post.md) | 📝 |

### Backend — Servicios
| Servicio | Spec | Estado |
| :--- | :--- | :---: |
| Generación del Excel consolidado | [reporte-consolidado-service.md](/.knowledge/3-features/ReporteConsolidadoSeguimientos/backend/services/reporte-consolidado-service.md) | 📝 |

### Frontend — Componentes
| Componente | Spec | Estado |
| :--- | :--- | :---: |
| Botón + modal de reporte consolidado | [BotonReporteConsolidado.md](/.knowledge/3-features/ReporteConsolidadoSeguimientos/frontend/components/BotonReporteConsolidado.md) | 📝 |

### Testing
| Plan | Spec | Estado |
| :--- | :--- | :---: |
| Plan de pruebas | [reporte-consolidado-test-plan.md](/.knowledge/3-features/ReporteConsolidadoSeguimientos/testing/reporte-consolidado-test-plan.md) | 📝 |

## Historias de Usuario

| ID | Historia | Estado |
| :--- | :--- | :---: |
| US-RPT-01 | Como Supervisor/Administrador, quiero descargar un Excel consolidado de los seguimientos de los casos creados en un rango de fechas, para consultar la evolución de la atención. | 📝 Especificado |

## Fuentes de datos (código existente)

| Sección del Excel | Modelo / origen | Archivo |
| :--- | :--- | :--- |
| Casos (detalle, víctima, hechos) | `VictimCase` + `VictimCaseForm1` + `VictimCaseForm2` | `src/internal/models/victim_case.go` |
| Seguimientos | `FollowUpV2` (`follow_up_v2`) | `src/internal/models/follow_up_v2.go` |
| Respuestas de formulario | `FormSubmission` → `Answer` (vía `FollowUpV2.FormSubmissionID`) | `src/internal/models/form_submission.go`, `src/internal/models/answer.go` |
| Timeline | `CaseTimelineEvent` | `src/internal/models/case_timeline_event.go` |
| Listado + scoping por rol | Consulta de casos | `src/internal/repository/cases_list_repository.go` |

> [!NOTE]
> Estado: `📝` = especificado, pendiente de implementar. La implementación no inicia hasta que este spec sea **aprobado**.
