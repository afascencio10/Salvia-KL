## 2026-11-02 — REP-12: la fecha del reporte muestra el día de Bogotá

La columna «Fecha Reporte» del Listado de Reportes salía en «—» porque el listado no traía `creationDateCO`. Ahora `GET /api/v1/victim-contacts/list` devuelve el día de creación en hora de Bogotá (`AAAA-MM-DD`, sobre el `timestamp` ya guardado, sin `AT TIME ZONE`) y la celda lo muestra sin `new Date()`: el día no cambia por la zona horaria del navegador. Probado en 5 zonas, en los bordes de día y de año y con un reporte real del asistente público. Cierra H11 del plan de QA.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/report_list_item.go` | `ReportListItem` agrega `CreationDateCO` (`creationDateCO`); `creationDate` se conserva |
| `src/internal/repository/victim_contacts_list_repository.go` | Selecciona `TO_CHAR(victim_contact_creation_date, 'YYYY-MM-DD') AS creation_date_co` y lo mapea a la fila |
| `cloudDevelopment/e2e/wp2-lib.js · wp2-rep12-fecha-reporte.js` | Nuevos — prueba backend + BD + pantalla en 5 zonas horarias (19 comprobaciones) |
