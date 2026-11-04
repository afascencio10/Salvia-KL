# QA — CR-E01, CR-E02, CR-E04 y la interfaz de CR-E03

Ejecutado 2026-11-02 con Playwright (Chromium) contra la app del sandbox.

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C13 | Máscara del campo y botón | «k7qm4r2x» pasa a K7QM-4R2X; con 7 símbolos o con I, O, 0 o 1 el botón sigue deshabilitado; recorta a 8 y descarta signos | UI |
| C14 | Resultado, error y «Hacer otra consulta» | la tarjeta trae radicado, fecha, estado (badge `text-bg-*`) y descripción; un radicado inexistente muestra la alerta genérica que se limpia al escribir; «Hacer otra consulta» limpia el formulario | UI |
| C15 | Bloqueo por 429 (respuesta simulada de 3 s) | mensaje con la espera y botón deshabilitado hasta que pase el tiempo, aunque el formato sea válido | UI |
| — | Error del servidor (500) | «No fue posible consultar en este momento. Intente de nuevo.» | UI |
| — | Carga | título «Consultar mi reporte», foco en el campo, botón deshabilitado, enlace del pie al wizard V2, sin errores JS | UI |
