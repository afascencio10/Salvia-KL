## 2026-11-02 — Pantalla nueva: Consulta Reporte Pública (WP-1 T4, REP-02)

Nueva pantalla pública `/salvia/public-consulta` donde la persona consulta el estado de su reporte con el radicado (recibido, en gestión o cerrado), sin datos personales y con límite de 5 fallos por IP cada 15 minutos. Probada por API (17 casos) y Playwright; ver `QA/`.

### Archivos

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/victim_contact/consulta_reporte.html` | Nuevo — app Vue 3 `appConsulta` |
| `src/salvia/facades/ReportStatusFacade.go` | Nuevo — `ReportStatusGET_Public` |
| `src/salvia/controller/report_status_controller.go` | Nuevo — `POST /api/v1/public/consulta-reporte` |
| `src/salvia/service/report_status_service.go` | Nuevo — normalización, estados y limitador por IP (+ pruebas Go) |
| `src/internal/repository/report_status_repository.go` | Nuevo — `FindByRadicado` |
| `src/salvia/facades/MainRouter.go`, `src/salvia/config/Enums.go`, `src/salvia/config/Locale.go`, `src/main.go` | Solo adiciones — ruta, plantilla, título y cableado |
