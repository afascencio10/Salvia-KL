# Changelog — Detalle del Caso — septiembre 2026

## 2026-09-09 — Tarjeta Mecanismos en Gestión institucional

Nueva tarjeta para registrar la cadena de escalamiento del caso ante mecanismos articuladores. El modal asigna un usuario con rol `mc` filtrado por alcance y por el municipio del perfil (`town → city → department`). Cada registro admite comentarios (autor + fecha) y la creación escribe un evento en el Timeline.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/frontend/js/components/case-mecanismos.js` | Card + modal |
| `src/frontend/html/salvia/case_detail/get_case_detail_sv.html` | Montaje en el tab Gestión institucional |
| `src/internal/models/case_articulating_mechanism.go` | Tablas `case_articulating_mechanism` y comentarios |
| `src/salvia/controller/case_articulating_mechanism_controller.go` | REST |
| `src/main.go` | Seed rol `mc` + AutoMigrate + wiring |

## 2026-09-02 — Enlace Territorial puede abrir el detalle

El rol `en` obtiene `get_case_detail_sv` para que «Ver caso» (Mis barreras / Detalle de barrera) no responda 401. Las acciones de reasignar, seguimientos y similares siguen ocultas: el HTML las limita a `sv`/`op`/`ro`.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/salvia/config/Menu.go` | `"en": true` en `get_case_detail_sv` |
