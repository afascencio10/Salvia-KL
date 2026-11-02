# QA — CR-E03 (Cuando presiona «Consultar») · API

Ejecutado 2026-11-02 con `requests` y `psql` contra la app del sandbox (Postgres local). Cada caso usa una IP distinta (`X-Forwarded-For`) para no interferir con el limitador. Pruebas Go del servicio: `go test -vet=off -overlay <omitir followup_v2_service_test.go> ./salvia/service -run 'Radicado|ReportStatus|Lookup|FailureLimiter'` (el paquete tiene un test previo que no compila).

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C1 | Radicado de un reporte recién creado | 200, estado «recibido», fecha de hoy | API |
| C2 | Reporte con agente asignado y sin intentos | «en gestión» | API |
| C3 | Reporte con un intento | «en gestión» | API |
| C4 | Reporte cerrado (ni, nc, co, dp, ap) y reporte legado invalidado | «cerrado»; el cuerpo no trae el motivo | API |
| C5 | Contenido de la respuesta | solo radicado, estado, texto, descripción y fecha; sin icode, nombres, teléfono, tipo ni agente | API |
| C6 | Radicado en minúscula, sin guion, con espacios o con guion largo | se normaliza y se encuentra | API |
| C7 | Formato inválido, radicado inexistente y vacío | mismo 404 con cuerpo idéntico | API |
| C8 | Cinco fallos seguidos desde una IP | el sexto (aun válido) responde 429 con `reintentarEnSegundos` y `Retry-After`; otra IP no se ve afectada; pasada la ventana vuelve a funcionar (Go) | API + GO |
| C9 | Consultas válidas repetidas | no cuentan como fallo ni bloquean | API + GO |
| C10 | Sin cookie de sesión | la página y la consulta funcionan | API |
| C11 | Entradas hostiles (comilla, `%`, `_`, `DROP TABLE`, NUL, 5000 caracteres, cuerpo malformado, número en vez de texto) | 404 o 400, nunca 500; la tabla sigue existiendo | API |
| C12 | Reporte sin radicado (legado) | no es consultable (404) | API |
| C16 | Reporte con borrador de caso iniciado (caso vinculado en `bo` o `in`) | «en gestión» | API |
| C17 | Caso vinculado en `te` o `cd` (por SQL en el sandbox) | «cerrado»; uno en `ra` (otro flujo) no cambia el estado | API |
