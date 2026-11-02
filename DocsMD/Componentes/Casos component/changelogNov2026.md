## 2026-11-02 — REP-12: la columna «Fecha Reporte» de los reportes deja de salir en «—»

Sin cambios en el componente: la celda `registration_date` ya leía `creationDateCO` (`formatFechaCO`, sin `new Date()`). El cambio es de datos: `ReportListItem` agrega `creationDateCO` (día de creación en hora de Bogotá), el mismo nombre que ya traía `CaseListItem`, así que Listado Reportes y Mis Reportes dejan de mostrar «—». Se actualiza la documentación de la columna.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/internal/models/report_list_item.go · src/internal/repository/victim_contacts_list_repository.go` | El listado de reportes devuelve `creationDateCO` — ver changelog de Listado Reportes |
