# `gestion-reporte-modal` — Índice

Componente **nuevo** — concentra toda la máquina de estados de "Gestión reportes" (validación previa, esquema 3x3, cierre). Ver diseño completo, tomado literal de `Mockup_Reportes_Registro.html`, en [`gestion-reporte-modal-interface.md`](gestion-reporte-modal-interface.md).

> **Estado:** ✅ Implementado y probado — 2026-10-13. 155 pruebas de Playwright en verde (ver [`QA/`](../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md)); resumen del cambio en [`changelogOct2026.md`](changelogOct2026.md). La apertura de caso (M-10) queda diseñada pero **diferida, fuera de alcance de esta fase** (instrucción del usuario, 2026-10-11 — coincide con la nota de FigJam); el backend rechaza `close_reason='ap'` con 400 mientras tanto.

## Archivos de este componente

| Archivo | Descripción |
|---|---|
| [`gestion-reporte-modal-interface.md`](gestion-reporte-modal-interface.md) | Árbol de interfaz — 4 pasos (validación previa, esquema 3x3, post-contacto, cerrado) |
| [`gestion-reporte-modal-usage.md`](gestion-reporte-modal-usage.md) | Método `open()`, evento `updated`, ejemplo de uso |
| [`gestion-reporte-modal-events.md`](gestion-reporte-modal-events.md) | Inventario de eventos (M-01 a M-11, M-10 diferido) |
| [`related-tables.md`](related-tables.md) | Puntero — mismas tablas que `Listado Reportes` |
| [`QA/M01-abre-modal-qa-cases.md`](QA/M01-abre-modal-qa-cases.md) | Casos de M-01 (detalle, contadores, paso inicial) |
| [`QA/M02-M04-validacion-previa-qa-cases.md`](QA/M02-M04-validacion-previa-qa-cases.md) | Casos de M-02 y M-04 |
| [`QA/M05-M07-esquema-3x3-qa-cases.md`](QA/M05-M07-esquema-3x3-qa-cases.md) | Casos de M-05, M-06 y M-07 (esquema 3x3) |
| [`QA/M08-M09-cierre-post-contacto-qa-cases.md`](QA/M08-M09-cierre-post-contacto-qa-cases.md) | Casos de M-08 y M-09 |
| [`QA/M11-cierra-modal-qa-cases.md`](QA/M11-cierra-modal-qa-cases.md) | Casos de M-11 |
| [`QA/gestion-reportes-qa-plan.md`](../../Screens/Listado%20Reportes/QA/gestion-reportes-qa-plan.md) | Plan de QA compartido: resultados, bugs encontrados y decisiones |
| [`changelogOct2026.md`](changelogOct2026.md) | Registro del cambio de octubre 2026 |

## Resumen de eventos

| # | Evento | Flujo |
|---|---|---|
| M-01 | Cuando abre el modal | [flow-M01](./Flujos/flow-M01-cuando-abre-modal.md) |
| M-02 | Cuando presiona "No información mínima" | [flow-M02](./Flujos/flow-M02-cuando-cierra-sin-informacion-minima.md) |
| ~~M-03~~ | ~~Cuando presiona "Reporte de prueba"~~ — ❌ removido | — (no se construye, ver `gestion-reporte-modal-events.md`) |
| M-04 | Cuando presiona "Procede recontacto" | [flow-M04](./Flujos/flow-M04-cuando-procede-recontacto.md) |
| M-05 | Cuando registra un intento fallido | [flow-M05](./Flujos/flow-M05-cuando-registra-intento-fallido.md) |
| M-06 | Cuando registra un intento exitoso | [flow-M06](./Flujos/flow-M06-cuando-registra-intento-exitoso.md) |
| M-07 | Cuando presiona "Cerrar por no contacto" | [flow-M07](./Flujos/flow-M07-cuando-cierra-por-no-contacto.md) |
| M-08 | Cuando confirma cierre por no competencia | [flow-M08](./Flujos/flow-M08-cuando-confirma-cierre-no-competencia.md) |
| M-09 | Cuando confirma cierre por decisión de la persona | [flow-M09](./Flujos/flow-M09-cuando-confirma-cierre-decision-persona.md) |
| ~~M-10~~ | ~~Cuando confirma apertura de caso~~ — ⏸️ diferido | [flow-M10](./Flujos/flow-M10-cuando-confirma-apertura-caso.md) (diseño conservado) |
| M-11 | Cuando cierra el modal | [flow-M11](./Flujos/flow-M11-cuando-cierra-modal.md) |

## Checklist de completitud

- [x] ¿El ciclo de vida inicial (carga del detalle) está cubierto? → M-01
- [x] ¿Toda rama de cierre previo al contacto está cubierta? → M-02, M-03
- [x] ¿El tránsito a esquema 3x3 está cubierto? → M-04
- [x] ¿Registro de intentos (ambos resultados) está cubierto? → M-05, M-06
- [x] ¿Cierre tras completar 3x3 está cubierto? → M-07
- [x] ¿Toda rama post-contacto activa en esta fase está cubierta? → M-08, M-09 (M-10 diferido)
- [x] ¿Cierre del modal y aviso al padre está cubierto? → M-11
- [x] ¿Cambios de BD identificados? → `related-tables.md` (compartido con `Listado Reportes`)
- [x] ¿GAPs resueltos? → GAP-2 resuelto (2026-10-11). GAP-1 (orígenes del reporte) cubierto por la implementación: el detalle lee `form1`, `form2` y `form2_v2` según existan; confirmar con negocio si existe un 4.º origen. GAP-6/GAP-7 diferidos junto con M-10.
