## 2026-09-24 — Banner con el motivo de solo lectura

Mismo cambio que en `hacer-seguimiento` (ver `DocsMD/Screens/hacer-seguimiento/changelogSep2026.md`): cuando `canEdit = false`, se muestra un banner 🔒 al inicio de la pantalla con el motivo exacto (`noEditReason`), calculado en la misma cascada de prioridad que ya usa esta pantalla. Estilo compartido vía `.ms-readonly-banner` en el CSS global, para que se vea idéntico en ambas pantallas.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | `noEditReason` en cada rama de la cascada de `canEdit`; banner `ReadonlyBanner` al inicio de `#follow-up-detail-app` |
| `src/frontend/css/styles.css` | Nuevas clases `.ms-readonly-banner*` (compartidas con hacer-seguimiento) |
| `DocsMD/Screens/Detalle de Seguimiento/detalle-seguimiento-interface.md` | Tabla de `canEdit` expandida a 5 condiciones + `noEditReason`; nodo del banner en el árbol |
