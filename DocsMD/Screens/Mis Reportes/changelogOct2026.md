## 2026-10-13 — Implementación y QA de Mis Reportes

Pantalla nueva `/salvia/mis-reportes` (solo op): los reportes asignados al agente de la sesión, con los mismos filtros, búsqueda, orden y modal de gestión que Listado Reportes, sin reasignación. El scope sale de la sesión con la prop `agentId` de `casos-component`, así que ningún filtro lo ensancha. En la primera prueba respondía 400 porque `MyReportsGET` buscaba el template en `my_reportes/` (corregido a `reportes/`). Entrada en el menú lateral solo para op.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/reportes/get_my_reportes.html` | Nuevo — pantalla: casos-component con `:agent-id` y gestion-reporte-modal |
| `src/salvia/facades/VictimContactFacade.go` | Nuevo `MyReportsGET`; carpeta del template corregida a `reportes/` |
| `src/salvia/facades/MainRouter.go` | Solo adiciones — ruta `/salvia/mis-reportes` |
| `src/salvia/config/Menu.go` | Solo adiciones — permiso `get_my_reports` (solo op) |
| `src/frontend/html/sidebar/sidebar.html` | Entrada "Mis Reportes" (solo op) |
