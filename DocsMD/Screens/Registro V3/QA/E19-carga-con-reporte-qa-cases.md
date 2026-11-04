# QA — E19 (Cuando carga la pantalla con un registro de origen) y la conversión desde un reporte (WP-1 T7)

Ejecutado 2026-11-02 con `requests`, `psql` y Playwright contra la app del sandbox. La conversión vive en `POST /api/v1/victim-cases-v3/desde-reporte/:icode`; el detalle del modal en [`../../../Componentes/Gestion Reporte Modal/QA/M10-M12-registro-de-caso-qa-cases.md`](../../../Componentes/Gestion%20Reporte%20Modal/QA/M10-M12-registro-de-caso-qa-cases.md).

| Caso | Qué se prueba | Resultado esperado | Método |
|---|---|---|---|
| C8 | `GET /api/v1/victim-cases-v3/{caseIcode}` tras convertir | CI_2, CI_4, HV_1 y CI_8 en `answeredKeys` (CI_7 no); `answers` trae nombre completo, teléfono, hechos y el ajuste | API |
| C14a | V3 con `?caso=<icode>` | abre directo el formulario (sin la autorización) con el nombre y el teléfono prellenados | UI |
| C14b | V3 sin parámetro | el flujo normal arranca en «Autorización de datos personales» | UI |
| C14c | `?caso=no-es-un-uuid` | se ignora; flujo normal | UI |
| C13 | «Ir al borrador del caso» desde el modal | abre el mismo borrador en el formulario | UI |
