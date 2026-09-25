## 2026-09-24 — Banner con el motivo de solo lectura

Cuando `canEdit = false`, ambas pantallas de seguimiento (`hacer-seguimiento` y `Detalle de Seguimiento`) ahora muestran un banner 🔒 al inicio de la pantalla explicando la razón exacta: formulario migrado, caso cerrado, más de 5 días desde completado, o rol de solo lectura (ps/ts). El cálculo de `noEditReason` se agregó en paralelo a `canEdit`, en la misma cascada de prioridad, para que nunca queden desincronizados. El estilo (`.ms-readonly-banner`) se agregó al CSS global (`styles.css`) para que se vea idéntico en ambas pantallas.

### Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/frontend/html/salvia/follow_up_v2/hacer_seguimiento.html` | `noEditReason` en cada rama de la cascada de `canEdit`; banner `ReadonlyBanner` al inicio de `#seguimiento-app` |
| `src/frontend/html/salvia/follow_up_detail/get_follow_up_detail.html` | Mismo cambio (ver `DocsMD/Screens/Detalle de Seguimiento/changelogSep2026.md`) |
| `src/frontend/css/styles.css` | Nuevas clases `.ms-readonly-banner*` |
| `DocsMD/Screens/hacer-seguimiento/interfaz.md` | Tabla de `canEdit` expandida a 5 condiciones + `noEditReason`; nodo del banner en el árbol |
| `DocsMD/Screens/hacer-seguimiento/Flujos/flow-E01-cuando-carga-pantalla.md` | PASO 4 incluye las 5 condiciones (antes faltaban form migrado y override de rol); PASO 5 documenta el banner |
