# QA — E-11 (Cuando presiona "Volver al Inicio") · captcha nuevo (WP-1 T2, REP-14*)

Ejecutado 2026-11-02 con Playwright (UI) y `requests` (API) contra la app del sandbox.

El captcha no se puede resolver automáticamente: el éxito se prueba con el bypass del cliente móvil (`User-Agent: flutter-client` solo en el POST) y el consumo del id se demuestra aparte.

| Caso | Descripción | Resultado esperado | Método |
|---|---|---|---|
| C1 | Enviar un reporte y presionar «Volver al Inicio» | se hace 1 GET con `Accept: application/json` a `/salvia/public-v2/nuevo` y el wizard vuelve al paso inicial | UI |
| C2 | Segundo formulario en la misma página | el `captchaID` es distinto al del primero y la imagen del nuevo se sirve (200 image/png) | UI |
| C3 | Dos reportes seguidos en la misma página | el segundo se envía sin error | UI |
| C4 | El GET JSON falla (500) | el wizard vuelve al inicio sin romperse y deja un aviso en consola | UI |
| C5 | Consumo del id | un id verificado en el servidor queda consumido (imagen de 0 bytes) y el 400 trae uno nuevo con imagen | API |
