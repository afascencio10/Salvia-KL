# QA — TL-E01 (Cuando carga el componente `case-timeline`) con `report-id` (WP-1 T5)

Ejecutado 2026-11-02 contra la app del sandbox: API con `requests` y `psql`, UI con Playwright.

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C1 | Esquema tras arrancar | existen `victim_contact_id` e índice; `case_id` admite NULL (AutoMigrate emitió el `DROP NOT NULL`) | API |
| C2 | Evento con reporte y sin caso (insert por SQL) | `GET /api/v1/reportes/{icode}/timeline-events` lo devuelve con el actor resuelto | API |
| C3 | Orden y vacío | más reciente primero; un reporte sin eventos (o un icode inexistente) devuelve `[]` | API |
| C4 | Sin sesión o con rol distinto de op y sv (et, us, ro) | 401 | API |
| C5 | Timeline de un caso | el endpoint de casos no trae los eventos de reportes, salvo el «Borrador de caso iniciado», que lleva `case_id`; el del reporte lo incluye también | API |
| C6 | Componente con `report-id` | pinta los eventos (verificado en el modal de gestión); sin ningún id muestra «Falta el identificador del caso o del reporte» | UI |
| C7 | Regresión ligera | el Detalle del Caso sigue mostrando su timeline | UI |
