---
okf_version: "1.0"
type: Project
title: "Feature: Reporte consolidado de contactos (reportes)"
description: "Exportación a Excel de todos los reportes (victim_contact) creados en un rango de fechas, con la información completa de Form1 y Form2 en una sola hoja plana con cabeceras amigables, desde la pantalla de reportes del supervisor."
owner: "@salvia-squad"
status: active
tags: [feature, reportes, contactos]
dependencies:
  - Feature: 3-features/ReporteConsolidadoSeguimientos/index.md
last_updated: "2026-07-17"
---

# Feature: Reporte consolidado de contactos (reportes)

## Descripción

Un botón **"Descargar reportes"** en la **pantalla de reportes** del supervisor ("Reportes Salvia", `get_victim_cases_sv.html`), ubicado junto al botón de crear nuevo caso, abre un modal donde el usuario elige un **rango de fechas** (máx. 1 año, sobre la **fecha de creación** del reporte). El backend genera un **Excel de una sola hoja plana** con una fila por reporte (`VictimContact`), incluyendo **toda** la información asociada sin importar si el reporte es de **Formulario 1** (reporte propio de la víctima) o **Formulario 2** (reporte de un tercero): las columnas del formulario que no aplica quedan vacías. Solo disponible para **Supervisor (`sv`)** y **Administrador (`ad`)**.

Permite a SALVIA consultar y analizar los reportes sin depender de acceso directo a la base de datos. Sigue el patrón del feature [Reporte consolidado de seguimientos](/.knowledge/3-features/ReporteConsolidadoSeguimientos/index.md) ya implementado.

## Documentos Relacionados

### Backend — Endpoints
| Endpoint | Spec | Estado |
| :--- | :--- | :---: |
| `POST /api/v1/reportes/contactos-consolidado` | [reporte-contactos-consolidado-post.md](/.knowledge/3-features/ReporteConsolidadoContactos/backend/endpoints/reporte-contactos-consolidado-post.md) | ✅ |

### Backend — Servicios
| Servicio | Spec | Estado |
| :--- | :--- | :---: |
| Generación del Excel de contactos | [reporte-contactos-service.md](/.knowledge/3-features/ReporteConsolidadoContactos/backend/services/reporte-contactos-service.md) | ✅ |

### Frontend — Componentes
| Componente | Spec | Estado |
| :--- | :--- | :---: |
| Botón + modal de descarga de reportes | [BotonReporteContactos.md](/.knowledge/3-features/ReporteConsolidadoContactos/frontend/components/BotonReporteContactos.md) | ✅ |

### Testing
| Plan | Spec | Estado |
| :--- | :--- | :---: |
| Plan de pruebas | [reporte-contactos-test-plan.md](/.knowledge/3-features/ReporteConsolidadoContactos/testing/reporte-contactos-test-plan.md) | ✅ |

## Historias de Usuario

| ID | Historia | Estado |
| :--- | :--- | :---: |
| US-RPC-01 | Como Supervisor/Administrador, quiero descargar un Excel con todos los reportes creados en un rango de fechas y su información completa (Form1 o Form2), para analizarlos sin acceso directo a la base de datos. | 📝 Especificado |

## Decisiones de diseño (entrevista 2026-07-17)

| Tema | Decisión |
| :--- | :--- |
| Estructura del Excel | **Una sola hoja plana**: una fila por reporte con todas las columnas `contacto_*` + `tipo_formulario_reporte` + `f1_*` + `f2_*`. |
| Roles | `sv` + `ad` (permiso nuevo `report_contacts_consolidated`). |
| Rango máximo | 1 año (366 días), validado en cliente y backend. |
| Filtro | Fecha de **creación** del reporte (`victim_contact_creation_date`), zona `America/Bogota`. |
| Estados | Se incluyen reportes en **todos** los estados. |
| Rango sin reportes | **No** se descarga archivo: el backend responde `400 VALIDATION_FAILED` y el modal muestra el mensaje. |
| Valores amigables | Enums traducidos a español (Locale), estado como descripción legible, fechas `dd/mm/yyyy hh:mm`, sí/no en vez de claves; municipio como código + nombre. |
| Auditoría | Sin registro especial de descargas (igual que el reporte de seguimientos). |
| Frescura | El Excel se construye consultando la BD al momento del request (sin caché). |

## Fuentes de datos (código existente)

| Sección del Excel | Modelo / origen | Archivo |
| :--- | :--- | :--- |
| Datos del contacto | `VictimContactDTO` (`victim_contact`) | `src/salvia/dao/VictimContactDAO.go` |
| Formulario 1 (reporte propio) | `VictimContactForm1DTO` | `src/salvia/dao/VictimContactForm1DAO.go` |
| Formulario 2 (reporte de tercero) | `VictimContactForm2DTO` (+ enums `VictimCaseForm2EnumsDTO`) | `src/salvia/dao/VictimContactForm2DAO.go` |
| Nombre de municipio | Tabla `towns` (`town_code` → `town_name`) | `src/internal/repository/location_repository.go` |
| Patrón de Excel existente | Builder con `excelize` | `src/salvia/service/followups_report_excel.go` |

## Pendientes / fuera de alcance

- 🔴 **Origen web/app**: el requerimiento original pedía diferenciar si el reporte se hizo desde el portal web o desde la app, pero `victim_contact` no tiene ningún campo que registre el canal de origen. **Decisión del usuario (2026-07-17): no incluir esa columna por ahora.** Para una iteración futura se evaluará derivarlo de `VictimContactGeneralUser` (app autenticada) vs. captcha (web anónimo), o agregar una columna de origen con migración.

> [!NOTE]
> Estado: `✅` = implementado. Spec aprobado el 2026-07-21; implementación completada el mismo día (endpoint activo en `ReportController`, permiso `report_contacts_consolidated` en `Menu.go`, botón + modal en la pantalla de reportes). Validado end-to-end el 2026-07-21: API con sesión real de supervisor (`TC-RPC-API-01/02/03/04/05/06/08`) y suite **E2E Playwright** en Chromium (`TC-RPC-UI-01/03/04/05/06`, 5/5 PASS) con confirmación visual del botón, el modal y el estado sin-reportes. Evidencia completa en el [plan de pruebas](/.knowledge/3-features/ReporteConsolidadoContactos/testing/reporte-contactos-test-plan.md).
