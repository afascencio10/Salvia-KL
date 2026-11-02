# QA — Estado derivado «Borrador iniciado» / «Cerrado — válido» y hitos de asignación (WP-1 T6 · T7)

Ejecutado 2026-11-02 con `requests`, `psql` y Playwright contra la app del sandbox. Pruebas Go de la regla: `go test -vet=off ./internal/repository -run 'ReportStatusFrom|ValidateConversion'`.

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C1 | Reporte con borrador (caso en `bo`/`in`) | `GET /api/v1/victim-contacts/list` lo devuelve con `reportStatus = borrador_iniciado`; el filtro `filter_estado_reporte=borrador_iniciado` lo incluye | API |
| C2 | `asignado` y `sin_asignar` | ya no incluyen el reporte con borrador | API |
| C3 | Caso vinculado pasa a `te` (por SQL) | el listado lo muestra `cerrado_valido` y el filtro `cerrado_valido` lo incluye; el detalle devuelve cierre `ap` «Válido - Apertura de Caso Salvia» | API |
| C4 | Precedencia | motivo de cierre y legado invalidado ganan sobre el caso vinculado; un caso de otro flujo (`ra`) no cambia el estado | GO |
| C5 | Listado en el navegador | etiqueta «Borrador iniciado» (badge ámbar) y «Cerrado — válido» tras completar el caso | UI |
| C6 | Mis Reportes | carga con la opción de filtro nueva y `case-timeline.js` sin errores JS ni 4xx/5xx de API | UI |
| C7 | Reasignación (E-06) | «Reporte asignado» y «Reporte reasignado»; masiva de 3 (1 inexistente) = 2 eventos; si el evento falla se revierte la asignación | API |
