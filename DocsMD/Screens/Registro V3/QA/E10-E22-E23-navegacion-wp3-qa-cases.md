# QA — Registro V3 · Menú de módulos, flechas, avance y porcentaje (WP-3 · REG-01, REG-02)

Casos de [`flow-E10`](../Flujos/flow-E10-cuando-cambia-de-modulo.md), [`flow-E22`](../Flujos/flow-E22-cuando-presiona-siguiente.md), [`flow-E23`](../Flujos/flow-E23-cuando-presiona-anterior.md) y de la parte reactiva de [`flow-E12`](../Flujos/flow-E12-cuando-guarda-una-seccion.md).

> **Estado: ejecutado el 2026-11-02** — `scripts/e2e/wp3-navegacion-encabezado.js`, **29 de 29 ✅** (las comprobaciones T7, T8 y T9 de ese script están en [`E04-encabezado-estados-wp3-qa-cases.md`](E04-encabezado-estados-wp3-qa-cases.md)).

**Método.** Solo pantalla (Playwright): la API no tiene nada de navegación. Se ejecuta un recorrido real de un Borrador nuevo hasta completar la tríada guardando por la pantalla, y se leen menú, barra de navegación y encabezado del DOM.

## Casos de prueba

| Caso | REG | Descripción | Condición inicial | Resultado esperado | Método |
|---|---|---|---|---|---|
| T10a | 01 | El Caso se abre en CI | Borrador nuevo | Ítem «Contacto e identificación» activo en el menú | DOM |
| T10b | 01 | Menú construido con `modulos[]` | — | Un ítem por módulo (7), mismo orden que `modulos` | DOM |
| T10l | 01 | Contrato con WP-4: un módulo agregado a `modulos[]` | Se agrega `{id:'ZZ'}` desde la consola | Aparece solo en el menú (8 ítems) | DOM |
| T10c | 01 | Guardar con éxito avanza | CI completo | Pasa a «Datos personales» y muestra «Sección guardada. Pasaste de… a…» | DOM |
| T10i, T10j | 01 | Avance DP → CA, HV → AG | Módulos guardados por pantalla | Cada guardado pasa al siguiente módulo | DOM |
| T10k | 01 | Al completar la tríada, guardar AG avanza a Tamizaje | Tríada completa | El siguiente de AG es TZ (ya desbloqueado) | DOM |
| T10d | 01 | «Siguiente» salta lo bloqueado | Tríada incompleta, en AG | Va a Gestión Salvia (TZ bloqueado) | DOM |
| T10e / T10g | 01 | Extremos | En GS y en CI | «Siguiente» deshabilitado en GS; «Anterior» deshabilitado en CI | DOM |
| T10f | 01 | «Anterior» salta lo bloqueado | En GS, TZ bloqueado | Va a AG | DOM |
| T10h | 01 | Navegar no guarda | Valor escrito en CI, navegar y volver | 0 POST a `/seccion` y el valor sigue en pantalla | Red + DOM |
| T10m / T10n | 01 | Maquetación | 1280 px y 375 px | Menú a la izquierda en escritorio; apilado arriba en móvil | Geometría |
| T11a / T11b | 02 | Borrador vacío | — | CI, DP, HV, AG y GS en «0 %»; CA y TZ en «—»; avance global 0 % | DOM |
| T11c | 02 | Un condicional que aparece sube el total | CI_7 = Sí (aparece CI_8) | CI baja a «89 %» (8 de 9) | DOM |
| T11d | 02 | Cuadra con el recordatorio | CI_8 sin responder, guardar | El mensaje lista 1 faltante = total − respondidas | DOM |
| T11e | 02 | CI completo | CI_8 oculta | «100 %» | DOM |
| T11f | 02 | Con la tríada completa | — | Avance global > 0 y sin errores de consola | DOM |

(Las comprobaciones de bloqueo solo del Tamizaje — REG-13 — están en [`E12-E15-triada-…`](E12-E15-triada-guardado-y-cierre-wp3-qa-cases.md), casos U1a–d y N4.)

## Cobertura

| Flujo de prueba | Casos cubiertos |
|---|---|
| `wp3-navegacion-encabezado.js` | T10a–n, T11a–f |

## Resultados de ejecución

**Fecha:** 2026-11-02 · **Resultado:** 29/29 ✅ (incluye T7, T8 y T9). Hallazgo durante la prueba: Vue ignora los `<style>` dentro del template (la página se monta sobre `<body>`), por eso los estilos del menú y de `v3-msg-error` se inyectan desde `registro-v3-app.js`.
