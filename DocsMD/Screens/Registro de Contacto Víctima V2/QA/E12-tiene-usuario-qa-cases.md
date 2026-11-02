# QA — E-12 (Cuando responde "¿Tiene usuario?") y reportante con sesión (WP-1 T8, REP-07)

Ejecutado 2026-11-02 con `requests` (API), `psql` y Playwright (UI) contra la app del sandbox, con los usuarios `cloud.test.<rol>`.

| Caso | Descripción | Resultado esperado | Método |
|---|---|---|---|
| C1 | POST anónimo | 200; `victim_contact_general_user` vacío y `victim_contact_reporter_role` NULL | API |
| C2 | POST con sesión en, op y et | usuario y rol guardados | API |
| C3 | POST con sesión sv, ad y ro | también quedan asociados (usuario y rol); el wizard arranca en la autorización, sin el mensaje inicial | API + UI |
| C4 | Payload con `generalUser` falsificado | se ignora: anónimo → vacío; con sesión op → vale la sesión | API |
| C5 | Cookie inválida o sesión cerrada | anónimo, sin error | API |
| C6 | `GET /salvia/public-v2/nuevo` con sesión en, op y et, y con sv | con en, op y et la página trae nombre y rol; con sv no hay mensaje y arranca en la autorización; el modo JSON no cambia | API |
| C7 | Login con `next=/salvia/public-v2/nuevo` (roles en y op) | vuelve al reporte | UI |
| C8 | `next` malicioso (`//evil.com`, `https://…`, `/seguridad/logout`, `/salvia/casos`, `javascript:…`, `/salvia/public-v2/nuevo/../casos`) | se ignora; aterriza en su pantalla por rol | UI |
| C9 | Wizard con sesión | muestra «Está reportando como … (Enlace territorial)» con «Continuar»; el envío sigue exigiendo captcha (captcha incorrecto → 400, sin radicar) | UI |
| C10 | «No, continuar sin usuario» | flujo anónimo completo; «Volver al Inicio» regresa al paso inicial | UI |
| C11 | Hito «Reporte recibido» | lleva como actor al usuario con sesión | API |

## Copiar radicado (T3 · E-13)

Los casos C31 y C34 de [`E10-cuando-envia-reporte-qa-cases.md`](E10-cuando-envia-reporte-qa-cases.md) cubren la confirmación con el radicado y el botón (portapapeles permitido y denegado).
