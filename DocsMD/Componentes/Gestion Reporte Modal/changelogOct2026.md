## 2026-10-13 — Implementación y QA de gestion-reporte-modal

Componente nuevo: validación previa, esquema de recontacto 3x3 (intentos fallidos y exitosos), cierre por no información, no contacto, no competencia o decisión de la persona, y aviso `updated` al padre. Endpoints nuevos de detalle, intento y cierre. La QA con Playwright corrigió, entre otros: días de Bogotá mal contados, horas guardadas en la zona del servidor, observaciones de más de 500 caracteres, legados gestionables y log con la zona del navegador. `ap` (M-10) queda diferido y el backend lo rechaza con 400.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/gestion-reporte-modal.js` | Nuevo — máquina de estados del modal (4 pasos) |
| `src/frontend/html/salvia/reportes/gestion_reporte_modal.html` | Nuevo — template del modal |
| `src/frontend/css/gestion-reporte-modal.css` | Nuevo — estilos `.grm-*` |
| `src/salvia/controller/victim_contacts_management_controller.go · service/…_management_service.go · internal/repository/…_management_repository.go` | Nuevo — `GET /api/v1/victim-contacts/:icode`, `POST …/intento` y `POST …/cerrar` |
| `src/salvia/service/victim_contacts_labels.go` | Nuevo — etiquetas legibles para el detalle |
| `src/internal/models/report_list_item.go` | Nuevo — `ReportDetail` y `ReportAttemptItem` |
| `src/migrations/add_victim_contact_management.sql` | Nuevo — `victim_contact_attempt` y motivos de cierre (no está en git) |
