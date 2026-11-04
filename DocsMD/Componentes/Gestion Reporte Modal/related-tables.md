# `gestion-reporte-modal` — Tablas de BD

Este componente es el que realmente **escribe** en las tablas descritas en el documento compartido de `Listado Reportes` — ahí está el detalle completo (columnas nuevas en `victim_contact`, tabla nueva `victim_contact_attempt`, motivos de cierre, GAPS):

📄 [`../../Screens/Listado Reportes/related-tables.md`](../../Screens/Listado%20Reportes/related-tables.md)

## Mapa rápido evento → escritura

| Evento | Tabla(s) escritas |
|---|---|
| M-02, M-07, M-08, M-09 | `victim_contact` (`close_reason`, `close_observation`, `closed_date`) + `case_timeline_event` («Reporte cerrado como inválido») en la misma transacción *(WP-1 T6)* |
| M-05, M-06 | `victim_contact_attempt` (INSERT) + `case_timeline_event` («Intento de contacto fallido» / «Contacto logrado») en la misma transacción *(WP-1 T6)* — el de M-06 (exitoso) es lo que permite a M-01 reabrir en `postContacto` en vez de `esquema3x3` si el agente no terminó de cerrar |
| M-10 *(WP-1 T7)* | `victim_case` (+ `victim_case_form2`, `general_user`…, vía `CreateDraft`), vínculo `victim_case.victim_case_victim_contact`, `case_timeline_event` («Borrador de caso iniciado»). **El reporte no se escribe:** queda «Borrador iniciado» por estado derivado y se lee «Cerrado — válido» cuando el caso pasa a `te`/`cd`. El backend sigue rechazando `close_reason='ap'` con 400 (`deferredCloseReasonCodes`) |
| M-01, M-04, M-11 | Ninguna — solo lectura o navegación local |

Otros hitos (WP-1 T6): la reasignación masiva escribe «Reporte asignado» / «Reporte reasignado» (ver `Screens/Listado Reportes/`) y el reporte V2 recibido escribe «Reporte recibido» (ver `Screens/Registro de Contacto Víctima V2/`). El estado derivado del reporte está en `Screens/Listado Reportes/related-tables.md`.
