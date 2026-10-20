# Mis Reportes — Tablas de BD

Esta pantalla no introduce ninguna tabla o columna adicional. Usa exactamente las mismas que `Listado Reportes` — ver el documento completo (columnas nuevas en `victim_contact`, tabla nueva `victim_contact_attempt`, motivos de cierre, GAPS) en:

📄 [`../Listado Reportes/related-tables.md`](../Listado%20Reportes/related-tables.md)

La única diferencia funcional es que el backend de `GET /api/v1/victim-contacts/list` recibe siempre `filter_key=persona_asignada&filter_value={currentUserICode}` desde esta pantalla, sin que el usuario pueda cambiarlo (ver GAP-5 en [`mis-reportes-interface.md`](mis-reportes-interface.md)).
