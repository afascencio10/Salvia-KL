# QA — Historial del reporte, hitos en el timeline, encabezado y «Hechos del reporte» (WP-1 T6, REP-06 parte b)

Ejecutado 2026-11-02 con `requests`, `psql` y Playwright contra la app del sandbox (Postgres local). Los fallos de evento se simulan con un trigger temporal sobre `case_timeline_event` que lanza una excepción para los tipos indicados en una tabla bandera (se elimina al terminar).

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C1 | Crear un reporte V2 | un evento «Reporte recibido» con tipo y radicado en la descripción; `category = Reportes`, sin `case_id` | API |
| C2 | El evento «recibido» falla | el reporte se crea igual (200, best-effort) y no queda evento | API |
| C3 | Reasignar un reporte sin agente y luego con agente | «Reporte asignado» y «Reporte reasignado» con agente, equipo (si tiene) y supervisor; el actor es el supervisor | API |
| C4 | Reasignación masiva de 3 reportes (1 inexistente) | `{updated: 2, skipped: 1}`; 2 eventos; el omitido no genera | API |
| C5 | Intento fallido | evento «Intento de contacto fallido» con el agente y «(N hoy)» | API |
| C6 | Intento exitoso (dos veces) | un evento «Contacto logrado» por cada uno (mismo tipo); no existe un tipo aparte para un segundo exitoso | API |
| C7 | Cierre ni, co, dp y nc (3 días completos) | «Reporte cerrado como inválido» con motivo, observación y agente | API |
| C8 | Cierre `ap` por `/cerrar` | sigue 400 y no escribe evento | API |
| C9 | Atomicidad (evento que falla) | el intento (fallido o exitoso), el cierre y la asignación se revierten (500 / 400 con el error); no queda el cambio sin su evento | API |
| C10 | Reporte ya cerrado | 409 en intento y cierre, sin evento nuevo | API |
| C11 | Modal: historial | al abrir muestra los eventos; tras registrar un intento (fallido o exitoso) aparece el nuevo sin recargar | UI |
| C12 | Modal con reporte cerrado | muestra el historial completo (incluye «Reporte cerrado como inválido») y se mantiene al reabrir | UI |
| C13 | Encabezado del modal | muestra nombre, tipo, teléfono y radicado, sin el `icode`; un reporte sin radicado no rompe (sin etiqueta de radicado) | UI |
| C14 | Sección «Hechos del reporte» | se ve en todos los pasos (validación previa, 3x3, post-contacto y cerrado) con quién reporta, teléfono, ajuste razonable y descripción | UI |
| C15 | Detalle extendido | `GET /api/v1/victim-contacts/:icode` trae `radicado` y `reportTypeCode` (solo V2) | API |
| C16 | Concurrencia | 6 intentos simultáneos → 6 intentos y 6 eventos; un cierre y un intento simultáneos no dejan intentos sin evento | API |
| C17 | Páginas anfitrionas | «Mis Reportes» y «Reportes» cargan `case-timeline.js` sin errores JS ni 4xx/5xx de API | UI |
