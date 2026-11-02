## 2026-11-02 — REP-12: la fecha del reporte muestra el día de Bogotá

«Mis Reportes» comparte el endpoint y la columna con el Listado de Reportes: la columna «Fecha Reporte» salía en «—» porque `ReportListItem` no traía `creationDateCO`. Ahora el listado devuelve el día de creación en hora de Bogotá (`AAAA-MM-DD`, sin `AT TIME ZONE`) y la celda lo muestra sin `new Date()`, así que el día no cambia por la zona horaria del navegador. Cierra H11 del plan de QA. El detalle del cambio y los archivos están en el changelog de Listado Reportes.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/report_list_item.go` | `ReportListItem` agrega `CreationDateCO` (`creationDateCO`) — ver changelog de Listado Reportes |
| `src/internal/repository/victim_contacts_list_repository.go` | Selecciona `TO_CHAR(victim_contact_creation_date, 'YYYY-MM-DD')` como `creation_date_co` |
| `cloudDevelopment/e2e/wp2-rep12-fecha-reporte.js` | Nuevo — incluye la prueba del scope de agente y de la pantalla en otra zona horaria |
