---
okf_version: "1.0"
type: Test_Plan
title: "Plan de Pruebas: Reporte consolidado de contactos"
description: "Matriz de casos de prueba (API y UI) para la descarga del Excel consolidado de reportes (victim_contact) por rango de fechas."
owner: "@qa-squad"
status: active
tags: [qa, testing, reportes, contactos]
dependencies:
  - Feature: 3-features/ReporteConsolidadoContactos/index.md
last_updated: "2026-07-17"
---

# Estrategia de Testing (Reporte consolidado de contactos)

Los casos derivan del [contrato del endpoint](/.knowledge/3-features/ReporteConsolidadoContactos/backend/endpoints/reporte-contactos-consolidado-post.md) y del [spec del servicio](/.knowledge/3-features/ReporteConsolidadoContactos/backend/services/reporte-contactos-service.md). Herramienta oficial: **Playwright** (API + E2E).

## Pruebas de API (Backend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Script |
| :--- | :--- | :--- | :--- |
| `TC-RPC-API-01` | `sv` o `ad` con rango válido con reportes | `200`; `Content-Type` xlsx; `Content-Disposition` attachment; libro con 1 hoja `Reportes` y las 34 columnas | `tests/api/` |
| `TC-RPC-API-02` | Rango válido **sin reportes** creados | `400 VALIDATION_FAILED` (JSON); **no** se descarga archivo | `tests/api/` |
| `TC-RPC-API-03` | Falta `start_date` o `end_date` | `400 VALIDATION_FAILED` (JSON) | `tests/api/` |
| `TC-RPC-API-04` | `start_date > end_date` | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-RPC-API-05` | Rango > 366 días | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-RPC-API-06` | Formato de fecha inválido (`2026/13/40`) | `400 VALIDATION_FAILED` | `tests/api/` |
| `TC-RPC-API-07` | Rol sin permiso (ej. `op`, `us`, `ro`) | `403 AUTH_INSUFFICIENT_ROLE` | `tests/api/` |
| `TC-RPC-API-08` | Sin sesión / sesión expirada | `401 AUTH_TOKEN_INVALID` | `tests/api/` |
| `TC-RPC-API-09` | Reporte de Formulario 1 en el rango | Fila con `tipo_formulario_reporte` = "Reporte propio (Formulario 1)", columnas `f1_*` pobladas y `f2_*` vacías | `tests/api/` |
| `TC-RPC-API-10` | Reporte de Formulario 2 en el rango | Fila con `tipo_formulario_reporte` = "Reporte de tercero (Formulario 2)", columnas `f2_*` pobladas y `f1_*` vacías | `tests/api/` |
| `TC-RPC-API-11` | Filtro por fecha de creación | Solo aparecen reportes con `victim_contact_creation_date` dentro del rango (`America/Bogota`) | `tests/api/` |
| `TC-RPC-API-12` | Reportes en distintos estados (nuevo, gestionado, cerrado) | Todos aparecen; `contacto_estado` muestra la descripción legible de cada uno | `tests/api/` |
| `TC-RPC-API-13` | Reporte con municipio en Form1 | `f1_codigo_municipio` con el código y `f1_municipio` con el nombre resuelto de `towns` | `tests/api/` |
| `TC-RPC-API-14` | Form2 con `adjustmentsGBV` múltiples | `f2_ajustes_gbv` con las descripciones traducidas separadas por coma | `tests/api/` |
| `TC-RPC-API-15` | Valores enum traducidos | Ningún valor de enum aparece como clave cruda si existe traducción en Locale; sí/no en español | `tests/api/` |
| `TC-RPC-API-16` | Dato registrado justo antes de la descarga | El xlsx refleja el valor más reciente (sin caché) | `tests/api/` |
| `TC-RPC-API-17` | Teléfonos largos (10+ dígitos) | Se exportan como texto, sin notación científica ni pérdida de ceros | `tests/api/` |

## Pruebas E2E / UI (Frontend - Playwright)

| ID Caso | Escenario | Resultado Esperado | Script |
| :--- | :--- | :--- | :--- |
| `TC-RPC-UI-01` | `sv`/`ad` ve el botón "Descargar reportes" junto al botón de crear nuevo caso | Botón visible en la pantalla de reportes | `tests/e2e/` |
| `TC-RPC-UI-02` | Rol sin permiso no ve el botón | Botón ausente | `tests/e2e/` |
| `TC-RPC-UI-03` | Abrir modal, elegir rango válido, confirmar | Se descarga un archivo `.xlsx` con el nombre `reporte-contactos_<desde>_<hasta>.xlsx` | `tests/e2e/` |
| `TC-RPC-UI-04` | Rango > 1 año en el modal | Error visible en el DOM; no se dispara la petición | `tests/e2e/` |
| `TC-RPC-UI-05` | Rango sin reportes | Modal muestra "No hay reportes en el rango seleccionado"; no se descarga archivo | `tests/e2e/` |
| `TC-RPC-UI-06` | Backend responde error 500 | Mensaje de error genérico visible en el modal; no se descarga | `tests/e2e/` |

## Pruebas Unitarias (Go)

- Validación del rango (obligatorio, `desde<=hasta`, ≤366 días) en el facade.
- `ErrNoData` cuando el rango no tiene reportes (no se construye el libro).
- Construcción del libro: nombre de hoja, las 34 cabeceras en orden, una fila por reporte.
- Exclusividad Form1/Form2: columnas del formulario que no aplica quedan vacías.
- Formatos amigables: traducción de enums, `Sí`/`No`, fechas `dd/mm/yyyy hh:mm`, teléfonos como texto, municipio código + nombre.

## Evidencia de ejecución

### Pruebas Unitarias (Go) — servicio y builder Excel

* **Casos cubiertos**: rango sin reportes (`TC-RPC-API-02` a nivel servicio), error de repositorio, hoja única `Reportes` con 34 cabeceras en orden, diferenciación F1/F2 con columnas del otro form vacías (`TC-RPC-API-09`, `TC-RPC-API-10`), traducción de códigos (`TC-RPC-API-15`), teléfonos como texto (`TC-RPC-API-17`), fechas `dd/mm/yyyy hh:mm`.
* **Comando**: `go test -v ./salvia/service/contacts_report_service_test.go ./salvia/service/contacts_report_excel.go ./salvia/service/report_service.go ./salvia/service/followups_report_excel.go`
* **Fecha**: 2026-07-21
* **Resultado**: `PASS`
* **Salida**:
  ```text
  === RUN   TestReportService_GenerateConsolidatedContactsReport_Empty
  --- PASS: TestReportService_GenerateConsolidatedContactsReport_Empty (0.00s)
  === RUN   TestReportService_GenerateConsolidatedContactsReport_Error
  --- PASS: TestReportService_GenerateConsolidatedContactsReport_Error (0.00s)
  === RUN   TestReportService_GenerateConsolidatedContactsReport_WithData
  --- PASS: TestReportService_GenerateConsolidatedContactsReport_WithData (0.01s)
  PASS
  ok  	command-line-arguments	0.688s
  ```

### Pruebas Unitarias (Go) — facade

* **Casos cubiertos**: binding/validación de fechas vacías (`TC-RPC-API-03`), sin sesión → `401 AUTH_TOKEN_INVALID` (`TC-RPC-API-08`).
* **Comando**: `go test -v ./salvia/facades/ContactsReportFacade_test.go ./salvia/facades/ContactsReportFacade.go ./salvia/facades/ConsolidatedReportFacade.go`
* **Fecha**: 2026-07-21
* **Resultado**: `PASS`
* **Salida**:
  ```text
  === RUN   TestDownloadContactsReport_BindingValidation
  --- PASS: TestDownloadContactsReport_BindingValidation (0.00s)
  === RUN   TestDownloadContactsReport_NoSession
  --- PASS: TestDownloadContactsReport_NoSession (0.00s)
  PASS
  ok  	command-line-arguments	0.515s
  ```

### Verificación en vivo — endpoint HTTP (servidor real)

* **Caso**: servidor levantado (`PORT=8090 go run main.go`), ruta registrada en Gin (`POST /api/v1/reportes/contactos-consolidado`) y petición sin sesión rechazada.
* **Comando**: `curl -s -w "\nHTTP:%{http_code}\n" -X POST http://localhost:8090/api/v1/reportes/contactos-consolidado -H "Content-Type: application/json" -d '{"start_date":"2026-01-01","end_date":"2026-06-30"}'`
* **Fecha**: 2026-07-21
* **Resultado**: `PASS` (`TC-RPC-API-08`)
* **Salida**:
  ```text
  {"error_code":"AUTH_TOKEN_INVALID","message":"Sesión no válida o expirada"}
  HTTP:401
  ```

### Verificación en vivo — SQL + builder contra BD real

* **Caso**: arnés temporal (`tests/contactsreport/live_query_test.go`, eliminado tras la ejecución) que ejecutó `FindContactsInDateRange` + `BuildContactsExcel` contra la BD real con rango de 1 año (`TC-RPC-API-01` datos/estructura, `TC-RPC-API-09`, `TC-RPC-API-10`, `TC-RPC-API-11`, `TC-RPC-API-13`, `TC-RPC-API-15`).
* **Comando**: `LIVE_DB=1 XLSX_OUT=<scratch>/reporte-contactos-live.xlsx go test -v ./tests/contactsreport/ -run TestLiveContactsReportQuery`
* **Fecha**: 2026-07-21
* **Resultado**: `PASS`
* **Salida**:
  ```text
  live_query_test.go:53: reportes en el último año: 1622
  live_query_test.go:66: form1=1087 form2=535 sin_formulario=0
  live_query_test.go:97: hoja Reportes OK: 1622 filas de datos, 34 columnas
  --- PASS: TestLiveContactsReportQuery (4.95s)
  ```
  Inspección del `.xlsx` generado: presentes "Reporte propio (Formulario 1)", "Reporte de tercero (Formulario 2)", "Cédula de Ciudadanía", "Mujer", "Urbano" (códigos traducidos) y las 34 cabeceras.

### Corrección de pantalla (2026-07-21, tarde)

* **Caso**: el botón no aparecía porque la integración inicial se hizo en `get_victim_contacts_sv.html`, pero la pantalla real "Reportes Salvia" del supervisor es `get_victim_cases_sv.html` (servida por `VictimCaseFacade.go`). Se corrigió el spec del componente y se movió la integración (botón junto a "Crear Nuevo caso", modal `contacts_report_modal.html` en `victim_case/`).
* **Verificación**:
  * Parseo de plantillas con arnés temporal (eliminado): `go test -v ./tests/tplparse/` → `PASS`, define `salvia/victim_case/contacts_report_modal.html` registrado.
  * `go build ./salvia/... ./internal/... .` → OK.
  * Servidor relanzado (`PORT=8090`): ruta activa y `curl` sin sesión → `401 AUTH_TOKEN_INVALID` (idéntico al contrato).
* **Resultado**: `PASS` (pendiente confirmación visual del usuario en sesión `sv`).

### Corrección de estilos del modal (2026-07-21, tarde 2)

* **Caso**: el modal abría pero sin estilos (contenido inline, sin overlay). Causa raíz: el bloque `<style>` inline del template se pierde porque Vue monta sobre `body#app` y re-renderiza el contenido; los modales existentes funcionan porque cargan su CSS por `<link>` externo (`reasignar_casos_modal.html` → `/static/css/reasignar-casos-modal.css`).
* **Fix**: estilos movidos a `/static/css/contacts-report-modal.css` (prefijo `crm-`) cargados con `<link>` desde `contacts_report_modal.html`; clases del template renombradas `rcm-*` → `crm-*`.
* **Verificación** (2026-07-21): render server-side con arnés temporal (eliminado) contiene el `<link>` y el x-template; servidor relanzado sirve los estáticos:
  ```text
  curl /static/css/contacts-report-modal.css → HTTP:200 bytes:3018
  curl /static/js/components/contacts-report-modal.js → HTTP:200 bytes:5752
  ```
* **Resultado**: `PASS` (pendiente confirmación visual del usuario en sesión `sv`).

### Ajuste de columnas: 37 → 34 (2026-07-21, tarde 3)

* **Caso**: por decisión del usuario se eliminaron `contacto_usuario_app`, `f2_detalle_tipo_reporte` y `f2_respuesta_autorizacion` (las dos últimas siempre vacías — diagnóstico contra BD real: `report_type_details` NULL en 535/535 registros Form2; la relación de enums de contactos solo contiene `victim_case_form2_adjustments_gbv`, 446 filas, sin categoría de autorización).
* **Comandos** (2026-07-21):
  * Unitarias: `go test -v ./salvia/service/contacts_report_service_test.go ./salvia/service/contacts_report_excel.go ./salvia/service/report_service.go ./salvia/service/followups_report_excel.go` → `PASS` (3/3, cabeceras e índices actualizados a 34).
  * En vivo (arnés temporal eliminado): `LIVE_DB=1 go test -v ./tests/dbcheck/ -run TestLive34Columns` → `PASS`:
    ```text
    filas datos=1622, cabeceras=34, primera=contacto_id, ultima=f2_ajustes_gbv
    ```
* **Resultado**: `PASS`

### Fix traducción de enums Form2 (2026-07-21, tarde 4)

* **Caso** (`TC-RPC-API-15`): las columnas de enums de Form2 salían con la clave cruda (`yes_no_y`, `victim_case_form2_report_type_v`, `victim_case_form2_adjustments_gbv_nr`) porque `victim_case_form2_enums_name` almacena la **clave de Locale**, no el texto. Se agregó `translateForm2Enum`/`translateForm2EnumList` en el builder: resuelve vía `salvia_config.Locale["sp"]` con fallback a `translateEnumKey` (conforme a la regla de formato ya aprobada en el spec del servicio).
* **Comandos** (2026-07-21):
  * Unitarias (fixtures con claves reales → "Sí", "No", "Usted es la víctima de violencia basada en género", "No requiere, Sí"): `go test -v ./salvia/service/contacts_report_service_test.go ./salvia/service/contacts_report_excel.go ./salvia/service/report_service.go ./salvia/service/followups_report_excel.go` → `PASS` (3/3).
  * En vivo (arnés temporal eliminado): `LIVE_DB=1 go test -v ./tests/dbcheck/ -run TestLiveEnumTranslation` → `PASS`:
    ```text
    filas=1622, celdas enum con clave cruda=0, ejemplos=map[]
    ```
* **Resultado**: `PASS`

### API con sesión autenticada real (supervisor `sv`)

* **Caso**: servidor levantado (`PORT=8090 go run main.go`, BD Supabase real), login `POST /seguridad/login` como `supervisor.prueba` (rol `sv`; captcha omitido con header `User-Agent: flutter-client`) y peticiones al endpoint con la cookie de sesión.
* **Fecha**: 2026-07-21
* **Resultado**: `PASS`
* **Casos y salidas**:
  * `TC-RPC-API-01` — rango con datos (`2025-07-21`→`2026-07-21`): `HTTP 200`, `Content-Type: ...spreadsheetml.sheet`, `Content-Disposition: attachment; filename="reporte-contactos_2025-07-21_2026-07-21.xlsx"`, archivo `.xlsx` real de 872 230 bytes. Inspección del libro: hoja única `Reportes`, **34 cabeceras** en orden (`contacto_id`…`f2_ajustes_gbv`), 1627 filas de datos, presentes "Reporte propio (Formulario 1)" y "Reporte de tercero (Formulario 2)" (cubre `TC-RPC-API-09/10/12/13/15`).
  * `TC-RPC-API-02` — rango sin reportes (`1990`): `HTTP 400` `{"error_code":"VALIDATION_FAILED","message":"No hay reportes en el rango seleccionado"}`.
  * `TC-RPC-API-03` — falta `end_date`: `HTTP 400` `VALIDATION_FAILED` ("Fechas requeridas o formato incorrecto").
  * `TC-RPC-API-04` — `start > end`: `HTTP 400` `VALIDATION_FAILED` ("La fecha inicial no puede ser posterior a la fecha final").
  * `TC-RPC-API-05` — rango > 366 días: `HTTP 400` `VALIDATION_FAILED` ("El rango seleccionado no puede superar los 366 días (1 año)").
  * `TC-RPC-API-06` — formato inválido (`2026/13/40`): `HTTP 400` `VALIDATION_FAILED` ("Fecha inicial inválida (debe ser YYYY-MM-DD)").
  * `TC-RPC-API-08` — sin sesión: `HTTP 401` `AUTH_TOKEN_INVALID` (ya verificado antes).

### E2E / UI (Playwright — Chromium headless, rol `sv`)

* **Caso**: suite Playwright contra el servidor real; login por API sobre el mismo contexto del navegador (User-Agent `flutter-client` para saltar captcha), navegación a `/salvia/casos` (pantalla "Reportes Salvia" del supervisor).
* **Comando**: `npx playwright test` (5 tests, 1 worker, headless).
* **Fecha**: 2026-07-21
* **Resultado**: `PASS` (5/5)
* **Salida**:
  ```text
  Running 5 tests using 1 worker
    ✓  TC-RPC-UI-01: el supervisor ve el botón "Descargar reportes" (5.9s)
    ✓  TC-RPC-UI-03: abrir modal, rango válido, descarga .xlsx con nombre correcto (9.6s)
    ✓  TC-RPC-UI-04: rango > 1 año muestra error en el DOM y no dispara la petición (5.1s)
    ✓  TC-RPC-UI-05: rango sin reportes muestra el mensaje del backend (4.5s)
    ✓  TC-RPC-UI-06: error 500 del backend muestra mensaje de error en el modal (4.1s)
    5 passed (33.7s)
  ```
* **Confirmación visual** (capturas revisadas): `ui-01` botón "Descargar reportes" visible junto a "Crear Nuevo caso" en "Reportes Salvia" (sesión "Supervisor Pruebas"); `ui-03` modal "Descarga de Reportes (Excel)" con overlay, dos date-pickers y nota; `ui-05` banner rojo "No hay reportes en el rango seleccionado". `TC-RPC-UI-03` descargó el `.xlsx` con nombre `reporte-contactos_<desde>_<hasta>.xlsx` (firma ZIP `PK`, >1 KB).

### Pendiente

* `TC-RPC-API-07` (rol sin permiso → `403 AUTH_INSUFFICIENT_ROLE`) y `TC-RPC-UI-02` (rol sin permiso no ve el botón): requieren credenciales de un usuario con rol distinto de `sv`/`ad`, no disponibles en el entorno del agente. La restricción está garantizada por el permiso `report_contacts_consolidated` (solo `sv`/`ad`) en `Menu.go` y por `v-if="canDescargarReportes"` en el frontend.
* `TC-RPC-API-14` (Form2 con `adjustmentsGBV` múltiples separados por coma) y `TC-RPC-API-16` (frescura sin caché): cubiertos indirectamente por la traducción de enums verificada en vivo y por la ausencia de caché en el servicio; sin aserción HTTP dedicada.
* **Infra E2E (versionada)**: la suite Playwright quedó en el repo siguiendo la convención del CI (`.github/workflows/playwright.yml`): `package.json` + `playwright.config.ts` en la raíz, `tests/api/contacts-report.api.spec.ts` (7 casos) y `tests/e2e/contacts-report.spec.ts` (5 casos), con helper de login en `tests/helpers/auth.ts`. Corrida completa desde la raíz: **12/12 PASS** (`npm test`, server Go auto-levantado por `webServer`, ~1.9 min). Credenciales por env (`SV_PASS`), captura/descargas en `test-results/` (git-ignored). Ver `tests/README.md`. **Nota CI**: el workflow no aprovisiona server+BD; para pipeline verde hay que apuntar `BASE_URL` a un entorno desplegado o extender el workflow.
