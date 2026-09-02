# Changelog — Detalle del Caso — septiembre 2026

## 2026-09-02 — Enlace Territorial puede abrir el detalle

El rol `en` obtiene `get_case_detail_sv` para que «Ver caso» (Mis barreras / Detalle de barrera) no responda 401. Las acciones de reasignar, seguimientos y similares siguen ocultas: el HTML las limita a `sv`/`op`/`ro`.

**Archivos**

| Archivo | Cambio |
|---|---|
| `src/salvia/config/Menu.go` | `"en": true` en `get_case_detail_sv` |
