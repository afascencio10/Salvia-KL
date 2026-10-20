## 2026-10-13 — Extensión aditiva para Gestión de reportes

`casos-component` suma los props `listEndpoint`, `sortOptions` e `itemLabelPlural`, las columnas de reporte y los filtros `estado_reporte` y `asignacion` (con «Sin asignar»); `reasignar-casos-modal` suma `entityKind`. Las pantallas de casos no pasan esos props y se comportan igual. Se corrigió una regresión del cambio: el contador de selección decía «casos(s)» (ahora usa `itemLabelSingular`). Limitación conocida: con `REASSIGN_CROSS_TEAM_CONTINGENCY = false` los reportes no serían seleccionables.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/casos-component.js` | Props nuevos, filtro `asignacion` por `by-role`, `itemLabelSingular` y respuesta `reports` |
| `src/frontend/html/salvia/case_component/case_component.html` | Columnas de reporte, `sortOptions`, textos con `itemLabelPlural`; `victim_info` muestra el ID del reporte |
| `src/frontend/css/casos-component.css` | Solo adiciones — estilos de los badges de tipo y estado del reporte |
| `src/frontend/js/components/reasignar-casos-modal.js` | Prop `entityKind`: candidatos `op`, endpoint de reportes y nombre del agente asignado |
| `src/frontend/html/salvia/list-cases/reasignar_casos_modal.html` | Textos condicionales por `entityKind` (sin `Doc:` ni riesgo para reportes) |
