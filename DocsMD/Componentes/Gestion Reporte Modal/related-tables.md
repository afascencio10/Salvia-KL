# `gestion-reporte-modal` — Tablas de BD

Este componente es el que realmente **escribe** en las tablas descritas en el documento compartido de `Listado Reportes` — ahí está el detalle completo (columnas nuevas en `victim_contact`, tabla nueva `victim_contact_attempt`, motivos de cierre, GAPS):

📄 [`../../Screens/Listado Reportes/related-tables.md`](../../Screens/Listado%20Reportes/related-tables.md)

## Mapa rápido evento → escritura

| Evento | Tabla(s) escritas |
|---|---|
| M-02, M-07, M-08, M-09 | `victim_contact` (`close_reason`, `close_observation`, `closed_date`) |
| M-05, M-06 | `victim_contact_attempt` (INSERT) — el de M-06 (exitoso) es lo que permite a M-01 reabrir en `postContacto` en vez de `esquema3x3` si el agente no terminó de cerrar |
| ~~M-10~~ | ⏸️ Diferido, no se implementa en esta fase — ver `gestion-reporte-modal-interface.md`. El backend rechaza `close_reason='ap'` con 400 (`deferredCloseReasonCodes`) |
| M-01, M-04, M-11 | Ninguna — solo lectura o navegación local |
