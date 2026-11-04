# QA — E02 (Cuando presiona «Ingresar») · regreso con `next` (WP-1 T8)

Ejecutado 2026-11-02 con Playwright contra la app del sandbox (el POST del login usa `User-Agent: flutter-client`, que omite el captcha).

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C1 | Login con `?next=/salvia/public-v2/nuevo` (roles en y op) | vuelve al reporte V2 | UI |
| C2 | Login sin `next` | aterriza en la pantalla por rol (en → `/salvia/mis-barreras`) | UI |
| C3 | `next` malicioso: `//evil.com`, `https://evil.com/x`, `/seguridad/logout`, `/salvia/casos`, `javascript:alert(1)`, `/salvia/public-v2/nuevo/../casos` | se ignora; aterriza en su pantalla por rol | UI |
| C4 | Smoke de los 19 roles (`cloudDevelopment/smoke-login.js`) | todos inician sesión y aterrizan sin 5xx | UI |
